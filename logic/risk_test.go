package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"stan-project/data"
	"testing"
)

func TestNewRiskLogic(t *testing.T) {
	t.Run("successfully initialize risk logic", func(t *testing.T) {
		mockDB := &mockRiskDB{}
		actual := NewRiskLogic(mockDB)
		assert.Equal(t, &riskLogic{riskDB: mockDB}, actual)
	})
}

func TestRiskLogic_Add(t *testing.T) {
	t.Run("successfully add a new risk", func(t *testing.T) {
		rl := NewRiskLogic(mockRiskDB{})
		risk := data.Risk{
			Title:       "threat 1",
			Description: "DDOS threat",
			State:       "open",
		}
		actual, err := rl.Add(context.Background(), risk)
		assert.Nil(t, err)
		risk.ID = actual.ID
		assert.Equal(t, risk, actual)
	})
	t.Run("failed to add a new risk, invalid state", func(t *testing.T) {
		rl := NewRiskLogic(mockRiskDB{})
		risk := data.Risk{
			Title:       "threat 1",
			Description: "DDOS threat",
			State:       "converted",
		}
		_, err := rl.Add(context.Background(), risk)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("risk state is invalid: %v", risk), err)
	})
	t.Run("failed to add a new risk, some error from db", func(t *testing.T) {
		rl := NewRiskLogic(mockRiskDB{err: errors.New("some error from DB")})
		risk := data.Risk{
			Title:       "threat 1",
			Description: "DDOS threat",
			State:       "closed",
		}
		_, err := rl.Add(context.Background(), risk)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("some error from DB"), err)
	})
}

func TestRiskLogic_GetByID(t *testing.T) {
	t.Run("successfully get a risk by ID", func(t *testing.T) {
		rl := NewRiskLogic(mockRiskDB{risk: data.Risk{
			ID:          uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
			Title:       "threat 1",
			Description: "DDOS threat",
			State:       "open",
		}})

		expected := data.Risk{
			ID:          uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
			Title:       "threat 1",
			Description: "DDOS threat",
			State:       "open",
		}

		actual, err := rl.GetByID(context.Background(), uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"))
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestRiskLogic_GetAll(t *testing.T) {
	t.Run("successfully get all risks", func(t *testing.T) {
		rl := NewRiskLogic(mockRiskDB{paginatedRisk: data.PaginatedResponse{
			TotalCount: 1,
			Risks: []data.Risk{
				{
					ID:          uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
					Title:       "threat 1",
					Description: "DDOS threat",
					State:       "open",
				},
			},
		}})

		expected := data.PaginatedResponse{
			TotalCount: 1,
			Risks: []data.Risk{
				{
					ID:          uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
					Title:       "threat 1",
					Description: "DDOS threat",
					State:       "open",
				},
			},
		}

		actual, err := rl.GetAll(context.Background(), data.Options{Offset: 0, Limit: 5, SortBy: "title", SortOrder: "desc"})
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("successfully get all risks, offset and limit is less than 0 fallback to defaults", func(t *testing.T) {
		rl := NewRiskLogic(mockRiskDB{paginatedRisk: data.PaginatedResponse{
			TotalCount: 1,
			Risks: []data.Risk{
				{
					ID:          uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
					Title:       "threat 1",
					Description: "DDOS threat",
					State:       "open",
				},
			},
		}})

		expected := data.PaginatedResponse{
			TotalCount: 1,
			Risks: []data.Risk{
				{
					ID:          uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
					Title:       "threat 1",
					Description: "DDOS threat",
					State:       "open",
				},
			},
		}

		actual, err := rl.GetAll(context.Background(), data.Options{Offset: -1, Limit: -2, SortBy: "title", SortOrder: "desc"})
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

type mockRiskDB struct {
	risk          data.Risk
	err           error
	paginatedRisk data.PaginatedResponse
}

func (m mockRiskDB) Add(ctx context.Context, risk data.Risk) error {
	return m.err
}

func (m mockRiskDB) GetByID(ctx context.Context, ID uuid.UUID) (data.Risk, error) {
	return m.risk, m.err
}

func (m mockRiskDB) GetAll(ctx context.Context, options data.Options) (data.PaginatedResponse, error) {
	return m.paginatedRisk, m.err
}
