package main

import (
	"database/sql"
	"encoding/gob"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhamadagilf/whipped_noodle_online/handler"
	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
	"github.com/muhamadagilf/whipped_noodle_online/internal/server"
	"github.com/muhamadagilf/whipped_noodle_online/middlewares"
	"github.com/muhamadagilf/whipped_noodle_online/util"

	_ "modernc.org/sqlite"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("view/*.html")),
	}
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("use container environment or .env not found")
	}

	gob.Register(util.MenuOrder{})
	gob.Register(util.Cart{})
	gob.Register(sql.NullString{})

	e := echo.New()
	e.Renderer = newTemplates()
	e.HTTPErrorHandler = util.HTTPErrorHandling
	e.Static("/static", "static")

	s, err := server.NewServer()
	if err != nil {
		log.Fatal("New s init Failed: ", err)
	}
	defer s.DB.Close()
	defer s.RDB.Close()

	if err := server.InitOAuth(); err != nil {
		log.Fatal("OAuth Init Failed: ", err)
	}

	if err := database.Migrate(s.DB); err != nil {
		log.Fatal("Migration Failed: ", err)
	}

	if err := s.DB.Ping(); err != nil {
		log.Fatal("DB Ping Failed: ", err)
	}

	if err := s.Queries.InsertMainMenuToDB(); err != nil {
		log.Fatal("INSERT ROW DB Failed: ", err)
	}

	mdl := middlewares.NewMiddlewares(s)
	h := handler.NewHandler(s)

	r := e.Group("")
	r.Use(middleware.RequestLogger())
	r.Use(mdl.Session)
	r.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		TokenLength:    32,
		TokenLookup:    "form:_csrf,header:HX-CSRF-TOKEN",
		ContextKey:     "csrf",
		CookieName:     "_csrf",
		CookieMaxAge:   86400,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteLaxMode,
	}))
	r.Use(mdl.Authentication)
	r.Use(mdl.VerifyRedirectURL)

	// ROUTES
	r.GET("/home", h.Homepage)
	r.GET("/login", h.Loginpage)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout", h.Logout)
	r.GET("/auth/oauth/callback", h.OauthCallback)

	r.POST("/cart/add", h.AddToCartSession)
	r.DELETE("/cart/delete/:menuID", h.DeleteFromCartSession)

	r.GET("/checkout", h.Checkoutpage)

	// BG_WORKER
	go util.DBSessionCleanUp(s.Queries)

	log.Println("SERVER RUNNING ON :8000")
	e.Logger.Fatal(e.Start(":8000"))
}
