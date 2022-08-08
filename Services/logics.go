package Services

const LimitMonthShop float64 = 5000
const Limit4sales float64 = 500

type Product struct {
	ProductID    int
	ProductName  string
	ProductPrice float64
	ProductVAT   int
}

type Basket struct {
	BasketID   int
	CustomerID int
	Product
	ProductNum        int
	ProductTotalPrice float64
}

type Sale struct {
	CustomerID int
	Product
	ProductNum        int
	ProductTotalPrice float64
	SaleDate          string
	CampaignOrderNum  int
}

type PStorage interface {
	ListProducts() (*[]Product, error)
	ShowBasket(id int) (*[]Basket, float64, error)
	HaveProductNumber(idCustomer, productID int) (bool, int, error)
	FindProductPrice(idProduct int) (float64, error)
	AddBasket(idCustomer, idProduct, productNum int, productTotalPrice float64) error
	AddCartItem(idCustomer, idProduct, productNum int) error
	DeleteRow(idCustomer, idProduct int) error
	DeleteCartItem(idCustomer, idProduct, productNum int) error
	Sale(idCustomer int) (*[]Sale, float64, error)
}

// campaign 1 = check 4 sales and VAT
func (basket *Basket) Campaign1(campaignOrderNumber int) float64 {
	var campaign1 float64
	if campaignOrderNumber == 4 && basket.ProductVAT != 1 {
		if basket.ProductVAT == 18 {
			campaign1 = basket.ProductPrice * float64(basket.ProductNum) * 0.85
		} else {
			campaign1 = basket.ProductPrice * float64(basket.ProductNum) * 0.90
		}
	} else {
		campaign1 = basket.ProductPrice * float64(basket.ProductNum)
	}
	return campaign1
}

// campaign 2 = same product
func (basket *Basket) Campaign2() float64 {

	var campaignTotal2 float64

	if basket.ProductNum > 3 {
		campaignTotal2 = 3 * basket.ProductPrice
		b := (float64(basket.ProductNum) - 3) * basket.ProductPrice
		campaignTotal2 = campaignTotal2 + b*0.92
	} else {
		campaignTotal2 = basket.ProductPrice * float64(basket.ProductNum)
	}

	return campaignTotal2
}

// campaign3 = month subscriber
func (basket *Basket) Campaign3(lastSales float64) float64 {
	var campaignTotal3 float64
	if lastSales > LimitMonthShop {
		campaignTotal3 = basket.ProductPrice * float64(basket.ProductNum) * 0.9
	} else {
		campaignTotal3 = basket.ProductPrice * float64(basket.ProductNum)
	}

	return campaignTotal3
}
