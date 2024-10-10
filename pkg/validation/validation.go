package validation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	res "github.com/luiidev/go/pkg/response"
)

type Errors map[string][]string

type Validated struct {
	message string
	errors  Errors
}

func (v *Validated) Message() string {
	return v.message
}

func (v *Validated) Errors() Errors {
	return v.errors
}

func (v *Validated) Fails() bool {
	return len(v.errors) > 0
}

func (v *Validated) Passes() bool {
	return !v.Fails()
}

func (v *Validated) Error(field string) string {
	if v.Fails() {
		if errors, ok := v.errors[field]; ok {
			return errors[0]
		}
	}

	return ""
}

func (v *Validated) Response(w http.ResponseWriter) {
	res.JSON(w, res.H{"message": v.message, "errors": v.errors}, http.StatusUnprocessableEntity)
}

var validate = validator.New()

// Mapa de mensajes de error personalizados
var errorMessages = map[string]string{
	"required": "The %s is required.",
	"min":      "The %s must be at least %s characters.",
	"max":      "The %s may not be greater than %s characters.",
	"email":    "The %s must be a valid email address.",
	"gte":      "The %s must be greater than or equal to %s.",
	"lte":      "The %s must be less than or equal to %s.",
	"len":      "The %s must be exactly %s characters.",
	"eqfield":  "The %s must be equal to %s.",
	"nefield":  "The %s must not be equal to %s.",
	"url":      "The %s must be a valid URL.",
	"numeric":  "The %s must be a numeric value.",
}

// Función para obtener el nombre del campo según el tag `json`
func getJSONFieldName(entity interface{}, field string) string {
	t := reflect.TypeOf(entity).Elem()
	fieldStruct, found := t.FieldByName(field)
	if !found {
		return field // Si no se encuentra, devolver el nombre del campo original
	}

	jsonTag := fieldStruct.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return field // Si no tiene tag `json`, devolver el nombre del campo
	}

	return jsonTag
}

// Decodifica y valida el JSON
func Make[T any](body io.ReadCloser, model *T) Validated {
	// Decodificar el JSON
	err := json.NewDecoder(body).Decode(model)
	if err != nil {
		return Validated{message: "Invalid JSON", errors: Errors{"json": []string{"Invalid JSON"}}}
	}

	// Validar los campos
	err = validate.Struct(model)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(Errors)

			// Recorrer los errores de validación y adaptarlos al formato de Laravel
			for _, err := range validationErrors {
				field := getJSONFieldName(model, err.Field())
				tag := err.Tag() // Obtener el campo que falló

				// Inicializar el mensaje base
				message := errorMessages[tag]

				formatFieldName := strings.Join(strings.Split(field, "_"), " ")

				// Usar fmt.Sprintf con el número correcto de argumentos
				if tag == "min" || tag == "max" || tag == "gte" || tag == "lte" || tag == "len" || tag == "eqfield" || tag == "nefield" {
					message = fmt.Sprintf(message, formatFieldName, err.Param())
				} else {
					message = fmt.Sprintf(message, formatFieldName)
				}

				// Agregar el mensaje a la lista
				errors[field] = append(errors[field], message)
			}

			// Retornar el formato de error
			return Validated{message: "Validation Failed", errors: errors}
		}
	}

	return Validated{}
}
