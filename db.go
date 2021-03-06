/*

	cyle's SIGIL database, version 0.1

*/

package main

import "fmt"
import "math"
import "net/http"
import "os"
import "io/ioutil"
import "bufio"
import "runtime"
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
	Extra map[string]string
}

type Connection struct {
	Id int
	Name string
	Source int
	Target int
	Distance float64
	DistanceMultiplier float64
	Extra map[string]string
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
var theService GraphService

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
	gorest.RegisterService(&theService)
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
	
	// path/spatial stuff
	getDistanceBetweenNodesHandler gorest.EndPoint `method:"GET" path:"/distance/from/{Source:int}/to/{Target:int}" output:"string"`
	//getPathsBetweenNodes gorest.EndPoint `method:"GET" path:"/paths/from/{Source:int}/to/{Target:int}" output:"[][]Connection"`
	getAstarBetweenNodes gorest.EndPoint `method:"GET" path:"/shortest/from/{Source:int}/to/{Target:int}" output:"[]Connection"`
	getNearbyNodes gorest.EndPoint `method:"GET" path:"/nodes/nearby/{Source:int}/radius/{Radius:float64}" output:"[]Node"`
	getClosestNode gorest.EndPoint `method:"GET" path:"/node/closest/{Source:int}" output:"Node"`
	
	// database-wide operations
	saveDatabaseHandler gorest.EndPoint `method:"GET" path:"/save" output:"string"`
	deleteAllNodes gorest.EndPoint `method:"DELETE" path:"/nodes"`
	deleteAllConnections gorest.EndPoint `method:"DELETE" path:"/connections"`
	
	// get memory info
	memoryInfoHandler gorest.EndPoint `method:"GET" path:"/meminfo" output:"string"`
}

func (serv GraphService) RootHandler() string {
	return "SIGIL, v0.1"
}

func (serv GraphService) SaveDatabaseHandler() string {
	fmt.Println("Saving database to file")
	saveAllTheData();
	fmt.Println("Saved database to file")
	return "okay"
}

func (serv GraphService) DeleteAllNodes() {
	fmt.Println("Deleting all nodes...")
	theData.Nodes = make([]Node, 0)
	// also deletes all connections
	theService.DeleteAllConnections()
	return
}

func (serv GraphService) DeleteAllConnections() {
	fmt.Println("Deleting all connections...")
	theData.Connections = make([]Connection, 0)
	return
}

func (serv GraphService) MemoryInfoHandler() (memstring string) {
	memstats := new(runtime.MemStats)
	runtime.ReadMemStats(memstats)
	memstring = fmt.Sprintf("Bytes Allocated: %d; Bytes Reserved from System: %d\nMegabytes Allocated: %f; Megabytes Reserved from System: %f", memstats.Alloc, memstats.Sys, float32(memstats.Alloc)/1024/1024, float32(memstats.Sys)/1024/1024)
	return 
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
	
	if n.Id == 0 { // could not find it
		serv.ResponseBuilder().SetResponseCode(404).Overide(true)  //Overide causes the entity returned by the method to be ignored. Other wise it would send back zeroed object
	    return
	}
	
	fmt.Printf("Giving: %+v \n", n)
	
    return
}

