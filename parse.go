package arriba

import (
	"encoding/xml"
	"fmt"
	"bytes"
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



const html1 = (`<html><body><div data-lift="ChangeName"><p name="name">Diego</p><p data-lift="ChangeLastName">Medina</p></div></body></html>`)

func Children(tokenArray []xml.Token, initialIndex int) ([]xml.Token, xml.Token) {
	openTags := 1
	closedTags := 0
	stack := new(Stack)
	var closingTag = new(xml.Token)
	for i:= initialIndex; i < len(tokenArray); i++ {
		tok := tokenArray[i]
		switch innerTok := tok.(type) {
		case xml.EndElement:
			closedTags++
			if openTags == closedTags {
				closingTag = innerTok
				break
			}
			stack.Push(innerTok)

		case xml.StartElement:
			openTags++
			stack.Push(innerTok)

		case xml.CharData:
			stack.Push(innerTok)
		}
	}
	var ret = make([]xml.Token, stack.Size())
	var i = stack.Size()
	for i > 0 {
		ret[i - 1] = stack.Pop()
		i--
	}
	return ret, closingTag
}

func toTokenArray(decoder *xml.Decoder) []xml.Token {
	stack := new(Stack)
	for {
		tok, err := decoder.Token()
		if err != nil {
			//We are done processing tokens, let's end.
			break
		}
		stack.Push(tok)
	}
	var array = make([]xml.Token, stack.Size())
	var i = stack.Size()
	for i > 0 {
		array[i - 1] = stack.Pop()
		i--
	}
	return array
}

func Princ() {

	decoder := xml.NewDecoder(bytes.NewBufferString(html1))

	var tokenArray = toTokenArray(decoder)

	fmt.Println(tokenArray)

	var children = Children(tokenArray, 1)


	fmt.Println(children)
//
//	token2, _ := decoder.Token()
//
//	fmt.Println(Children(&token2))




}
