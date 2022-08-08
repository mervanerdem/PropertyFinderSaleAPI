package Server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mervanerdem/PropertyFinderSaleAPI/Services"
	"net/http"
	"strconv"
)

func NewServer(storage Services.PStorage) http.Handler {
	router := gin.New()
	//get list of products
	router.GET("/api/products", func(ctx *gin.Context) {

		products, err := storage.ListProducts()
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, map[string]string{
			"Product List": "Successful",
		})
		jsonProduct, err := json.Marshal(*products)
		_, err = ctx.Writer.Write(jsonProduct)
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, map[string]string{
				"error": err.Error(),
			})
		}

	})

	//get list of basket
	router.GET("/api/:idCustomer/basket", func(ctx *gin.Context) {
		id_str := ctx.Param("idCustomer")
		idCustomer, err := strconv.Atoi(id_str)
		if err != nil || idCustomer <= 0 {
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Unsuitable ID number",
			})
			return
		}

		ShowBasket(ctx, idCustomer, storage, 200)
	})

	//Add to cart
	router.POST("/api/:idCustomer/basket/add", func(ctx *gin.Context) {
		id_str := ctx.Param("idCustomer")
		idCustomer, err := strconv.Atoi(id_str)
		if err != nil || idCustomer <= 0 {
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Unsuitable ID number",
			})
			return
		}

		var data = struct {
			ProductID     int
			ProductNumber int
		}{}
		err = ctx.BindJSON(&data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		productPrice, err := storage.FindProductPrice(data.ProductID)
		if err != nil {
			ctx.JSON(http.StatusInsufficientStorage, map[string]string{
				"error": err.Error(),
			})
			return
		}

		if data.ProductNumber <= 0 {
			ctx.JSON(http.StatusBadRequest, map[string]any{
				"ProductNumber": "Amount Unsuitable",
			})
			return
		}
		haveProductNum, pNum, err := storage.HaveProductNumber(idCustomer, data.ProductID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		if haveProductNum {
			pNum = pNum + data.ProductNumber
			err = storage.AddCartItem(idCustomer, data.ProductID, pNum)
			if err != nil {
				ctx.JSON(http.StatusNotFound, map[string]string{
					"error": err.Error(),
				})
				return
			}
		} else {
			productTotalPrice := productPrice * float64(data.ProductNumber)
			err = storage.AddBasket(idCustomer, data.ProductID, data.ProductNumber, productTotalPrice)
			if err != nil {
				ctx.JSON(http.StatusNotFound, map[string]string{
					"error": err.Error(),
				})
				return
			}
		}
		ShowBasket(ctx, idCustomer, storage, 201)
	})

	//delete cart
	router.POST("/api/:idCustomer/basket/delete", func(ctx *gin.Context) {
		id_str := ctx.Param("idCustomer")
		idCustomer, err := strconv.Atoi(id_str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Unsuitable ID number",
			})
			return
		}

		var data = struct {
			ProductID     int
			ProductNumber int
		}{}
		err = ctx.BindJSON(&data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		_, err = storage.FindProductPrice(data.ProductID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		if data.ProductNumber <= 0 {
			ctx.JSON(http.StatusBadRequest, map[string]any{
				"ProductNumber": "Amount Unsuitable",
			})
			return
		}

		haveProductNum, pNum, err := storage.HaveProductNumber(idCustomer, data.ProductID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}

		if haveProductNum && (pNum-data.ProductNumber) > 0 {
			err = storage.DeleteCartItem(idCustomer, data.ProductID, data.ProductNumber)
			if err != nil {
				ctx.JSON(http.StatusNotFound, map[string]string{
					"error": err.Error(),
				})
				return
			}
		} else {
			err = storage.DeleteRow(idCustomer, data.ProductID)
			if err != nil {
				ctx.JSON(http.StatusNotFound, map[string]string{
					"error": err.Error(),
				})
				return
			}
		}
		ShowBasket(ctx, idCustomer, storage, 201)
	})

	//sale
	router.POST("/api/:idCustomer/Sale", func(ctx *gin.Context) {
		id_str := ctx.Param("idCustomer")
		idCustomer, err := strconv.Atoi(id_str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Unsuitable ID number",
			})
			return
		}
		saleShow, totalPay, err := storage.Sale(idCustomer)
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, map[string]string{
			"Sale": "Successful",
		})
		jsonProduct, err := json.Marshal(*saleShow)
		_, err = ctx.Writer.Write(jsonProduct)
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, map[string]any{
			"Total Pay": totalPay,
		})

	})

	return router
}

// show cart
func ShowBasket(ctx *gin.Context, idCustomer int, storage Services.PStorage, status int) {
	basket2, totalPay, err := storage.ShowBasket(idCustomer)
	if err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(status, map[string]string{
		"Message": "Successful",
	})
	jsonProduct, err := json.Marshal(*basket2)
	_, err = ctx.Writer.Write(jsonProduct)
	if err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(status, map[string]any{
		"Total Pay": totalPay,
	})
}
