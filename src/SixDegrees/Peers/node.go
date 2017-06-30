package node

import "fmt"

func Init(isMaster bool, masterIp string) {
	fmt.Println("Node initalizing...")
	if isMaster {
		fmt.Println("Node is master")
	} else {
		fmt.Println("Node is slave")
	}
	fmt.Println("Master ip is:", masterIp)
}
