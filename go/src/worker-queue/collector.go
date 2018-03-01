package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

func Collector(w http.ResponseWriter, r *http.Request) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Tree part, now works by just hardcoded json file
	// But need to fix / use the json from http post in the future
	file, e := ioutil.ReadFile("./tree-example.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	myJson := string(file)

	tree, err := parseJson(myJson)
	if err != nil {
		fmt.Println(err)
	}
	root := parseBinaryTree(tree)
	/* -- */

	work := WorkRequest{ID: 1, Tree: &root}
	fmt.Println("Workrequest tree created")

	// Push the work onto the queue.
	WorkQueue <- work
	fmt.Println("Work request queued")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	return
}
