package main

import (
	"os"

	"bookstore-api/db"
	"bookstore-api/handlers"

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

	e.POST("/signup", h.Signup)
	e.POST("/account", h.Account)
	e.GET("/market", h.Market)
	e.POST("/market/deal", h.Deal)

	e.Logger.Fatal(e.Start((os.Getenv("WEBURL") + ":" + os.Getenv("WEBPORT"))))
}