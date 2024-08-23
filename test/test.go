package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func main() {
	app := gin.New()
	app.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"hello": "world"})
	})

	http.ReadRequest()
	http.ReadResponse()

	r, _ := http.NewRequest("GET", "http://127.0.0.1:8080/test", nil)

	r.Header
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)

	app.ServeHTTP()

	_ = app.Run(":808")
}
