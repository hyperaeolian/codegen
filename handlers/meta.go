package handlers

import (
	"net/http"
	"encoding/json"
)

func Meta(w http.ResponseWriter, r *http.Request) {
	supported_codgen_variations := []string{"nodejs-fetch"}
	data := map[string]interface{}{"supported_options": supported_codgen_variations}
	json.NewEncoder(w).Encode(data)
}