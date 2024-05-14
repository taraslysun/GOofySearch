package main

import (
    "dcs/core"
    "os"
	"fmt"
	"strconv"
)

func main() {
    nodeType := os.Args[1]
	masterIp := os.Args[2]
	fmt.Println("Starting")
    switch nodeType {
    case "master":
		numWorkers, _ := strconv.Atoi(os.Args[3])
        core.GetMasterNode(masterIp, numWorkers).Start()
    case "worker":
		id, _ := strconv.Atoi(os.Args[3])
        core.GetWorkerNode(masterIp, id).Start()
    default:
        panic("invalid node type")
    }
}