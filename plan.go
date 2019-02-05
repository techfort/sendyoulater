package sendyoulater

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// Plan is a representation of usage plan
type Plan struct {
	Name      string
	MaxEmails int64
	MaxSMS    int64
}

var (
	// Plans is the supported plans
	Plans = map[string]Plan{
		"basic":      Plan{"basic", 100, 100},
		"enterprise": Plan{"enterprise", 100, 100},
	}
)

type planRepo struct {
	store
}

func (s store) NewPlanRepo() PlanRepo {
	return planRepo{s}
}

// PlanRepo interface that wwraps plan related functions of Store
type PlanRepo interface {
	ByName(plan string) (Plan, error)
	SavePlan(name string, maxemails, maxsms int64) (Plan, error)
}

// PlanByName retrieves a plan by its name
func (p planRepo) ByName(name string) (Plan, error) {
	var (
		ret       map[string]string
		plan      Plan
		err       error
		maxemails int64
		maxsms    int64
	)
	ret, err = p.HGetAll(name).Result()
	if maxemails, err = strconv.ParseInt(ret["MaxEmails"], 10, 64); err != nil {
		return plan, errors.Wrap(err, "failed to parse maxemails value")
	}
	if maxsms, err = strconv.ParseInt(ret["MaxSMS"], 10, 64); err != nil {
		return plan, errors.Wrap(err, "failed to parse MaxSMS value")
	}
	plan.Name = ret["Name"]
	plan.MaxEmails = maxemails
	plan.MaxSMS = maxsms
	return plan, err
}

// SavePlan saves a plan to cache
func (p planRepo) SavePlan(name string, maxemails, maxsms int64) (Plan, error) {
	plan := Plan{name, maxemails, maxsms}
	planKey := KeyPlan(name)
	pipe := p.TxPipeline()
	pipe.HSet(planKey, "Name", plan.Name)
	pipe.HSet(planKey, "MaxEmails", maxemails)
	pipe.HSet(planKey, "MaxSMS", plan.MaxSMS)
	_, err := pipe.Exec()
	if err != nil {
		return plan, errors.Wrap(err, fmt.Sprintf("failed to save plan: %+v", plan))
	}
	return plan, err
}
