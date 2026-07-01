package analyzer

import (
	"afv/internal/model"
	"strings"
)

func ExtractEndpoints(root model.ASTNode) []model.Endpoint {
	endpoints := []model.Endpoint{}

	Walk(root, nil, func(node model.ASTNode, parent *model.ASTNode) {
		if !isEndpoint(node) {
			return
		}

		endpoint := extractEndpoint(node)
		endpoints = append(endpoints, endpoint)
	})

	return endpoints
}

func isEndpoint(node model.ASTNode) bool {
	if node.Type != "CallExpression" {
		return false
	}

	for _, child := range node.Children {

		if child.Role != "callee" {
			continue
		}

		// fetch
		if child.Type == "Identifier" && child.Name == "fetch" {
			return true
		}

		// axios.xxx / xhr.open
		if child.Type == "MemberExpression" {

			hasAxios := false
			hasHTTPMethod := false
			hasOpen := false

			for _, child2 := range child.Children {

				if child2.Role == "object" && child2.Name == "axios" {
					hasAxios = true
				}

				if child2.Role == "property" && isHTTPMethod(child2.Name) {
					hasHTTPMethod = true
				}

				if child2.Role == "property" && child2.Name == "open" {
					hasOpen = true
				}
			}

			if hasAxios && hasHTTPMethod {
				return true
			}

			if hasOpen {
				return true
			}
		}
	}

	return false
}

func isHTTPMethod(name string) bool {
	switch name {
	case "get", "post", "put", "delete", "patch", "head", "options":
		return true
	default:
		return false
	}
}

func extractEndpoint(node model.ASTNode) model.Endpoint {

	endpoint := model.Endpoint{}
	endpoint.Location = node.Location

	for _, child := range node.Children {

		if child.Role != "callee" {
			continue
		}

		// fetch
		if child.Type == "Identifier" && child.Name == "fetch" {
			endpoint.Method = parseFetchMethod(node)
			endpoint.URL = parseFetchURL(node)
			return endpoint
		}

		// axios.xxx / xhr.open
		if child.Type == "MemberExpression" {

			object, property := getMemberParts(child)

			if object == "axios" {
				endpoint.Method = strings.ToUpper(property)
				endpoint.URL = parseAxiosURL(node)
				return endpoint
			}

			if property == "open" {
				endpoint.Method = parseXHRMethod(node)
				endpoint.URL = parseXHRURL(node)
				return endpoint
			}
		}
	}

	return endpoint
}

func parseFetchMethod(node model.ASTNode) string {

	for _, arg := range node.Children {

		if arg.Role != "arguments" || arg.Type != "ObjectExpression" { // search for an aurguments that is an objects. ({method: "POST"} (an object expression))
			continue
		}

		for _, prop := range arg.Children {

			if prop.Type != "ObjectProperty" { // only real properies (key-values) not every child.
				continue
			}

			var key string
			var value string

			for _, child := range prop.Children { // iterate on key-value pairs

				if child.Role == "key" {
					key = child.Name
				}

				if child.Role == "value" {
					if v, ok := child.Value.(string); ok {
						value = v
					}
				}
			}

			if key == "method" {
				return strings.ToUpper(value)
			}
		}
	}

	return "GET"
}

func parseFetchURL(node model.ASTNode) string {

	argIndex := 0

	for _, arg := range node.Children {

		if arg.Role != "arguments" {
			continue
		}

		if argIndex == 0 {

			if v, ok := arg.Value.(string); ok {
				return v
			}

			return "<dynamic>"
		}

		argIndex++
	}

	return ""
}

func parseAxiosURL(node model.ASTNode) string {

	argIndex := 0

	for _, arg := range node.Children {

		if arg.Role != "arguments" {
			continue
		}

		if argIndex == 0 {

			if v, ok := arg.Value.(string); ok {
				return v
			}

			return "<dynamic>"
		}

		argIndex++
	}

	return ""
}

func parseXHRMethod(node model.ASTNode) string {

	argIndex := 0

	for _, arg := range node.Children {

		if arg.Role != "arguments" {
			continue
		}

		if argIndex == 0 {

			if v, ok := arg.Value.(string); ok {
				return strings.ToUpper(v)
			}

			return "<dynamic>"
		}

		argIndex++
	}

	return ""
}

func parseXHRURL(node model.ASTNode) string {

	argIndex := 0

	for _, arg := range node.Children {

		if arg.Role != "arguments" {
			continue
		}

		if argIndex == 1 {

			if v, ok := arg.Value.(string); ok {
				return v
			}

			return "<dynamic>"
		}

		argIndex++
	}

	return ""
}

func getMemberParts(node model.ASTNode) (object string, property string) {

	for _, child := range node.Children {

		if child.Role == "object" {
			object = child.Name
		}

		if child.Role == "property" {
			property = child.Name
		}
	}

	return object, property
}
