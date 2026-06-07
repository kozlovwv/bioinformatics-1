package main

import (
    // Import SciPipe, aliased to sp
    sp "github.com/scipipe/scipipe"
)

func main() {
    // Init workflow and max concurrent tasks
    wf := sp.NewWorkflow("hello_world", 4)

    // Initialize processes, and file extensions
    hello := wf.NewProc("hello", "echo 'Hello ' > {o:out|.txt}")
    world := wf.NewProc("world", "echo $(cat {i:in}) World > {o:out|.txt}")

    // Define data flow
    world.In("in").From(hello.Out("out"))

    // Run workflow
    wf.Run()
}