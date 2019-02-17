package sendyoulater

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type janitor struct {
	*redis.Client
}

// Janitor interface
type Janitor interface {
	CleanUp(threshold int64, pattern string) (int64, error)
}

// NewJanitor returns a janitor that performs a cleanup
func NewJanitor(r *redis.Client) Janitor {
	return janitor{r}
}

// CleanUp deletes all the keys in redis matching a pattern that have not been read for a time greater than the provided threshold in seconds
func (j janitor) CleanUp(threshold int64, pattern string) (int64, error) {
	maxtime := time.Duration(threshold) * time.Second
	toBeDeleted := []string{}
	cursor := uint64(0)
	for {
		keys, c, err := j.Scan(cursor, pattern, 10).Result()
		if err != nil {
			return 0, err
		}
		for _, k := range keys {
			t, err := j.ObjectIdleTime(k).Result()
			if err != nil {
				fmt.Printf("Error retrieving idletime for key %v", k)
				continue
			}
			if t > maxtime {
				toBeDeleted = append(toBeDeleted, k)
			}
		}
		if c == 0 {
			break
		}
		cursor = c
	}
	fmt.Printf("Deleting: %+v", toBeDeleted)
	if len(toBeDeleted) == 0 {
		return 0, nil
	}
	return j.Del(toBeDeleted...).Result()
}
