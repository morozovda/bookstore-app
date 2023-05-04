package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"bookstore-api/models"

	"github.com/labstack/echo/v4"
)

func (h *DBH) Market (c echo.Context) error {
	books := []models.Book{}
	var market models.Market
	var e models.Error
	rows, err := h.DB.Query("SELECT \"id\", \"title\", \"author\", \"price\", \"amount\" FROM \"book\";")
	switch err {
		case sql.ErrNoRows:
			return c.JSON(http.StatusOK, books)
            
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
				if book.Amount == 0 {
					continue
				}
				book.Price = book.Price/100
				books = append(books, book)
			}
    
			err = rows.Err()
			if err != nil {
				log.Fatal("service unavailable")
				e.Message = "service unavailable"
				return c.JSON(http.StatusServiceUnavailable, e)        
			}
	
			market.Books = books
			return c.JSON(http.StatusOK, market)
            
		default:
			log.Fatal("service unavailable")
			e.Message = "service unavailable"
			return c.JSON(http.StatusBadRequest, e)
	}
}