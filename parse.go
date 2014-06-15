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

type snippetAndNode struct {
	FunctionName string
	HTML         string
}

var functionMap = struct {
	sync.RWMutex
	m map[string]HTMLTransform
}{m: make(map[string]HTMLTransform)}

var ch = make(chan snippetAndNode)

func MarshallElem(in string) string {
	//fmt.Println("\n\n\n\n\n1")
	functionMap.Lock()
	functionMap.m["ChangeName"] = ChangeName
	functionMap.m["ChangeLastName"] = ChangeLastName
	functionMap.Unlock()
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
			functionName := ""
			for _, value := range element.Attr {
				err, res := processSnippet(value, decoder, element.Name.Local)
				if err != nil {
					return err.Error()
				}
				completeHTML = completeHTML + res
			}
			if !strings.HasSuffix(completeHTML, ">") {
				completeHTML = completeHTML + ">"
			}
			if functionName != "" {
				fmt.Printf("functionName: %v\n", functionName)
			}
		case xml.CharData:
			fmt.Printf("CharData: %+v\n", string(element))
		case xml.EndElement:
			completeHTML = completeHTML + "</" + element.Name.Local + ">"
		case xml.Comment:
			fmt.Printf("Comment: %+v\n", element)
		case xml.Directive:
			fmt.Printf("Directive: %+v\n", string(element))
		case xml.Token:
			fmt.Printf("4: %+v\n", element)

		default:
			fmt.Errorf("\nIf yo uare here, you are missing a type: %v\n", element)
		}

	}
	return completeHTML
}

func processSnippet(value xml.Attr, decoder *xml.Decoder, parentTag string) (error, string) {
	snippetHTML := ""
	if parentTag != "" {
		snippetHTML = "<" + parentTag + ">"
	}
	open := 1
	closingTags := 0

	if value.Name.Local == "data-lift" {
		for {
			tok, err := decoder.Token()
			if err != nil {
				//We are done processing tokens, let's end.
				close(ch)
				return err, snippetHTML
			}
			switch innerTok := tok.(type) {
			case xml.StartElement:
				snippetHTML = snippetHTML + "<" + innerTok.Name.Local
				for _, attr := range innerTok.Attr {
					if attr.Name.Local != "data-lift" {
						snippetHTML = snippetHTML + " " + attr.Name.Local + "=\"" + attr.Value + "\""
					}
					err, super := processSnippet(attr, decoder, "")
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
					functionMap.RLock()
					f, ok := functionMap.m[value.Value]
					functionMap.RUnlock()
					if ok {
						return nil, f(snippetHTML)
					} else {
						return errors.New("Did not find function " + value.Value), ""
					}

				}
			}
		}

	}
	return nil, ""
}

func ChangeName(html string) string {
	return strings.Replace(html, "Diego", "Gabriel", 1)
}

func ChangeLastName(html string) string {
	return strings.Replace(html, "Medina", "Bauman", 1)
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
				//functionMap.RLock()
				//f := functionMap.m[buffer.FunctionName]
				//functionMap.RUnlock()
				//fmt.Printf("result is : %v\n", f(buffer.HTML))
			}
		}
	}

}
*/
