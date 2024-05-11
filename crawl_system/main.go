package main

import (
    "dcs/core"
    "os"
	"fmt"
)

func main() {
    nodeType := os.Args[1]
	masterIp := os.Args[2]
	fmt.Println("Starting")
    switch nodeType {
    case "master":
        core.GetMasterNode(masterIp).Start()
    case "worker":
        core.GetWorkerNode(masterIp).Start()
    default:
        panic("invalid node type")
    }
}