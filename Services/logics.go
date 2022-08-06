package Services

type Product struct {
	ProductID    int
	ProductName  string
	ProductStock int
	ProductPrice int
	ProductVAT   int
}

type Basket struct {
	BasketID   int
	CustomerID int
	Product
	ProductNum        int
	ProductTotalPrice int
}

type PStorage interface {
	ListProducts() (*[]Product, error)
	ShowBasket(id int) (*[]Basket, int, error)
	IsHaveProductNumber(idCustomer, productID int) (bool, int, error)
	IsHaveProductID(idProduct int) (int, error)
	AddBasket(idCustomer, idProduct, productNum, productTotalPrice int) error
	AddCartItem(idCustomer, idProduct, productNum int) error
	DeleteRow(idCustomer, idProduct int) error
	DeleteCartItem(idCustomer, idProduct, productNum int) error
	Sale(idCustomer int) error
}

func (p *Product) List() Product {
	return *p
}

func (b *Basket) ListBasket() Basket {
	return *b
}
