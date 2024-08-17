package handlers

import (
	"net/http"
	"encoding/json"
)

func Root(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"name": "Codegen", "status": "OK", "version": "1.0.0"}
	json.NewEncoder(w).Encode(data)
}