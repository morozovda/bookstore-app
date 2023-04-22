package models

type (
	Ids struct {
		Id int `json:"id"`
	}
	
	Book struct {
		Id int `json:"id"`
		Title string `json:"title"`
		Author string `json:"author"`
		Price int `json:"price"`
		Amount int `json:"amount"`
	}
	
	Market struct {
		Books []Book `json:"books"`
	}
	
	Customer struct {
		Name string
		Email string
		Passwd string
	}
	
	Deal struct {
		Book_id int
		Order_amount int 
		Customer_id int
	}
	
	Account struct {
		Books []Book `json:"books"`
		Balance int `json:"balance"`
	}
	
	Error struct {
		Message string `json:"message"`
	}
)