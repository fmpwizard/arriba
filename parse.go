package arriba

import (
	"encoding/xml"
	"fmt"
)

//Node holds the function name taken from the data-lift attribute and
// the NodeSeq is the raw inner html that we will pass to the function on the data-lift attribute
type Node struct {
	FunctionName string `xml:"data-lift,attr"`
	NodeSeq      string `xml:",innerxml"`
}

func push2Stack(array []Node, stack *Stack) {
	for _, item := range array {
		stack.Push(item)
	}
}

//Result holds a slice of Node values
type Result struct {
	XMLName   xml.Name
	Functions []Node `xml:",any"`
}

//GetFunctions takes the complete html of a page and returns a map of
//function names => html that we should pass to those functions
func GetFunctions(html string) map[string]string {
	err, v := marshalNode(html)
	if err != nil {
		fmt.Printf("Error 1: %v\n\n", err)
		return nil
	}
	return loop(v)
}

func marshalNode(html string) (error, Result) {
	v := Result{}
	//horrible hack to get the complete html that is inside the node, otherwise we only get child nodes and miss data
	err := xml.Unmarshal([]byte("<p>"+html+"</p>"), &v)
	if err != nil {
		return err, v
	}
	return nil, v
}

func loop(v Result) map[string]string {
	var stack = new(Stack)

	var functionMap = make(map[string]string)
	// Stack initialization with array elements
	push2Stack(v.Functions, stack)
	for stack.Size() > 0 {
		innerNode := stack.Pop().(Node)
		if innerNode.FunctionName != "" {
			functionMap[innerNode.FunctionName] = innerNode.NodeSeq
		}

		err, node := marshalNode(innerNode.NodeSeq)
		if err != nil {
			fmt.Printf("Error 2: %v ==>> %v\n\n", innerNode.NodeSeq, err)
		}
		// we have more html, add the pending nodes to the stack

		push2Stack(node.Functions, stack)
	}
	return functionMap
}
