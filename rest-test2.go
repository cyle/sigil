//
// using go-rest: https://github.com/googollee/go-rest
//

package main

import (
	"net/http"
	"testing"
	"time"
	"github.com/googollee/go-rest"
)

type RestExample struct {
	rest.Service `prefix:"/cyle" mime:"application/json" charset:"utf-8"`

	CreateHello rest.Processor `method:"POST" path:"/hello"`
	GetHello    rest.Processor `method:"GET" path:"/hello/:to" func:"HandleHello"`
	Watch       rest.Streaming `method:"GET" path:"/hello/:to/streaming"`

	post  map[string]string
	watch map[string]chan string
}

type HelloArg struct {
	To   string `json:"to"`
	Post string `json:"post"`
}

// Post example:
// > curl "http://127.0.0.1:8080/prefix/hello" -d '{"to":"rest", "post":"rest is powerful"}'
//
// No response
func (r RestExample) HandleCreateHello(arg HelloArg) {
	r.post[arg.To] = arg.Post
	c, ok := r.watch[arg.To]
	if ok {
		select {
		case c <- arg.Post:
		default:
		}
	}
}

// Get example:
// > curl "http://127.0.0.1:8080/prefix/hello/rest"
//
// Response:
//   {"to":"rest","post":"rest is powerful"}
func (r RestExample) HandleHello() HelloArg {
	to := r.Vars()["to"]
	post, ok := r.post[to]
	if !ok {
		r.Error(http.StatusNotFound, r.DetailError(2, "can't find hello to %s", to))
		return HelloArg{}
	}
	return HelloArg{
		To:   to,
		Post: post,
	}
}

// Streaming example:
// > curl "http://127.0.0.1:8080/prefix/hello/rest/streaming"
//
// It create a long-live connection and will receive post content "rest is powerful"
// when running post example.
func (r RestExample) HandleWatch(s rest.Stream) {
	to := r.Vars()["to"]
	if to == "" {
		r.Error(http.StatusBadRequest, r.DetailError(3, "need to"))
		return
	}
	r.WriteHeader(http.StatusOK)
	c := make(chan string)
	r.watch[to] = c
	for {
		post := <-c
		s.SetWriteDeadline(time.Now().Add(time.Second))
		err := s.Write(post)
		if err != nil {
			close(c)
			delete(r.watch, to)
			return
		}
	}
}

// in unit test file
func TestExample(t *testing.T) {
	instance := &RestExample{
		post:  make(map[string]string),
		watch: make(map[string]chan string),
	}

	instance.HandleCreateHello(HelloArg{
		To:   "rest",
		Post: "rest is powerful",
	})

	resp, err := rest.SetTest(instance, map[string]string{"to": "rest"}, nil)
	if err != nil {
		panic(err)
	}
	arg := instance.HandleHello()
	if resp.Code != http.StatusOK {
		t.Error("should return ok")
	}
	if arg.To != "rest" {
		t.Error("arg.To should be rest")
	}
	if arg.Post != "rest is powerful" {
		t.Error("arg.Post should be 'rest is powerful'")
	}

	resp, err = rest.SetTest(instance, map[string]string{"to": "123"}, nil)
	if err != nil {
		panic(err)
	}
	arg = instance.HandleHello()
	if resp.Code != http.StatusNotFound {
		t.Error("should return not found")
	}
}

// The usage of rest.
func ExampleRest() {
	instance := &RestExample{
		post:  make(map[string]string),
		watch: make(map[string]chan string),
	}
	rest, err := rest.New(instance)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":8080", rest)
}