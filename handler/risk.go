package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"stan-project/data"
	"strconv"
)

const (
	offset    = "offset"
	limit     = "limit"
	sortBy    = "sortBy"
	sortOrder = "sortOrder"
	title     = "title"
	asc       = "asc"
	desc      = "desc"
)

type (
	riskLogic interface {
		Add(ctx context.Context, risk data.Risk) (data.Risk, error)
		GetByID(ctx context.Context, ID uuid.UUID) (data.Risk, error)
		GetAll(ctx context.Context, options data.Options) (data.PaginatedResponse, error)
	}

	riskHandler struct {
		riskLogic riskLogic
	}
)

func NewRiskHandler(riskLogic riskLogic) *riskHandler {
	return &riskHandler{riskLogic: riskLogic}
}

func (rh *riskHandler) Add(w http.ResponseWriter, r *http.Request) {

	requestID := r.Context().Value("requestID").(string)
	log.Printf("received a request to create a new risk with requestID: %s, req: %v", requestID, r)

	ctx := r.Context()
	risk, err := decodeReq(r)
	if err != nil {
		log.Printf("error unmarshallling risk request: %s", err)
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "error decoding risk request"})
		return
	}

	risk, err = rh.riskLogic.Add(ctx, risk)
	if err != nil {
		log.Printf("error adding risk: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "error processing the risk add request"})
		return
	}

	log.Printf("successfully added a new risk with ID: %s", risk.ID)
	respondWithJSON(w, http.StatusCreated, risk)
}

func (rh *riskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("requestID").(string)
	log.Printf("received a request to fetch a risk with requestID: %s, req: %v", requestID, r)
	ctx := r.Context()

	ID := mux.Vars(r)["id"]

	log.Printf("received riskID: %s", ID)

	riskID, err := uuid.Parse(ID)
	if err != nil {
		log.Printf("invalid riskID: %s", ID)
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Errorf("invalid riskID, expected a UUID but received: %s", ID).Error()})
		return
	}

	risk, err := rh.riskLogic.GetByID(ctx, riskID)
	if err != nil {
		log.Printf("error fetching risk with ID: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Errorf("error fetching risk with ID: %s", riskID).Error()})
	}

	log.Printf("successfully fetched risk with ID: %s, risk: %v", riskID, risk)
	respondWithJSON(w, http.StatusOK, risk)
}

func (rh *riskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("requestID").(string)
	log.Printf("received a request to fetch all risks with requestID: %s, req: %v", requestID, r)
	ctx := r.Context()

	var options data.Options
	options.SortBy = title
	options.SortOrder = asc

	offsetOpt := getQueryParam(offset, r)
	offsetVal, err := strconv.Atoi(offsetOpt)
	if err != nil {
		log.Printf("invalid offset in request options so setting it to 0, err: %s", err)
		offsetVal = 0
	}
	limitOpt := getQueryParam(limit, r)
	limitVal, err := strconv.Atoi(limitOpt)
	if err != nil {
		log.Printf("invalid limit option, setting to default value 5, err: %s", err)
		limitVal = 10
	}
	sortByVal := getQueryParam(sortBy, r)
	sortOrderVal := getQueryParam(sortOrder, r)

	options.Offset = offsetVal
	options.Limit = limitVal
	if sortByVal != "" {
		options.SortBy = sortByVal
	}
	if sortOrderVal == desc {
		options.SortOrder = desc
	}

	log.Printf("fetching risks with options: %v", options)

	risks, err := rh.riskLogic.GetAll(ctx, options)
	if err != nil {
		log.Printf("error fetchiing all risks, %s", err)
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "error fetching risks"})
		return
	}

	log.Printf("successfully fetched all risks: %v", risks)
	respondWithJSON(w, http.StatusOK, risks)
}

func decodeReq(req *http.Request) (data.Risk, error) {
	var risk data.Risk
	err := json.NewDecoder(req.Body).Decode(&risk)
	if err != nil {
		return risk, err
	}
	return risk, err
}
