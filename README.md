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
