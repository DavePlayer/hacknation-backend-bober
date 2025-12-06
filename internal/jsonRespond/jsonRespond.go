package jsonRespond

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusFail    Status = "fail"
	StatusError   Status = "error"
)

type Response struct {
	Status  Status      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ogólny helper
func Write(c *gin.Context, httpStatus int, status Status, data any) {
	c.JSON(httpStatus, Response{
		Status: status,
		Data:   data,
	})
}

// jeśli chcesz surowe JSON bez wrappera Response
func SendJSON(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Status: StatusSuccess,
		Data:   data,
	})
}

func Fail(c *gin.Context, httpStatus int, message string, data any) {
	c.JSON(httpStatus, Response{
		Status:  StatusFail,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus int, message string, err error) {
	log.Printf("http error: %d %s: %v", httpStatus, message, err)

	c.JSON(httpStatus, Response{
		Status:  StatusError,
		Message: message,
		Data: map[string]string{
			"error": err.Error(),
		},
	})
}
