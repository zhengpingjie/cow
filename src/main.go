package main

import (
	"./serveGate"
	"net/http"
	_"net/http/pprof"
)



func main(){
	//http://localhost:6060/debug/pprof/
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
	serveGate.Start()
}