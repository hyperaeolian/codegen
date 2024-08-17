package handlers

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/gorilla/mux"
	"codegen/generators"
	"log"
)

type CodegenRequest struct {
	ApiKey string `json:"api_key"`
	OpenAPISpec map[string]interface{} `json:"spec"`
	LanguageOptions map[string]interface{} `json:"language_options"`
}

func CodegenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	language := vars["language"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var codegenRequest CodegenRequest
	if err = json.Unmarshal(body, &codegenRequest); err != nil {
		http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
		return
	}

	if codegenRequest.OpenAPISpec == nil {
		http.Error(w, "OpenAPI spec is required", http.StatusBadRequest)
		return
	}
	if codegenRequest.ApiKey == "" {
		http.Error(w, "API key is required", http.StatusBadRequest)
		return
	}
	if codegenRequest.LanguageOptions == nil {
		http.Error(w, "Language options are required", http.StatusBadRequest)
		return
	}

	sdk, err := generators.GenerateSDK(codegenRequest.OpenAPISpec, codegenRequest.ApiKey, language, codegenRequest.LanguageOptions)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error generating code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=sdk.zip")
	w.Header().Set("Content-Type", "application/zip")
	w.WriteHeader(http.StatusOK)
	w.Write(sdk)
}
