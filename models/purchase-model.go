package models

import "time"

type Purchase struct {
	Id           int64   `db:"id" form:"id"`
	ProductCode  string  `db:"product_code" form:"product_code"`
	ItemQuantity float64 `db:"item_quantity" form:"item_quantity"`
	UnitPrice    float64 `db:"item_price" form:"item_price"`
}

type PurchaseDataDetail struct {
	Id            int64     `db:"id" form:"id"`
	InvoiceNumber string    `db:"invoice_number" form:"invoice_number"`
	InvoiceDate   time.Time `db:"invoice_date" form:"invoice_date"`
	TotalAmount   float64   `db:"total_amount" form:"total_amount"`
	MiscExpense   float64   `db:"misc_expense" form:"misc_expense"`
	CreatedAt     time.Time `db:"created_at" form:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" form:"updated_at"`
}

type Invoice struct {
	InvoiceLineItems []Purchase
	InvoiceDetail    PurchaseDataDetail
	UpdatedAt        time.Time
}
