package apidoc

import (
	"sort"
	"strings"

	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/version"
)

type Doc struct {
	Version string
	Paths   []Path
}

func NewDoc() Doc {
	return Doc{
		Version: version.GetVersion().Version,
		Paths:   make([]Path, 0),
	}
}

func (d *Doc) Add(p ...Path) {
	d.Paths = append(d.Paths, p...)
}

func (d *Doc) LoadBackend(prefix string, backend *framework.Backend) {
	for _, p := range backend.Paths {
		paths := procLogicalPath(prefix, p)
		d.Paths = append(d.Paths, paths...)
	}
}
func (d *Doc) LoadBackend2(prefix string, backendPaths []*framework.Path) {
	for _, p := range backendPaths {
		paths := procLogicalPath(prefix, p)
		d.Paths = append(d.Paths, paths...)
	}
}

func (d *Doc) SortPaths() {
	sort.Slice(d.Paths, func(i, j int) bool {
		return d.Paths[i].Pattern < d.Paths[j].Pattern
	})
}

type Path struct {
	Pattern string
	Methods map[string]Method
}

func NewPath(pattern string) Path {
	return Path{
		Pattern: pattern,
		Methods: make(map[string]Method),
	}
}

type Method struct {
	Summary    string
	Tags       []string
	Parameters []Parameter
	BodyProps  []Property
	Responses  []Response
}

func NewMethod(summary string, tags ...string) Method {
	return Method{
		Summary: summary,
		Tags:    tags,
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
