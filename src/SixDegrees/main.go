package main

import (
	routes "SixDegrees/Api"
	node "SixDegrees/Peers"
	"fmt"
	"os"
	"time"
)

var (
	HTTP_SERVER_START_DELAY = time.Duration(0) //stall to allow rpc to connect first
	DEFAULT_PORT            = "4000"
)

func main() {
	//parse args
	isMaster, masterIp, err := parseArgs()
	if err != nil {
		panic(err)
	}

	//make a channel that does not close to force the main to stay open
	done := make(chan int)

	//Init the nodes on the server (master or slave)
	go node.Init(isMaster, masterIp)

	//Give the nodes time to init before setting up http server on master
	if isMaster {
		go func() {
			time.Sleep(HTTP_SERVER_START_DELAY * time.Second)
			routes.Init()
		}()
	}

	// holds the main open
	<-done
}

func parseArgs() (isMaster bool, masterIp string, err error) {
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Println("Incorrect number of args. Indicate master/slave with -m/-s and the IP of the master node")
		err = fmt.Errorf("Incorrect number of args", err)
		return
	}

	isMaster = args[0] == "-m"

	masterIp = args[1]

	return
}
