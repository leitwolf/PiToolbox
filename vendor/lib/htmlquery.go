package lib

//
// 查询html标签
//
import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// HTMLQuery html相关的标签查询
type HTMLQuery struct {
	// 根节点
	Doc *html.Node
}

// Read 从Reader读取
func (q *HTMLQuery) Read(r io.Reader) (err error) {
	d, err := html.Parse(r)
	if err != nil {
		return
	}
	q.Doc = d
	return
}

// GetNodeByID 获取具有id=value的节点
func (q *HTMLQuery) GetNodeByID(id string) (node *html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if node != nil {
			return
		}
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if strings.ToLower(attr.Key) == "id" && attr.Val == id {
					node = n
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(q.Doc)
	return
}

// GetNodeAttr 获取一个节点指定的属性值
func (q *HTMLQuery) GetNodeAttr(node *html.Node, attrName string) (attr string) {
	for _, a := range node.Attr {
		if a.Key == attrName {
			attr = a.Val
			return
		}
	}
	return
}

// GetNodesByName 获取具有name=value的所有节点
func (q *HTMLQuery) GetNodesByName(name string) (nodes []*html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if strings.ToLower(attr.Key) == "name" && attr.Val == name {
					nodes = append(nodes, n)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(q.Doc)
	return
}

// NewQueryFromReader 新建Query
func NewQueryFromReader(r io.Reader) (q *HTMLQuery, err error) {
	q = &HTMLQuery{}
	err = q.Read(r)
	return
}
