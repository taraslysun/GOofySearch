package core

import (
    "context"
    "google.golang.org/grpc"
	"fmt"
    "strings"
    "dcs/crawler"
    "net/http"
    "log"
    "strconv"
    "sync"
)

type WorkerNode struct {
    conn *grpc.ClientConn  // grpc client connection
    c    NodeServiceClient // grpc client
    masterIP string // ip address of master node
    ID int // int of this worker node
}

func (n *WorkerNode) Init(masterIp string, id int) (err error) {

    n.masterIP = masterIp
    n.ID = id

    // connect to master node
    n.conn, err = grpc.Dial(n.masterIP+":50051", grpc.WithInsecure())
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
    es := crawler.Setup()

    // report status
    _, _ = n.c.ReportStatus(context.Background(), &Request{})

    // assign task
    stream, _ := n.c.AssignTask(context.Background(), &Request{})
    for {
        
        resp, err := http.Get("http://"+ n.masterIP+":9092/notify/"+strconv.Itoa(n.ID))
        if err != nil {
          log.Fatal(err)
        }
    
        err = resp.Body.Close()
        if err != nil {
          log.Fatal(err)
        }
        
        var wg sync.WaitGroup
        for i := 1; i <= 5; i++ {
            wg.Add(1)
            // receive links from master node
            res, err := stream.Recv()
            if err != nil {
                return
            }
        
            // run CrawlerMain with received links
            links := strings.Split(res.Data, " ")

            go func() {
                defer wg.Done()
                crawler.CrawlerMain(links, len(links), es, n.masterIP)
            }()

        }
        wg.Wait()

    
      }
}

var workerNode *WorkerNode

func GetWorkerNode(masterIp string, id int) *WorkerNode {
    if workerNode == nil {
        // node
        workerNode = &WorkerNode{}

        // initialize node
        if err := workerNode.Init(masterIp, id); err != nil {
            panic(err)
        }
    }

    return workerNode
}
