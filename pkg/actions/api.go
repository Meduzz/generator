package actions

import (
	"fmt"
	"strings"

	"github.com/Meduzz/helper/fp/slice"
)

type (
	Kind    string
	SubKind string
	Source  string

	API struct {
		Name        string
		Method      string
		Path        string
		Params      []*Param
		ContentType string
	}

	// TODO have a required property?
	Param struct {
		Name    string
		Kind    Kind
		Subkind SubKind
		Source  Source
		Value   any
	}
)

const (
	Path   Kind = "path"
	Query  Kind = "query"
	Header Kind = "header"

	String SubKind = "string"
	Int    SubKind = "int"
	Map    SubKind = "map"
	Array  SubKind = "array"

	Flag Source = "flag"
	Env  Source = "env"
)

var (
	QueryHelpers = map[SubKind]func(*Param) string{
		// helps turn map of string -> string into gin query maps.
		Map: func(p *Param) string {
			valueMap, ok := p.Value.(map[string]string)

			if !ok {
				return ""
			}

			pairs := make([]string, 0)
			for k, v := range valueMap {
				pairs = append(pairs, fmt.Sprintf("%s[%s]=%s", p.Name, k, v))
			}

			return strings.Join(pairs, "&")
		},
		// helps turn arrays into gin query array.
		Array: func(p *Param) string {
			valueArray, ok := p.Value.([]string)

			if !ok {
				return ""
			}

			pairs := slice.Map(valueArray, func(v string) string {
				return fmt.Sprintf("%s[]=%s", p.Name, v)
			})

			return strings.Join(pairs, "&")
		},
	}

	HeaderHelpers = map[SubKind]func(*Param) string{
		Map: func(p *Param) string {
			mapData, ok := p.Value.(map[string]string)

			if !ok {
				return ""
			}

			pairs := make([]string, 0)

			for k, v := range mapData {
				pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
			}

			return strings.Join(pairs, ";")
		},
		Array: func(p *Param) string {
			arrayData, ok := p.Value.([]string)

			if !ok {
				return ""
			}

			return strings.Join(arrayData, ",")
		},
	}
)

func (p *Param) Query() string {
	helper, ok := QueryHelpers[p.Subkind]

	if !ok {
		helper = func(p *Param) string {
			return fmt.Sprintf("%s=%v", p.Name, p.Value)
		}
	}

	return helper(p)
}

func (p *Param) Header() string {
	helper, ok := HeaderHelpers[p.Subkind]

	if !ok {
		helper = func(p *Param) string {
			return fmt.Sprintf("%v", p.Value)
		}
	}

	return helper(p)
}
