// Package middlewares
package middlewares

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
	"github.com/muhamadagilf/whipped_noodle_online/util"
)

func (m *Middlewares) Session(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		query := m.Server.Queries
		if c.Request().URL.Path == "/" {
			return c.Redirect(http.StatusFound, "/home")
		}
		session, err := m.Server.SessionStore.Get(c.Request(), m.Server.SessionName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if session.IsNew {
			log.Println("CREATE NEW ONE")
			sessionID := fmt.Sprintf("sessid_%v_%v", uuid.New().String(), time.Now().Local().UnixMilli())
			cartID := fmt.Sprintf("orderid_%v_%v", uuid.New().String(), time.Now().Local().Unix())
			session.Values["session_id"] = sessionID
			session.Values["user_id"] = sql.NullString{Valid: false}
			session.Values["cart"] = util.Cart{
				ID:    cartID,
				Menus: make(map[string]util.MenuOrder),
				Total: 0,
			}

			if err := query.Transaction(c.Request().Context(), m.Server.DB, func(qtx *database.Queries) error {
				if err := qtx.CreateSession(c.Request().Context(), database.CreateSessionParams{
					SessionID: sessionID,
					ExpiredAt: time.Now().Local().Add(24 * time.Hour).UnixMilli(),
					UserID:    sql.NullString{Valid: false},
				}); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			log.Println("[SESSION_CREATED]# ", sessionID)
		}

		log.Println("[SESSION_DEBUG]# ", session.Values["session_id"].(string))
		c.Set("session", session)
		if err := session.Save(c.Request(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return next(c)
	}
}
