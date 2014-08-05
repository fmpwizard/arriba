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
