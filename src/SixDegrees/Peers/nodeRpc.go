package node

import (
	"fmt"
)

type NodeRpc int

func (nr *NodeRpc) TestMethod(req string, resp *int) error {
	fmt.Println("RPC works")
	CrawlInit(req)
	return nil
}
