package models

import (
	"gl/db"
)

type Product struct {
	Id                 *int64   `db:"id" form:"id"`
	ProductCode        string   `db:"product_code" form:"product_code"`
	ProductDescription string   `db:"product_description" form:"product_description"`
	Cogs               *float64 `db:"cogs" form:"cogs"`
	CurrentStock       *float64 `db:"current_stock" form:"current_stock"`
	InventoryAccountId *int64   `db:"inventory_account_id" form:"inventory_account_id"`
	SalesAccountId     *int64   `db:"sales_account_id" form:"sales_account_id"`
	COGSAccountId      *int64   `db:"cogs_account_id" form:"cogs_account_id"`
	Unit               *string  `db:"unit" form:"unit"`
	UserId             int64    `db:"user_id" form:"user_id"`
}

type SubAccountType struct {
	Id             int64
	AccountCode    string
	SubAccountName string
	Subtype        string
}

type Unit struct {
	Id       string
	UnitName string
}

var Units []Unit = []Unit{
	{Id: "pcs", UnitName: "Piece"},
	{Id: "lbs", UnitName: "Pound"},
	{Id: "kg", UnitName: "Kilogram"},
	{Id: "Dz", UnitName: "Dozen"},
}

func GetAllSubAccountType() (map[string][]SubAccountType, error) {

	grouped := make(map[string][]SubAccountType)
	var all []SubAccountType

	rows, err := db.Conn.Query(`
		SELECT id, account_code, account_name, subtype
		FROM general_ledger.accounts
		WHERE subtype IS NOT NULL
		ORDER BY account_name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var subAc SubAccountType

		if err := rows.Scan(
			&subAc.Id,
			&subAc.AccountCode,
			&subAc.SubAccountName,
			&subAc.Subtype,
		); err != nil {
			return nil, err
		}

		// add to "all"
		all = append(all, subAc)

		// add to grouped
		grouped[subAc.Subtype] = append(grouped[subAc.Subtype], subAc)
	}

	return grouped, nil
}

func GetProduct(id int64) (*Product, error) {
	var product Product
	err := db.Conn.QueryRow(`SELECT id, product_code, product_description, inventory_account_id, sales_account_id, cogs_account_id, unit, cogs, current_stock FROM general_ledger.products_cogs where id = $1`,
		id).Scan(&product.Id, &product.ProductCode, &product.ProductDescription, &product.InventoryAccountId,
		&product.SalesAccountId, &product.COGSAccountId, &product.Unit, &product.Cogs, &product.CurrentStock)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func SaveProduct(product Product) (*int64, error) {
	cogs := 0.0
	currentStock := 0.0
	product.Cogs = &cogs
	product.CurrentStock = &currentStock

	if product.Id != nil {
		return UpdateProduct(product)
	}
	var id int64
	err := db.Conn.QueryRow(
		`INSERT INTO general_ledger.products_cogs (product_code, product_description, cogs, current_stock, 
		inventory_account_id, sales_account_id, cogs_account_id, unit, user_id) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		product.ProductCode,
		product.ProductDescription,
		product.Cogs,
		product.CurrentStock,
		product.InventoryAccountId,
		product.SalesAccountId,
		product.COGSAccountId,
		product.Unit,
		product.UserId,
	).Scan(&id)

	if err != nil {
		return nil, err
	}
	return &id, nil
}

func UpdateProduct(product Product) (*int64, error) {
	var id int64
	err := db.Conn.QueryRow(
		`UPDATE general_ledger.products_cogs SET product_code = $1, product_description = $2, 
		inventory_account_id = $3, sales_account_id = $4, cogs_account_id = $5, unit = $6 WHERE id = $7 RETURNING id`,
		product.ProductCode,
		product.ProductDescription,
		product.InventoryAccountId,
		product.SalesAccountId,
		product.COGSAccountId,
		product.Unit,
		*product.Id,
	).Scan(&id)

	if err != nil {
		return nil, err
	}
	return &id, nil
}

func ListAllProducts(productCode string) ([]Product, error) {
	products := []Product{}
	rows, err := db.Conn.Query(`SELECT id, product_code, product_description, inventory_account_id, 
	sales_account_id, cogs_account_id, unit, cogs, current_stock FROM general_ledger.products_cogs
	where product_code like $1 ORDER BY product_code ASC`, "%"+productCode+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Id, &product.ProductCode, &product.ProductDescription, &product.InventoryAccountId,
			&product.SalesAccountId, &product.COGSAccountId, &product.Unit, &product.Cogs, &product.CurrentStock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func FindProductByCode(productCode string) ([]Product, error) {
	var products []Product
	rows, err := db.Conn.Query(`SELECT id, product_code, product_description, inventory_account_id, sales_account_id, cogs_account_id, unit, cogs, current_stock FROM general_ledger.products_cogs WHERE product_code like $1`, "%"+productCode+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Id, &p.ProductCode, &p.ProductDescription, &p.InventoryAccountId,
			&p.SalesAccountId, &p.COGSAccountId, &p.Unit, &p.Cogs, &p.CurrentStock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
