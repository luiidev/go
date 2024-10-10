package controllers

import (
	"net/http"

	"github.com/luiidev/go/pkg/logger"
	res "github.com/luiidev/go/pkg/response"
	"gorm.io/gorm"
)

type ExampleController struct {
	l  logger.Logger
	db gorm.DB
}

func NewExampleController(l logger.Logger, db gorm.DB) *ExampleController {
	return &ExampleController{l: l, db: db}
}

func (c ExampleController) Helloworld(w http.ResponseWriter, r *http.Request) {
	res.JSON(w, res.H{"message": "Hello World!"})
}
