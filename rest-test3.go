//
// using go-json-rest: http://godoc.org/github.com/ant0ine/go-json-rest
//

package main

import (
        "github.com/ant0ine/go-json-rest"
        "net/http"
)

type User struct {
        Id   string
        Name string
}

func GetUser(w *rest.ResponseWriter, req *rest.Request) {
        user := User{
                Id:   req.PathParam("id"),
                Name: "Antoine",
        }
        w.WriteJson(&user)
}

func GetRoot(w *rest.ResponseWriter, req *rest.Request) {
	data := struct {
		Wut string
		Huh string
	}{
		"hello",
		"thar",
	}
	w.WriteJson( data )
}

func main() {
        handler := rest.ResourceHandler{}
        handler.SetRoutes(
                rest.Route{"GET", "/users/:id", GetUser},
				rest.Route{"GET", "/", GetRoot},
        )
        http.ListenAndServe(":8080", &handler)
}