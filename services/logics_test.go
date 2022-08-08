package services

import "testing"

func TestBasket(t *testing.T) {
	assertCorrectMessage := func(t testing.TB, got, want float64) {
		t.Helper()
		if got != want {
			t.Errorf("got %f want %f", got, want)
		}
	}

	t.Run("Campaing 1 VAT 1 order number lower 4", func(t *testing.T) {
		dummyBasket := getDummyBasket(1, 4, 1)
		got := dummyBasket.Campaign1(1)
		want := dummyBasket.ProductTotalPrice
		assertCorrectMessage(t, got, want)
	})
	t.Run("Campaing 1 VAT 8 order number lower 4", func(t *testing.T) {
		dummyBasket := getDummyBasket(2, 40, 8)
		got := dummyBasket.Campaign1(2)
		want := dummyBasket.ProductTotalPrice
		assertCorrectMessage(t, got, want)
	})
	t.Run("Campaing 1 VAT 18 order number lower 4", func(t *testing.T) {
		dummyBasket := getDummyBasket(3, 50, 18)
		got := dummyBasket.Campaign1(2)
		want := dummyBasket.ProductTotalPrice
		assertCorrectMessage(t, got, want)
	})

	t.Run("Campaing 1 VAT 1 4th order", func(t *testing.T) {
		dummyBasket := getDummyBasket(1, 100, 1)
		got := dummyBasket.Campaign1(3)
		want := dummyBasket.ProductTotalPrice
		assertCorrectMessage(t, got, want)
	})

	t.Run("Campaing 1 VAT 8 4th order", func(t *testing.T) {
		dummyBasket := getDummyBasket(2, 100, 8)
		want := dummyBasket.ProductTotalPrice * 0.9
		got := dummyBasket.Campaign1(4)
		assertCorrectMessage(t, got, want)
	})

	t.Run("Campaing 1 VAT 18 4th order", func(t *testing.T) {
		dummyBasket := getDummyBasket(3, 100, 18)
		want := dummyBasket.ProductTotalPrice * 0.85
		got := dummyBasket.Campaign1(4)
		assertCorrectMessage(t, got, want)
	})

	t.Run("Campaing 2 product number lower than 4", func(t *testing.T) {
		dummyBasket := getDummyBasket(1, 100, 1)
		want := 100.0
		got := dummyBasket.Campaign2()
		assertCorrectMessage(t, got, want)
	})

	t.Run("Campaing 2 product number after 3 product", func(t *testing.T) {
		dummyBasket := getDummyBasket(10, 100, 1)
		want := 944.0
		got := dummyBasket.Campaign2()
		assertCorrectMessage(t, got, want)
	})

	t.Run("Campaing 3 if last one month sales lower than given amount", func(t *testing.T) {
		dummyBasket := getDummyBasket(10, 100, 1)
		want := dummyBasket.ProductTotalPrice
		got := dummyBasket.Campaign3(1000)
		assertCorrectMessage(t, got, want)
	})

	t.Run("Campaing 3 if last one month sales higher than given amount", func(t *testing.T) {
		dummyBasket := getDummyBasket(10, 100, 1)
		want := dummyBasket.ProductTotalPrice * 0.9
		got := dummyBasket.Campaign3(6000)
		assertCorrectMessage(t, got, want)
	})
}

func getDummyBasket(ProductNum int, ProductPrice float64, ProductVAT int) Basket {
	var basket Basket
	basket.ProductNum = ProductNum
	basket.ProductPrice = ProductPrice
	basket.ProductVAT = ProductVAT
	basket.ProductTotalPrice = float64(basket.ProductNum) * basket.ProductPrice
	basket.BasketID = 1
	basket.CustomerID = 1
	basket.ProductID = 1
	basket.ProductName = "test product"

	return basket
}
