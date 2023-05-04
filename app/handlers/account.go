package handlers

import (
	"log"
	"net/http"

	"bookstore-api/models"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *DBH) Account (c echo.Context) error {
	bs := []models.Book{}
	var customerAcc models.Account
	var customerId uuid.UUID
	var balance int
	var cExists bool
	var e models.Error
	
	token := c.Get("customer").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	customerId, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		e.Message = "request failed"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	_ = h.DB.QueryRow("SELECT EXISTS (SELECT \"id\" FROM \"customer\" WHERE \"id\" = $1);", customerId).Scan(&cExists)
	if cExists == false {
		log.Printf("customer %s doesnt exist", customerId)
		e.Message = "invalid credentials"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	rows, err := h.DB.Query("SELECT \"book\".\"id\", \"book\".\"title\", \"book\".\"author\", \"book\".\"price\", \"order_amount\" FROM \"deal\" INNER JOIN \"book\" ON \"deal\".\"book_id\" = \"book\".\"id\" WHERE \"customer_id\"=$1;", customerId)
	switch err {
		case nil:
			defer rows.Close()
			for rows.Next() {
				var book models.Book
				err = rows.Scan(&book.Id, &book.Title, &book.Author, &book.Price, &book.Amount)
				if err != nil {
					log.Fatal("service unavailable")
					e.Message = "service unavailable"
					return c.JSON(http.StatusServiceUnavailable, e)
				}
				book.Price = book.Price/100
				bs = append(bs, book)
			}

			err = rows.Err()
			if err != nil {
				log.Fatal("service unavailable")
				e.Message = "service unavailable"
				return c.JSON(http.StatusServiceUnavailable, e)        
			}
	
			_ = h.DB.QueryRow("SELECT \"balance\" FROM \"customer\" WHERE \"id\" = $1;", customerId).Scan(&balance)
			
			customerAcc.Books = bs
			customerAcc.Balance = balance/100
			
			return c.JSON(http.StatusOK, customerAcc)
            
		default:
			log.Fatal("service unavailable")
			e.Message = "service unavailable"
			return c.JSON(http.StatusBadRequest, e)
	}
}