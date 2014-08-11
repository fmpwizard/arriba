package arriba

import (
	"bytes"
	"errors"

	//"bytes"
	//"encoding/xml"
	//"errors"
	"fmt"
	"github.com/fmpwizard/arriba/vendor/code.google.com/p/go-html-transform/h5"
	"github.com/fmpwizard/arriba/vendor/code.google.com/p/go.net/html"
	//"strings"
	"sync"
)

type HTMLTransform func(*html.Node) *html.Node

var FunctionMap = struct {
	sync.RWMutex
	M map[string]HTMLTransform
}{M: make(map[string]HTMLTransform)}

/*type snippetAndNode struct {
	FunctionName string
	HTML         string
}
*/
/*var ch = make(chan snippetAndNode)*/

/*func MarshallElem(in string) string {
	//go readChan(ch)

	completeHTML := ""
	decoder := xml.NewDecoder(bytes.NewBufferString(in))
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch element := token.(type) {
		case xml.StartElement:
			completeHTML = completeHTML + "<" + element.Name.Local
			for _, attr := range element.Attr {
				if attr.Name.Local != "data-lift" {
					completeHTML = completeHTML + " " + attr.Name.Local + "=\"" + attr.Value + "\""
				}
			}
			completeHTML = completeHTML + ">"

			err, res := processSnippet(decoder, element.Name.Local, "")
			if err != nil {
				return err.Error()
			}
			completeHTML = completeHTML + res
			if !strings.HasSuffix(completeHTML, ">") {
				completeHTML = completeHTML + ">"
			}
		case xml.CharData:
			completeHTML = completeHTML + string(element)
		case xml.EndElement:
			completeHTML = completeHTML + "</" + element.Name.Local + ">"
		case xml.Comment:
			fmt.Printf("Comment: %+v\n", element)
		case xml.Directive:
			fmt.Printf("Directive: %+v\n", string(element))
		case xml.Token:
			fmt.Printf("4: %+v\n", element)

		default:
			fmt.Errorf("\nIf you are here, you are missing a type: %v\n", element)
		}

	}
	return completeHTML
}

func processSnippet(decoder *xml.Decoder, parentTag string, scopeFunction string) (error, string) {
	snippetHTML := ""

	open := 1
	closingTags := 0

	for {
		tok, err := decoder.Token()
		if err != nil {
			//We are done processing tokens, let's end.
			//close(ch)
			if err.Error() == "EOF" {
				return nil, snippetHTML
			} else {
				return err, snippetHTML
			}

		}
		switch innerTok := tok.(type) {
		case xml.StartElement:
			snippetHTML = snippetHTML + "<" + innerTok.Name.Local
			for _, attr := range innerTok.Attr {
				if attr.Name.Local != "data-lift" {
					snippetHTML = snippetHTML + " " + attr.Name.Local + "=\"" + attr.Value + "\""
				} else {
					scopeFunction = attr.Value
				}
			}
			snippetHTML = snippetHTML + ">"
			err, super := processSnippet(decoder, "", scopeFunction)
			if err != nil {
				return err, ""
			}
			if strings.HasSuffix(snippetHTML, ">") {
				snippetHTML = snippetHTML + super
			} else {
				snippetHTML = snippetHTML + ">" + super
			}
			open++
		case xml.CharData:
			snippetHTML = snippetHTML + string(innerTok)
		case xml.EndElement:
			snippetHTML = snippetHTML + "</" + innerTok.Name.Local + ">"
			closingTags++
			if open == closingTags { //do we have our matching closing tag? //This fails with autoclose tags I think
				//ch <- snippetAndNode{scopeFunction, snippetHTML}

				if scopeFunction != "" {
					FunctionMap.RLock()
					f, found := FunctionMap.M[scopeFunction]
					FunctionMap.RUnlock()
					if found {
						scopeFunction = ""
						return nil, f(snippetHTML)
					} else {
						return errors.New("Did not find function: '" + scopeFunction + "'"), ""
					}
				} else {
					return nil, snippetHTML
				}

			}
		default:
			fmt.Errorf("\n1- If you are here, you are missing a type: %v\n", innerTok)
		}
	}

}*/

//readChan receives the snippet name and the html we will work on.
//So far is an alternative way to process snippets. We wil lcompare speeds once
//we are further along
/*func readChan(ch chan snippetAndNode) {
	var buffer snippetAndNode
	for {
		select {
		case data, ok := <-ch:
			if ok == true { //if is false when you close the channel
				buffer = data
				fmt.Printf("Found snippet: %v\n", buffer.FunctionName)
				//FunctionMap.RLock()
				//f := FunctionMap.m[buffer.FunctionName]
				//FunctionMap.RUnlock()
				//fmt.Printf("result is : %v\n", f(buffer.HTML))
			}
		}
	}

}
*/

func Process(in string) string {
	//in is the html we get from the template
	node, _ := h5.NewFromString(in)
	//because the html may not be a full page, we use the Partial* function
	//node, _ := h5.PartialFromString(in)
	node.Walk(walkTree)
	return in
}

func walkTree(n *html.Node) {

	if len(n.Attr) == 0 {
		buf := bytes.NewBufferString("")
		html.Render(buf, n)
		fmt.Println("node1 " + string(buf.Bytes()))
	}
	for _, attr := range n.Attr {
		if attr.Key == "data-lift" {
			transformedNode, err := do(attr.Val, n)
			if err != nil {
				fmt.Errorf("We got error %+v", err.Error())
			}
			buf := bytes.NewBufferString("")
			html.Render(buf, transformedNode)
			fmt.Println("node2 " + string(buf.Bytes()))
		}
	}
}

func do(scopeFunction string, snippetHTML *html.Node) (*html.Node, error) {
	FunctionMap.RLock()
	f, found := FunctionMap.M[scopeFunction]
	FunctionMap.RUnlock()
	if found {
		return f(snippetHTML), nil
	} else {
		return &html.Node{}, errors.New("Did not find function: '" + scopeFunction + "'")
	}
}

/*for _, value1 := range node {
	//apply the css selector to get a []*html.Node of matching nodes
	ret := dataLiftSelector.Find(value1)
	for _, value2 := range ret {
		//if we wanted to have a *Tree of the nodes, use this
		t := h5.NewTree(value2)
		//here we loop over the attributes of the matching node
		for _, attr := range value2.Attr {
			fmt.Println("Function Name: " + attr.Val)
			fmt.Println("html to process: " + h5.RenderNodesToString([]*html.Node{value2}))
			err, result := do(attr.Val, value2)
			if err != nil {
				fmt.Errorf("Error was %+v", err.Error())
			}
			fmt.Println("result is " + h5.RenderNodesToString([]*html.Node{result}))
		}
		return t.String()
	}
}*/
