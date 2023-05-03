package models

import (
	"github.com/google/uuid"
)

type (
	Jwts struct {
		Jwtstr string `json:"token"`
	}
	
	Book struct {
		Id uuid.UUID `json:"id"`
		Title string `json:"title"`
		Author string `json:"author"`
		Price int `json:"price"`
		Amount int `json:"amount"`
	}
	
	Market struct {
		Books []Book `json:"books"`
	}
	
	Regcustomer struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Passwd string `json:"passwd"`
	}
	
	Logincustomer struct {
		Email string `json:"email"`
		Passwd string `json:"passwd"`
	}
	
	Deal struct {
		Book_id uuid.UUID `json:"book"`
		Order_amount int `json:"amount"`
	}
	
	Account struct {
		Books []Book `json:"books"`
		Balance int `json:"balance"`
	}
	
	Error struct {
		Message string `json:"message"`
	}
)