package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"stan-project/data"
	"testing"
)

func TestNewRisksDB(t *testing.T) {
	t.Run("successfully initialize risks DB", func(t *testing.T) {
		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initializing DB for test: %s", err)
		}
		defer pDB.client.Close(ctx)

		rDB := NewRisksDB(pDB)

		assert.Equal(t, &risksDB{db: pDB}, rDB)
	})
}

func TestRisksDB_Add(t *testing.T) {
	t.Run("successfully add a new risk", func(t *testing.T) {
		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initializing DB for test: %s", err)
		}
		defer pDB.client.Close(ctx)

		rDB := NewRisksDB(pDB)

		riskID := uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909")
		err = rDB.Add(ctx, data.Risk{
			ID:          riskID,
			Title:       "threat 1",
			Description: "DDOS threat",
			State:       "open",
		})

		assert.Nil(t, err)

		defer func() {
			//clean up
			deleteEr := rDB.DeleteByID(ctx, riskID)
			if deleteEr != nil {
				t.Logf("error cleaning up test data: %s", err)
			}
		}()
	})
}

func TestRisksDB_GetByID(t *testing.T) {
	t.Run("successfully get a risk by ID", func(t *testing.T) {

		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initializing DB for test: %s", err)
		}
		defer pDB.client.Close(ctx)

		rDB := NewRisksDB(pDB)

		//Add test data
		riskID := uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909")
		addErr := rDB.Add(ctx, data.Risk{
			ID:          riskID,
			Title:       "threat 1",
			Description: "DDOS threat",
			State:       "open",
		})

		if addErr != nil {
			t.Fatalf("error adding test data: %s", err)
		}

		defer func() {
			//clean up
			deleteEr := rDB.DeleteByID(ctx, riskID)
			if deleteEr != nil {
				t.Logf("error cleaning up test data: %s", err)
			}
		}()

		expected := data.Risk{
			ID:          riskID,
			Title:       "threat 1",
			Description: "DDOS threat",
			State:       "open",
		}

		actual, err := rDB.GetByID(ctx, riskID)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestRisksDB_GetAll(t *testing.T) {
	t.Run("successfully get all risks", func(t *testing.T) {

		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initializing DB for test: %s", err)
		}
		defer pDB.client.Close(ctx)

		rDB := NewRisksDB(pDB)

		//Add test data
		riskID := uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909")
		addErr := rDB.Add(ctx, data.Risk{
			ID:          riskID,
			Title:       "threat 1",
			Description: "DDOS threat",
			State:       "open",
		})

		if addErr != nil {
			t.Fatalf("error adding test data: %s", err)
		}

		defer func() {
			//clean up
			deleteEr := rDB.DeleteByID(ctx, riskID)
			if deleteEr != nil {
				t.Logf("error cleaning up test data: %s", err)
			}
		}()

		expected := data.PaginatedResponse{
			TotalCount: 1,
			Risks: []data.Risk{
				{
					ID:          riskID,
					Title:       "threat 1",
					Description: "DDOS threat",
					State:       "open",
				},
			},
		}

		actual, err := rDB.GetAll(ctx, data.Options{Offset: 0, Limit: 3, SortBy: "title", SortOrder: "asc"})

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}
