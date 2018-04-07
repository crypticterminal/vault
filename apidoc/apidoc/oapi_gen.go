package apidoc

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

var optRe = regexp.MustCompile(`(?U)\(.*\)\?`)
var cleanRe = regexp.MustCompile("[()$?]")
var reqdRe = regexp.MustCompile(`\(\?P<(\w+)>[^)]*\)`)

type pathlet struct {
	pattern string
	params  []string
}

func deRegex(s string) string {
	return cleanRe.ReplaceAllString(s, "")
}

func convertType(t framework.FieldType) (string, string) {
	var ret, sub string

	switch t {
	case framework.TypeString, framework.TypeNameString, framework.TypeKVPairs:
		ret = "string"
	case framework.TypeInt, framework.TypeDurationSecond:
		ret = "number"
	case framework.TypeBool:
		ret = "boolean"
	case framework.TypeMap:
		ret = "object"
	case framework.TypeSlice, framework.TypeStringSlice, framework.TypeCommaStringSlice:
		ret = "array"
		sub = "string"
	case framework.TypeCommaIntSlice:
		ret = "array"
		sub = "number"
	default:
		panic(fmt.Sprintf("Unsupported type %d", t))
	}

	return ret, sub
}

// expandPattern expands a regex pattern by generating permutations of any optional parameters
// and changing named parameters into their {open_api} style equivalents.
func expandPattern(root, pat string) []pathlet {
	// This construct is added by GenericNameRegex and is much easier to remove now
	// than compensate for in the other regexes.
	pat = strings.Replace(pat, `\w(([\w-.]+)?\w)?`, "", -1)

	toppaths := []string{pat}

	// expand all optional elements into two paths. This apporach really only useful up to 2 optional
	// groups, but we probably don't want to deal with the exponential increase for the general case.
	for i := 0; i < len(toppaths); i++ {
		p := toppaths[i]
		match := optRe.FindStringIndex(p)
		if match != nil {
			toppaths[i] = p[0:match[0]] + p[match[0]+1:match[1]-2] + p[match[1]:]
			toppaths = append(toppaths, p[0:match[0]]+p[match[1]:])
			i--
		}

	}

	sort.Strings(toppaths)

	paths := make([]pathlet, 0)
	for _, pat := range toppaths {
		var params []string
		result := reqdRe.FindAllStringSubmatch(pat, -1)
		if result != nil {
			for _, p := range result {
				par := p[1]
				params = append(params, par)
				pat = strings.Replace(pat, p[0], fmt.Sprintf("{%s}", par), 1)
			}
		}
		pat = fmt.Sprintf("/%s/%s", root, pat)
		pat = deRegex(pat)
		paths = append(paths, pathlet{pat, params})
	}
	return paths
}

// TODO: this is conservative. Should omit surrounding quotes if not needed.
func prepareString(syn string) string {
	syn = strings.TrimSpace(syn)
	if idx := strings.Index(syn, "\n"); idx != -1 {
		syn = syn[0:idx] + "…"
	}
	syn = strings.Replace(syn, `"`, `\"`, -1)
	return fmt.Sprintf(`"%s"`, syn)
}

func procLogicalPath(prefix string, p *framework.Path) []Path {
	var docPaths []Path
	var httpMethod string

	paths := expandPattern(prefix, p.Pattern)

	for _, path := range paths {
		methods := make(map[string]Method)

		for opType := range p.Callbacks {
			m := Method{
				Summary: prepareString(p.HelpSynopsis),
				//Tags:    []string{prefix},
			}
			switch opType {
			case logical.CreateOperation:
				httpMethod = "post"
				m.Responses = []Response{StdRespNoContent}
			case logical.UpdateOperation:
				httpMethod = "put"
				m.Responses = []Response{StdRespNoContent}
			case logical.DeleteOperation:
				httpMethod = "delete"
				m.Responses = []Response{StdRespNoContent}
			case logical.ReadOperation, logical.ListOperation:
				httpMethod = "get"
				m.Responses = []Response{StdRespOK}
			default:
				panic(fmt.Sprintf("unknown operation type %v", opType))
			}

			d := make(map[string]bool)
			for _, param := range path.params {
				d[param] = true
				typ, sub := convertType(p.Fields[param].Type)
				m.Parameters = append(m.Parameters, Parameter{
					In: "path",
					Property: Property{
						Name:        param,
						Type:        typ,
						SubType:     sub,
						Description: prepareString(p.Fields[param].Description),
					},
				})
			}

			// It's assumed that any fields not present in the path can be part of
			// the body for POST methods.
			if httpMethod == "post" || httpMethod == "put" {
				for name, field := range p.Fields {
					if _, ok := d[name]; !ok {
						typ, sub := convertType(field.Type)
						m.BodyProps = append(m.BodyProps, Property{
							Name:        name,
							Description: prepareString(field.Description),
							Type:        typ,
							SubType:     sub,
						})
					}
				}
			}
			methods[httpMethod] = m
		}
		if len(methods) > 0 {
			pd := Path{
				Pattern: path.pattern,
				Methods: methods,
			}
			docPaths = append(docPaths, pd)
		}
	}

	return docPaths
}