func (serv GraphService) PostNodeHandler(n Node) {
	fmt.Printf("Just got: %+v \n", n)
	current_max_id := 0
	// check if this already exists. if so, update it.
	for key, value := range theData.Nodes {
		if value.Id == n.Id {
			fmt.Printf("Updating node ID %d \n", n.Id)
			theData.Nodes[key] = n
			serv.ResponseBuilder().SetResponseCode(200)
			return
		}
		if value.Id > current_max_id {
			current_max_id = value.Id
		}
	}
	// doesn't exist? create it.
	fmt.Println("Creating new node based on input")
	n.Id = current_max_id + 1
	theData.Nodes = append(theData.Nodes, n)
	//serv.ResponseBuilder().SetResponseCode(201)
	fmt.Printf("Created node ID %d \n", n.Id)
	serv.ResponseBuilder().Created("http://localhost:8777/node/"+fmt.Sprintf("%d", n.Id))
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
	for _, conn := range theData.Connections {
		if conn.Source == Id || conn.Target == Id {
			theData.Connections = deleteConnectionFromSlice(conn, theData.Connections)
		}
	}
	
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
	if c.Source == 0 || c.Target == 0 {
		fmt.Println("Cannot create connection where SOURCE or TARGET are zero")
		serv.ResponseBuilder().SetResponseCode(400).Overide(true)
		return
	}
	// check to see if connection already exists
	current_max_id := 0
	for key, value := range theData.Connections {
		if value.Id == c.Id {
			fmt.Printf("Updating connection ID %d \n", c.Id)
			c.Distance = getDistanceBetweenNodes(c.Source, c.Target) // update distance
			theData.Connections[key] = c
			serv.ResponseBuilder().SetResponseCode(200)
			return
		}
		if value.Id > current_max_id {
			current_max_id = value.Id
		}
	}
	// does not exist! create a new connection.
	fmt.Println("Creating new connection based on input")
	c.Id = current_max_id + 1
	c.Distance = getDistanceBetweenNodes(c.Source, c.Target) // make sure distance is set
	theData.Connections = append(theData.Connections, c)
	//serv.ResponseBuilder().SetResponseCode(201)
	fmt.Printf("Created connection ID %d \n", c.Id)
	serv.ResponseBuilder().Created("http://localhost:8777/connection/"+fmt.Sprintf("%d", c.Id))
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
					//newHScore := getDistanceBetweenNodes(tmpPathStep.NodeId, Target)
					//fmt.Printf("New HScore should be, from %d to %d: %f \n", tmpPathStep.NodeId, Target, newHScore)
					tmpPathStep.RecalcHScore(Target)
					tmpPathStep.RecalcFScore()
					//fmt.Printf("Added to open path list: %+v \n", tmpPathStep)
					openPathList = append(openPathList, tmpPathStep)
				} else {
					//fmt.Println("Already in open list, checking it out...")
					// this adjacent node is already in the open path list
					//fmt.Printf("Checking out: %+v \n", tmpPathStep)
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
	nextNodeId := 0
	for i := len(closedPathList) - 1; i >= 0; i-- {
		if i == len(closedPathList) - 1 || nextNodeId == closedPathList[i].NodeId {
			reverseFinalPath = append(reverseFinalPath, closedPathList[i].NodeId)
			nextNodeId = closedPathList[i].ParentNodeId
		}
		if closedPathList[i].NodeId == Source {
			break
		}
	}
	
	finalPath := []int{}
	for i := len(reverseFinalPath) - 1; i >= 0; i-- {
		finalPath = append(finalPath, reverseFinalPath[i])
	}
	
	fmt.Printf("Final list: %+v \n", finalPath)
	
	// get the connections that make up this path
	for i := 0; i < len(finalPath) - 1; i++ {
		for _, conn := range theData.Connections {
			if (conn.Source == finalPath[i] && conn.Target == finalPath[i+1]) || (conn.Source == finalPath[i+1] && conn.Target == finalPath[i]) {
				connections = append(connections, conn)
			}
		}
	}
	
	// make sure the destination is in the list. if not, 404
	foundTarget := false
	for _, conn := range connections {
		if conn.Target == Target || conn.Source == Target {
			foundTarget = true
		}
	}
	
	if foundTarget == false {
		serv.ResponseBuilder().SetResponseCode(404).Overide(true)
	}
	
	return
}

func (serv GraphService) GetNearbyNodes(Source int, Radius float64) (nodes []Node) {
	// get nodes that are within X of source node
	
	sourceNode := getNode(Source)
	boundMinX := float64(sourceNode.X) - Radius
	boundMinY := float64(sourceNode.Y) - Radius
	boundMaxX := float64(sourceNode.X) + Radius
	boundMaxY := float64(sourceNode.Y) + Radius
	
	for _, node := range theData.Nodes {
		if node.Id == Source {
			continue
		}
		if float64(node.X) < boundMaxX && float64(node.X) > boundMinX && float64(node.Y) < boundMaxY && float64(node.Y) > boundMinY {
			// it might be within the radius
			if getDistanceBetweenNodes(node.Id, Source) <= Radius {
				nodes = append(nodes, node)
			}
		}
	}
	
	return
}

func (serv GraphService) GetClosestNode(Source int) (n Node) {
	// get nodes that are within X of source node
	
	var shortestDistance float64
	
	for _, node := range theData.Nodes {
		if node.Id == Source {
			continue
		}
		distance := getDistanceBetweenNodes(Source, node.Id)
		if distance < shortestDistance || shortestDistance == 0.0 {
			shortestDistance = distance
			n = node
		}
	}
	
	return
}

/*
func (serv GraphService) GetPathsBetweenNodes(Source int, Target int) (paths [][]Connection) {
	fmt.Printf("Get connection paths between nodes %d and %d \n", Source, Target)
	
	paths = make([][]Connection, 1)
	for i := range paths {
		paths[i] = make([]Connection, 1)
	}
	
	return paths
}
*/

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
}
