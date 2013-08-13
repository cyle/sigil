/*

	cyle's simple graph database, version 0.1

*/

package main

import "fmt"
import "math"
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
	X int
	Y int
	Z int
	ExtraJSONBytes []byte
	ExtraJSON []interface{}
}

type Connection struct {
	Id int
	Name string
	Source int
	Target int
	Distance float64
	DistanceMultiplier float64
}

type PathStep struct {
	GScore float64
	HScore float64
	FScore float64
	NodeId int
	ParentNodeId int
}

func (p *PathStep) GetAdjacentNodes() []Node {
	return getAdjacentNodes(p.NodeId)
}

func (p *PathStep) RecalcFScore() {
	p.FScore = p.GScore + p.HScore
}

func (p *PathStep) RecalcHScore( destNodeId int ) {
	p.HScore = getDistanceBetweenNodes(p.NodeId, destNodeId)
}

func getLowestFScore(path []PathStep) (step PathStep) {
	notSetYet := true
	lastKey := 0
	for stepKey, tmpStep := range path {
		if notSetYet || tmpStep.FScore < step.FScore {
			step = tmpStep
			if notSetYet { notSetYet = false }
		} else if tmpStep.FScore == step.FScore {
			if stepKey > lastKey {
				lastKey = stepKey
				step = tmpStep
			} else {
				lastKey = stepKey
			}
		}
	}
	return
}

func doesPathExistAlready(needle PathStep, haystack []PathStep) bool {
	for _, val := range haystack {
		if needle.NodeId == val.NodeId {
			return true
		}
	}
	return false
}

func removeFromPath(needle PathStep, haystack []PathStep) (newPath []PathStep) {
	thekey := -1
	for key, value := range haystack {
		if value == needle {
			thekey = key
		}
	}
	// look at all of this bullshit we have to do because of memory management
	if thekey > -1 {
		if thekey == 0 {
			lastPartOfSlice := haystack[1:] // copy everything AFTER the node
			for _, value := range lastPartOfSlice {
				newPath = append(newPath, value)
			}
		} else {
			firstPartOfSlice := haystack[:thekey]
			copy(newPath, firstPartOfSlice) // copy everything BEFORE the node
			theNextKey := thekey + 1
			lastPartOfSlice := haystack[theNextKey:] // copy everything AFTER the node
			for _, value := range lastPartOfSlice {
				newPath = append(newPath, value)
			}
		}
		return newPath
	} else {
		// not found -- send haystack back unchanged
		return haystack
	}
}

var db_filename string = "ALLTHEDATA.json"
var theData AllTheData

