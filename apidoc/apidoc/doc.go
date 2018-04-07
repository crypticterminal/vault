package apidoc

import (
	"strings"

	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/version"
)

// Documente is a set a API documentation. The structure of it and its descendants roughly
// follow the organization of OpenAPI, but it is not rigidly tied to that. Additional elements
// can be added, and many OpenAPI constructs are omitted. It is meant as an itermediate format
// from which OpenAPI or other targets can be generated.
type Document struct {
	Version string
	//Paths   []Path
	Mounts map[string][]Path
}

// NewDoc initialized a new, empty Document.
func NewDoc() Document {
	return Document{
		Version: version.GetVersion().Version,
		//Paths:   make([]Path, 0),
		Mounts: make(map[string][]Path),
	}
}

func (d *Document) Add(mount string, p ...Path) {
	d.Mounts[mount] = append(d.Mounts[mount], p...)
}

func (d *Document) PathsLikeBefore() []Path {
	var paths []Path
	for _, pathe := range d.Mounts {
		paths = append(paths, pathe...)
	}
	return paths
}

func (d *Document) LoadBackend(prefix string, backend *framework.Backend) {
	for _, p := range backend.Paths {
		paths := procLogicalPath(prefix, p)
		//d.Paths = append(d.Paths, paths...)
		d.Mounts[prefix] = append(d.Mounts[prefix], paths...)
	}
}

func (d *Document) SortPaths() {
	//sort.Slice(d.Paths, func(i, j int) bool {
	//	return d.Paths[i].Pattern < d.Paths[j].Pattern
	//})
}

// Path is the structure for a single path, including all of its methods.
// The path description is kept split as: /<prefix>/<pattern>
type Path struct {
	Pattern string
	Methods map[string]Method
}

// NewPath returns a new Path.
func NewPath(pattern string) Path {
	return Path{
		Pattern: pattern,
		Methods: make(map[string]Method),
	}
}

type Method struct {
	Summary    string
	Parameters []Parameter
	BodyProps  []Property
	Responses  []Response
	//Tags       []string
}

//func NewMethod(summary string, tags ...string) Method {
func NewMethod(summary string) Method {
	return Method{
		Summary: summary,
		//	Tags:    tags,
	}
}

func (m *Method) AddResponse(code int, example string) {
	var description string
	switch code {
	case 200:
		description = "OK"
	case 204:
		description = "empty body"
	}
	m.Responses = append(m.Responses, NewResponse(code, description, example))
}

type Property struct {
	Name        string
	Type        string
	SubType     string
	Description string
}

func NewProperty(name, typ, description string) Property {
	p := Property{
		Name:        name,
		Description: description,
	}
	parts := strings.Split(typ, "/")
	p.Type = parts[0]
	if len(parts) > 1 && parts[0] == "array" {
		p.SubType = parts[1]
	}

	return p
}

type Parameter struct {
	Property Property
	In       string
}

type Response struct {
	Code        int
	Description string
	Example     string
}

func NewResponse(code int, description, example string) Response {
	example = strings.TrimSpace(example)
	example = strings.Replace(example, "\t", "  ", -1)
	return Response{
		Code:        code,
		Description: description,
		Example:     example,
	}
}

var StdRespOK = Response{
	Code:        200,
	Description: "OK",
}

var StdRespNoContent = Response{
	Code:        204,
	Description: "empty body",
}
