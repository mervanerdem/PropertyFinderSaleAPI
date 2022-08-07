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

// Database Configurations
func NewMStorage(dsn string) (*MStorage, *sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	errHandle(err, "database configurations")
	return &MStorage{client: db}, db, nil
}

// List Products
func (m *MStorage) ListProducts() (*[]Services.Product, error) {
	var product Services.Product
	var products []Services.Product
	list, err := m.client.Query("SELECT * FROM products")
	errHandle(err, "list product query")
	for list.Next() {

		err := list.Scan(&product.ProductID, &product.ProductName, &product.ProductPrice, &product.ProductVAT)
		errHandle(err, "list product scan")

		products = append(products, product)
	}
	return &products, nil
}

// Show Cart
func (m *MStorage) ShowBasket(idCustomer int) (*[]Services.Basket, float64, error) {
	queryString := "Select baskets.idBasket,baskets.idCustomer, products.idProduct, products.productName, " +
		"products.productPrice,products.productVat,baskets.proNum, baskets.productTotalPrice " +
		"From baskets " +
		"INNER JOIN products ON baskets.idProduct = products.idProduct " +
		"Where baskets.idCustomer = ?"

	show, err := m.client.Query(queryString, idCustomer)
	errHandle(err, "show cart list query")

	var basket Services.Basket
	var basketProduct []Services.Basket
	var totalPay float64
	var campaign1 float64
	var campaign2 float64
	var campaign3 float64
	var campaignTotal1 float64
	var campaignTotal2 float64
	var campaignTotal3 float64
	var campaignOrderNumber int
	var lastSales float64

	for show.Next() {

		err := show.Scan(&basket.BasketID, &basket.CustomerID, &basket.Product.ProductID, &basket.Product.ProductName,
			&basket.Product.ProductPrice, &basket.Product.ProductVAT, &basket.ProductNum, &basket.ProductTotalPrice)
		errHandle(err, "show cart list scan")

		rowCheck, err := m.client.Exec("select campaignOrderNumber from sales where idCustomer = ? ORDER BY saleDate DESC LIMIT 1;", idCustomer)
		errHandle(err, "campaign order number is have any row")
		affectRow, err := rowCheck.RowsAffected()
		errHandle(err, "campaign order number is have any row 2")
		if affectRow > 0 {
			//campaign 1
			campOrdNum, err := m.client.Query("select campaignOrderNumber from sales where idCustomer = ? ORDER BY saleDate DESC LIMIT 1;", idCustomer)
			errHandle(err, "campaign 1 query")
			campOrdNum.Next()
			err = campOrdNum.Scan(&campaignOrderNumber)
			errHandle(err, "campaign 1 scan 1")

			//campaign 3
			currentTime := time.Now()
			now := currentTime.AddDate(0, -1, 0)
			lastMonth := now.Format("2006.01.02 15:04:05")

			proTotPri, err := m.client.Query("select Sum(productTotalPrice)  from sales where idCustomer = ? and saleDate > ?", basket.CustomerID, lastMonth)
			errHandle(err, "product total price query 1")
			for proTotPri.Next() {
				err = proTotPri.Scan(&lastSales)
				errHandle(err, "product total price scan 1")
			}
		}

		campaign1 = basket.Campaign1(campaignOrderNumber)
		campaign2 = basket.Campaign2()
		campaign3 = basket.Campaign3(lastSales)

		campaignTotal1 += campaign1
		campaignTotal2 += campaign2
		campaignTotal3 += campaign3

	}

	finalCampaign := compareCampaign(campaignTotal1, campaignTotal2, campaignTotal3)

	log.Println("campaign1:", campaignTotal1)
	log.Println("campaign2:", campaignTotal2)
	log.Println("campaign3:", campaignTotal3)
	log.Println("final:", campaignTotal3)

	if finalCampaign == campaignTotal3 {
		//campaign 3
		show, err := m.client.Query(queryString, idCustomer)
		errHandle(err, "campaign3 show query")
		for show.Next() {

			err := show.Scan(&basket.BasketID, &basket.CustomerID, &basket.Product.ProductID, &basket.Product.ProductName,
				&basket.Product.ProductPrice, &basket.Product.ProductVAT, &basket.ProductNum, &basket.ProductTotalPrice)
			errHandle(err, "campaign3 show scan")

			currentTime := time.Now()
			now := currentTime.AddDate(0, -1, 0)
			lastMonth := now.Format("2006.01.02 15:04:05")

			proTotPri, err := m.client.Query("select Sum(productTotalPrice)  from sales where idCustomer = ? and saleDate > ?", basket.CustomerID, lastMonth)
			errHandle(err, "product total price query 2")
			for proTotPri.Next() {
				err = proTotPri.Scan(&lastSales)
				errHandle(err, "product total price scan 2")
			}

			if lastSales > 2000 {
				basket.ProductTotalPrice = basket.ProductPrice * float64(basket.ProductNum) * 0.9
			}

			_, err = m.client.Query("UPDATE baskets SET productTotalPrice = ? WHERE idProduct = ? and idcustomer = ?", basket.ProductTotalPrice, basket.ProductID, basket.CustomerID)
			errHandle(err, "campaign 3 update query")
			totalPay = totalPay + basket.ProductTotalPrice

			basketProduct = append(basketProduct, basket)
		}
	} else if finalCampaign == campaignTotal2 {
		//campaign 2
		show, err := m.client.Query(queryString, idCustomer)
		errHandle(err, "campaign2 show query")
		for show.Next() {

			err := show.Scan(&basket.BasketID, &basket.CustomerID, &basket.Product.ProductID, &basket.Product.ProductName,
				&basket.Product.ProductPrice, &basket.Product.ProductVAT, &basket.ProductNum, &basket.ProductTotalPrice)
			errHandle(err, "campaign2 show scan")

			campaign2 = basket.Campaign2()
			basket.ProductTotalPrice = campaign2

			_, err = m.client.Query("UPDATE baskets SET productTotalPrice = ? WHERE idProduct = ? and idcustomer = ?", basket.ProductTotalPrice, basket.ProductID, basket.CustomerID)
			errHandle(err, "campaign 2 update query")
			totalPay = totalPay + basket.ProductTotalPrice

			basketProduct = append(basketProduct, basket)
		}

	} else if finalCampaign == campaignTotal1 {
		show, err := m.client.Query(queryString, idCustomer)
		errHandle(err, "campaign1 show query")
		for show.Next() {

			err := show.Scan(&basket.BasketID, &basket.CustomerID, &basket.Product.ProductID, &basket.Product.ProductName,
				&basket.Product.ProductPrice, &basket.Product.ProductVAT, &basket.ProductNum, &basket.ProductTotalPrice)
			errHandle(err, "campaign1 show scan")

			campOrdNum, err := m.client.Query("select campaignOrderNumber from sales where idCustomer = ? ORDER BY saleDate DESC LIMIT 1;", idCustomer)
			errHandle(err, "campaign 1 query")
			campOrdNum.Next()
			err = campOrdNum.Scan(&campaignOrderNumber)
			errHandle(err, "campaign 1 scan 2")

			basket.ProductTotalPrice = campaign1

			totalPay = totalPay + basket.ProductTotalPrice
			_, err = m.client.Query("UPDATE baskets SET productTotalPrice = ? WHERE idProduct = ? and idcustomer = ?", basket.ProductTotalPrice, basket.ProductID, basket.CustomerID)
			errHandle(err, "campaign 1 update query")

			basketProduct = append(basketProduct, basket)
		}
	}

	if basket.CustomerID <= 0 {
		err = fmt.Errorf("this customer does not have anything in her/him basket")
	}
	return &basketProduct, totalPay, err
}

