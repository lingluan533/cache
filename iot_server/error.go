package iot_server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func NewParameterError(param string) error {
	return fmt.Errorf("%s format is error", param)
}

// ErrorHandler handle the apiError
func ErrorHandler(err error, c echo.Context) {
	var (
		code   = http.StatusInternalServerError
		result = NewResult()
	)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		// TODO: use config here
		if he.Internal != nil && c.Logger().Level() == log.DEBUG {
			result.Err = fmt.Sprintf("%v, %v", err, he.Internal)
		} else {
			result.Err = he.Message
		}
	} else {
		// TODO: custom bind and db error
		code = http.StatusOK
		result.Err = err.Error()
	}
	// Send response
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD { // Issue #608
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, result)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
