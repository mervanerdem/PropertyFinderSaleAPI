package sqlConnection

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mervanerdem/PropertyFinderSaleAPI/services"
	"log"
	"sort"
	"time"
)

type PStorage interface {
	ListProducts() (*[]services.Product, error)
	ShowBasket(id int) (*[]services.Basket, float64, error)
	HaveProductNumber(idCustomer, productID int) (bool, int, error)
	FindProductPrice(idProduct int) (float64, error)
	AddBasket(idCustomer, idProduct, productNum int, productTotalPrice float64) error
	AddCartItem(idCustomer, idProduct, productNum int) error
	DeleteRow(idCustomer, idProduct int) error
	DeleteCartItem(idCustomer, idProduct, productNum int) error
	Sale(idCustomer int) (*[]services.Sale, float64, error)
}

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
func (m *MStorage) ListProducts() (*[]services.Product, error) {
	var product services.Product
	var products []services.Product
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
func (m *MStorage) ShowBasket(idCustomer int) (*[]services.Basket, float64, error) {
	queryString := "Select baskets.idBasket,baskets.idCustomer, products.idProduct, products.productName, " +
		"products.productPrice,products.productVat,baskets.proNum, baskets.productTotalPrice " +
		"From baskets " +
		"INNER JOIN products ON baskets.idProduct = products.idProduct " +
		"Where baskets.idCustomer = ?"

	show, err := m.client.Query(queryString, idCustomer)
	errHandle(err, "show cart list query")

	var basket services.Basket
	var basketProduct []services.Basket
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

		//campaign 1
		campOrdNum, err := m.client.Query("select campaignOrderNumber from sales where idCustomer = ? ORDER BY saleDate DESC LIMIT 1;", idCustomer)
		errHandle(err, "campaign 1 query")
		checkHaveRow := campOrdNum.Next()
		if checkHaveRow {
			err = campOrdNum.Scan(&campaignOrderNumber)
			errHandle(err, "campaign 1 scan 1")
		}

		//campaign 3
		currentTime := time.Now()
		now := currentTime.AddDate(0, -1, 0)
		lastMonth := now.Format("2006.01.02 15:04:05")

		proTotPri, err := m.client.Query("select Sum(productTotalPrice)  from sales where idCustomer = ? and saleDate > ?", basket.CustomerID, lastMonth)
		errHandle(err, "product total price query 1")
		checkHaveRow = proTotPri.Next()
		if checkHaveRow {
			err = proTotPri.Scan(&lastSales)
			errHandle(err, "product total price scan 1")
		}

		campaign1 = basket.Campaign1(campaignOrderNumber)
		campaign2 = basket.Campaign2()
		campaign3 = basket.Campaign3(lastSales)

		campaignTotal1 += campaign1
		campaignTotal2 += campaign2
		campaignTotal3 += campaign3

	}

	finalCampaign := compareCampaign(campaignTotal1, campaignTotal2, campaignTotal3)

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
func (m *MStorage) HaveProductNumber(idCustomer, idProduct int) (bool, int, error) {
	var basket services.Basket
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
func (m *MStorage) FindProductPrice(idProduct int) (float64, error) {
	var product services.Product
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
		return product.ProductPrice, err
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
	var basket services.Basket
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
	var basket services.Basket

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
func (m *MStorage) Sale(idCustomer int) (*[]services.Sale, float64, error) {
	var sumOrder float64
	var basket services.Basket
	var sale services.Sale
	var saleProduct []services.Sale
	var totalPay float64
	var campaignOrderNumber int

	campOrder, err := m.client.Query("select campaignOrderNumber from sales where idCustomer = ? ORDER BY saleDate DESC LIMIT 1;", idCustomer)
	errHandle(err, "sale campaign order number query")
	campOrder.Scan()

	checkHaveRow := campOrder.Next()
	if checkHaveRow {
		err = campOrder.Scan(&campaignOrderNumber)
		errHandle(err, "sale campaign order number scan")
	}

	saleTotalPrice, err := m.client.Query("select Sum(productTotalPrice)  from baskets where idCustomer = ?", idCustomer)
	errHandle(err, "sale find product total price query")
	saleTotalPrice.Next()
	err = saleTotalPrice.Scan(&sumOrder)
	errHandle(err, "sale find product total price scan")

	Limit4sales := services.GetLimit4sales()

	if sumOrder > float64(Limit4sales) {
		campaignOrderNumber++
	}
	if campaignOrderNumber == 5 {
		campaignOrderNumber = 1
	}

	currentTime := time.Now()
	saleDate := currentTime.Format("2006.01.02 15:04:05")

	saleFind, err := m.client.Query("select idCustomer,idProduct,proNum,productTotalPrice from baskets where idCustomer = ?", idCustomer)
	errHandle(err, "sale find basket query")

	for saleFind.Next() {

		err := saleFind.Scan(&basket.CustomerID, &basket.Product.ProductID, &basket.ProductNum, &basket.ProductTotalPrice)
		errHandle(err, "sale find basket ")

		_, err = m.client.Query("insert into sales (idCustomer, idProduct, proNum,productTotalPrice,saleDate,campaignOrderNumber) "+
			"values (?, ?, ?,?,?,?);", basket.CustomerID, basket.Product.ProductID, basket.ProductNum, basket.ProductTotalPrice, saleDate, campaignOrderNumber)
		errHandle(err, "insert sales data query")
		_, err = m.client.Query("DELETE FROM baskets where idCustomer = ?", idCustomer)
		errHandle(err, "delete sales data from basket")

	}

	if basket.CustomerID != idCustomer {
		return nil, 0.0, fmt.Errorf("this customer does not have any product in her or him basket")
	}

	showSale, err := m.client.Query("SELECT sales.idCustomer,sales.idProduct,products.productName,products.productVat, "+
		"sales.proNum, products.productPrice, sales.productTotalPrice,sales.saleDate,campaignOrderNumber "+
		"FROM sales "+
		"INNER JOIN products ON sales.idProduct = products.idProduct "+
		"where saleDate = ? and idCustomer = ?", saleDate, idCustomer)
	errHandle(err, "show sale list query")

	for showSale.Next() {
		err = showSale.Scan(&sale.CustomerID, &sale.ProductID, &sale.Product.ProductName, &sale.Product.ProductVAT,
			&sale.ProductNum, &sale.Product.ProductPrice, &sale.ProductTotalPrice, &sale.SaleDate, &sale.CampaignOrderNum)
		errHandle(err, "show sale list scan")

		saleProduct = append(saleProduct, sale)
		totalPay = totalPay + sale.ProductTotalPrice
	}

	return &saleProduct, totalPay, err
}

// Handle SQL Errors
func errHandle(err error, errorName string) {
	if err != nil {
		log.Println(errorName, "\n", err)
	}
}

func compareCampaign(campaignTotal1, campaignTotal2, campaignTotal3 float64) float64 {
	var finalCampaign float64
	x := []float64{campaignTotal1, campaignTotal2, campaignTotal3}
	sort.Float64s(x)
	for i := 0; i < 3; i++ {
		if x[i] >= 0 {
			finalCampaign = x[i]
			break
		}
	}
	return finalCampaign
}
