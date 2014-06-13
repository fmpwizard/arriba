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
	snipetHTML := ""

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
			open := 1
			close := 0
			functionName := ""
		Loop:
			for _, value := range element.Attr {
				//fmt.Printf("Att: %v ==> value: %v\n", value.Name.Local, value.Value)
				parentTag = parentTag + " " + value.Name.Local + "=\"" + value.Value + "\""
				if value.Name.Local == "data-lift" {
					functionName = value.Value

					for {
						tok, err := decoder.Token()
						if err != nil {
							return err.Error()
						}
						switch innerTok := tok.(type) {
						case xml.StartElement:
							if snipetHTML == "" {
								snipetHTML = parentTag + ">" //we found first inner node, so close the parent
							}

							snipetHTML = snipetHTML + "<" + innerTok.Name.Local
							for _, attr := range innerTok.Attr {
								snipetHTML = snipetHTML + " " + attr.Name.Local + "=\"" + attr.Value + "\""
							}
							snipetHTML = snipetHTML + ">"
							open++
						case xml.CharData:
							snipetHTML = snipetHTML + string(innerTok)

						case xml.EndElement:
							snipetHTML = snipetHTML + "</" + innerTok.Name.Local + ">"
							close++
							if open == close { //do we have our matching closing tag? //This fails with autoclose tags I think
								rawHTML := snipetHTML
								completeHTML = completeHTML + ChangeName(rawHTML) //hard coded for now
								snipetHTML = ""
								parentTag = ""
								break Loop
							}
						}
						fmt.Printf(" ==>> snipetHTML  %v\n", snipetHTML)
					}

				}
			}
			/*if snipetHTML != "" {
				rawHTML := snipetHTML
				completeHTML = completeHTML + ChangeName(rawHTML) //hard coded for now
			}

			*/
			//fmt.Printf("Start: %v\n", element.Name.Local)
			if !strings.HasSuffix(completeHTML, ">") {
				completeHTML = completeHTML + ">"
			}

			//fmt.Printf("completeHTML: %v\n", completeHTML)
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

func ChangeTime(html string) string {
	return strings.Replace(html, "Time goes here", time.Now().Format("2006-01-02T15:04:05.999999999Z07:00"), 1)
}

func ChangeName(html string) string {
	return strings.Replace(html, "Diego", "Gabriel", 1)
}
