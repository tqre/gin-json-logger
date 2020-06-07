package main

import "github.com/gin-gonic/gin"

func main() {
	server := gin.New()

	// Default settings mostly
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#New
	server.RedirectTrailingSlash = false // not default
	server.RedirectFixedPath = false
	server.HandleMethodNotAllowed = false
	server.ForwardedByClientIP = true

	// https://golang.org/src/net/url/url.go
	server.UseRawPath = false
	server.UnescapePathValues = true

	// Use the middleware
	server.Use(Logger_JSON("server.log", true))

	// Skeleton request handler
	server.GET("/", index)

	// TODO: RunTLS
	server.Run("localhost:8001")
}

func index(context *gin.Context) {
	context.String(200, "Hello, world!")
}