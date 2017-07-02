package node

import (
	"fmt"
	"net"
	"net/rpc"
	"strings"
)

var (
	defaultPort   = "4000"
	connectPort   string
	workerClients []*rpc.Client
)

func Init(isMaster bool, masterIp string) {
	fmt.Println("Node initalizing...")
	ipWithPort := formatIpWithPort(masterIp)
	if isMaster {
		fmt.Println("Node is master")
		err := startServerListen(ipWithPort)
		if err != nil {
			fmt.Println("Server start error:", err)
			panic(err)
		}
	} else {
		fmt.Println("Node is slave")
		conn, err := dialServer(ipWithPort)
		if err != nil {
			fmt.Println("Could not connect to server:", err)
		} else {
			go rpcListener(conn)
		}
	}
	fmt.Println("Master ip is:", ipWithPort)
}

func startServerListen(ipAddr string) error {
	ipFullAddr := ipAddr
	listen, err := net.Listen("tcp", ipFullAddr)
	if err != nil {
		fmt.Println("Server could not start listening on", ipFullAddr)
		err = fmt.Errorf("%v", err.Error())
		return err
	}

	mServer := new(NodeRpc)
	server := rpc.NewServer()
	server.Register(mServer)

	for {
		conn, err := listen.Accept()

		if err != nil {
			fmt.Println("Worker could not connect:", err)
		} else {
			workerAddr := conn.RemoteAddr().String()
			fmt.Println("Worker connected on:", workerAddr)
			workerClients = append(workerClients, rpc.NewClient(conn))
		}
	}
}

func dialServer(masterIp string) (conn *net.TCPConn, err error) {
	addr, err := net.ResolveTCPAddr("tcp", masterIp)
	if err != nil {
		fmt.Println("Could not resolve ip", err)
		err = fmt.Errorf("Worker connect error: %v", err.Error())
		return
	}

	conn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("Could not connect to master node", err)
		err = fmt.Errorf("Worker connect error: %v", err.Error())
		return
	}

	return conn, nil
}

func rpcListener(conn *net.TCPConn) {
	server := rpc.NewServer()
	rpcType := new(NodeRpc)
	server.Register(rpcType)

	server.ServeConn(conn)
	fmt.Println("Connection to server broken")
}

func formatIpWithPort(ip string) string {
	if strings.Contains(ip, ":") {
		return ip
	} else {
		//assume its only the ip that was passed without portDefault
		return ip + ":" + defaultPort
	}
}
