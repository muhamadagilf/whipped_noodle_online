// Package middlewares
package middlewares

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
)

func (m *Middlewares) Session(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().URL.Path == "/" {
			return c.Redirect(http.StatusFound, "/home")
		}

		if c.Request().URL.Path == "/favicon.ico" {
			return next(c)
		}
		query := m.Server.Queries
		session, err := m.Server.SessionStore.Get(c.Request(), m.Server.SessionName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		sid, ok := session.Values["session_id"].(string)
		if sid == "" && !ok {
			sessionID := fmt.Sprintf("sessid_%v_kuncisesi4422", time.Now().Local().UnixMilli()*time.Now().Local().UnixMicro()*123)
			session.Values["session_id"] = sessionID
			session.Values["user_id"] = ""
			session.Values["cart"] = ""
			if err := query.Transaction(c.Request().Context(), m.Server.DB, func(qtx *database.Queries) error {
				if err := qtx.CreateSession(c.Request().Context(), database.CreateSessionParams{
					SessionID: sessionID,
					ExpiredAt: time.Now().Local().Add(12 * time.Hour).UnixMilli(),
					UserID:    sql.NullString{Valid: false},
				}); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error()+"; session")
			}

		}
		c.Set("session", session)
		if err := session.Save(c.Request(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return next(c)
	}
}
