package core

import (
    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
    "net"
    "net/http"
    "fmt"
    "encoding/json"
    "strings"
    "bytes"
    "log"
	"strconv"
    "io"
    "dcs/crawler"

)

// MasterNode is the node instance
type MasterNode struct {
    api     *gin.Engine            // api server
    ln      net.Listener           // listener
    svr     *grpc.Server           // grpc server
    nodeSvr *NodeServiceGrpcServer // node service
    masterIP string // ip address of master node
}

func (n *MasterNode) Init(masterIp string) (err error) {

    n.masterIP = masterIp

    // grpc server listener with port as 50051
    n.ln, err = net.Listen("tcp", n.masterIP+":50051")
    if err != nil {
        return err
    }

    // grpc server
    n.svr = grpc.NewServer()

    // node service
    n.nodeSvr = GetNodeServiceGrpcServer()

    // register node service to grpc server
    RegisterNodeServiceServer(node.svr, n.nodeSvr)

    // api
    n.api = gin.Default()
    n.api.POST("/links", func(c *gin.Context) {
        // parse payload
        var payload struct {
            Links string `json:"links"`
        }
        if err := c.ShouldBindJSON(&payload); err != nil {
            c.AbortWithStatus(http.StatusBadRequest)
            return
        }

        // post to task manager
        links := strings.Split(payload.Links, " ")
        
        jsonLinks, err := json.Marshal(links)
        if err != nil {
            log.Fatal(err)
        }

        client := &http.Client{}

        req, err := http.NewRequest("POST", "http://localhost:8080/links", bytes.NewBuffer(jsonLinks))
	    req.Header.Set("Content-Type", "application/json")

        client.Do(req)

        c.AbortWithStatus(http.StatusOK)
        
    })

    n.api.GET("/notify/:id", func(c *gin.Context) {
        id, _ := strconv.Atoi(c.Param("id"))
        n.DistributeLinks(id)
        c.AbortWithStatus(http.StatusOK)

    })

    n.DistributeLinks(1)
    go crawler.MasterCrawler(nil, n.masterIP)

    return nil
}


func (n *MasterNode) Start() {
    // start grpc server
    go n.svr.Serve(n.ln)

    // start api server
    _ = n.api.Run(n.masterIP+":9092")

    // run master crawler

    // create es client
    // es := crawler.Setup()    

    // wait for exit
    n.svr.Stop()
}

func (n *MasterNode) DistributeLinks(id int) {

    // get links from task manager
    client := &http.Client{}
    res, err := http.NewRequest("GET", "http://localhost:8080/links", nil)
    if err != nil {
        log.Fatal(err)
    }

    q := res.URL.Query()
    q.Add("CID", strconv.Itoa(id))
    res.URL.RawQuery = q.Encode()


    resp, err := client.Do(res)
    if err != nil {
        log.Fatal(err)
    }
    defer func(Body io.ReadCloser) {
        err := Body.Close()
        if err != nil {
            log.Fatal(err)
        }
    }(resp.Body)

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    var links []string
    err = json.Unmarshal(body, &links)
    if err != nil {
        log.Fatal(err)
    }

    if len(links) == 0 {
        return
    }
    joined := strings.Join(links, " ")

    // move links to LinksChannel
    fmt.Println("moving links to links channel, len: ", len(links))
    n.nodeSvr.LinksChannel <- joined
}

var node *MasterNode

// GetMasterNode returns the node instance
func GetMasterNode(masterIp string) *MasterNode {
    if node == nil {
        // node
        node = &MasterNode{}

        // initialize node
        if err := node.Init(masterIp); err != nil {
            panic(err)
        }
    }

    return node
}
