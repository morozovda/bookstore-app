package handlers

import (
	"log"
	"net/http"

	"bookstore-api/models"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Deal (c echo.Context) error {
	d := new(models.Deal)
	var e models.Error
	var cid uuid.UUID
	var bExists bool
	var cExists bool
	var bookAmount int
	var bookPrice int
	var customerBalance int
	var balanceAfer int
	
	token := c.Get("customer").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	cid, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		e.Message = "request failed"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	err = c.Bind(d)
	if err != nil {
		e.Message = "request failed"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	_ = h.DB.QueryRow("SELECT EXISTS (SELECT \"id\" FROM \"customer\" WHERE \"id\" = $1);", cid).Scan(&cExists)
	if cExists == false {
		log.Printf("customer %s doesnt exist", cid.String())
		e.Message = "invalid credentials"
		return c.JSON(http.StatusBadRequest, e)
	}
	
	_ = h.DB.QueryRow("SELECT EXISTS (SELECT \"id\" FROM \"book\" WHERE \"id\" = $1);", d.Book_id).Scan(&bExists)
	if bExists == false {
		log.Printf("book %s doesnt exist", d.Book_id.String())
		e.Message = "book doesnt exist"
		return c.JSON(http.StatusNotFound, e)
	}
	
	_ = h.DB.QueryRow("SELECT \"amount\" FROM \"book\" WHERE \"id\" = $1;", d.Book_id).Scan(&bookAmount)
	if bookAmount < d.Order_amount {
		log.Printf("not enough %d books in market", d.Book_id)
		e.Message = "not enough books in market"
		return c.JSON(http.StatusConflict, e)
	}
	
	_ = h.DB.QueryRow("SELECT \"balance\" FROM \"customer\" WHERE \"id\" = $1;", cid).Scan(&customerBalance)
	_ = h.DB.QueryRow("SELECT \"price\" FROM \"book\" WHERE \"id\" = $1;", d.Book_id).Scan(&bookPrice)
	if customerBalance < bookPrice*d.Order_amount {
		log.Printf("not enough funds in balance %d", cid)
		e.Message = "not enough funds in balance"
		return c.JSON(http.StatusConflict, e)
	}
	
	balanceAfer = customerBalance - bookPrice
	_, err = h.DB.Exec("UPDATE \"customer\" SET \"balance\" = $1 WHERE \"id\" = $2;", balanceAfer, cid)
	if err != nil {
		log.Fatal("transaction failed")
		e.Message = "transaction unavailable"
		return c.JSON(http.StatusServiceUnavailable, e)
	}
	
	_, err = h.DB.Exec("INSERT INTO \"deal\" (\"book_id\", \"order_amount\", \"customer_id\") VALUES ($1, $2, $3);", d.Book_id, d.Order_amount, cid)
	if err != nil {
		log.Fatal("request failed")
		e.Message = "service unavailable"
		return c.JSON(http.StatusServiceUnavailable, e)
	}
	
	return c.NoContent(http.StatusCreated)
}