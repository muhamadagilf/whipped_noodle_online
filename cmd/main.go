package main

import (
	"html/template"
	"io"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhamadagilf/whipped_noodle_online/handler"
	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
	"github.com/muhamadagilf/whipped_noodle_online/internal/server"
	"github.com/muhamadagilf/whipped_noodle_online/middlewares"

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

	e := echo.New()
	e.Renderer = newTemplates()
	e.Static("/static", "static")

	s, err := server.NewServer()
	if err != nil {
		log.Fatal("New s init Failed: ", err)
	}
	defer s.DB.Close()

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

	public := e.Group("")
	public.Use(middleware.RequestLogger())
	public.Use(mdl.Session)
	public.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		TokenLength:    32,
		TokenLookup:    "form:_csrf",
		ContextKey:     "csrf",
		CookieName:     "_csrf",
		CookieMaxAge:   86400,
		CookieHTTPOnly: true,
	}))
	public.Use(mdl.Authentication)
	public.Use(mdl.VerifyRedirectURL)

	public.GET("/home", h.Homepage)
	public.GET("/login", h.Loginpage)
	public.POST("/auth/login", h.Login)
	public.GET("/auth/oauth/callback", h.OauthCallback)
	public.GET("/checkout", h.Checkout)

	e.Logger.Fatal(e.Start(":8000"))
}
