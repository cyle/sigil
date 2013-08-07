# The Simple Graph Database REST API

## Querying

To get all nodes, simply send a GET request to `/nodes`.

To get all connections, simply send a GET request to `/connections`.

To get a specific node, simply send a GET request to `/node/:id`, where `:id` is the unique ID number of the node you're looking for.

To get a specific connection, simply send a GET request to `/connection/:id`, where `:id` is the unique ID number of the connection you're looking for.

To get all connections attached to a given node, simply send a GET request to `/node/:id/connections`, where `:id` is the unique ID number of the node you want connections for.

## Creating

To create a new node, send a POST request to `/node` with a JSON object in the body with the following attributes:

    { "Name": "A name for your node" }

A unique ID for your new node will be generated automatically.

To create a new connection, send a POST request to `/connection` with a JSON object in the body with the following attributes: 

    { "Name": "A name for your connection", "Source": 1, "Target": 2 }

A unique ID for your new connection will be generated automatically. You **must** supply the `Source` and `Target` attributes, and they must not be the same.

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

You can delete a node by sending a DELETE request to `/node/:id` where `:id` is the unique node ID you'd like to have deleted.

You can delete a connection by sending a DELETE request to `/connection/:id` where `:id` is the unique connection ID you'd like to have deleted.