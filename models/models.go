package models

import( "time" )

type Product struct{
	Name string `json:"product_name"`
	Image []string `json:"image"`
	Price float64 `json:"price"`
	Desc string `json:"description"`
	Sizes []string `json:"sizes"`
	Category string `json:"category"`
	CategoryId string `json:"category_id"`
	Instock bool `json:"instock"`
	IsNewArrival bool `json:"newarrival"`

	Id string `json:"id,omitempty"`
	DateAdded time.Time `json:"updated_on"`
}

type Category struct{
	Name string `json:"product_name"`
	Image string `json:"image"`
	Description string `json:"description"`
	NewArrivals int `json:"newarrivals"`

	Id string `json:"id"`
	DateAdded time.Time `json:"created_on"`
}

type SubCategory struct{
	Name string `json:"product_name"`
	Image string `json:"image"`
	Description string `json:"description"`
	ParentId string `json:"parent_category"`
	NewArrivals int `json:"newarrivals"`

	Id string `json:"id"`
	DateAdded time.Time `json:"created_on"`
}

type Admin struct{
	UserID string `json:"userid" example:"admin"`
	Id string `json:"id"`
	Pass string `json:"pass" example:"pass"`
}

type Login struct{
	UserID string `json:"userid" example:"admin"`
	Pass string `json:"pass" example:"pass"`
}

type Contact struct{
	Contact string `json:"contact"`
	Whatsapp string `json:"whatsapp"`
}

type HidePrice struct{
	HidePrice bool `bson:"hideprice"`
}

type Promotion struct{
	Image string `json:"image"`
	FromDate time.Time `json:"from"`
	ToDate time.Time `json:"to"`

	Id string `json:"id"`
	DateAdded time.Time `json:"created_on"`
}
