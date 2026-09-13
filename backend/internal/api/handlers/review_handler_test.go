package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/pkg/config"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

type stubReviewService struct {
	reviewErr error
}

func (s *stubReviewService) GetByOriginalID(ctx context.Context, dataSource, originalID string) (interface{}, error) {
	return nil, nil
}

func (s *stubReviewService) GetAll(ctx context.Context, params domain.QueryParams) (*domain.PaginatedResult, error) {
	return nil, nil
}

func (s *stubReviewService) ReviewItem(ctx context.Context, dataSource, originalID string, update domain.ReviewStatusUpdate) error {
	return s.reviewErr
}

func (s *stubReviewService) GetFilterOptions(ctx context.Context, dataSource string) (map[string][]string, error) {
	return nil, nil
}

func (s *stubReviewService) GetDataSources(ctx context.Context) ([]string, error) {
	return nil, nil
}

func setupReviewRouter(svc *stubReviewService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewReviewHandler(svc, logger.New(config.LoggerConfig{Level: "error", Format: "text"}))
	r := gin.New()
	r.POST("/api/review/:dataSource/items/:id/review", h.ReviewItem)
	return r
}

func TestReviewItem_InternalErrorOmitsDetails(t *testing.T) {
	sqlLeak := errors.New(`pq: relation "sat_oneprep_items" does not exist`)
	r := setupReviewRouter(&stubReviewService{reviewErr: sqlLeak})

	body, _ := json.Marshal(domain.ReviewStatusUpdate{Status: domain.ReviewStatusApproved})
	req := httptest.NewRequest(http.MethodPost, "/api/review/sat_oneprep/items/item-1/review", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp["error"] != "Failed to update review status" {
		t.Fatalf("error = %v, want generic message", resp["error"])
	}
	if _, ok := resp["details"]; ok {
		t.Fatalf("response must not include details field, got %v", resp["details"])
	}
	if bytes.Contains(w.Body.Bytes(), []byte("sat_oneprep_items")) || bytes.Contains(w.Body.Bytes(), []byte("pq:")) {
		t.Fatalf("response leaked internal error text: %s", w.Body.String())
	}
}

func TestReviewItem_NotFoundReturns404(t *testing.T) {
	r := setupReviewRouter(&stubReviewService{
		reviewErr: fmt.Errorf("no item found with original_id %s", "missing-id"),
	})

	body, _ := json.Marshal(domain.ReviewStatusUpdate{Status: domain.ReviewStatusApproved})
	req := httptest.NewRequest(http.MethodPost, "/api/review/sat_oneprep/items/missing-id/review", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp["error"] != "Item with ID missing-id not found" {
		t.Fatalf("error = %v, want stable not-found message", resp["error"])
	}
	if _, ok := resp["details"]; ok {
		t.Fatalf("response must not include details field, got %v", resp["details"])
	}
}
