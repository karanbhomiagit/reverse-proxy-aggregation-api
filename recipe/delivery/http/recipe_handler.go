package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/karanbhomiagit/reverse-proxy-aggregation-api/recipe"
)

type HttpRecipeHandler struct {
	RUsecase recipe.Usecase
}

func NewRecipeHttpHandler(ru recipe.Usecase) {
	handler := &HttpRecipeHandler{
		RUsecase: ru,
	}

	http.HandleFunc("/recipes", handler.RecipesHandler)
}

//RecipesHandler is the entrypoint for any requests received for the path "/recipes"
func (h *HttpRecipeHandler) RecipesHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query()["ids"]
	//Check if "ids" are passed as query parameters
	if idParam != nil {
		ids := strings.Split(idParam[0], ",")
		h.getRecipesByIds(w, r, ids)
	} else {
		h.getAllRecipes(w, r)
	}
}

//getRecipesByIds fetches recipes of given ids by contacting the usecase layer.
func (h *HttpRecipeHandler) getRecipesByIds(w http.ResponseWriter, r *http.Request, idsParam []string) {
	method := r.Method
	//Check the http request method
	switch method {
	case http.MethodGet:
		fmt.Println("Request GET /recipes?ids=", idsParam)
		//Make call to usecase layer to fetch the recipes by ids (sorted by PrepTime)
		res, err := h.RUsecase.FetchByIdsSortedByPrepTime(idsParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}
		//Encode the response to JSON
		b, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		//Return 405 in case of other methods
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Unsupported Request Method")
	}
}

func (h *HttpRecipeHandler) getAllRecipes(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	//Check the http request method
	switch method {
	case http.MethodGet:
		fmt.Println("Request GET /recipes")
		//Check if $skip and $top params were passed
		skipParam := r.URL.Query()["$skip"]
		topParam := r.URL.Query()["$top"]
		//Use default values if not passed
		skipParamVal := "0"
		topParamVal := "10"
		if skipParam != nil {
			skipParamVal = skipParam[0]
		}
		if topParam != nil {
			topParamVal = topParam[0]
		}
		//Convert values of skip and top to integer
		skip, _ := strconv.Atoi(skipParamVal) //receive err and return 422
		top, _ := strconv.Atoi(topParamVal)   //receive err and return 422
		//Call helper method to get the recipes in specified range
		h.getRecipesInRange(skip, top, w)
	default:
		//Return 405 in case of other methods
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Unsupported Request Method")
	}
}

//getRecipesInRange is a helper method to get the recipes in specified range
func (h *HttpRecipeHandler) getRecipesInRange(skip int, top int, w http.ResponseWriter) {
	//Prepare the ids to fetch using values of skip and top
	ids := []string{}
	for i := skip + 1; i <= skip+top; i++ {
		ids = append(ids, strconv.Itoa(i))
	}
	//Make call to usecase layer to fetch the recipes by ids
	res, err := h.RUsecase.FetchByIds(ids)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}
	//Marshal the json
	b, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
