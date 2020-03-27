kmud-2020
===========

A refactor and update of [kmud](https://github.com/Cristofori/kmud) written in Go.

Development
===========

The current development priority is splitting the project into services that can be run independently in docker/kubernetes.  

For the week of March 26th, 2020:
    
* Frontend Service
    
    The focus of this service is a passthrough frontend that can communicate with backend services. "Frontend" delivers all telnet interaction with clients. 
    
    


Installation
============

go get github.com/yamamushi/kmud
go install github.com/yamamushi/kmud


Dependencies
============
Google Go v1.2
MongoDB: www.mongodb.org 

mgo: http://labix.org/mgo
go get gopkg.in/mgo.v2

go check: http://labix.org/gocheck
go get gopkg.in/check.v1


Roadmap
============

* [x] Add config file 
* [ ] Add proper logging output
* [ ] Refactor to clustered services model
* [ ] Create embedded builder tools 
* [ ] Refactor menu system
* [ ] Refactor admin tools
* [ ] Add [mccp2](https://mudhalla.net/tintin/protocols/mccp/) support
* [ ] Add color support
* [ ] Add unicode support 
* [ ] Add tls support
* [ ] Add [mxp](http://www.zuggsoft.com/zmud/mxp.htm) support
* [ ] Add [mmcp](https://mudhalla.net/tintin/protocols/mmcp/) support
* [ ] Add boat/ship engine support
* [ ] Add player house support
* [ ] Add global z-level support
* [ ] Refactor world generator
* [ ] Add world editing support

