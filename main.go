package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/Qalifah/shop/routes"
	"github.com/Qalifah/shop/docs/docs.go"
)

// @title SHOP API
// @version 1.0
// @description This is a sample E-Commerce server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host shop.swagger.io
// @BasePath /v1
func main() {
	port := ":8080"
	http.Handle("/", routes.Routers())
	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}