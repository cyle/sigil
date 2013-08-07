# simple graph database

Version: 0.1

This is supposed to be just a very simple graph database, accessible via a REST interface.

## warning

This is still in very, very early development. It's not really usable yet in any substantial way.

## how to use

First, set up your Go working directory:

    mkdir simple-graph-database
    cd simple-graph-database

Next, set your $GOPATH to this working directory (you'll need this to install the `gorest` package later):

    export GOPATH=`pwd`

Next, create the necessary working directory subdirectories (based on the [how to write Go code](http://golang.org/doc/code.html#Workspaces) document):

    mkdir src pkg bin

Okay cool -- now clone in this repo to the `src/` directory:

    git clone git@github.com:cyle/simple-graph-database.git src/

Right now this primarily relies on the `gorest` third-party package. To install it, do:

    go get code.google.com/p/gorest

That should be it for installing things. Now to run this database:

    cd src/
    go run db.go

That's it... go to `http://localhost:8777/` to see if it's working. That'll be your API endpoint.

## API documentation

To use the graph database in an application, visit the [API doc](API.md).

## ideas / lists

documents:

- nodes are just JSON documents with a unique ID attribute
- connections are just JSON documents with a unique ID attribute, and unique SOURCE and TARGET attributes
- nodes cannot be connected to themselves

add/update/remove actions:

- add new node
- add new connection between two nodes
- update node based on ID
- update connection based on ID
- upsert connection
	- if connection already exists between the two nodes, just update it
	- if one does not exist, just insert it
- delete node based on ID or other matching expression
	- delete connections going to/from the deleted node
- delete connection based on ID or other matching expression

query actions:

- select node based on ID or expression
	- get node ID #1
	- get node with { name: "cyle" }
- select connection based on ID or expression
	- get connection ID #4
	- get connection with { type: "works with" }
- get all connections tied to node
	- get all connections for node ID #1
- get connections tied to node with matching expression
	- get connections for node ID #1 with { type: "works with" }
- get all nodes based on connection query
	- get node connected via { type: "works with" }
- get paths between X node and Y node
	- return whole path between nodes
	- need "max depth" parameter to limit ridiculousness
	- options for either "get all paths" or "get shortest path"
- get all nodes connected to node
	- list of nodes and connections they're based on

backend stuff:

- config file
	- whether to save to disk or just keep in memory
	- how often to save to disk
	- where to save to disk
- flush/save data to disk
- load from disk on startup, if file exists

## why?

I liked the ideas of neo4j but I hate Java. And I want to build something that needs a simple graph database component.
