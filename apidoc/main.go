package main

import (
	"os"

	"github.com/hashicorp/vault/apidoc/apidoc"
	"github.com/hashicorp/vault/builtin/logical/aws"
	"github.com/hashicorp/vault/vault"
)

func main() {
	aws_be := aws.Backend()

	doc := apidoc.NewDoc()
	doc.LoadBackend("aws", aws_be.Backend)
	doc.LoadBackend("sys", vault.Backend())
	doc.Add("sys", vault.ManualPaths()...)

	r := apidoc.OAPIRenderer{
		Output:   os.Stdout,
		Template: apidoc.Tmpl,
		Version:  2,
	}
	r.Render(doc)
}
