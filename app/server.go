package main

import (
	"net/http"
	"os"

	"bookstore-api/db"
	"bookstore-api/handlers"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dbc := db.Start()
	defer dbc.Close()
	h := &handlers.Handler{DB: dbc}
	
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(handlers.Key),
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/signup" || c.Path() == "/signin" || c.Path() == "/market" {
				return true
			}
			return false
		},
		ContextKey: "customer",
	}))
	
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		//AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderContentLength},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))
	
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.NewString()
		},
	}))

	e.POST("/signup", h.Signup)
	e.POST("/signin", h.Signin)
	e.POST("/account", h.Account)
	e.GET("/market", h.Market)
	e.POST("/market/deal", h.Deal)

	e.Logger.Fatal(e.Start((os.Getenv("WEBURL") + ":" + os.Getenv("WEBPORT"))))
}