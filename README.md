```
  _____     _______ ____  __  __ ___ _   _ ____
 / _ \ \   / / ____|  _ \|  \/  |_ _| \ | |  _ \
| | | \ \ / /|  _| | |_) | |\/| || ||  \| | | | |
| |_| |\ V / | |___|  _ <| |  | || || |\  | |_| |
 \___/  \_/  |_____|_| \_\_|  |_|___|_| \_|____/
```

A minimalistic group of microservices for controlling a zerg swarm.

It's used to demonstrate deployment of microservices onto a Kubernetes cluster.

Components
==========

Overmind
--------
The overmind is the service that the user communicate with to control the swarm. Overmind can be instructed to spawn zerglings through a viper and control the spawned zerglings.

Brain
-----
The brain of the overmind service is a [CouchDB](https://couchdb.apache.org) database. CouchDB was chosen for its simplicity and because of its RESTish-ness, we can treat CouchDB as if it were another web service (albeit a stateful one).

API Spec
========

The marsrovers microservice is designed as a REST-API. The endpoints are:

* `GET /_health` - The health of the overmind and its subordinates
* `GET /zerglings/` - All zerglings the overmind is aware of and their locations
* `GET /zerglings/<zergling_id>` - Get the status of the specified zergling
* `POST /zerglings/<zergling_id>/_action` - Move the zergling
* `POST /zerglings/` - Spawn a zergling

Build
=====

## Build on host

    make overmind

This will generate `overmind` binary.

## Run on host

    ./overmind

The overmind service will be on `0.0.0.0:8080` by default. To run it on a different ip or port, use `-http.addr=`:

    ./overmind -http.addr=0.0.0.0:9999

## Run inside Docker container

First you will need to build the docker image:

    make image

Run it:

    docker run overmind

To run it on a different port:

    docker run -e OVERMIND_HTTP_ADDR=0.0.0.0:9999 overmind

Local Dev
=========

Startup CouchDB
---------------

When developing locally, use the `docker-compose.yaml` file to start up a local CouchDB instance.

    docker-compose up -d

Bootstrap
---------

CouchDB user and database need to be bootstrapped. We can use the `bootstrap.sh` script but first, we need to find out the ip address of the CouchDB instance we just started.

Get the id of the container:

    $ docker ps | grep couchdb
    3303970683ac        couchdb:1.6.1       "tini -- /docker-e..."   2 minutes ago       Up About a minute   5984/tcp            overmind_couchdb_1

Inspect the container:

    docker inspect 3303 | grep IPAdd
    [...]
    "IPAddress": "172.18.0.2",

There's the ip address.

Modify the `dev.env` and change `COUCHDB_SERVICE_HOST` to the ip address above.

Then source the env file and source the bootstrap file:

    source dev.env
    source bootstrap.sh

This will create the couchdb user and the `zergling` database.

Compile and Start the overmind service
--------------------------------------

    make overmind

Start it:

    ./overmind

You will see:

    ts=2017-05-27T16:49:42.418479742Z caller=service.go:187 COUCHDB_USERNAME=admin COUCHDB_PASSWORD=*** COUCHDB_SERVICE_HOST=172.18.0.2 COUCHDB_SERVICE_PORT=5984
    ts=2017-05-27T16:49:42.41943168Z caller=main.go:50 transport=HTTP addr=:8080

Verify its health:

    $ curl localhost:8080/_health
    {"health":{"version":"1.0.0","brain":"ok"}}

Sample Usage
============

Verify health:

    $ curl localhost:8080/_health
    {"health":{"version":"1.0.0","brain":"ok"}}

Spawn a new zergling:

    $ curl -XPOST localhost:8080/zerglings/
    {"zergling":{"id":"5fb03720-c5c4-4961-8f1f-6d40bfd3f46c","location":{"x":0,"y":0},"facing":"N"}}

Get all zerglings:

	$ curl -s localhost:8080/zerglings/ | jq .
	{
	  "zerglings": [
		{
		  "id": "20214980-b68c-4379-8971-1633822930e5",
		  "location": {
			"x": 0,
			"y": 0
		  },
		  "facing": ""
		},
		{
		  "id": "5fb03720-c5c4-4961-8f1f-6d40bfd3f46c",
		  "location": {
			"x": 0,
			"y": 0
		  },
		  "facing": ""
		},
		{
		  "id": "612ab522-9c52-493d-94f1-8243cd06b581",
		  "location": {
			"x": 0,
			"y": 0
		  },
		  "facing": ""
		},
		{
		  "id": "d0fd6371-d37b-4bd6-805e-a45519d00ecc",
		  "location": {
			"x": 0,
			"y": 0
		  },
		  "facing": ""
		}
	  ]
	}

Get a zergling by id:

	$ curl localhost:8080/zerglings/20214980-b68c-4379-8971-1633822930e5
	{"zergling":{"id":"20214980-b68c-4379-8971-1633822930e5","location":{"x":0,"y":0},"facing":"N","_rev":"1-943f87a80c81dbc9e6260bbfb5df8e2a"}}


Move a zergling:

	$ curl -XPOST localhost:8080/zerglings/20214980-b68c-4379-8971-1633822930e5 -d'"L"'
	{"zergling":{"id":"20214980-b68c-4379-8971-1633822930e5","location":{"x":0,"y":0},"facing":"W","commandHistory":["L"],"_rev":"1-943f87a80c81dbc9e6260bbfb5df8e2a"}}

	$ curl -XPOST localhost:8080/zerglings/20214980-b68c-4379-8971-1633822930e5 -d'"M"'
	{"zergling":{"id":"20214980-b68c-4379-8971-1633822930e5","location":{"x":-1,"y":0},"facing":"W","commandHistory":["L","M"],"_rev":"2-68acfc992e04e5d75ca7a14f2783cb87"}}

	$ curl -XPOST localhost:8080/zerglings/20214980-b68c-4379-8971-1633822930e5 -d'"M"'
	{"zergling":{"id":"20214980-b68c-4379-8971-1633822930e5","location":{"x":-2,"y":0},"facing":"W","commandHistory":["L","M","M"],"_rev":"3-1a8ae9d0b3828901f936b2aba818b8ca"}}

Deploy on Kubernetes
====================

TBD
