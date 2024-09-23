package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRouter() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	useMetrics(r)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	rg := r.Group("/api/v1")
	rg.GET("stores/:storeId/products/:id", getProduct)
	rg.DELETE("stores/:storeId/products/item/:id", deleteProductItem)

	return r
}

func getProduct(ctx *gin.Context) {
	strStoreID := ctx.Param("storeId")
	storeID, _ := strconv.Atoi(strStoreID)
	strID := ctx.Param("id")
	id, _ := strconv.Atoi(strID)
	res, err := GetProduct(uint(storeID), uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"product": res,
	})
}

func deleteProductItem(ctx *gin.Context) {
	strStoreID := ctx.Param("storeId")
	storeID, _ := strconv.Atoi(strStoreID)
	strID := ctx.Param("id")
	id, _ := strconv.Atoi(strID)
	err := DeleteProduct(uint(storeID), uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, "")
}

func useMetrics(e *gin.Engine) {
	metricsMiddleware := func(c *gin.Context) {
		start := time.Now()
		reqSz := computeApproximateRequestSize(c.Request)

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		resSz := float64(c.Writer.Size())

		url := c.Request.URL.Path

		PromMetrics.Rest.ReqCnt.WithLabelValues(status, c.Request.Method, c.Request.Host, url).Inc()
		PromMetrics.Rest.ReqDur.WithLabelValues(status, c.Request.Method, url).Observe(elapsed)
		PromMetrics.Rest.ReqSz.Observe(float64(reqSz))
		PromMetrics.Rest.ResSz.Observe(float64(resSz))
	}

	e.Use(metricsMiddleware)
}

// From https://github.com/DanielHeckrath/gin-prometheus/blob/master/gin_prometheus.go
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
