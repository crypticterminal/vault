package apidoc

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"text/template"
)

type OAPIRenderer struct {
	Output   io.Writer
	Template string
	Version  int
}

func (r *OAPIRenderer) Render(doc Document) {
	funcs := map[string]interface{}{
		"indent": funcIndent,
	}

	doc.SortPaths()
	//fmt.Println(doc.PathsLikeBefore())
	tmpl, _ := template.New("root").Funcs(funcs).Parse(r.Template)
	tmpl.Execute(r.Output, &doc)
}

func funcIndent(count int, text string) string {
	var buf bytes.Buffer
	prefix := strings.Repeat(" ", count)
	scan := bufio.NewScanner(strings.NewReader(text))
	for scan.Scan() {
		buf.WriteString(prefix + scan.Text() + "\n")
	}

	return strings.TrimRight(buf.String(), "\n")
}
