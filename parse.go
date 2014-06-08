package template

import (
	"encoding/xml"
	"fmt"
)

type Node struct {
	FunctionName string `xml:"data-lift,attr"`
	NodeSeq      string `xml:",innerxml"`
}
type Result struct {
	XMLName   xml.Name
	Functions []Node `xml:",any"`
}

func GetFunctions(html string) {
	err, v := MarshalNode(html)
	if err != nil {
		fmt.Printf("Error 1: %v\n\n", err)
		return
	}
	Loop(v)
}

func MarshalNode(html string) (error, Result) {
	v := Result{}
	err := xml.Unmarshal([]byte("<p>"+html+"</p>"), &v)
	if err != nil {
		return err, v
	}
	return nil, v
}

func Loop(v Result) {
	for _, innerNode := range v.Functions {
		if innerNode.FunctionName != "" {
			//fmt.Println("Found Function: " + innerNode.FunctionName)
			fmt.Printf("Calling: %v( %v )\n\n", innerNode.FunctionName, innerNode.NodeSeq)
		}

		err, node := MarshalNode(innerNode.NodeSeq)
		if err != nil {
			fmt.Printf("Error 2: %v ==>> %v\n\n", innerNode.NodeSeq, err)
		}
		Loop(node)
	}
}
