/*

	cyle's simple graph database, version 0.1

*/

package main

import "fmt"
import "net/http"
import "os"
import "io/ioutil"
import "bufio"
import "encoding/json" // documentation: http://golang.org/pkg/encoding/json/

// using gorest: https://code.google.com/p/gorest/wiki/GettingStarted?tm=6
import "code.google.com/p/gorest"

type AllTheData struct {
	Name string
	Nodes []Node
	Connections []Connection
}

type Node struct {
	Id int
	Name string
	ExtraJSONBytes []byte
	ExtraJSON []interface{}
}

type Connection struct {
	Id int
	Name string
	Source int
	Target int
}

var db_filename string = "ALLTHEDATA.json"
var theData AllTheData

func main() {
	
	theData.Name = "The Graph Database"
	
	fmt.Println("Oh dear, a graph database...")
	
	// if the database file exists, load it
	check, _ := doesFileExist(db_filename);
	if check {
		loadAllTheData()
	} else {
		// create some dummy nodes!
		for i := 1; i <= 10; i++ {
			tmpNode := Node{ i, "Node "+fmt.Sprintf("%d", i), nil, nil }
			theData.Nodes = append(theData.Nodes, tmpNode)
		}
		// create some dummy connections!
		connOne := Connection{ 1, "Node 1 to 2", 1, 2 }
		connTwo := Connection{ 2, "Node 2 to 3", 2, 3 }
		connThree := Connection{ 3, "Node 3 to 4", 3, 4 }
		connFour := Connection{ 4, "Node 4 to 5", 4, 5 }
		connFive := Connection{ 5, "Node 5 to 6", 5, 6 }
		connSix := Connection{ 5, "Node 3 to 9", 3, 9 }
		connSeven := Connection{ 5, "Node 9 to 8", 9, 8 }
		connEight := Connection{ 5, "Node 8 to 3", 8, 3 }
		connNine := Connection{ 5, "Node 3 to 7", 3, 7 }
		connTen := Connection{ 5, "Node 7 to 5", 7, 5 }
		// add connections to the big data pool
		theData.Connections = append(theData.Connections, connOne, connTwo, connThree, connFour, connFive)
		// save this dummy data for future use
		saveAllTheData()
	}
		
	// start the REST service to access the data
	gorest.RegisterService(new(GraphService))
	http.Handle("/", gorest.Handle())    
	http.ListenAndServe(":8777", nil)
}

type GraphService struct{
    // service level config
    gorest.RestService `root:"/" consumes:"application/json" produces:"application/json"`
	// define routes
	
	// deal with the root
    rootHandler gorest.EndPoint `method:"GET" path:"/" output:"string"`
	
	// node stuff
	getNodesHandler gorest.EndPoint `method:"GET" path:"/nodes" output:"[]Node"`
	getNodeHandler gorest.EndPoint `method:"GET" path:"/node/{Id:int}" output:"Node"`
	postNodeHandler gorest.EndPoint `method:"POST" path:"/node" postdata:"Node"`
	deleteNodeHandler gorest.EndPoint `method:"DELETE" path:"/node/{Id:int}"`
	getConnectionsForNodeHandler gorest.EndPoint `method:"GET" path:"/node/{Id:int}/connections" output:"[]Connection"`
	
	// connections stuff
	getConnectionsHandler gorest.EndPoint `method:"GET" path:"/connections" output:"[]Connection"`
	getConnectionHandler gorest.EndPoint `method:"GET" path:"/connection/{Id:int}" output:"Connection"`
	postConnectionHandler gorest.EndPoint `method:"POST" path:"/connection" postdata:"Connection"`
	deleteConnectionHandler gorest.EndPoint `method:"DELETE" path:"/connection/{Id:int}"`
	
	// save the database
	saveDatabaseHandler gorest.EndPoint `method:"GET" path:"/save" output:"string"`
}

func (serv GraphService) RootHandler() string {
	return "Simple Graph Database, v0.1"
}

func (serv GraphService) SaveDatabaseHandler() string {
	fmt.Println("Saving database to file")
	saveAllTheData();
	fmt.Println("Saved database to file")
	return "okay"
}

