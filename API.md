# The SIGIL REST API

## Consider the Following

This document is probably not accurate while I'm developing the database.

## Result Types

There are two main result data types that will return from these functions as JSON objects.

### Node Objects

They look like this:

    {
	    Id (integer)
	    Name (string)
	    X (integer)
	    Y (integer)
	    Z (integer)
	    ExtraJSON (object)
    }

The attributes should be pretty self-explanatory. The `ExtraJSON` attribute can be an object of anything you want, if you'd like to attach additional data to the node.

## Connection Objects

They look like this:

    {
        Id (integer)
        Name (string)
        Source (integer)
        Target (integer)
        Distance (float)
        DistanceMultiplier (float)
        ExtraJSON (object)
    }

The attributes should be pretty self-explanatory. While connections do have *directed* links, they are never used as such within the database itself (at least, not yet). The `DistanceMultiplier` attribute is also not used yet. The `ExtraJSON` attribute can be an object of anything you want, if you'd like to attach additional data to the node.

## Querying

To get all nodes, simply send a GET request to `/nodes`.

To get all connections, simply send a GET request to `/connections`.

To get a specific node, simply send a GET request to `/node/:id`, where `:id` is the unique ID number of the node you're looking for.

To get a specific connection, simply send a GET request to `/connection/:id`, where `:id` is the unique ID number of the connection you're looking for.

To get all connections attached to a given node, simply send a GET request to `/node/:id/connections`, where `:id` is the unique ID number of the node you want connections for. This will return a list of connections.

To get all nodes adjacent/connected to a given node, simply send a GET request to `/node/:id/adjacent`, where `:id` is the unique ID number of the node you want adjacent nodes for. This will return a list of nodes.

To get the shortest path between two nodes, simply send a GET request to `/shortest/from/:source/to/:target`, where `:source` is the origin node unique ID and `:target` is the destination node unique ID. This will return a list of connections, from the source to the target.

To get the straight distance between two nodes, simply send a GET request to `/distance/from/:source/to/:target`, where `:source` is the origin node unique ID and `:target` is the destination node unique ID. This will return a floating point number.

## Creating

To create a new node, send a POST request to `/node` with a JSON object with as many of the above Node-type attributes as you want. Any attributes left out will be zero'd. For example, if you send along:

    { "Name": "A name for your node" }

A unique ID for your new node will be generated automatically, and the `X`, `Y`, and `Z` attributes will all be 0, and the `ExtraJSON` attribute will be `null`.

To create a new connection, send a POST request to `/connection` with a JSON object with as many of the above Connection-type attributes as you want. You **must** supply the `Source` and `Target` attributes, and they must not be the same. Any attributes left out will be zero'd or auto-filled. For example, if you send along:

    { "Name": "A name for your connection", "Source": 1, "Target": 2 }

A unique ID for your new connection will be generated automatically, and the `Distance` attribute will be generated for you, and the `ExtraJSON` attribute will be `null`.

## Updating

To update a node, simply send a POST request to `/node` with a JSON object in the body with what you'd like the updated node to be. This time, also set an `Id` attribute, of the node ID you'd like to update, like so:

    { "Id": 4, "Name": "A renamed node!" }

That will update the node with ID #4 to the new name.

To update a connection, simply send a POST request to `/connection` with a JSON object in the body with what you'd like the updated connection to be. This time, also set an `Id` attribute, of the connection ID you'd like to update, like so:

    { "Id": 4, "Name": "A renamed connection!", "Source": 5, "Target": 4 }

That will update the connection with ID #4 to the new name, and potentially new source/target if you modified them. Note that the `Source` and `Target` attributes still cannot be the same! And you *do* need to include them **even if you are not modifying them**. (For now, at least.)

## Upserting

As a special note, if you do an above "update" type command, if the ID you're trying to update does not already exist, a new entry will be created automatically, much like how [MongoDB does upserts](http://docs.mongodb.org/manual/reference/method/db.collection.update/).

## Deleting

You can delete a node by sending a DELETE request to `/node/:id` where `:id` is the unique node ID you'd like to have deleted. Any connections that were made to that node will also be deleted.

You can delete a connection by sending a DELETE request to `/connection/:id` where `:id` is the unique connection ID you'd like to have deleted.

## Saving

You can save the current database to disk by sending a GET request to `/save`, you should get an "okay" in return.