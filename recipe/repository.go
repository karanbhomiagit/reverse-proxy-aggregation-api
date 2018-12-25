package recipe

import "github.com/karanbhomiagit/reverse-proxy-aggregation-api/models"

// Repository represents the recipe's storage/retrieval as an interface
type Repository interface {
	FetchByID(id string) (*models.Recipe, error)
}
