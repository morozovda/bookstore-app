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
		Name string `json:"name"`
		Email string `json:"email"`
		Passwd string `json:"passwd"`
	}
	
	Deal struct {
		Book_id int `json:"book"`
		Order_amount int `json:"amount"`
		Customer_id int `json:"id"`
	}
	
	Account struct {
		Books []Book `json:"books"`
		Balance int `json:"balance"`
	}
	
	Error struct {
		Message string `json:"message"`
	}
)