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
	getNodeHandler gorest.EndPoint `method:"GET" path:"/node/{Id:int}" output:"Node"`
	getConnectionHandler gorest.EndPoint `method:"GET" path:"/connection/{Id:int}" output:"Connection"`
}

func(serv GraphService) RootHandler() string {
	return "Simple Graph Database, v0.1"
}

func(serv GraphService) GetNodeHandler(Id int) (n Node){

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

func(serv GraphService) GetConnectionHandler(Id int) (c Connection){
	
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
