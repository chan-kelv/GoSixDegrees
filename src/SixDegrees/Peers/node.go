package node

/**
 * This package deals with node to node communication to do the crawling
 */

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"strings"
)

var (
	defaultPort = "4000"
	connectPort string
	workerNodes []*rpc.Client
)

func Init(isMaster bool, masterIp string) {
	fmt.Println("Node initalizing...")

	//make sure the ip has an assigned port
	ipWithPort := formatIpWithPort(masterIp)

	if isMaster {
		fmt.Println("Node is master")
		//listen for nodes tring to connect
		err := startServerListen(ipWithPort)
		if err != nil {
			fmt.Println("Server start error:", err)
			panic(err)
		}
	} else {
		fmt.Println("Node is slave")
		//try to connect to master node
		conn, err := dialServer(ipWithPort)
		if err != nil {
			fmt.Println("Could not connect to server:", err)
		} else {
			//hold the connection open
			go slaveServeConn(conn)
		}
	}
}

/**
  * Uses the net package to listen for slave nodes trying to connect.
	* Uses a tcp connection. When a node connects, save the connection.
	* Registers NodeRpc to allow for rpc methods (TODO needed?)
*/
func startServerListen(ipFullAddr string) error {
	//give the ip:port to listen on
	listen, err := net.Listen("tcp", ipFullAddr)
	if err != nil {
		fmt.Println("Server could not start listening on", ipFullAddr)
		err = fmt.Errorf("%v", err.Error())
		return err
	}

	//Register class for rpc
	mServer := new(NodeRpc)
	server := rpc.NewServer()
	server.Register(mServer)

	//Wait for a connection and handle nodes connecting
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Worker could not connect:", err)
		} else {
			//Save the connection of the slave node
			workerAddr := conn.RemoteAddr().String()
			fmt.Println("Worker connected on:", workerAddr)
			workerNodes = append(workerNodes, rpc.NewClient(conn))
		}
	}
}

/**
  * Slave node trying to connect to the master. If successful, will return
	* the connection object. It connects via tcp given an ip:port
*/
func dialServer(masterIp string) (conn *net.TCPConn, err error) {
	//Turn the string ip into tcp addr
	addr, err := net.ResolveTCPAddr("tcp", masterIp)
	if err != nil {
		fmt.Println("Could not resolve ip", err)
		err = fmt.Errorf("Worker connect error: %v", err.Error())
		return
	}

	// Attempt to connect to the master node
	conn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("Could not connect to master node", err)
		err = fmt.Errorf("Worker connect error: %v", err.Error())
		return
	}
	return conn, nil
}

/**
 * From the connection given, establishes a serveconn to listen for master
 * node calling rpc methods to this slave node. This method blocks on ServeConn
 */
func slaveServeConn(conn *net.TCPConn) {
	//Register NodeRpc to allow for rpc methods (TODO needed?)
	server := rpc.NewServer()
	rpcType := new(NodeRpc)
	server.Register(rpcType)

	server.ServeConn(conn)
	fmt.Println("Connection to server broken")
}

/**
  * Formats the given ip to make sure it has an attached port
	* This only conforms with ipv4 addresses for now...
*/
func formatIpWithPort(ip string) string {
	if strings.Contains(ip, ":") {
		return ip
	} else {
		//assume its only the ip that was passed without portDefault
		return ip + ":" + defaultPort
	}
}

/**
 * Crawl the initial page
 */
func CrawlInit(crawlTerm string) {
	fmt.Println("Begin crawl on page:", crawlTerm)
	resp, err := http.Get(crawlTerm)
	if err != nil {
		fmt.Println("Crawl init error:", err)
		return
	}

	respBody := resp.Body
	if err != nil {
		fmt.Println("Crawl init body error:", err)
		return
	}

	// z := html.NewTokenizer(respBody)
	// for {
	// 	tokenType := z.Next()
	// 	if tokenType == html.ErrorToken {
	// 		fmt.Println("Error token or done")
	// 		return
	// 	}
	// 	switch tokenType {
	// 	case html.StartTagToken: // <tag>
	// 		// type Token struct {
	// 		//     Type     TokenType
	// 		//     DataAtom atom.Atom
	// 		//     Data     string
	// 		//     Attr     []Attribute
	// 		// }
	// 		//
	// 		// type Attribute struct {
	// 		//     Namespace, Key, Val string
	// 		// }
	// 		token := z.Token()
	// 		isAnchor := token.Data == "a"
	// 		if isAnchor {
	// 			for _, a := range token.Attr {
	// 				if a.Key == "href" {
	// 					fmt.Println("Link found:", a.Val)
	// 				}
	// 			}
	// 		}
	// 	case html.TextToken: // text between start and end tag
	// 	case html.EndTagToken: // </tag>
	// 	case html.SelfClosingTagToken: // <tag/>
	// 	}
	// }
	fmt.Println("Done crawl")
}

func MakeTestRpc(crawlTerm string) {
	// workerNode := workerNodes[0]
	// var reply int
	// workerNode.Call("NodeRpc.TestMethod", "", &reply)
	CrawlInit(crawlTerm)
}
