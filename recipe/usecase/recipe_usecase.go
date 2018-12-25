package usecase

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/karanbhomiagit/reverse-proxy-aggregation-api/models"
	"github.com/karanbhomiagit/reverse-proxy-aggregation-api/recipe"
)

type recipeUsecase struct {
	recipeRepo recipe.Repository
}

func NewRecipeUsecase(r recipe.Repository) recipe.Usecase {
	return &recipeUsecase{
		recipeRepo: r,
	}
}

//FetchByIds returns a list of recipes for the given ids
func (r *recipeUsecase) FetchByIds(ids []string) ([]*models.Recipe, error) {
	return r.fetchByIdsHelper(ids)
}

//FetchByIds returns a list of recipes for the given ids sorted by PrepTime
func (r *recipeUsecase) FetchByIdsSortedByPrepTime(ids []string) ([]*models.Recipe, error) {
	res, err := r.fetchByIdsHelper(ids)
	if err != nil {
		return nil, err
	}
	return SortRecipesByPrepTime(res)
}

//fetchByIdsHelper contacts the repository layer to fetch the recipes for provided ids
func (r *recipeUsecase) fetchByIdsHelper(ids []string) ([]*models.Recipe, error) {
	//Make separate channels to receive responses and errors
	respc, errc := make(chan *models.Recipe), make(chan error)
	//For all ids passed, start a goroutine which contacts the repository layer and fetches the recipe
	for _, id := range ids {
		go func(idToFetch string) {
			body, err := r.recipeRepo.FetchByID(idToFetch)
			if err != nil {
				errc <- err
				return
			}
			respc <- body
		}(id)
	}

	recipes := []*models.Recipe{}

	//Receive the response/error in the channels
	for range ids {
		select {
		case resp := <-respc:
			recipes = append(recipes, resp)
		case err := <-errc:
			fmt.Println("Response Error : ", err)
			//If response is 404, ignore. Else, return error to delivery layer.
			if err.Error() != "Not found" {
				return nil, err
			}
		}
	}

	return recipes, nil
}

//SortRecipesByPrepTime receives a slice of Recipe and sorts based on PrepTime
func SortRecipesByPrepTime(recipes []*models.Recipe) ([]*models.Recipe, error) {
	sort.Slice(recipes, func(i, j int) bool {
		//If the PrepTime is zero value for any recipe, push it to the last of slice
		if recipes[i].PrepTime == "" {
			return false
		}
		if recipes[j].PrepTime == "" {
			return true
		}

		si := recipes[i].PrepTime
		sj := recipes[j].PrepTime
		//If the PrepTime is in the format "PT3-5M", take the lower value out, i.e PT3M
		if strings.Contains(si, "-") {
			temp := strings.Split(si, "-")
			si = temp[0] + "M"
		}
		if strings.Contains(sj, "-") {
			temp := strings.Split(sj, "-")
			sj = temp[0] + "M"
		}

		//Remove prefixes and suffixes
		si = strings.TrimPrefix(si, "PT")
		si = strings.TrimSuffix(si, "M")

		sj = strings.TrimPrefix(sj, "PT")
		sj = strings.TrimSuffix(sj, "M")

		//Convert to int and compare the values
		vali, _ := strconv.Atoi(si)
		valj, _ := strconv.Atoi(sj)
		return vali < valj

	})
	return recipes, nil
}
