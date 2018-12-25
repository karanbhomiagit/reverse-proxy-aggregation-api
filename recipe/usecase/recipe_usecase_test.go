package usecase_test

import (
	"testing"

	"github.com/karanbhomiagit/reverse-proxy-aggregation-api/models"
	"github.com/karanbhomiagit/reverse-proxy-aggregation-api/recipe/usecase"
)

func TestSortRecipesByPrepTime(t *testing.T) {

	recipe1 := models.Recipe{
		ID:       "1",
		Name:     "Bacon",
		PrepTime: "PT25M",
	}

	recipe2 := models.Recipe{
		ID:       "2",
		Name:     "Omelette",
		PrepTime: "PT5M",
	}

	recipe3 := models.Recipe{
		ID:       "3",
		Name:     "Eggs",
		PrepTime: "PT2-3M",
	}

	recipe4 := models.Recipe{
		ID:       "4",
		Name:     "Random",
		PrepTime: "",
	}

	recipe5 := models.Recipe{
		ID:       "5",
		Name:     "Chicken",
		PrepTime: "PT120M",
	}

	recipes := []*models.Recipe{&recipe5, &recipe2, &recipe3, &recipe1, &recipe4}
	sortedRecipes, err := usecase.SortRecipesByPrepTime(recipes)
	if err != nil {
		t.Error("Expected sorted slice of struct pointers. Got err", err)
	}
	// Expected order is 3,2,1,5,4
	if sortedRecipes[0].ID != "3" || sortedRecipes[1].ID != "2" || sortedRecipes[2].ID != "1" || sortedRecipes[3].ID != "5" || sortedRecipes[4].ID != "4" {
		t.Error("Unable to sort based on prepTime.")
	}
}
