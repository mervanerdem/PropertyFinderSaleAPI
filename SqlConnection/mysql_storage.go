package SqlConnection

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mervanerdem/PropertyFinderSaleAPI/Services"
	"log"
	"time"
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

// Show Cart
func (m *MStorage) ShowBasket(idCustomer int) (*[]Services.Basket, int, error) {
	var basket Services.Basket
	var basketProduct []Services.Basket
	var totalPay int
	res, err := m.client.Query("Select baskets.idBasket,baskets.idCustomer, products.idProduct, products.productName,products.productStock, "+
		"products.productPrice,products.productVat,baskets.proNum, baskets.productTotalPrice "+
		"From baskets "+
		"INNER JOIN products ON baskets.idProduct = products.idProduct "+
		"Where baskets.idCustomer = ?", idCustomer)
	if err != nil {
		log.Fatal(err)
	}

	for res.Next() {

		err := res.Scan(&basket.BasketID, &basket.CustomerID, &basket.Product.ProductID, &basket.Product.ProductName, &basket.Product.ProductStock,
			&basket.Product.ProductPrice, &basket.Product.ProductVAT, &basket.ProductNum, &basket.ProductTotalPrice)
		if err != nil {
			log.Fatal(err)
		}

		totalPay += basket.ProductPrice * basket.ProductNum

		basketProduct = append(basketProduct, basket)
	}
	if basket.CustomerID <= 0 {
		err = fmt.Errorf("this customer does not have any thing in her/him basket")
	}
	return &basketProduct, totalPay, err
}

// read from baskets table
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

// read from products table
func (m *MStorage) IsHaveProductID(idProduct int) (int, error) {
	var product Services.Product
	haveProductID := false
	res, err := m.client.Query("select idProduct from products")
	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		err := res.Scan(&product.ProductID)
		if err != nil {
			log.Fatal(err)
		}
		if product.ProductID == idProduct {
			haveProductID = true
		}
	}
	if haveProductID {
		res, err = m.client.Query("select productPrice from products Where idProduct = ?", idProduct)
		for res.Next() {
			err := res.Scan(&product.ProductPrice)
			if err != nil {
				log.Fatal(err)
			}
		}
		return product.ProductPrice, nil
	} else {
		return 0, fmt.Errorf("the product id does not exist")
	}

}

// add to cart as new
func (m *MStorage) AddBasket(idCustomer, idProduct, productNum, productTotalPrice int) error {
	_, err := m.client.Query("insert into baskets (idCustomer, idProduct, proNum,productTotalPrice) values (?, ?, ?,?)", idCustomer, idProduct, productNum, productTotalPrice)
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
	productTotalPrice := productNum * basket.ProductPrice

	_, err = m.client.Query("UPDATE baskets SET proNum = ?, productTotalPrice = ? WHERE idProduct = ? and idcustomer = ?", productNum, productTotalPrice, idProduct, idCustomer)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// DeleteRow delete cart Row
func (m *MStorage) DeleteRow(idCustomer, idProduct int) error {
	_, err := m.client.Query("DELETE FROM baskets where idProduct = ? and idCustomer = ?", idProduct, idCustomer)
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

func (m *MStorage) Sale(idCustomer int) error {
	var basket Services.Basket
	res, err := m.client.Query("select idCustomer,idProduct,proNum,productTotalPrice from baskets where idCustomer = ?", idCustomer)
	if err != nil {
		log.Fatal(err)
	}
	currentTime := time.Now()
	saleDate := currentTime.Format("2022-01-02")
	for res.Next() {

		err := res.Scan(&basket.CustomerID, &basket.Product.ProductID, &basket.ProductNum, &basket.ProductTotalPrice)
		if err != nil {
			log.Fatal(err)
		}
		_, err = m.client.Query("insert into sales (idCustomer, idProduct, proNum,productTotalPrice,saleDate) "+
			"values (?, ?, ?,?,?);", basket.CustomerID, basket.Product.ProductID, basket.ProductNum, basket.ProductTotalPrice, saleDate)
		if err != nil {
			log.Fatal(err)
		}
		_, err = m.client.Query("DELETE FROM baskets where idCustomer = ?", idCustomer)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
