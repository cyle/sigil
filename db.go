package main

import "fmt"
import "net/http"
import "encoding/json" // documentation: http://golang.org/pkg/encoding/json/

//
// using gorest: https://code.google.com/p/gorest/wiki/GettingStarted?tm=6
//
import "code.google.com/p/gorest"

func main() {
	fmt.Println("Oh dear, a graph database...")
	gorest.RegisterService(new(GraphService))
	http.Handle("/", gorest.Handle())    
	http.ListenAndServe(":8777", nil)
}

type Node struct {
	Id int
	Name string
	ExtraJSONBytes []byte
	ExtraJSON []interface{}
}

type Edge struct {
	Id int
	Name string
	Source int
	Target int
	ExtraJSON []byte
}

type GraphService struct{
    // service level config
    gorest.RestService    `root:"/" consumes:"application/json" produces:"application/json"`
	// define routes
    rootHandler gorest.EndPoint `method:"GET" path:"/" output:"string"`
	getNodeHandler gorest.EndPoint `method:"GET" path:"/node/{Id:int}" output:"Node"`
	getEdgeHandler gorest.EndPoint `method:"GET" path:"/edge/{Id:int}" output:"Edge"`
}

func(serv GraphService) RootHandler() string {
	return "oh hello"
}

func(serv GraphService) GetNodeHandler(Id int) (n Node){
	fmt.Printf("Asking for node ID: %d \n", Id)
	
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
	
	fmt.Printf("Giving %+v \n", n)
	return
    //serv.ResponseBuilder().SetResponseCode(404).Overide(true)  //Overide causes the entity returned by the method to be ignored. Other wise it would send back zeroed object
    //return
}

func(serv GraphService) GetEdgeHandler(Id int) (e Edge){
	fmt.Printf("Asking for edge ID: %d \n", Id)
	e.Id = Id
	e.Name = "Some Edge"
	fmt.Printf("Giving %+v \n", e)
	return
    //serv.ResponseBuilder().SetResponseCode(404).Overide(true)  //Overide causes the entity returned by the method to be ignored. Other wise it would send back zeroed object
    //return
}
