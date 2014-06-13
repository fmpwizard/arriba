package arriba

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

//Node holds the function name taken from the data-lift attribute and
// the NodeSeq is the raw inner html that we will pass to the function on the data-lift attribute
type Node struct {
	FunctionName  string `xml:"data-lift,attr"`
	NodeSeq       string `xml:",innerxml"`
	ProcessedHTML string
}

//Result holds a slice of Node values
type Result struct {
	XMLName       xml.Name
	Functions     []Node `xml:",any"`
	ParentNodeSeq string
}

type PartialHTML struct {
	RawHTML      string
	FunctionName string
	FunctionHtml string
}
type HTMLTransform func(string) string

type Elem struct {
	XMLName  xml.Name
	Comment  xml.Comment
	Attr     xml.Attr
	InnerXML string `xml:",innerxml"`
	Kids     []Elem `xml:",any"`
}

var funcMap = make(map[string]HTMLTransform)

func MarshallElem(in string) string {
	funcMap["ChangeTime"] = ChangeTime
	funcMap["ChangeName"] = ChangeName
	//node := Elem{}
	//err := xml.Unmarshal([]byte(in), &node)

	decoder := xml.NewDecoder(bytes.NewBufferString(in))

	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		switch startElement := token.(type) {
		case xml.Attr:
			fmt.Printf("\n\n1: %+v\n", startElement)
		case xml.StartElement:
			//fmt.Printf("\n\n2: %+v\n", startElement.Attr)
			for _, value := range startElement.Attr {
				fmt.Printf("Att: %v ==> value: %v\n", value.Name.Local, value.Value)
				if value.Name.Local == "data-lift" {
					snipetHTML := ""
					open := 0
					close := 0
					for {
						tok, err := decoder.Token()
						if err != nil {
							return err.Error()
						}
						switch innerTok := tok.(type) {
						case xml.StartElement:
							snipetHTML = snipetHTML + "<" + innerTok.Name.Local
							for _, attr := range innerTok.Attr {
								snipetHTML = snipetHTML + " " + attr.Name.Local + "=" + attr.Value
							}
							snipetHTML = snipetHTML + ">"
							open++
							fmt.Printf("1111: %v\n", innerTok.Name.Local)
						case xml.CharData:
							snipetHTML = snipetHTML + string(innerTok)

						case xml.EndElement:
							fmt.Printf("2222 %v\n", innerTok.Name.Local)
							snipetHTML = snipetHTML + "</" + innerTok.Name.Local + ">"
							close++
							if open == close {
								return snipetHTML
							}
						}
						fmt.Printf(" ==>> snipetHTML  %v\n\n", snipetHTML)
					}

				}
			}
			fmt.Printf("Start: %v\n", startElement.Name.Local)
			fmt.Printf("End: %v\n", startElement.End().Name.Local)
		case xml.CharData:
			fmt.Printf("\n\nCharData: %+v\n", string(startElement))
		case xml.EndElement:
			fmt.Printf("\n\nEnd: %+v\n", startElement.Name.Local)
		case xml.Comment:
			fmt.Printf("\n\nComment: %+v\n", startElement)
		case xml.Directive:
			fmt.Printf("\n\nDirective: %+v\n", string(startElement))
		case xml.Token:
			fmt.Printf("\n\n4: %+v\n", startElement)

		default:
			fmt.Printf("\nedd\n")
		}

	}
	return "diego"
}

//GetFunctions takes the complete html of a page and returns a map of
//function names => html that we should pass to those functions
func GetFunctions(html string) []PartialHTML {
	funcMap["ChangeTime"] = ChangeTime

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
	v.ParentNodeSeq = parentHtml
	return nil, v
}

func ProcessHTML(html, parent string) (error, string) {
	funcMap["ChangeTime"] = ChangeTime
	//fmt.Printf("html is %v\n", html)
	temp := parent
	err, v := marshalNode(html, parent)
	if err != nil {
		fmt.Printf("Error A: %v\n\n", err)
		return err, ""
	}
	//fmt.Printf("v.Functions %v\n==============\n", v.Functions)
	for key, innerNode := range v.Functions {
		if innerNode.FunctionName != "" {
			//fmt.Printf("inner: %v\n=============\n", funcMap[innerNode.FunctionName](innerNode.NodeSeq))
			//fmt.Println("2")
			v.Functions[key].ProcessedHTML = funcMap[innerNode.FunctionName](innerNode.NodeSeq)
			_, temp = ProcessHTML(funcMap[innerNode.FunctionName](innerNode.NodeSeq), v.ParentNodeSeq)
			return nil, temp
		}
		//fmt.Printf("3 %v\n=========\n", innerNode.NodeSeq)
		//fmt.Printf("4 %v\n=========\n", v.ParentNodeSeq)
		temp = v.ParentNodeSeq
		//return ProcessHTML(innerNode.NodeSeq, v.ParentNodeSeq)
		_, h := ProcessHTML(innerNode.NodeSeq, v.ParentNodeSeq)
		h = h
		//fmt.Printf("5 %v\n============\n", h)
	}

	fmt.Println("here " + temp)

	return nil, temp

}

func loop(v Result) []PartialHTML {
	var template []PartialHTML
	var funccMap = make(map[string]string)
	for _, innerNode := range v.Functions {
		if innerNode.FunctionName != "" {
			funccMap[innerNode.FunctionName] = innerNode.NodeSeq
			fmt.Printf("inner: \n%v\n\n=============\n\n", funcMap[innerNode.FunctionName](innerNode.NodeSeq))

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

func ChangeTime(html string) string {
	return strings.Replace(html, "Time goes here", time.Now().Format("2006-01-02T15:04:05.999999999Z07:00"), 1)
}

func ChangeName(html string) string {
	return strings.Replace(html, "Diego", "Gabriel", 1)
}
