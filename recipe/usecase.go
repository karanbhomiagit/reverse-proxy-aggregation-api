package recipe

import (
	"github.com/karanbhomiagit/reverse-proxy-aggregation-api/models"
)

// Usecase represents the recipe's business logic as an interface
type Usecase interface {
	FetchByIds(ids []string) ([]*models.Recipe, error)
	FetchByIdsSortedByPrepTime(ids []string) ([]*models.Recipe, error)
}
