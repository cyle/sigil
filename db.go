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
	for i := 0; i < 10; i++ {
		tmpNode := Node{ i, "Node "+fmt.Sprintf("%d", i), nil, nil }
		theData.nodes = append(theData.nodes, tmpNode)
	}
	
	// create some dummy connections!
	connOne := Connection{ 0, "Node 1 to 2", 1, 2 }
	connTwo := Connection{ 1, "Node 2 to 4", 2, 4 }
	connThree := Connection{ 2, "Node 4 to 5", 4, 5 }
	connFour := Connection{ 3, "Node 2 to 3", 2, 3 }
	connFive := Connection{ 4, "Node 3 to 5", 3, 5 }
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
	getNodeHandler gorest.EndPoint `method:"GET" path:"/node/{Id:int}" output:"Node"`
	postNodeHandler gorest.EndPoint `method:"POST" path:"/node" postdata:"Node"`
	deleteNodeHandler gorest.EndPoint `method:"DELETE" path:"/node/{Id:int}"`
	
	// connections stuff
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
	n.Id = len(theData.nodes)
	theData.nodes = append(theData.nodes, n)
	serv.ResponseBuilder().SetResponseCode(200)
	return
}

func (serv GraphService) DeleteNodeHandler(Id int) {
	serv.ResponseBuilder().SetResponseCode(200)
	return
}

/*

	connection functions

*/

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
	c.Id = len(theData.connections)
	theData.connections = append(theData.connections, c)
	serv.ResponseBuilder().SetResponseCode(200)
	return
}

func (serv GraphService) DeleteConnectionHandler(Id int) {
	serv.ResponseBuilder().SetResponseCode(200)
	return
}
