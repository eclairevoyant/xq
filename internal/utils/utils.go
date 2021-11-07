package utils

import (
	"encoding/xml"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/fatih/color"
	"strings"
)

func FormatXml(str string) string {
	decoder := xml.NewDecoder(strings.NewReader(str))
	level := 0
	hasContent := false
	result := new(strings.Builder)

	tagColor := color.New(color.FgYellow).SprintFunc()
	attrColor := color.New(color.FgGreen).SprintFunc()
	commentColor := color.New(color.FgHiBlue).SprintFunc()

	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		switch typedToken := token.(type) {
		case xml.ProcInst:
			_, _ = fmt.Fprintf(result, "%s%s %s%s\n", tagColor("<?"), typedToken.Target, string(typedToken.Inst), tagColor("?>"))
		case xml.StartElement:
			if level > 0 {
				_, _ = fmt.Fprint(result, "\n", strings.Repeat("  ", level))
			}
			var attrs []string
			for _, attr := range typedToken.Attr {
				attrs = append(attrs, attr.Name.Local + attrColor("=\"" + attr.Value + "\""))
			}
			attrsStr := strings.Join(attrs, " ")
			if attrsStr != "" {
				attrsStr = " " + attrsStr
			}
			_, _ = fmt.Fprint(result, tagColor("<" + typedToken.Name.Local) + attrsStr + tagColor(">"))
			level++
		case xml.CharData:
			str := string(typedToken)
			str = strings.TrimSpace(str)
			_, _ = fmt.Fprint(result, str)
			hasContent = str != ""
		case xml.Comment:
			if !hasContent && level > 0 {
				_, _ = fmt.Fprint(result, "\n", strings.Repeat("  ", level))
			}
			_, _ = fmt.Fprint(result, commentColor("<!--" + string(typedToken) + "-->"))
			if level == 0 {
				fmt.Fprint(result, "\n")
			}
		case xml.EndElement:
			level--
			if !hasContent {
				_, _ = fmt.Fprint(result, "\n", strings.Repeat("  ", level))
			}
			_, _ = fmt.Fprint(result, tagColor("</" + typedToken.Name.Local + ">"))
			hasContent = false
		default:
		}
	}

	return result.String()
}

func XPathQuery(str string, query string) string {
	result := new(strings.Builder)

	doc, err := xmlquery.Parse(strings.NewReader(str))
	if err != nil {
		panic(err)
	}

	for _, n := range xmlquery.Find(doc, query) {
		_, _ = fmt.Fprintf(result, "%s\n", n.InnerText())
	}

	return result.String()
}