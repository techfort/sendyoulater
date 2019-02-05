package sendyoulater

import "github.com/go-redis/redis"

type store struct {
	*redis.Client
}

// Store interface exposes repo creation functions
type Store interface {
	NewUserRepo() UserRepo
	NewPlanRepo() PlanRepo
	NewEmailActionRepo() EmailActionRepo
	NewEmailUseCase(u UserRepo, p PlanRepo, e EmailActionRepo) EmailUseCase
}

// NewStore returns a Store interface to obtain repos
func NewStore(r *redis.Client) Store {
	return store{r}
}
