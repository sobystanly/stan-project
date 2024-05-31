package logic

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"stan-project/data"
)

type (
	riskDB interface {
		Add(ctx context.Context, risk data.Risk) error
		GetByID(ctx context.Context, ID uuid.UUID) (data.Risk, error)
		GetAll(ctx context.Context, options data.Options) (data.PaginatedResponse, error)
	}
	riskLogic struct {
		riskDB riskDB
	}
)

func NewRiskLogic(riskDB riskDB) *riskLogic {
	return &riskLogic{riskDB: riskDB}
}

func (r *riskLogic) Add(ctx context.Context, risk data.Risk) (data.Risk, error) {
	if !risk.State.IsValid() {
		log.Printf("given risk: %v is invalid", risk)
		return data.Risk{}, fmt.Errorf("risk is invalid: %v", risk)
	}

	risk.ID = uuid.New()

	err := r.riskDB.Add(ctx, risk)
	if err != nil {
		log.Printf("error adding new risk: %s", err)
		return data.Risk{}, err
	}

	return risk, nil
}

func (r *riskLogic) GetByID(ctx context.Context, ID uuid.UUID) (data.Risk, error) {
	return r.riskDB.GetByID(ctx, ID)
}

func (r *riskLogic) GetAll(ctx context.Context, options data.Options) (data.PaginatedResponse, error) {
	if options.Offset < 0 {
		options.Offset = 0
	}
	return r.riskDB.GetAll(ctx, options)
}
