/*

	cyle's simple graph database, version 0.1

*/

package main

import "fmt"
import "net/http"
//import "encoding/json" // documentation: http://golang.org/pkg/encoding/json/

// using gorest: https://code.google.com/p/gorest/wiki/GettingStarted?tm=6
import "code.google.com/p/gorest"

type AllTheData struct {
	nodes []Node
	connections []Connection
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

var theData AllTheData

func main() {
	
	fmt.Println("Oh dear, a graph database...")
	
	// create some dummy nodes!
	for i := 1; i <= 10; i++ {
		tmpNode := Node{ i, "Node "+fmt.Sprintf("%d", i), nil, nil }
		theData.nodes = append(theData.nodes, tmpNode)
	}
	
	// create some dummy connections!
	connOne := Connection{ 1, "Node 1 to 2", 1, 2 }
	connTwo := Connection{ 2, "Node 2 to 4", 2, 4 }
	connThree := Connection{ 3, "Node 4 to 5", 4, 5 }
	connFour := Connection{ 4, "Node 2 to 3", 2, 3 }
	connFive := Connection{ 5, "Node 3 to 5", 3, 5 }
	// add connections to the big data pool
	theData.connections = append(theData.connections, connOne, connTwo, connThree, connFour, connFive)
	
	//fmt.Printf("%+v \n", theData)
	
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
}

func (serv GraphService) RootHandler() string {
	return "Simple Graph Database, v0.1"
}

/*

	node functions

*/

func (serv GraphService) GetNodesHandler() []Node {
	fmt.Println("Sending along current list of nodes")
	return theData.nodes
}

func (serv GraphService) GetNodeHandler(Id int) (n Node){

	fmt.Printf("Asking for node ID: %d \n", Id)
	
	for _, value := range theData.nodes {
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
	for key, value := range theData.nodes {
		if value.Id == n.Id {
			fmt.Printf("Updating node ID %d \n", n.Id)
			theData.nodes[key] = n
			serv.ResponseBuilder().SetResponseCode(200)
			return
		}
	}
	// doesn't exist? create it.
	fmt.Println("Creating new node based on input")
	n.Id = len(theData.nodes) + 1 // +1 because it's 1-based instead of 0-based
	theData.nodes = append(theData.nodes, n)
	serv.ResponseBuilder().SetResponseCode(201)
	return
}

func (serv GraphService) DeleteNodeHandler(Id int) {
	fmt.Printf("Trying to delete node ID %d \n", Id)
	thekey := -1
	for key, value := range theData.nodes {
		if value.Id == Id {
			thekey = key
		}
	}
	// look at all of this bullshit we have to do because of memory management
	if thekey > -1 {
		//fmt.Printf("Found the node to delete: %d \n", thekey)
		var tmpWhatever []Node
		if thekey == 0 {
			tmpWhatever = make([]Node, len(theData.nodes) - 1)
			lastPartOfSlice := theData.nodes[1:] // copy everything AFTER the node
			for _, value := range lastPartOfSlice {
				//fmt.Printf("Copying node: %+v \n", value)
				tmpWhatever = append(tmpWhatever, value)
			}
		} else {
			tmpWhatever = make([]Node, thekey)
			firstPartOfSlice := theData.nodes[:thekey]
			copy(tmpWhatever, firstPartOfSlice) // copy everything BEFORE the node
			//fmt.Printf("Nodes so far: %+v \n", tmpWhatever)
			theNextKey := thekey + 1
			lastPartOfSlice := theData.nodes[theNextKey:] // copy everything AFTER the node
			for _, value := range lastPartOfSlice {
				//fmt.Printf("Copying node: %+v \n", value)
				tmpWhatever = append(tmpWhatever, value)
			}
		}
		//fmt.Printf("Nodes so far: %+v \n", tmpWhatever)
		theData.nodes = tmpWhatever
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
		
	for _, conn := range theData.connections {
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
	return theData.connections
}

func (serv GraphService) GetConnectionHandler(Id int) (c Connection){
	
	fmt.Printf("Asking for connection ID: %d \n", Id)
	
	for _, value := range theData.connections {
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
	for key, value := range theData.connections {
		if value.Id == c.Id {
			fmt.Printf("Updating connection ID %d \n", c.Id)
			theData.connections[key] = c
			serv.ResponseBuilder().SetResponseCode(200)
			return
		}
	}
	// does not exist! create a new connection.
	fmt.Println("Creating new connection based on input")
	c.Id = len(theData.connections) + 1 // +1 because it's 1-based instead of 0-based
	theData.connections = append(theData.connections, c)
	serv.ResponseBuilder().SetResponseCode(201)
	return
}

func (serv GraphService) DeleteConnectionHandler(Id int) {
	fmt.Printf("Trying to delete connection ID %d", Id)
	thekey := -1
	for key, value := range theData.connections {
		if value.Id == Id {
			thekey = key
		}
	}
	if thekey > -1 {
		var tmpWhatever []Connection
		if thekey == 0 {
			tmpWhatever = make([]Connection, len(theData.connections) - 1)
			lastPartOfSlice := theData.connections[1:] // copy everything AFTER
			for _, value := range lastPartOfSlice {
				tmpWhatever = append(tmpWhatever, value)
			}
		} else {
			tmpWhatever = make([]Connection, thekey)
			firstPartOfSlice := theData.connections[:thekey]
			copy(tmpWhatever, firstPartOfSlice) // copy everything BEFORE
			theNextKey := thekey + 1
			lastPartOfSlice := theData.connections[theNextKey:] // copy everything AFTER
			for _, value := range lastPartOfSlice {
				tmpWhatever = append(tmpWhatever, value)
			}
		}
		theData.connections = tmpWhatever
		fmt.Println("Connection deleted")
	} else {
		fmt.Println("Could not find that connection ID to delete, weird")
	}
	serv.ResponseBuilder().SetResponseCode(200)
	return
}
