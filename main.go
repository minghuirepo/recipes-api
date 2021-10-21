// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
//	Schemes: http
//  Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//	Contact: Mohamed Labouardy <mohamed@labouardy.com> https://labouardy.com
//
//  Security:
//  - api_key:
//
//   SecurityDefinitions:
//   api_key:
//        type: apiKey
//        name: Authorization
//        in: header
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"study.recipes.api/handlers"
	"study.recipes.api/middlewares"
)

var (
	recipesHandler *handlers.RecipesHandler
)

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping(ctx)
	fmt.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

func main() {
	router := gin.Default()

	router.Use(middlewares.CORSMiddleware())

	router.GET("/recipes", recipesHandler.ListRecipesHandler)

	router.POST("/recipes", recipesHandler.NewRecipeHandler)

	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)

	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)

	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)

	router.GET("/recipes/:id", recipesHandler.GetRecipeHandler)

	router.Run()
}
