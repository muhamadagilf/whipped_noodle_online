// Package handler
package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
	"github.com/muhamadagilf/whipped_noodle_online/internal/server"
	"github.com/muhamadagilf/whipped_noodle_online/middlewares"
	"github.com/muhamadagilf/whipped_noodle_online/util"
	"github.com/redis/go-redis/v9"
)

type UserData struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func (h *Handler) Homepage(c echo.Context) error {
	rdb := h.Server.RDB
	query := h.Server.Queries
	csrf, ok := c.Get("csrf").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoCSRFError)
	}

	var menu []database.Menu
	cached, err := rdb.Get(c.Request().Context(), "menu:all").Result()
	if err != nil {
		if err == redis.Nil {
			menu, err := query.GetAllMenu(c.Request().Context())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			menuJSON, err := json.Marshal(menu)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			rdb.Set(c.Request().Context(), "menu:all", menuJSON, 24*time.Hour)
		}
	}
	if err := json.Unmarshal([]byte(cached), &menu); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	cart, ok := c.Get("cart").(*util.Cart)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoCartError)
	}
	userCred, ok := c.Get("userCred").(middlewares.UserCred)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoUserIDError)
	}

	return c.Render(http.StatusOK, "home", Data{
		"csrf_token": csrf,
		"menu":       menu,
		"cart":       cart.Menus,
		"total":      cart.Total,
		"cred":       userCred,
		"email":      userCred.Email,
	})
}

func (h *Handler) Loginpage(c echo.Context) error {
	csrf, ok := c.Get("csrf").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoCSRFError)
	}
	return c.Render(http.StatusOK, "login", Data{
		"csrf_token": csrf,
	})
}

func (h *Handler) Login(c echo.Context) error {
	stateID := fmt.Sprintf(
		"oauthstate%v%v",
		uuid.New().String(),
		time.Now().Local().UnixMilli()*time.Now().Local().UnixMicro(),
	)
	stateCookie := http.Cookie{
		Name:     "oauth_state",
		HttpOnly: true,
		Expires:  time.Now().Add(120 * time.Second),
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Value:    stateID,
	}
	c.SetCookie(&stateCookie)
	oauthRedirectURL := server.GoogleOAuthConfig.AuthCodeURL(stateID)
	return c.Redirect(http.StatusFound, oauthRedirectURL)
}

func (h *Handler) OauthCallback(c echo.Context) error {
	query := h.Server.Queries
	oauthCode := c.QueryParam("code")
	oauthStateID := c.QueryParam("state")
	if oauthCode == "" || oauthStateID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "oauth_code or oauth_state_id not found")
	}

	stateCookie, err := c.Request().Cookie("oauth_state")
	if err != nil || oauthStateID != stateCookie.Value {
		return echo.NewHTTPError(http.StatusBadRequest, "state_id parameter mismatched")
	}

	stateCookie.MaxAge = -1
	c.SetCookie(stateCookie)

	token, err := server.GoogleOAuthConfig.Exchange(c.Request().Context(), oauthCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	client := server.GoogleOAuthConfig.Client(c.Request().Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	userData := &UserData{}
	if err := json.Unmarshal(data, userData); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	session, ok := c.Get("session").(*sessions.Session)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	sessionID, ok := session.Values["session_id"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionIDError)
	}

	if err = query.Transaction(c.Request().Context(), h.Server.DB, func(qtx *database.Queries) error {
		if err := qtx.CreateUser(c.Request().Context(), database.CreateUserParam{
			ID:            userData.ID,
			Email:         userData.Email,
			VerifiedEmail: userData.VerifiedEmail,
			Name:          userData.Name,
			Picture:       userData.Picture,
		}); err != nil {
			return err
		}
		if err := qtx.UpdateSessionAuthentication(c.Request().Context(), database.UpdateSessionParams{
			SessionID: sessionID,
			UserID:    sql.NullString{String: userData.ID, Valid: true},
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	session.Values["user_id"] = sql.NullString{String: userData.ID, Valid: true}
	session.Values["user_email"] = userData.Email
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	returnToURL, ok := session.Values["return_to"].(string)
	if !ok && returnToURL == "" {
		returnToURL = "/home"
	}
	return c.Redirect(http.StatusFound, returnToURL)
}

func (h *Handler) Logout(c echo.Context) error {
	time.Sleep(300 * time.Millisecond)
	query := h.Server.Queries
	session, ok := c.Get("session").(*sessions.Session)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	sid := session.Values["session_id"].(string)
	session.Options.MaxAge = -1
	err := query.Transaction(c.Request().Context(), h.Server.DB, func(qtx *database.Queries) error {
		if err := qtx.DeleteSessionBySessionID(c.Request().Context(), sid); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	c.Response().Header().Set("HX-Redirect", "/home")
	return c.NoContent(http.StatusOK)
}
