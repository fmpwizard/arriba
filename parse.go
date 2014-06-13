package arriba

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

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
	completeHTML := ""
	parentTag := ""

	decoder := xml.NewDecoder(bytes.NewBufferString(in))

	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch element := token.(type) {
		case xml.StartElement:

			if completeHTML != "" {
				completeHTML = completeHTML + "<" + element.Name.Local
			} else {
				completeHTML = "<" + element.Name.Local
			}
			functionName := ""
			for _, value := range element.Attr {
				parentTag = parentTag + " " + value.Name.Local + "=\"" + value.Value + "\""
				if value.Name.Local == "data-lift" {
					_, res := processSnippet(value, decoder, parentTag)
					parentTag = ""
					completeHTML = completeHTML + res

				}
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
	open := 1
	closingTags := 0

	for {
		tok, err := decoder.Token()
		if err != nil {
			return err, ""
		}
		switch innerTok := tok.(type) {
		case xml.StartElement:
			if snippetHTML == "" {
				snippetHTML = parentTag + ">" //we found first inner node, so close the parent
			}

			snippetHTML = snippetHTML + "<" + innerTok.Name.Local
			for _, attr := range innerTok.Attr {
				snippetHTML = snippetHTML + " " + attr.Name.Local + "=\"" + attr.Value + "\""
			}
			snippetHTML = snippetHTML + ">"
			open++
		case xml.CharData:
			snippetHTML = snippetHTML + string(innerTok)

		case xml.EndElement:
			snippetHTML = snippetHTML + "</" + innerTok.Name.Local + ">"
			closingTags++
			if open == closingTags { //do we have our matching closing tag? //This fails with autoclose tags I think
				return nil, ChangeName(snippetHTML)
			}
		}
		fmt.Printf(" ==>> snippetHTML  %v\n", snippetHTML)
	}
}

func ChangeTime(html string) string {
	return strings.Replace(html, "Time goes here", time.Now().Format("2006-01-02T15:04:05.999999999Z07:00"), 1)
}

func ChangeName(html string) string {
	return strings.Replace(html, "Diego", "Gabriel", 1)
}
