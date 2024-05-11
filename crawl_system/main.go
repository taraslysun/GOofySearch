package main

import (
    "dcs/core"
    "os"
	"fmt"
)

func main() {
    nodeType := os.Args[1]
	fmt.Println("Starting")
    switch nodeType {
    case "master":
        core.GetMasterNode().Start()
    case "worker":
        core.GetWorkerNode().Start()
    default:
        panic("invalid node type")
    }
}