// read from baskets table
func (m *MStorage) IsHaveProductNumber(idCustomer, idProduct int) (bool, int, error) {
	var basket Services.Basket
	haveProduct := false
	checkProduct, err := m.client.Query("Select proNum From baskets Where idProduct = ? and idcustomer = ?", idProduct, idCustomer)
	errHandle(err, "check product from basket query")

	for checkProduct.Next() {
		err := checkProduct.Scan(&basket.ProductNum)
		errHandle(err, "check product from basket scan")
		if basket.ProductNum >= 1 {
			haveProduct = true
		} else {
			haveProduct = false
		}
	}
	return haveProduct, basket.ProductNum, nil
}

// read from products table
func (m *MStorage) IsHaveProductID(idProduct int) (float64, error) {
	var product Services.Product
	haveProductID := false
	checkID, err := m.client.Query("select idProduct from products")
	errHandle(err, "check product id from products, query")
	for checkID.Next() {
		err := checkID.Scan(&product.ProductID)
		errHandle(err, "check product id from products scan")
		if product.ProductID == idProduct {
			haveProductID = true
		}
	}
	if haveProductID {
		proPrice, err := m.client.Query("select productPrice from products Where idProduct = ?", idProduct)
		errHandle(err, "check product price from products, query")
		for proPrice.Next() {
			err := proPrice.Scan(&product.ProductPrice)
			errHandle(err, "check product price from products, scan")
		}
		return product.ProductPrice, nil
	} else {
		return 0, fmt.Errorf("the product id does not exist")
	}

}

