package node

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

type NodeRpc int

type WikiPage struct {
	Title     string
	ChildPage []string
	Depth     int
}

func (nr *NodeRpc) TestMethod(req string, resp *int) error {
	fmt.Println("RPC works")
	CrawlInit(req)
	return nil
}

func crawlWiki(respBody io.Reader) (page WikiPage, err error) {
	if respBody == nil {
		fmt.Println("No Resp error:", err.Error())
		err = fmt.Errorf("Crawl Error", err)
		return
	}
	z := html.NewTokenizer(respBody)
	for {
		tokenType := z.Next()
		if tokenType == html.ErrorToken {
			fmt.Println("Error token or done")
			return
		}
		switch tokenType {
		case html.StartTagToken: // <tag>
			// type Token struct {
			//     Type     TokenType
			//     DataAtom atom.Atom
			//     Data     string
			//     Attr     []Attribute
			// }
			//
			// type Attribute struct {
			//     Namespace, Key, Val string
			// }
			token := z.Token()
			isAnchor := token.Data == "a"
			if isAnchor {
				for _, a := range token.Attr {
					if a.Key == "href" {
						fmt.Println("Link found:", a.Val)
					}
				}
			}
		case html.TextToken: // text between start and end tag
		case html.EndTagToken: // </tag>
		case html.SelfClosingTagToken: // <tag/>
		}
	}
}
