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
		if err := c.Render(
			http.StatusBadRequest,
			"error-message",
			map[string]any{"message": HTTPErr.Message},
		); err != nil {
			log.Print(err)
		}
	case http.StatusInternalServerError:
		if err := c.String(
			http.StatusInternalServerError,
			HTTPErr.Error(),
		); err != nil {
			log.Print(err)
		}
	}
}
