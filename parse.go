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

//Result holds a slice of Node values
type Result struct {
	XMLName       xml.Name
	Functions     []Node `xml:",any"`
	ParentNodeSeq string
}

//GetFunctions takes the complete html of a page and returns a map of
//function names => html that we should pass to those functions
func GetFunctions(html string) []PartialHTML {
	//func GetFunctions(html string) map[string]string {
	err, v := marshalNode(html, html)
	if err != nil {
		fmt.Printf("Error 1: %v\n\n", err)
		return nil
	}
	return loop(v)
}

func marshalNode(html string, parentHtml string) (error, Result) {
	v := Result{}
	//horrible hack to get the complete html that is inside the node, otherwise we only get child nodes and miss data
	err := xml.Unmarshal([]byte("<p>"+html+"</p>"), &v)
	if err != nil {
		return err, v
	}
	v.ParentNodeSeq = html
	return nil, v
}

func loop(v Result) []PartialHTML {
	var template []PartialHTML
	var funccMap = make(map[string]string)
	for _, innerNode := range v.Functions {
		if innerNode.FunctionName != "" {
			funccMap[innerNode.FunctionName] = innerNode.NodeSeq
			//fmt.Printf("function name: %v\n\n============\n\n", innerNode.FunctionName)
		}
		//fmt.Printf("parent html: \n%v\n\n============\n\n", v.ParentNodeSeq)
		//fmt.Printf("html: \n%v\n\n============\n\n", innerNode.NodeSeq)
		pre := PartialHTML{v.ParentNodeSeq, innerNode.FunctionName, innerNode.NodeSeq}
		template = append(template, pre)
		//fmt.Printf("pre:\n%+v\n", pre)
		err, node := marshalNode(innerNode.NodeSeq, v.ParentNodeSeq)
		if err != nil {
			fmt.Printf("Error 2: %v ==>> %v\n\n", innerNode.NodeSeq, err)
		}
		//we have more html, so we recurse, but we need to keep the old map of functions
		//for k, v := range loop(node) {
		//	funccMap[k] = v
		//}
		template = append(template, loop(node)...)
	}
	//return funccMap
	return template
}

type PartialHTML struct {
	RawHTML      string
	FunctionName string
	FunctionHtml string
}
