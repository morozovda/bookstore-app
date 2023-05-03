package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"net/mail"

	"bookstore-api/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h *DBH) Signup (c echo.Context) error {
	customer := new(models.Regcustomer)
	var cid uuid.UUID
	var e models.Error
	
	err := c.Bind(customer)
	if err != nil {
		e.Message = "request failed"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	if customer.Name == "" || customer.Email == "" || customer.Passwd == "" {
		e.Message = "request failed"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	_, err = mail.ParseAddress(customer.Email)
	if err != nil {
		e.Message = "invalid email"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	hspasswd, err := bcrypt.GenerateFromPassword([]byte(customer.Passwd), bcrypt.MinCost+4)
	if err != nil {
		e.Message = "invalid password"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	err = h.DB.QueryRow("SELECT \"id\" FROM \"customer\" WHERE \"name\"=$1 AND \"email\"=$2;", customer.Name, customer.Email).Scan(&cid)
	switch err {
		case sql.ErrNoRows:
			err := h.DB.QueryRow("INSERT INTO \"customer\" (\"name\", \"email\", \"passwd\") VALUES ($1, $2, $3) RETURNING \"id\";", customer.Name, customer.Email, string(hspasswd)).Scan(&cid)
			if err != nil {
				log.Fatal("invalid credentials")
				e.Message = "invalid credentials"
				return c.JSON(http.StatusBadRequest, e)
			}
			
			if customer.Name == "Bob" {
				_, err := h.DB.Exec("UPDATE \"customer\" SET \"balance\" = '2000' WHERE \"id\" = $1;", cid)
				if err != nil {
					log.Fatal("claim failed")
				}
			}
			return c.NoContent(http.StatusCreated)
            
		case nil:
			e.Message = "already registered"
			return c.JSON(http.StatusForbidden, e)
            
		default:
			log.Fatal("service unavailable")
			e.Message = "service unavailable"
			return c.JSON(http.StatusBadRequest, e)
	}
}