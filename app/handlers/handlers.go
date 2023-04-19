package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"net/mail"

	"bookstore-api/models"

	"golang.org/x/crypto/bcrypt"
	"github.com/labstack/echo/v4"
)

type (
	Handler struct {
		DB *sql.DB
	}
)

func (h *Handler) Market (c echo.Context) error {
	books := []models.Book{}
	var market models.Market
	var e models.Error
	rows, err := h.DB.Query("SELECT \"id\", \"title\", \"author\", \"price\", \"amount\" FROM \"book\";")
	switch err {
		case sql.ErrNoRows:
			return c.NoContent(http.StatusNoContent)
            
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

func (h *Handler) Signup (c echo.Context) error {
	customer := new(models.Customer)
	var cid models.Ids
	var e models.Error
	
	err := c.Bind(customer)
	if err != nil {
		log.Fatal("request failed")
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
	
	err = h.DB.QueryRow("SELECT \"id\" FROM \"customer\" WHERE \"name\"=$1 AND \"email\"=$2;", customer.Name, customer.Email).Scan(&cid.Id)
	switch err {
		case sql.ErrNoRows:
			err := h.DB.QueryRow("INSERT INTO \"customer\" (\"name\", \"email\", \"passwd\") VALUES ($1, $2, $3) RETURNING \"id\";", customer.Name, customer.Email, string(hspasswd)).Scan(&cid.Id)
			if err != nil {
				log.Fatal("request failed")
				e.Message = "service unavailable"
				return c.JSON(http.StatusServiceUnavailable, e)
			}

			if cid.Id <= 3 {
				_, err := h.DB.Exec("UPDATE \"customer\" SET \"balance\" = '20' WHERE \"id\" = $1;", cid.Id)
				if err != nil {
					log.Fatal("claim failed")
				}
			}

			return c.JSON(http.StatusCreated, cid)
            
		case nil:
			e.Message = "already registered"
			return c.JSON(http.StatusForbidden, e)
            
		default:
			log.Fatal("service unavailable")
			e.Message = "service unavailable"
			return c.JSON(http.StatusBadRequest, e)
	}
}

func (h *Handler) Deal (c echo.Context) error {
	d := new(models.Deal)
	var e models.Error
	var bExists bool
	var cExists bool
	var bookAmount int
	var bookPrice int
	var customerBalance int
	var balanceAfer int
	
	err := c.Bind(d)
	if err != nil {
		log.Fatal("request failed")
		e.Message = "request failed"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	_ = h.DB.QueryRow("SELECT EXISTS (SELECT \"id\" FROM \"customer\" WHERE \"id\" = $1);", d.Customer_id).Scan(&cExists)
	if cExists == false {
		log.Printf("customer %d doesnt exist", d.Customer_id)
		e.Message = "customer doesnt exist"
		return c.JSON(http.StatusNotFound, e)
	}
	
	_ = h.DB.QueryRow("SELECT EXISTS (SELECT \"id\" FROM \"book\" WHERE \"id\" = $1);", d.Book_id).Scan(&bExists)
	if bExists == false {
		log.Printf("book %d doesnt exist", d.Book_id)
		e.Message = "book doesnt exist"
		return c.JSON(http.StatusNotFound, e)
	}
	
	_ = h.DB.QueryRow("SELECT \"amount\" FROM \"book\" WHERE \"id\" = $1;", d.Book_id).Scan(&bookAmount)
	if bookAmount < d.Order_amount {
		log.Printf("not enough %d books in market", d.Book_id)
		e.Message = "not enough books in market"
		return c.JSON(http.StatusConflict, e)
	}
	
	_ = h.DB.QueryRow("SELECT \"balance\" FROM \"customer\" WHERE \"id\" = $1;", d.Customer_id).Scan(&customerBalance)
	_ = h.DB.QueryRow("SELECT \"price\" FROM \"book\" WHERE \"id\" = $1;", d.Book_id).Scan(&bookPrice)
	if customerBalance < bookPrice*d.Order_amount {
		log.Printf("not enough funds in balance %d", d.Customer_id)
		e.Message = "not enough funds in balance"
		return c.JSON(http.StatusConflict, e)
	}
	
	balanceAfer = customerBalance - bookPrice
	_, err = h.DB.Exec("UPDATE \"customer\" SET \"balance\" = $1 WHERE \"id\" = $2;", balanceAfer, d.Customer_id)
	if err != nil {
		log.Fatal("transaction failed")
		e.Message = "transaction unavailable"
		return c.JSON(http.StatusServiceUnavailable, e)
	}
	
	_, err = h.DB.Exec("INSERT INTO \"deal\" (\"book_id\", \"order_amount\", \"customer_id\") VALUES ($1, $2, $3);", d.Book_id, d.Order_amount, d.Customer_id)
	if err != nil {
		log.Fatal("request failed")
		e.Message = "service unavailable"
		return c.JSON(http.StatusServiceUnavailable, e)
	}
	
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Account (c echo.Context) error {
	i := new(models.Ids)
	bs := []models.Book{}
	var cd models.Customerdeal
	var balance int
	var cExists bool
	var e models.Error
	
	err := c.Bind(i)
	if err != nil {
		log.Fatal("request failed")
		e.Message = "request failed"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	_ = h.DB.QueryRow("SELECT EXISTS (SELECT \"id\" FROM \"customer\" WHERE \"id\" = $1);", i.Id).Scan(&cExists)
	if cExists == false {
		log.Printf("customer %d doesnt exist", i.Id)
		e.Message = "customer doesnt exist"
		return c.JSON(http.StatusNotFound, e)
	}
	
	rows, err := h.DB.Query("SELECT \"book\".\"id\", \"book\".\"title\", \"book\".\"author\", \"book\".\"price\", \"order_amount\" FROM \"deal\" INNER JOIN \"book\" ON \"deal\".\"book_id\" = \"book\".\"id\" WHERE \"customer_id\"=$1;", i.Id)
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
				bs = append(bs, book)
			}

			err = rows.Err()
			if err != nil {
				log.Fatal("service unavailable")
				e.Message = "service unavailable"
				return c.JSON(http.StatusServiceUnavailable, e)        
			}
	
			_ = h.DB.QueryRow("SELECT \"balance\" FROM \"customer\" WHERE \"id\" = $1;", i.Id).Scan(&balance)
			
			cd.Books = bs
			cd.Balance = balance
			
			return c.JSON(http.StatusOK, cd)
            
		default:
			log.Fatal("service unavailable")
			e.Message = "service unavailable"
			return c.JSON(http.StatusBadRequest, e)
	}
}