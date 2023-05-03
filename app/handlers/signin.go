package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"net/mail"
	"time"

	"bookstore-api/models"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Signin (c echo.Context) error {
	lc := new(models.Logincustomer)
	var jwts models.Jwts
	var e models.Error
	var dbpasswd string
	var dbid uuid.UUID
	
	err := c.Bind(lc)
	if err != nil {
		e.Message = "request failed"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	if lc.Email == "" || lc.Passwd == "" {
		e.Message = "request failed"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	_, err = mail.ParseAddress(lc.Email)
	if err != nil {
		e.Message = "invalid email"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	err = h.DB.QueryRow("SELECT \"id\", \"passwd\" FROM \"customer\" WHERE \"email\"=$1", lc.Email).Scan(&dbid, &dbpasswd)
	switch err {
		case sql.ErrNoRows:
			e.Message = "request failed"
			return c.JSON(http.StatusBadRequest, e)
			
		case nil:
			err = bcrypt.CompareHashAndPassword([]byte(dbpasswd), []byte(lc.Passwd))
			if err != nil {
				e.Message = "request failed"
				return c.JSON(http.StatusBadRequest, e)
			}
			
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["id"] = dbid
			claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
			jwts.Jwtstr, err = token.SignedString([]byte(Key))
			if err != nil {
				log.Fatal("service unavailable")
				e.Message = "service unavailable"
				return c.JSON(http.StatusServiceUnavailable, e)
			}
			return c.JSON(http.StatusOK, jwts)
            
		default:
			log.Fatal("service unavailable")
			e.Message = "service unavailable"
			return c.JSON(http.StatusServiceUnavailable, e)
	}
}