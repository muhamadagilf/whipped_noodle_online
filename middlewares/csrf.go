package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (m *Middlewares) CSRF(csrfHandler func(http.Handler) http.Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			handler := csrfHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)
				next(c)
			}))

			handler.ServeHTTP(c.Response(), c.Request())
			return nil
		}
	}
}

func (m *Middlewares) CrossProtectionHTTP(COPHandler *http.CrossOriginProtection) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var handlerError error
			h := COPHandler.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)
				handlerError = next(c)
			}))

			h.ServeHTTP(c.Response().Writer, c.Request())
			return handlerError
		}
	}
}
