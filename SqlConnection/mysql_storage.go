package SqlConnection

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mervanerdem/PropertyFinderSaleAPI/Services"
	"log"
)

type MStorage struct {
	client *sql.DB
}

// NewMStorage Database Configurations
func NewMStorage(dsn string) (*MStorage, *sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	return &MStorage{client: db}, db, nil
}

// ListProducts List Products
func (m *MStorage) ListProducts() (*[]Services.Product, error) {
	var product Services.Product
	var products []Services.Product
	res, err := m.client.Query("SELECT * FROM products")
	if err != nil {
		log.Fatal(err)
	}

	for res.Next() {

		err := res.Scan(&product.ProductID, &product.ProductName, &product.ProductStock, &product.ProductPrice, &product.ProductVAT)

		if err != nil {
			log.Fatal(err)
		}

		products = append(products, product)
	}
	return &products, nil
}

// ShowBasket Show Cart
func (m *MStorage) ShowBasket(idCustomer int) (*[]Services.Basket, int, error) {
	var basket Services.Basket
	var basketProduct []Services.Basket
	var totalPay int
	res, err := m.client.Query("Select baskets.idBasket,baskets.idCustomer, products.idProduct, products.productName,products.productStock,products.productPrice,products.productVat,baskets.proNum "+
		"From baskets "+
		"INNER JOIN products ON baskets.idProduct = products.idProduct "+
		"INNER JOIN customers ON customers.idCustomer = baskets.idCustomer "+
		"Where baskets.idCustomer = ?", idCustomer)
	if err != nil {
		log.Fatal(err)
	}

	for res.Next() {

		err := res.Scan(&basket.BasketID, &basket.CustomerID, &basket.Product.ProductID, &basket.Product.ProductName, &basket.Product.ProductStock, &basket.Product.ProductPrice, &basket.Product.ProductVAT, &basket.ProductNum)
		if err != nil {
			log.Fatal(err)
		}

		totalPay += basket.ProductPrice * basket.ProductNum

		basketProduct = append(basketProduct, basket)
	}

	return &basketProduct, totalPay, nil
}

// IsHaveProductNumber in rows
func (m *MStorage) IsHaveProductNumber(idCustomer, idProduct int) (bool, int, error) {
	var basket Services.Basket
	haveProduct := false
	res, err := m.client.Query("Select proNum From baskets Where idProduct = ? and idcustomer = ?", idProduct, idCustomer)
	if err != nil {
		log.Fatal(err)
	}

	for res.Next() {
		err := res.Scan(&basket.ProductNum)
		if err != nil {
			log.Fatal(err)
		}
		if basket.ProductNum >= 1 {
			haveProduct = true
		} else {
			haveProduct = false
		}
	}

	return haveProduct, basket.ProductNum, nil
}

// AddBasket add to cart as new
func (m *MStorage) AddBasket(idCustomer, idProduct int, productNum int) error {
	_, err := m.client.Query("insert into baskets (idCustomer, idProduct, proNum) values (?, ?, ?)", idCustomer, idProduct, productNum)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// AddCartItem add to cart in row
func (m *MStorage) AddCartItem(idCustomer, idProduct, productNum int) error {
	var basket Services.Basket
	res, err := m.client.Query("Select proNum From baskets Where idProduct = ? and idcustomer = ?", idProduct, idCustomer)
	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		err = res.Scan(&basket.ProductNum)
		if err != nil {
			log.Fatal(err)
		}
	}
	productNum = productNum + basket.ProductNum
	_, err = m.client.Query("UPDATE baskets SET proNum = ? WHERE idProduct = ? and idcustomer = ?", productNum, idProduct, idCustomer)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// DeleteRow delete cart Row
func (m *MStorage) DeleteRow(idCustomer, idProduct int) error {
	_, err := m.client.Query("DELETE FROM baskets where idProduct = ? and idCustomer = ?", idCustomer, idProduct)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// DeleteCartItem delete cart in row
func (m *MStorage) DeleteCartItem(idCustomer, idProduct, productNum int) error {
	var basket Services.Basket

	res, err := m.client.Query("Select proNum From baskets Where idProduct = ? and idcustomer = ?", idProduct, idCustomer)
	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		err = res.Scan(&basket.ProductNum)
		if err != nil {
			log.Fatal(err)
		}

		if basket.ProductNum < productNum {
			return fmt.Errorf("product number can not be higher than given amount")
		}

	}
	basket.ProductNum = basket.ProductNum - productNum
	_, err = m.client.Query("UPDATE baskets SET proNum = ? WHERE idProduct = ? and idcustomer = ?", basket.ProductNum, idProduct, idCustomer)

	if err != nil {
		log.Fatal(err)
	}

	return err
}
