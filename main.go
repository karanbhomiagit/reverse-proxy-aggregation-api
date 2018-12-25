package main

import (
	"log"
	"net/http"
	"os"

	httpDeliver "github.com/karanbhomiagit/reverse-proxy-aggregation-api/recipe/delivery/http"
	recipeRepo "github.com/karanbhomiagit/reverse-proxy-aggregation-api/recipe/repository"
	recipeUsecase "github.com/karanbhomiagit/reverse-proxy-aggregation-api/recipe/usecase"
)

func main() {
	//Initializing the repository
	rr := recipeRepo.NewS3RecipeRepository()
	//Initializing the usecase
	ru := recipeUsecase.NewRecipeUsecase(rr)

	//Initializing the delivery
	httpDeliver.NewRecipeHttpHandler(ru)
	log.Fatal(http.ListenAndServe(port(), nil))
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}
