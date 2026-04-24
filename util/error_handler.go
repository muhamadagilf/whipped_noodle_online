package util

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandling(err error, c echo.Context) {
	HTTPErr, ok := err.(*echo.HTTPError)
	if !ok {
		if err := c.String(
			http.StatusInternalServerError,
			HTTPErrorAssertErr,
		); err != nil {
			log.Print(err)
		}
	}

	switch HTTPErr.Code {
	case http.StatusBadRequest:
		c.Render(HTTPErr.Code, "error-message", map[string]any{"message": HTTPErr.Message.(string)})
	case http.StatusNotFound, http.StatusInternalServerError:
		c.String(HTTPErr.Code, HTTPErr.Message.(string))
	}

}