func main() {
		
	fmt.Println("Oh dear, a graph database...")
	
	// if the database file exists, load it
	check, _ := doesFileExist(db_filename);
	if check {
		fmt.Println("Loading data from file...")
		loadAllTheData()
		fmt.Println("Loaded.")
	} else {
		// wat do if no data? leave it all blank, i guess
		fmt.Println("No data, providing blank database.")
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
	getAdjacentNodesHandler gorest.EndPoint `method:"GET" path:"/node/{Id:int}/adjacent" output:"[]Node"`
	
	// connections stuff
	getConnectionsHandler gorest.EndPoint `method:"GET" path:"/connections" output:"[]Connection"`
	getConnectionHandler gorest.EndPoint `method:"GET" path:"/connection/{Id:int}" output:"Connection"`
	postConnectionHandler gorest.EndPoint `method:"POST" path:"/connection" postdata:"Connection"`
	deleteConnectionHandler gorest.EndPoint `method:"DELETE" path:"/connection/{Id:int}"`
	
	// paths stuff
	getDistanceBetweenNodesHandler gorest.EndPoint `method:"GET" path:"/distance/from/{Source:int}/to/{Target:int}" output:"string"`
	getAstarBetweenNodes gorest.EndPoint `method:"GET" path:"/astar/from/{Source:int}/to/{Target:int}" output:"[]Connection"`
	getPathBetweenNodes gorest.EndPoint `method:"GET" path:"/path/from/{Source:int}/to/{Target:int}" output:"[]Connection"`
	//getPathsBetweenNodes gorest.EndPoint `method:"GET" path:"/paths/from/{Source:int}/to/{Target:int}" output:"[][]Connection"`
	getShortestPathBetweenNodes gorest.EndPoint `method:"GET" path:"/shortest/from/{Source:int}/to/{Target:int}" output:"[]Connection"`
	
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

func getNode(Id int) (n Node) {
	// go through all the nodes and return the one with the given ID
	for _, value := range theData.Nodes {
		if value.Id == Id {
			n = value
		}
	}
	return
}

func (serv GraphService) GetNodeHandler(Id int) (n Node){
	
	fmt.Printf("Asking for node ID: %d \n", Id)
	
	n = getNode(Id)
	
	fmt.Printf("Giving: %+v \n", n)
	
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
	tmpNode := Node{}
	for _, value := range theData.Nodes {
		if value.Id == Id {
			tmpNode = value
			break
		}
	}
	theData.Nodes = deleteNodeFromSlice(tmpNode, theData.Nodes)
	fmt.Println("Node deleted")
	
	
	// also delete any connections that were connected to the node
	
	
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
		return
	} else {
		// could not find any! send 404
	    serv.ResponseBuilder().SetResponseCode(404).Overide(true)  //Overide causes the entity returned by the method to be ignored. Other wise it would send back zeroed object
	    return
	}
	
}

func getAdjacentNodes(Id int) (nodes []Node) {
	// go through connections, get ones that are connected to this node
	for _, conn := range theData.Connections {
		if conn.Source == Id {
			nodes = append(nodes, getNode(conn.Target))
		} else if conn.Target == Id {
			nodes = append(nodes, getNode(conn.Source))
		}
	}
	return
}

func (serv GraphService) GetAdjacentNodesHandler(Id int) (nodes []Node) {
	
	nodes = getAdjacentNodes(Id);
	
	if len(nodes) > 0 {
		return
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
	tmpConn := Connection{}
	for _, value := range theData.Connections {
		if value.Id == Id {
			tmpConn = value
			break
		}
	}
	theData.Connections = deleteConnectionFromSlice(tmpConn, theData.Connections)
	fmt.Println("Deleted connection")
	serv.ResponseBuilder().SetResponseCode(200)
	return
}

/*

	paths functions

*/

func getDistanceBetweenNodes(Source int, Target int) (distance float64) {

	if Source == Target {
		return 0.0
	}
	
	var x1 int
	var y1 int
	var z1 int
	
	var x2 int
	var y2 int
	var z2 int
	
	// get X and Y of the Source node
	// get X and Y of the Target node
	for _, node := range theData.Nodes {
		if (node.Id == Source) {
			x1 = node.X
			y1 = node.Y
			z1 = node.Z
		} else if (node.Id == Target) {
			x2 = node.X
			y2 = node.Y
			z2 = node.Z
		}
	}
	
	// return distance as a string
	xD := x2 - x1
	yD := y2 - y1
	zD := z2 - z1	
	toSqrt := float64((xD * xD) + (yD * yD) + (zD * zD))	
	distance = math.Sqrt( toSqrt )
	
	return
}

func (serv GraphService) GetDistanceBetweenNodesHandler(Source int, Target int) string {
	
	fmt.Printf("Getting raw distance between nodes %d and %d \n", Source, Target)
	
	distance := getDistanceBetweenNodes(Source, Target)
	
	fmt.Printf("Raw distance between nodes %d and %d is %f \n", Source, Target, distance)
	
	return fmt.Sprintf("%f", distance)
	
}

func (serv GraphService) GetAstarBetweenNodes(Source int, Target int) (connections []Connection) {
	fmt.Printf("Using A* to get a connection path between nodes %d and %d \n", Source, Target)
	
	// references:
	//   http://www.raywenderlich.com/4946/introduction-to-a-pathfinding
	//   http://www.raywenderlich.com/4970/how-to-implement-a-pathfinding-with-cocos2d-tutorial
	//   http://theory.stanford.edu/~amitp/GameProgramming/
	
	// make:
	// open list - all nodes being considered for the path
	openPathList := make([]PathStep, 1)
	// closed list -- all the nodes definitely not to consider again
	closedPathList := []PathStep{}
	
	// current node goes into the closed list, of course
	// all nodes connected to current node goes into the open list
	firstPathStep := PathStep{}
	firstPathStep.NodeId = Source
	openPathList[0] = firstPathStep
	
	lastPathStep := PathStep{}
	lastPathStep.NodeId = Target
	
	// each node's score is F, which is G + H
	// G is the distance from the current node (always current G score + 1)
	// H is the (estimated) distance to the destination node
	
	// the loop:
	// get the node in the open list that has the lowest score.
	//   what if there are more than one? take the most recent one added
	// remove that node from the open list and add it to the closed list
	// if the node we just added to the closed list is the destination, then we're done!
	// if not...
	// for each node connected to that node:
	//   if it's already in the closed list, ignore it
	//   if it's not in the open list, add it to the open list and compute its F score
	//   if it's already in the open list, 
	//     check if its G score is lower than the current node's G score + 1
	//       if so, update its G score to be current node's G score + 1
	
	for len(openPathList) > 0 {
		currentPathStep := getLowestFScore(openPathList)
		//fmt.Printf("Current node ID: %d \n", currentPathStep.NodeId)
		closedPathList = append(closedPathList, currentPathStep)
		openPathList = removeFromPath(currentPathStep, openPathList)
		if doesPathExistAlready(lastPathStep, closedPathList) {
			// we just added the destination! we're done!
			//fmt.Println("Destination found! All done!")
			break
		} else {
			// get adjacent paths
			adjacentNodes := currentPathStep.GetAdjacentNodes();
			//fmt.Println("Going through adjacent nodes...")
			for _, tmpNode := range adjacentNodes {
				//fmt.Printf("Checking adjacent node ID: %d \n", tmpNode.Id)
				tmpPathStep := PathStep{}
				tmpPathStep.NodeId = tmpNode.Id
				tmpPathStep.ParentNodeId = currentPathStep.NodeId
				if doesPathExistAlready(tmpPathStep, closedPathList) {
					//fmt.Println("Already been here, continuing...")
					continue // keep going if we've already been there
				}
				if doesPathExistAlready(tmpPathStep, openPathList) == false {
					// this adjacent node is not yet in the open path list
					//fmt.Println("Not yet in open path list, adding...")
					tmpPathStep.GScore = currentPathStep.GScore + 1
					tmpPathStep.RecalcHScore(Target)
					tmpPathStep.RecalcFScore()
					openPathList = append(openPathList, tmpPathStep)
				} else {
					//fmt.Println("Already in open list, checking it out...")
					// this adjacent node is already in the open path list
					if tmpPathStep.GScore > currentPathStep.GScore + 1 {
						//fmt.Println("It's better, go there!")
						tmpPathStep.GScore = currentPathStep.GScore + 1
						tmpPathStep.RecalcFScore()
					}
				}
			}
		}
	}
	
	//fmt.Printf("Closed path list: %+v \n", closedPathList)
	//for _, val := range closedPathList {
	//	fmt.Printf("Going to %d from %d \n", val.NodeId, val.ParentNodeId)
	//}
	
	reverseFinalPath := []int{}
	for i := len(closedPathList) - 1; i >= 0; i-- {
		if i == len(closedPathList) - 1 {
			reverseFinalPath = append(reverseFinalPath, closedPathList[i].NodeId)
		}
		reverseFinalPath = append(reverseFinalPath, closedPathList[i].ParentNodeId)
		if closedPathList[i].ParentNodeId == Source {
			break
		}
	}
	
	finalPath := []int{}
	for i := len(reverseFinalPath) - 1; i >= 0; i-- {
		finalPath = append(finalPath, reverseFinalPath[i])
	}
	
	fmt.Printf("Final list: %+v \n", finalPath)
	
	// need:
	// open list with nodes and their H and G scores
	// AstarNode { properties: hScore, gScore, parentNode; methods: fScore, adjacentNodes, calcHScore }
	// OpenList is a slice of these AstarNodes
	// closedList is just a slice of Node ID ints
	// method to get node in open list with lowest F score (which is H + G)
	// method to get adjacent nodes (already have that)
	// method to get "parent" node in working set
	
	return
}

func (serv GraphService) GetPathBetweenNodes(Source int, Target int) (connections []Connection) {
	fmt.Printf("Get a connection path between nodes %d and %d \n", Source, Target)
	
	// trying: http://en.wikipedia.org/wiki/Graph_traversal
	
	tries := 0
	
	foundTarget := false
	
	nodeQueue := make([]int, 1)
	nodeQueue[0] = Source

	nodeMarked := make([]int, 1)
	nodeMarked[0] = Source
	
	for len(nodeQueue) != 0 {
		//fmt.Printf("Queue looks like: %+v \n", nodeQueue)
		//fmt.Printf("List of marked nodes looks like: %+v \n", nodeMarked)
		if tries > 100 { break }
		tmpNode := nodeQueue[0] // take first element
		nodeQueue = nodeQueue[1:] // remove first element
		//fmt.Printf("current node is %d \n", tmpNode)
		if tmpNode == Target {
			foundTarget = true
			break
		} else {
			//fmt.Println("not directly connected, getting connections...")
			for _, conn := range theData.Connections {
				if conn.Source == tmpNode || conn.Target == tmpNode {
					nextNode := 0
					if (conn.Source == tmpNode) {
						nextNode = conn.Target
					} else {
						nextNode = conn.Source
					}
					//fmt.Printf("seeing if %d is marked... \n", nextNode)
					if doesIntExist(nextNode, nodeMarked) == false {
						nodeMarked = append(nodeMarked, nextNode)
						nodeQueue = append(nodeQueue, nextNode)
						//fmt.Println("not marked, going deeper...")
					} else {
						//fmt.Println("marked, skipping along")
					}
				}
			}
		}
		tries += 1
	}
	
	if foundTarget {
		fmt.Printf("found the target, took %d iterations! \n", tries)
		fmt.Printf("final list of queued nodes: %+v \n", nodeQueue)
		fmt.Printf("final list of marked nodes: %+v \n", nodeMarked)
	} else {
		fmt.Println("could not find route to target within 200 iterations")
	}
	
	return
}

func (serv GraphService) GetPathsBetweenNodes(Source int, Target int) (paths [][]Connection) {
	fmt.Printf("Get connection paths between nodes %d and %d \n", Source, Target)
	
	paths = make([][]Connection, 1)
	for i := range paths {
		paths[i] = make([]Connection, 1)
	}
	
	return paths
}

func (serv GraphService) GetShortestPathBetweenNodes(Source int, Target int) (connections []Connection) {
	fmt.Printf("Get shortest connection path between nodes %d and %d \n", Source, Target)
	// run GetPathsBetweenNodes and just pick the shortest one and send along the list of connections
	
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

// doesFileExist returns whether the given file or directory exists or not
// from: http://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-denoted-by-a-path-exists-in-golang
func doesFileExist(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

// checks if int is inside of a slice
func doesIntExist(needle int, haystack []int) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}
	return false
}

// removes a string from a slice
func deleteStringFromSlice(needle string, haystack []string) (newSlice []string) {
	thekey := -1
	for key, value := range haystack {
		if value == needle {
			thekey = key
		}
	}
	if thekey > -1 {
		if thekey == 0 {
			lastPartOfSlice := haystack[1:] // copy everything AFTER
			for _, value := range lastPartOfSlice {
				newSlice = append(newSlice, value)
			}
		} else {
			firstPartOfSlice := haystack[:thekey]
			copy(newSlice, firstPartOfSlice) // copy everything BEFORE
			theNextKey := thekey + 1
			lastPartOfSlice := haystack[theNextKey:] // copy everything AFTER
			for _, value := range lastPartOfSlice {
				newSlice = append(newSlice, value)
			}
		}
		// deleted
		return newSlice
	} else {
		// could not find it in the slice, return original
		return haystack
	}
	return
}


// removes a node from a slice
func deleteNodeFromSlice(needle Node, haystack []Node) (newSlice []Node) {
	thekey := -1
	for key, value := range haystack {
		if value.Id == needle.Id {
			thekey = key
		}
	}
	if thekey > -1 {
		if thekey == 0 {
			lastPartOfSlice := haystack[1:] // copy everything AFTER
			for _, value := range lastPartOfSlice {
				newSlice = append(newSlice, value)
			}
		} else {
			firstPartOfSlice := haystack[:thekey]
			copy(newSlice, firstPartOfSlice) // copy everything BEFORE
			theNextKey := thekey + 1
			lastPartOfSlice := haystack[theNextKey:] // copy everything AFTER
			for _, value := range lastPartOfSlice {
				newSlice = append(newSlice, value)
			}
		}
		// deleted
		return newSlice
	} else {
		// could not find it in the slice, return original
		return haystack
	}
	return
}

// removes a connection from a slice
func deleteConnectionFromSlice(needle Connection, haystack []Connection) (newSlice []Connection) {
	thekey := -1
	for key, value := range haystack {
		if value == needle {
			thekey = key
		}
	}
	if thekey > -1 {
		if thekey == 0 {
			lastPartOfSlice := haystack[1:] // copy everything AFTER
			for _, value := range lastPartOfSlice {
				newSlice = append(newSlice, value)
			}
		} else {
			firstPartOfSlice := haystack[:thekey]
			copy(newSlice, firstPartOfSlice) // copy everything BEFORE
			theNextKey := thekey + 1
			lastPartOfSlice := haystack[theNextKey:] // copy everything AFTER
			for _, value := range lastPartOfSlice {
				newSlice = append(newSlice, value)
			}
		}
		// deleted
		return newSlice
	} else {
		// could not find it in the slice, return original
		return haystack
	}
}
