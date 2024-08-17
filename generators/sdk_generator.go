package generators

import (
	"fmt"
	"archive/zip"
	"bytes"
	"strings"
	"log"
)

func GenerateSDK(spec map[string]interface{}, apiKey string, language string, languageOptions map[string]interface{}) ([]byte, error) {
	var sdk []byte
	var err error

	log.Println("Generating SDK for language:", language)

	switch language {
	case "nodejs-fetch":
		if languageOptions == nil {
			return nil, fmt.Errorf("languageOptions is nil for language: %s", language)
		}
		sdk, err = generateNodeJSFetchSDK(spec, apiKey, languageOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to generate NodeJS Fetch SDK: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	sdkFile, err := zipWriter.Create("sdk.js")
	if err != nil {
		return nil, err
	}
	_, err = sdkFile.Write(sdk)
	if err != nil {
		return nil, err
	}
	zipWriter.Close()
	return buf.Bytes(), nil
}

func generateNodeJSFetchSDK(openAPISpec map[string]interface{}, apiKey string, languageOptions map[string]interface{}) ([]byte, error) {
    var sb strings.Builder

    // SDK class header
    baseUrl, ok := languageOptions["baseUrl"].(string)
    if !ok || baseUrl == "" {
        return nil, fmt.Errorf("invalid or missing baseUrl in languageOptions")
    }
    className := "DeclarativeClient"
    sb.WriteString("const fetch = require('node-fetch');\n\n")
    sb.WriteString(fmt.Sprintf("class %s {\n", className))
    sb.WriteString("\tconstructor(apiKey) {\n")
    sb.WriteString("\t\tthis.apiKey = apiKey;\n")
    sb.WriteString(fmt.Sprintf("\t\tthis.baseUrl = '%s';\n", baseUrl))
    sb.WriteString("\t}\n\n")

    // Iterate over paths in the OpenAPI spec
    paths, ok := openAPISpec["paths"].(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("invalid or missing paths in OpenAPI spec")
    }
    
    for path, methods := range paths {
        resourceName := extractResourceName(path)

        for method, _ := range methods.(map[string]interface{}) {
            funcName := generateFunctionName(resourceName, method, path)
            sb.WriteString(generateFunctionTemplate(funcName, method, path))
            sb.WriteString("\n\n")
        }
    }

    // SDK class footer
    sb.WriteString("}\n\n")
    sb.WriteString(fmt.Sprintf("module.exports = %s;\n", className))

    return []byte(sb.String()), nil
}

// extractResourceName extracts the resource name from the path, e.g., "/users/{userId}" -> "users".
func extractResourceName(path string) string {
    parts := strings.Split(path, "/")
    for _, part := range parts {
        if len(part) > 0 && !strings.Contains(part, "{") {
            return part
        }
    }
    return "resource"
}

func generateFunctionName(resourceName, method, path string) string {
    switch strings.ToUpper(method) {
    case "GET":
        if strings.Contains(path, "{") {
            return fmt.Sprintf("get%sById", capitalize(resourceName))
        }
        return fmt.Sprintf("list%s", capitalize(resourceName))
    case "POST":
        return fmt.Sprintf("create%s", capitalize(resourceName))
    case "PUT":
        return fmt.Sprintf("update%s", capitalize(resourceName))
    case "DELETE":
        return fmt.Sprintf("delete%s", capitalize(resourceName))
    default:
        return fmt.Sprintf("handle%s%s", capitalize(method), capitalize(resourceName))
    }
}

func generateFunctionTemplate(funcName, method, path string) string {
    var sb strings.Builder

    sb.WriteString(fmt.Sprintf("\tasync %s(data) {\n", funcName))
    sb.WriteString(fmt.Sprintf("\t\tconst url = `${this.baseUrl}%s`;\n", convertPathToTemplate(path)))
    sb.WriteString(fmt.Sprintf("\t\tconst response = await fetch(url, {\n"))
    sb.WriteString(fmt.Sprintf("\t\t\tmethod: '%s',\n", strings.ToUpper(method)))
    sb.WriteString("\t\t\theaders: {\n")
    sb.WriteString("\t\t\t\t'X-API-KEY': `${this.apiKey}`,\n")

    if method == "POST" || method == "PUT" {
        sb.WriteString("\t\t\t\t'Content-Type': 'application/json',\n")
    }

    sb.WriteString("\t\t\t},\n")

    if method == "POST" || method == "PUT" {
        sb.WriteString("\t\t\tbody: JSON.stringify(data),\n")
    }

    sb.WriteString("\t\t});\n")
    sb.WriteString("\t\treturn response.json();\n")
    sb.WriteString("\t}")

    return sb.String()
}

// convertPathToTemplate converts an OpenAPI path to a template string, e.g., "/users/{userId}" -> "/users/${userId}".
func convertPathToTemplate(path string) string {
    return strings.Replace(path, "{", "${", -1)
}

// capitalize capitalizes the first letter of a string.
func capitalize(s string) string {
    if len(s) == 0 {
        return s
    }
    return strings.ToUpper(s[:1]) + s[1:]
}