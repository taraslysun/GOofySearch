package core

import (
    "context"
    "google.golang.org/grpc"
	"fmt"
    "strings"
    "dcs/crawler"
)

type WorkerNode struct {
    conn *grpc.ClientConn  // grpc client connection
    c    NodeServiceClient // grpc client
}

func (n *WorkerNode) Init() (err error) {
    // connect to master node
    n.conn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        return err
    }

    // grpc client
    n.c = NewNodeServiceClient(n.conn)

    return nil
}

func (n *WorkerNode) Start() {
    // log
    fmt.Println("worker node started")

    // create es client
    // es := crawler.Setup()

    // report status
    _, _ = n.c.ReportStatus(context.Background(), &Request{})

    // assign task
    stream, _ := n.c.AssignTask(context.Background(), &Request{})
    for {
        // receive links from master node
        res, err := stream.Recv()
        if err != nil {
            return
        }


        // run CrawlerMain with recieved links
        links := strings.Split(res.Data, " ")
        // log command
        fmt.Println("worker received links: ", len(links))

        // TODO: divide links, run multiple CrawlerMain
        crawler.CrawlerMain(links, len(links), nil)


    }
}

var workerNode *WorkerNode

func GetWorkerNode() *WorkerNode {
    if workerNode == nil {
        // node
        workerNode = &WorkerNode{}

        // initialize node
        if err := workerNode.Init(); err != nil {
            panic(err)
        }
    }

    return workerNode
}