// add to cart as new
func (m *MStorage) AddBasket(idCustomer, idProduct, productNum int, productTotalPrice float64) error {
	_, err := m.client.Query("insert into baskets (idCustomer, idProduct, proNum,productTotalPrice) values (?, ?, ?,?)", idCustomer, idProduct, productNum, productTotalPrice)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// add to cart in row
func (m *MStorage) AddCartItem(idCustomer, idProduct, productNum int) error {
	var basket Services.Basket
	addCartPro, err := m.client.Query("Select proNum From baskets Where idProduct = ? and idcustomer = ?", idProduct, idCustomer)
	errHandle(err, "add cart item pronum query")
	for addCartPro.Next() {
		err = addCartPro.Scan(&basket.ProductNum)
		errHandle(err, "add cart item pronum scan")
	}

	addCartPrice, err := m.client.Query("Select productPrice From products Where idProduct = ?", idProduct)
	errHandle(err, "add cart item product price query")
	for addCartPrice.Next() {
		err = addCartPrice.Scan(&basket.ProductPrice)
		errHandle(err, "add cart item product price scan")
	}

	productTotalPrice := float64(productNum) * basket.ProductPrice

	_, err = m.client.Query("UPDATE baskets SET proNum = ?, productTotalPrice = ? WHERE idProduct = ? and idcustomer = ?", productNum, productTotalPrice, idProduct, idCustomer)
	errHandle(err, "update add cart item query")

	return err
}

// delete cart Row
func (m *MStorage) DeleteRow(idCustomer, idProduct int) error {
	rowCheck, err := m.client.Exec("DELETE FROM baskets where idProduct = ? and idCustomer = ?", idProduct, idCustomer)
	errHandle(err, "delete row query")
	delRow, err := rowCheck.RowsAffected()
	errHandle(err, "delete row , row affected")
	if delRow == 0 {
		err = fmt.Errorf("no items have been deleted")
	}

	return err
}

// delete cart in row
func (m *MStorage) DeleteCartItem(idCustomer, idProduct, productNum int) error {
	var basket Services.Basket

	deleteProNum, err := m.client.Query("Select proNum From baskets Where idProduct = ? and idcustomer = ?", idProduct, idCustomer)
	errHandle(err, "delete cart item pronum query")
	for deleteProNum.Next() {
		err = deleteProNum.Scan(&basket.ProductNum)
		errHandle(err, "delete cart item pronum scan")

		if basket.ProductNum < productNum {
			return fmt.Errorf("product number can not be higher than given amount")
		}
	}
	basket.ProductNum = basket.ProductNum - productNum
	_, err = m.client.Query("UPDATE baskets SET proNum = ? WHERE idProduct = ? and idcustomer = ?", basket.ProductNum, idProduct, idCustomer)
	errHandle(err, "delete cart item update")

	return err
}

// sale
func (m *MStorage) Sale(idCustomer int) error {
	var sumOrder float64
	var basket Services.Basket
	var lastSaleDate string

	saleFind, err := m.client.Query("select idCustomer,idProduct,proNum,productTotalPrice from baskets where idCustomer = ?", idCustomer)
	errHandle(err, "sale find basket query")

	rowCheck, err := m.client.Exec("select campaignOrderNumber from sales where idCustomer = ? ORDER BY saleDate DESC LIMIT 1;", idCustomer)
	errHandle(err, "sale campaign order number is have any row")
	affectRow, err := rowCheck.RowsAffected()
	errHandle(err, "sale campaign order number is have any row")
	campaignOrderNumber := 0
	if affectRow > 0 {
		campOrder, err := m.client.Query("select campaignOrderNumber,saleDate from sales where idCustomer = ? ORDER BY saleDate DESC LIMIT 1;", idCustomer)
		errHandle(err, "sale campaign order number query")
		campOrder.Next()
		err = campOrder.Scan(&campaignOrderNumber, &lastSaleDate)
		errHandle(err, "sale campaign order number scan")

		saleTotalPrice, err := m.client.Query("select Sum(productTotalPrice)  from sales where idCustomer = ? and saleDate = ?", idCustomer, lastSaleDate)
		errHandle(err, "sale find product total price query")
		saleTotalPrice.Next()
		err = saleTotalPrice.Scan(&sumOrder)
		errHandle(err, "sale find product total price scan")
	}

	if sumOrder > Services.LimitMonthShop {
		campaignOrderNumber++
	}
	if campaignOrderNumber == 4 {
		campaignOrderNumber = 0
	}

	currentTime := time.Now()
	saleDate := currentTime.Format("2006.01.02 15:04:05")

	for saleFind.Next() {

		err := saleFind.Scan(&basket.CustomerID, &basket.Product.ProductID, &basket.ProductNum, &basket.ProductTotalPrice)
		errHandle(err, "sale find basket ")
		_, err = m.client.Query("insert into sales (idCustomer, idProduct, proNum,productTotalPrice,saleDate,campaignOrderNumber) "+
			"values (?, ?, ?,?,?,?);", basket.CustomerID, basket.Product.ProductID, basket.ProductNum, basket.ProductTotalPrice, saleDate, campaignOrderNumber)
		errHandle(err, "insert sales data query")
		_, err = m.client.Query("DELETE FROM baskets where idCustomer = ?", idCustomer)
		errHandle(err, "delete sales data from basket")
	}
	return nil
}

// Handle SQL Errors
func errHandle(err error, errorName string) {
	if err != nil {
		log.Fatal(errorName, "\n", err)
	}
}

func compareCampaign(campaignTotal1, campaignTotal2, campaignTotal3 float64) float64 {
	var finalCampaign float64

	e := make([]float64, 0)
	e = append(e, campaignTotal1)
	e = append(e, campaignTotal2)
	e = append(e, campaignTotal3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if (e[i] < e[j] && e[i] > 0) || (e[i] == e[j] && e[i] != 0) {
				finalCampaign = e[i]
				break
			}
		}
	}

	return finalCampaign
}
