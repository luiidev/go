package controllers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func JsonResponse(w http.ResponseWriter, response Response, statusCode ...int) {
	// Si no se proporciona `statusCode`, usar 200 como valor por defecto
	status := http.StatusOK
	if len(statusCode) > 0 {
		status = statusCode[0]
	}

	// Configurar el encabezado de la respuesta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Codificar la respuesta a JSON
	json.NewEncoder(w).Encode(response)
}
