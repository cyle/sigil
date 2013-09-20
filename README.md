# SIGIL

A simple graph and spatial database, version 0.1.

This is supposed to be just a very simple database of graphs and spatial information, accessible via a REST interface.

## warning

This is still in very, very early development. It's usable, but the internals of the project may change.

## how to use

**Note:** these instructions were written with the intention of being used on a Mac via Terminal, or on a Debian/Ubuntu box via Bash. I don't yet know how to do all of this on Windows.

First, set up your Go working directory:

    mkdir sigil
    cd sigil

Next, set your $GOPATH to this working directory (you'll need this to install the `gorest` package later):

    export GOPATH=`pwd`

Next, create the necessary working directory subdirectories (based on the [how to write Go code](http://golang.org/doc/code.html#Workspaces) document):

    mkdir src pkg bin

Okay cool -- now clone in this repo to the `src/` directory:

    git clone git@github.com:cyle/sigil.git src/

Right now this primarily relies on the `gorest` third-party package. To install it, do:

    go get code.google.com/p/gorest

That should be it for installing things. Now to run this database:

    cd src/
    go run db.go

That's it... go to `http://localhost:8777/` to see if it's working. That'll be your API endpoint.

## populate it with a demo database

If you have the PHP CLI installed, you can run:

    php make-node-map.php

which will populate the database with a simple test graph. To save it, go to a browser and hit `http://localhost:8777/save`

## check out how the database looks

If you have a web server with PHP running on your development box, you can put the `visualize.php` page somewhere and use that to get an HTML5 Canvas-drawn visualization of your graph database. (Adding 3D support with something like three.js may come someday.)

## API documentation

To use the graph database in an application, visit the [API doc](API.md). Please note that I'll probably be updating this database faster than the API documentation.

## clients

Currently, I have only written two clients:

- [sigil-php](https://github.com/cyle/sigil-php/)
- [sigil-node](https://github.com/cyle/sigil-node/)

## why make this?

I liked the ideas of neo4j but I hate Java. And I want to build something that needs a simple graph database component.

## why "SIGIL"?

I tried coming up with some kind of acronym for "simple graph and spatial database" (SGSDB? lame) and I just thought of calling it SIGIL instead. No reason other than that.