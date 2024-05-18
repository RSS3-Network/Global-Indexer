package docs

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/a-h/rest"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type RouteOption struct {
	// HTTP Method
	Method string
	// the list of group names of the echo router, which is used as tags
	GroupNames []string
	// the path of the route in echo's format
	Path string
	// the handler function for echo
	Handler func(echo.Context) error
}

// GetCompletePath outputs the complete path for rest.API
func (r *RouteOption) GetCompletePath() string {
	var builder strings.Builder
	for _, name := range r.GroupNames {
		builder.WriteString(name)
	}

	builder.WriteString(r.Path)

	return replacePathParams(builder.String())
}

func (r *RouteOption) GetTags() []string {
	tags := make([]string, len(r.GroupNames))
	caser := cases.Title(language.English)

	caps := []string{"dsl", "nta", "rss"}

	for i, name := range r.GroupNames {
		name = strings.TrimPrefix(name, "/")
		if lo.Contains(caps, name) {
			tags[i] = strings.ToUpper(name)
		} else {
			tags[i] = caser.String(name)
		}
	}

	return tags
}

// replacePathParams uses regex to replace `:abc` with `{abc}`
func replacePathParams(input string) string {
	path := regexp.MustCompile(`:([a-zA-Z0-9_]+)`).ReplaceAllString(input, `{$1}`)

	//return strings.TrimSuffix(path, "/")
	return path
}

// ConvertStructToParams converts a struct type to a Params struct
func ConvertStructToParams[T any]() rest.Params {
	params := rest.Params{
		Path:  make(map[string]rest.PathParam),
		Query: make(map[string]rest.QueryParam),
	}

	// Get the reflect type of the generic type T
	typ := reflect.TypeOf((*T)(nil)).Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Process query parameters
		queryTag := field.Tag.Get("query")
		queryDescription := field.Tag.Get("description")

		if queryTag != "" {
			if queryDescription == "" {
				queryDescription = strings.ReplaceAll(queryTag, "_", " ")
			}

			params.Query[queryTag] = rest.QueryParam{
				Description:       queryDescription,
				ApplyCustomSchema: getCustomSchemaFunc(queryTag, reflect.TypeOf(rest.QueryParam{})),
			}
		}

		// Process path parameters
		paramTag := field.Tag.Get("param")
		paramDescription := field.Tag.Get("description")

		if paramTag != "" {
			if paramDescription == "" {
				paramDescription = strings.ReplaceAll(paramTag, "_", " ")
			}

			params.Path[paramTag] = rest.PathParam{
				Description:       paramDescription,
				ApplyCustomSchema: getCustomSchemaFunc(paramTag, reflect.TypeOf(rest.PathParam{})),
			}
		}
	}

	return params
}

// getCustomSchemaFunc returns a function to apply custom schema based on the tag and parameter type
func getCustomSchemaFunc(tag string, paramType reflect.Type) func(*openapi3.Parameter) {
	return func(s *openapi3.Parameter) {
		s.Name = tag
		// Customize schema example based on field name or type
		// FIXME: add more examples, also add description to request struct fields
		switch tag {
		case "cursor":
			s.Example = "0xb7876dbaf18cbed6378fa5540681cedb6c36507df99fc4e9e3b2889fc8842cbc"
		case "sender", "receiver", "address":
			s.Example = "0x0000000000000000000000000000000000000000"
		case "type":
			s.Example = "deposit"
		case "limit":
			s.Example = "20"
		case "transaction_hash":
			if paramType == reflect.TypeOf(rest.PathParam{}) {
				s.Example = "0xb7876dbaf18cbed6378fa5540681cedb6c36507df99fc4e9e3b2889fc8842cbc"
			}
		}
	}
}
