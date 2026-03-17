package middlewares

import (
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/internal/server"
)

func (m *Middlewares) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		query := m.Server.Queries
		ctx := c.Request().Context()
		session, ok := c.Get("session").(*sessions.Session)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "cannot find session in context")
		}

		userID, ok := session.Values["user_id"].(string)
		if !ok {
			return c.String(http.StatusInternalServerError, "cannot find user_id in session")
		}

		if slices.Contains(server.PublicURL, c.Request().URL.Path) {
			if c.Request().URL.Path == "/home" {
				return next(c)
			}
			if userID == "" {
				return next(c)
			}
			return c.Redirect(http.StatusFound, "/home")
		}

		// private url
		sessionID, ok := session.Values["session_id"].(string)
		if !ok {
			return c.String(http.StatusInternalServerError, "cannot find user_id in session")
		}

		sessionData, err := query.GetSessionBySessionID(ctx, sessionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error()+"; auth")
		}

		if userID == "" && !sessionData.UserID.Valid {
			redirectURL := "/login?redirect=" + url.QueryEscape(c.Request().URL.Path)
			return c.Redirect(http.StatusFound, redirectURL)
		}

		if time.Now().Local().UnixMilli() > sessionData.ExpiredAt {
			session.Values["is_authenticated"] = false
			session.Values["user_id"] = ""
			session.Values["cart"] = ""
			if err = session.Save(c.Request(), c.Response()); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			if err := query.DeleteSessionBySessionID(ctx, sessionID); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

		}

		c.Set("user_id", sessionData.UserID.String)
		return next(c)
	}
}

func (m *Middlewares) VerifyRedirectURL(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		redirectURL := c.QueryParam("redirect")
		if redirectURL == "" {
			return next(c)
		}
		if strings.Contains(redirectURL, "https://") || strings.Contains(redirectURL, "http://") {
			return echo.NewHTTPError(http.StatusNotFound, "cannot find URL; auth")
		}
		if !slices.Contains(server.ProtectedURL, redirectURL) {
			return echo.NewHTTPError(http.StatusNotFound, "cannot find URL; auth")
		}
		return next(c)
	}
}
