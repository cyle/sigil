package main

import "fmt"
import "net/http"

//
// using gorest: https://code.google.com/p/gorest/wiki/GettingStarted?tm=6
//
import "code.google.com/p/gorest"

func main() {
	fmt.Println("Oh dear, a graph database...")
	gorest.RegisterService(new(CyleService))
	http.Handle("/", gorest.Handle())    
	http.ListenAndServe(":8777", nil)
}

type Welp struct {
	Haha, Wut int
}

type CyleService struct{
    // service level config
    gorest.RestService    `root:"/" consumes:"application/json" produces:"application/json"`
    theDeets gorest.EndPoint `method:"GET" path:"/welp/{Id:int}" output:"Welp"`
	rootHandler gorest.EndPoint `method:"GET" path:"/" output:"string"`
}

func(serv CyleService) RootHandler() string {
	return "oh hello"
}

func(serv CyleService) TheDeets(Id int) (u Welp){
	fmt.Println("Asking for: ", Id)
	u.Haha = Id // as a reflection test
	u.Wut = 300 // as a random test to see what happens
	fmt.Println("Giving ", u)
	return
    //serv.ResponseBuilder().SetResponseCode(404).Overide(true)  //Overide causes the entity returned by the method to be ignored. Other wise it would send back zeroed object
    //return
}