/*

	node functions

*/

func (serv GraphService) GetNodesHandler() []Node {
	fmt.Println("Sending along current list of nodes")
	return theData.Nodes
}

func (serv GraphService) GetNodeHandler(Id int) (n Node){

	fmt.Printf("Asking for node ID: %d \n", Id)
	
	for _, value := range theData.Nodes {
		if value.Id == Id {
			n = value
			fmt.Printf("Giving: %+v \n", n)
			return
		}
	}
	
	/*
	n.Id = Id
	n.Name = "Some Node"
	n.ExtraJSONBytes = []byte(`[{"Name": "Platypus", "Order": "Monotremata"}]`)
	// this technique of taking arbitrary JSON and turning it into something usable came from: http://blog.golang.org/json-and-go
	var tmp interface {}
	err := json.Unmarshal(n.ExtraJSONBytes, &tmp)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Printf("JSON: %T %+v \n", tmp, tmp)
	n.ExtraJSON = tmp.([]interface{})
	n.ExtraJSONBytes = nil
	*/
	
	// could not find it! send 404
    serv.ResponseBuilder().SetResponseCode(404).Overide(true)  //Overide causes the entity returned by the method to be ignored. Other wise it would send back zeroed object
    return
}

func (serv GraphService) PostNodeHandler(n Node) {
	fmt.Printf("Just got: %+v \n", n)
	// check if this already exists. if so, update it.
	for key, value := range theData.Nodes {
		if value.Id == n.Id {
			fmt.Printf("Updating node ID %d \n", n.Id)
			theData.Nodes[key] = n
			serv.ResponseBuilder().SetResponseCode(200)
			return
		}
	}
	// doesn't exist? create it.
	fmt.Println("Creating new node based on input")
	n.Id = len(theData.Nodes) + 1 // +1 because it's 1-based instead of 0-based
	theData.Nodes = append(theData.Nodes, n)
	serv.ResponseBuilder().SetResponseCode(201)
	return
}

func (serv GraphService) DeleteNodeHandler(Id int) {
	fmt.Printf("Trying to delete node ID %d \n", Id)
	thekey := -1
	for key, value := range theData.Nodes {
		if value.Id == Id {
			thekey = key
		}
	}
	// look at all of this bullshit we have to do because of memory management
	if thekey > -1 {
		//fmt.Printf("Found the node to delete: %d \n", thekey)
		var tmpWhatever []Node
		if thekey == 0 {
			tmpWhatever = make([]Node, len(theData.Nodes) - 1)
			lastPartOfSlice := theData.Nodes[1:] // copy everything AFTER the node
			for _, value := range lastPartOfSlice {
				//fmt.Printf("Copying node: %+v \n", value)
				tmpWhatever = append(tmpWhatever, value)
			}
		} else {
			tmpWhatever = make([]Node, thekey)
			firstPartOfSlice := theData.Nodes[:thekey]
			copy(tmpWhatever, firstPartOfSlice) // copy everything BEFORE the node
			//fmt.Printf("Nodes so far: %+v \n", tmpWhatever)
			theNextKey := thekey + 1
			lastPartOfSlice := theData.Nodes[theNextKey:] // copy everything AFTER the node
			for _, value := range lastPartOfSlice {
				//fmt.Printf("Copying node: %+v \n", value)
				tmpWhatever = append(tmpWhatever, value)
			}
		}
		//fmt.Printf("Nodes so far: %+v \n", tmpWhatever)
		theData.Nodes = tmpWhatever
		//fmt.Printf("Nodes should be copied now!\n")
		fmt.Println("Node deleted")
	} else {
		fmt.Println("Could not find that node ID to delete, weird")
	}
	serv.ResponseBuilder().SetResponseCode(200)
	return
}

func (serv GraphService) GetConnectionsForNodeHandler(Id int) (connections []Connection) {
	// get the connections attached to a given node based on the node's ID
	fmt.Printf("Asking for connections for node ID: %d \n", Id)
		
	for _, conn := range theData.Connections {
		if conn.Source == Id || conn.Target == Id {
			connections = append(connections, conn)
		}
	}
	
	if len(connections) > 0 {
		return connections
	} else {
		// could not find any! send 404
	    serv.ResponseBuilder().SetResponseCode(404).Overide(true)  //Overide causes the entity returned by the method to be ignored. Other wise it would send back zeroed object
	    return
	}
	
}

