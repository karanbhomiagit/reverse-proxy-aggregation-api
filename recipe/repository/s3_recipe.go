package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/karanbhomiagit/reverse-proxy-aggregation-api/models"
	"github.com/karanbhomiagit/reverse-proxy-aggregation-api/recipe"
)

type s3RecipeRepository struct {
}

func NewS3RecipeRepository() recipe.Repository {
	return &s3RecipeRepository{}
}

//FetchByID contacts S3 via http call to get the recipe for a particular id
func (*s3RecipeRepository) FetchByID(id string) (*models.Recipe, error) {
	if id == "" {
		return nil, errors.New("ID should not be null")
	}

	url := fmt.Sprintf("%s%s", getBaseURLForRecipes(), id)

	//Add a timeout of 1 sec to the call to S3 URL
	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error : ", err)
		return nil, err
	}

	//If the response status code is not 200, return appropriate error
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("Not found") //Move this string to ENV variables
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to fetch record")
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Error : ", err)
		return nil, err
	}
	r, err := recipeFromJSON(body)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func getBaseURLForRecipes() string {
	return "https://replaceme.com/"
}

func recipeFromJSON(data []byte) (*models.Recipe, error) {
	r := models.Recipe{}
	err := json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
