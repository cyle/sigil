# simple graph database

This is supposed to be just a very simple graph database, accessible via a REST interface.

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

backend stuff:

- flush/save data to disk
- load from disk on startup, if file exists

## why?

I liked the ideas of neo4j but I hate Java. And I want to build something that needs a simple graph database component.
