package arriba

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

type HTMLTransform func(string) string

var funcMap = make(map[string]HTMLTransform)
var ch = make(chan string)

func MarshallElem(in string) string {
	fmt.Println("\n\n\n\n\n")

	go readChan(ch)
	funcMap["ChangeName"] = ChangeName
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
				_, res := processSnippet(value, decoder, element.Name.Local)
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

	//functionName := value.Value

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
					ch <- snippetHTML
					_, super := processSnippet(attr, decoder, "")
					if strings.HasSuffix(snippetHTML, ">") {
						snippetHTML = snippetHTML + super
					} else {
						snippetHTML = snippetHTML + ">" + super
					}
					ch <- snippetHTML
				}
				open++
			case xml.CharData:
				snippetHTML = snippetHTML + string(innerTok)
			case xml.EndElement:
				snippetHTML = snippetHTML + "</" + innerTok.Name.Local + ">"
				closingTags++
				if open == closingTags { //do we have our matching closing tag? //This fails with autoclose tags I think
					ch <- snippetHTML
					return nil, ChangeName(snippetHTML)
				}
			}
		}

	}
	return nil, ""
}

func ChangeName(html string) string {
	return strings.Replace(html, "Diego", "Gabriel", 1)
}

func readChan(ch chan string) {
	var buffer string
	for {
		select {
		case data, ok := <-ch:
			if ok == true {
				fmt.Printf("got: %v\n", data)
				buffer = data
				//buffer =  buffer + data
			} else {
				fmt.Printf("sending: %v\n", buffer)
				return
			}
		}
	}

}
