package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"bookstore-api/models"

	"github.com/google/uuid"
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

func (h *DBH) Marketbook (c echo.Context) error {
	book := models.Book{}
	var bookId uuid.UUID
	var e models.Error
	
	bookId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		e.Message = "request failed"
		return c.JSON(http.StatusBadRequest, e)
	}

	err = h.DB.QueryRow("SELECT \"id\", \"title\", \"author\", \"price\", \"amount\" FROM \"book\" WHERE \"id\" = $1;", bookId).Scan(&book.Id, &book.Title, &book.Author, &book.Price, &book.Amount)
	switch err {
		case sql.ErrNoRows:
			log.Printf("book %s doesnt exist", bookId.String())
			e.Message = "book doesnt exist"
			return c.JSON(http.StatusNotFound, e)
			
		case nil:
			return c.JSON(http.StatusOK, book)
		default: 
			log.Fatal("service unavailable")
			e.Message = "service unavailable"
			return c.JSON(http.StatusServiceUnavailable, e)
	}
}