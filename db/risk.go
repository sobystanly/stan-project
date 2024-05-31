package db

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/google/uuid"
	"stan-project/data"
)

type risksDB struct {
	db *db
}

func NewRisksDB(db *db) *risksDB {
	return &risksDB{db: db}
}

//go:embed sql/insert_risk.sql
var insertRisk string

func (rdb *risksDB) Add(ctx context.Context, risk data.Risk) error {
	var err error
	_, err = rdb.db.client.Exec(ctx, insertRisk, risk.ID, risk.Title, risk.Description, risk.State)
	return err
}

//go:embed sql/get_risk_by_id.sql
var getRiskByID string

func (rdb *risksDB) GetByID(ctx context.Context, ID uuid.UUID) (data.Risk, error) {
	rows, err := rdb.db.client.Query(ctx, getRiskByID, ID)
	if err != nil {
		return data.Risk{}, err
	}
	defer rows.Close()

	var risk data.Risk

	for rows.Next() {
		err = rows.Scan(&risk.ID, &risk.Title, &risk.Description, &risk.State)
		if err != nil {
			return data.Risk{}, err
		}
	}

	return risk, nil
}

//go:embed sql/get_all_risks.sql
var getAllRisks string

//go:embed sql/count_all_risks.sql
var countAllRisks string

func (rdb *risksDB) GetAll(ctx context.Context, options data.Options) (data.PaginatedResponse, error) {

	var count int
	err := rdb.db.client.QueryRow(ctx, countAllRisks).Scan(&count)
	if err != nil {
		return data.PaginatedResponse{}, err
	}

	formattedQuery := fmt.Sprintf(getAllRisks, options.SortBy, options.SortOrder)

	rows, err := rdb.db.client.Query(ctx, formattedQuery, options.Limit, options.Offset)
	if err != nil {
		return data.PaginatedResponse{}, err
	}
	defer rows.Close()

	var risks []data.Risk

	for rows.Next() {
		var risk data.Risk
		err = rows.Scan(&risk.ID, &risk.Title, &risk.Description, &risk.State)
		if err != nil {
			return data.PaginatedResponse{}, err
		}
		risks = append(risks, risk)
	}

	return data.PaginatedResponse{TotalCount: count, Risks: risks}, nil
}

//go:embed sql/delete_risk_by_id.sql
var deleteRiskByID string

func (rdb *risksDB) DeleteByID(ctx context.Context, ID uuid.UUID) error {
	var err error
	_, err = rdb.db.client.Exec(ctx, deleteRiskByID, ID)
	return err
}
