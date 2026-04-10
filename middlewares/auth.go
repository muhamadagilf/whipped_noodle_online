package middlewares

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/internal/server"
	"github.com/muhamadagilf/whipped_noodle_online/util"
)

type UserCred struct {
	UserID sql.NullString
	Email  string
}

func (m *Middlewares) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		query := m.Server.Queries
		ctx := c.Request().Context()
		session, ok := c.Get("session").(*sessions.Session)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
		}

		userID, ok := session.Values["user_id"].(sql.NullString)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, util.NoUserIDError)
		}

		userEmail, ok := session.Values["user_email"].(string)
		if !ok {
			userEmail = ""
		}

		cart, ok := session.Values["cart"].(util.Cart)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, util.NoCartError)
		}

		c.Set("cart", &cart)
		c.Set("userCred", UserCred{
			UserID: userID,
			Email:  userEmail,
		})

		log.Println(c.Path())
		var requestURL string
		if c.Path() == "/cart/delete/:menuID" {
			requestURL = c.Path()
		} else {
			requestURL = c.Request().URL.Path
		}

		if slices.Contains(server.PublicURL, requestURL) {
			if userID.Valid && requestURL == "/login" {
				return c.Redirect(http.StatusFound, "/home")
			}
			return next(c)
		}

		// protected url
		sessionID, ok := session.Values["session_id"].(string)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionIDError)
		}

		sessionData, err := query.GetSessionBySessionID(ctx, sessionID)
		if err == sql.ErrNoRows {
			session.Options.MaxAge = -1
			if err = session.Save(c.Request(), c.Response()); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			// redirectURL := "/login?redirect=" + url.QueryEscape(c.Request().URL.Path)
			return c.Redirect(http.StatusFound, "/home")
		}

		if sessionData.IsAuthenticated == 0 && !sessionData.UserID.Valid {
			redirectURL := "/login?redirect=" + url.QueryEscape(c.Request().URL.Path)
			return c.Redirect(http.StatusFound, redirectURL)
		}

		c.Set("userCred", UserCred{
			UserID: sessionData.UserID,
			Email:  userEmail,
		})
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
		if redirectURL != "" {
			session, ok := c.Get("session").(*sessions.Session)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
			}
			session.Values["return_to"] = redirectURL
			if err := session.Save(c.Request(), c.Response()); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}
		return next(c)
	}
}
