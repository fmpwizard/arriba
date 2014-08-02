package arriba

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type HTMLTransform func(string) string

/*type snippetAndNode struct {
	FunctionName string
	HTML         string
}
*/
var FunctionMap = struct {
	sync.RWMutex
	M map[string]HTMLTransform
}{M: make(map[string]HTMLTransform)}

/*var ch = make(chan snippetAndNode)*/

func MarshallElem(in string) string {
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
			if len(element.Attr) == 0 {
				completeHTML = completeHTML + "<" + element.Name.Local
			}
			for _, attr := range element.Attr {
				err, res := processSnippet(attr, decoder, element.Name.Local, "")
				if err != nil {
					return err.Error()
				}
				completeHTML = completeHTML + res
			}
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

func processSnippet(currentAttr xml.Attr, decoder *xml.Decoder, parentTag string, scopeFunction string) (error, string) {
	snippetHTML := ""
	if parentTag != "" && currentAttr.Name.Local != "data-lift" {
		snippetHTML = "<" + parentTag + " " + currentAttr.Name.Local + "=\"" + currentAttr.Value + "\">"
	} else if parentTag != "" {
		snippetHTML = "<" + parentTag + ">"
	}

	open := 1
	closingTags := 0

	//if currentAttr.Name.Local != "data-lift" {
	for {
		if currentAttr.Name.Local == "data-lift" {
			//fmt.Println("0 " + currentAttr.Value)
			scopeFunction = currentAttr.Value
		}

		//fmt.Println("1 current attr " + currentAttr.Name.Local)
		//fmt.Println("2 scopeFunction: ===> " + scopeFunction)
		tok, err := decoder.Token()
		if err != nil {
			//fmt.Println("6 " + err.Error())
			//We are done processing tokens, let's end.
			/*close(ch)*/
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
				}
				err, super := processSnippet(attr, decoder, "", scopeFunction)
				if err != nil {
					return err, ""
				}
				if strings.HasSuffix(snippetHTML, ">") {
					snippetHTML = snippetHTML + super
				} else {
					snippetHTML = snippetHTML + ">" + super
				}
			}
			open++
		case xml.CharData:
			snippetHTML = snippetHTML + string(innerTok)
		case xml.EndElement:
			snippetHTML = snippetHTML + "</" + innerTok.Name.Local + ">"
			closingTags++
			if open == closingTags { //do we have our matching closing tag? //This fails with autoclose tags I think
				//ch <- snippetAndNode{value.Value, snippetHTML}

				//if value.Name.Local == "data-lift" {
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

	//}
	//return nil, ""
}

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
