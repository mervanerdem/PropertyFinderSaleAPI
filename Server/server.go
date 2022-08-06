package Server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mervanerdem/PropertyFinderSaleAPI/Services"
	"log"
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
		jsonProduct, err := json.Marshal(*products)
		_, err = ctx.Writer.Write(jsonProduct)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

	})

	//get list of basket
	router.GET("/api/:idCustomer/basket", func(ctx *gin.Context) {
		id_str := ctx.Param("idCustomer")
		idCustomer, err := strconv.Atoi(id_str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Unsuitable ID number",
			})
			return
		}

		ShowBasket(ctx, idCustomer, storage)
	})

	//Add to cart
	router.POST("/api/:idCustomer/basket/add", func(ctx *gin.Context) {
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
			return
		}

		productPrice, err := storage.IsHaveProductID(data.ProductID)
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
		haveProductNum, pNum, err := storage.IsHaveProductNumber(idCustomer, data.ProductID)

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
			productTotalPrice := productPrice * data.ProductNumber
			err = storage.AddBasket(idCustomer, data.ProductID, data.ProductNumber, productTotalPrice)
			if err != nil {
				ctx.JSON(http.StatusNotFound, map[string]string{
					"error": err.Error(),
				})
				return
			}
		}
		ctx.JSON(http.StatusOK, map[string]string{
			"Message": "Successful",
		})
		ShowBasket(ctx, idCustomer, storage)
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

		_, err = storage.IsHaveProductID(data.ProductID)
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

		haveProductNum, pNum, err := storage.IsHaveProductNumber(idCustomer, data.ProductID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}
		if haveProductNum && (data.ProductNumber-pNum) > 0 {
			log.Printf("birinci")
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

		ctx.JSON(http.StatusOK, map[string]string{
			"Message": "Successful",
		})
		ShowBasket(ctx, idCustomer, storage)
	})

	//sale //düzenlenecek
	router.POST("/api/:idCustomer/Sale", func(ctx *gin.Context) {
		id_str := ctx.Param("idCustomer")
		idCustomer, err := strconv.Atoi(id_str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Unsuitable ID number",
			})
			return
		}
		err = storage.Sale(idCustomer)
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, map[string]string{
			"Sale": "Successful",
		})

	})

	return router
}

// show cart
func ShowBasket(ctx *gin.Context, idCustomer int, storage Services.PStorage) {
	basket2, totalPay, err := storage.ShowBasket(idCustomer)
	if err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}
	jsonProduct, err := json.Marshal(*basket2)
	_, err = ctx.Writer.Write(jsonProduct)
	if err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{
		"Total Pay": totalPay,
	})
}
