package main

import "C"
import (
	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"golang-angular/handlers"
	"gopkg.in/square/go-jose.v2"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var (
	audience string
	domain   string
)

func main() {

	setAuth0Variables()
	r := gin.Default()
	r.Use(CORSMiddleware())

	// This will ensure that Angular files are served correctly
	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui/" + path.Join(dir, file))
		}
	})

	// Handle non-mapped routes
	//r.NoRoute(func(c *gin.Context) {
	//	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	//})

	authorized := r.Group("/")
	authorized.Use(authRequired())

	authorized.GET("/dish", handlers.GetDishListHandler)
	authorized.POST("/dish", handlers.AddDishHandler)
	authorized.DELETE("/dish/:id", handlers.DeleteDishHandler)

	authorized.GET("/hero", handlers.GetHeroListHandler)
	authorized.POST("/hero", handlers.AddHeroHandler)
	authorized.DELETE("/hero/:id", handlers.DeleteHeroHandler)

	authorized.GET("/todo", handlers.GetTodoListHandler)
	authorized.POST("/todo", handlers.AddTodoHandler)
	authorized.DELETE("/todo/:id", handlers.DeleteTodoHandler)
	authorized.PUT("/todo", handlers.CompleteTodoHandler)
	r.OPTIONS("/todo", func(c *gin.Context) {
		c.Status(http.StatusOK)
		return
	})

	err := r.Run(":5005")
	if err != nil {
		panic(err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, OPTIONS, POST, PUT")
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			return
		}
	}
}

// ValidateRequest will verify that a token received from an http request is valid and signed by auth0
func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth0Domain = "https://" + domain + "/"
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: auth0Domain + ".well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{audience}, auth0Domain, jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		jwt, err := validator.ValidateRequest(c.Request)
		_ = jwt

		if err != nil {
			log.Println(err)
			terminateWithError(http.StatusUnauthorized, "token is not valid", c)
			return
		}

		c.Next()
	}
}

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}

func setAuth0Variables() {
	audience = os.Getenv("AUTH0_API_IDENTIFIER")
	domain = os.Getenv("AUTH0_DOMAIN")
	log.Println("audience: " + audience)
	log.Println("domain: " + domain)
}