/*

	connection functions

*/

func (serv GraphService) GetConnectionsHandler() []Connection {
	fmt.Println("Sending along current list of connections")
	return theData.Connections
}

func (serv GraphService) GetConnectionHandler(Id int) (c Connection){
	
	fmt.Printf("Asking for connection ID: %d \n", Id)
	
	for _, value := range theData.Connections {
		if value.Id == Id {
			c = value
			fmt.Printf("Giving: %+v \n", c)
			return
		}
	}
	
	// could not find it! send 404
    serv.ResponseBuilder().SetResponseCode(404).Overide(true)  //Overide causes the entity returned by the method to be ignored. Other wise it would send back zeroed object
    return
}

func (serv GraphService) PostConnectionHandler(c Connection) {
	fmt.Printf("Just got: %+v \n", c)
	// make sure it's not invalid
	if c.Source == c.Target {
		fmt.Println("Cannot create connection where SOURCE and TARGET are the same")
		serv.ResponseBuilder().SetResponseCode(400).Overide(true)
		return
	}
	// check to see if connection already exists
	for key, value := range theData.Connections {
		if value.Id == c.Id {
			fmt.Printf("Updating connection ID %d \n", c.Id)
			theData.Connections[key] = c
			serv.ResponseBuilder().SetResponseCode(200)
			return
		}
	}
	// does not exist! create a new connection.
	fmt.Println("Creating new connection based on input")
	c.Id = len(theData.Connections) + 1 // +1 because it's 1-based instead of 0-based
	theData.Connections = append(theData.Connections, c)
	serv.ResponseBuilder().SetResponseCode(201)
	return
}

func (serv GraphService) DeleteConnectionHandler(Id int) {
	fmt.Printf("Trying to delete connection ID %d", Id)
	thekey := -1
	for key, value := range theData.Connections {
		if value.Id == Id {
			thekey = key
		}
	}
	if thekey > -1 {
		var tmpWhatever []Connection
		if thekey == 0 {
			tmpWhatever = make([]Connection, len(theData.Connections) - 1)
			lastPartOfSlice := theData.Connections[1:] // copy everything AFTER
			for _, value := range lastPartOfSlice {
				tmpWhatever = append(tmpWhatever, value)
			}
		} else {
			tmpWhatever = make([]Connection, thekey)
			firstPartOfSlice := theData.Connections[:thekey]
			copy(tmpWhatever, firstPartOfSlice) // copy everything BEFORE
			theNextKey := thekey + 1
			lastPartOfSlice := theData.Connections[theNextKey:] // copy everything AFTER
			for _, value := range lastPartOfSlice {
				tmpWhatever = append(tmpWhatever, value)
			}
		}
		theData.Connections = tmpWhatever
		fmt.Println("Connection deleted")
	} else {
		fmt.Println("Could not find that connection ID to delete, weird")
	}
	serv.ResponseBuilder().SetResponseCode(200)
	return
}

/*

	save and load the data

*/

func saveAllTheData() {
	// spit data out to JSON into a file
	// open output file
    fo, err := os.Create(db_filename)
    if err != nil { panic(err) }
    // close fo on exit and check for its returned error
    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }()
    // make a write buffer
    w := bufio.NewWriter(fo)
	allTheDataJSON, err := json.Marshal(theData)
	if err != nil { panic(err) }
	_, err = w.Write(allTheDataJSON)
	if (err != nil) { panic(err) }
	if err = w.Flush(); err != nil { panic(err) }
}

func loadAllTheData() {
	// ingest data via JSON from a file
	allJSON, err := ioutil.ReadFile(db_filename)
	if err != nil { panic(err) }
	unmarshal_err := json.Unmarshal(allJSON, &theData)
	if unmarshal_err != nil { panic(unmarshal_err) }
}

/*

	helper functions

*/

// exists returns whether the given file or directory exists or not
// from: http://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-denoted-by-a-path-exists-in-golang
func doesFileExist(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}