package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"plutus/internal/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	// Set DbRepository middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rr := r.WithContext(context.WithValue(r.Context(), GetRepo,
				func(r context.Context) (model.DbRepository, error) {
					return DbRepoMock{}, nil
				}))
			next.ServeHTTP(w, rr)
		})
	})

	SetRoutes(router)
	return router
}

func TestApi(t *testing.T) {
	r := setupRouter()

	tests := []struct {
		title       string
		url         string
		method      string
		body        any
		expectCode  int
		expectError *ResponseError
		expectModel map[string]any
	}{
		{
			title:       "Not Found",
			url:         "/unknown",
			method:      http.MethodGet,
			expectCode:  http.StatusNotFound,
			expectError: &ResponseError{Error: "Not Found"},
		},
		{
			title:  "Transaction write",
			url:    "/webhooks/transaction",
			method: http.MethodPost,
			body: model.Transaction{
				TransactionId:   "tqZi6QapS41zcEHy2",
				TransactionType: model.Sale,
				OrderId:         "c66oxMaisTwJQXjD",
				Amount:          10,
				Currency:        model.Eur,
				Description:     "Test transaction",
				AccountId:       "001",
			},
			expectCode:  http.StatusOK,
			expectModel: map[string]any{"transaction_id": "tqZi6QapS41zcEHy2"},
		},
		{
			title:  "Transaction write - bad data",
			url:    "/webhooks/transaction",
			method: http.MethodPost,
			body: model.Transaction{
				TransactionId:   "tqZi6QapS41zcEHy2",
				TransactionType: model.Sale,
				// OrderId:         "c66oxMaisTwJQXjD",
				Amount:      10,
				Currency:    model.Eur,
				Description: "Test transaction",
				AccountId:   "001",
			},
			expectCode:  http.StatusBadRequest,
			expectError: &ResponseError{Error: "Missing or invalid fields"},
		},
		{
			title:       "Missing account",
			url:         "/account/unknown",
			method:      http.MethodGet,
			expectCode:  http.StatusNotFound,
			expectError: &ResponseError{Error: "No transactions found"},
		},
		{
			title:       "Account balance",
			url:         "/account/001",
			method:      http.MethodGet,
			expectCode:  http.StatusOK,
			expectModel: map[string]any{"balance": "20.00"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			var body io.Reader
			if tt.body != nil {
				bb := new(bytes.Buffer)
				assert.NoError(t, json.NewEncoder(bb).Encode(tt.body))
				body = bb
			}
			req, _ := http.NewRequest(tt.method, tt.url, body)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectCode, rr.Code)

			if tt.expectError != nil {
				var res = ResponseError{}
				b, e := io.ReadAll(rr.Body)
				assert.NoError(t, e)
				assert.NoError(t, json.Unmarshal(b, &res))

				assert.Contains(t, res.Error, tt.expectError.Error)
			}

			if tt.expectModel != nil {
				var res map[string]any
				b, e := io.ReadAll(rr.Body)
				assert.NoError(t, e)
				assert.NoError(t, json.Unmarshal(b, &res))

				assert.Equal(t, tt.expectModel, res)
			}
		})
	}
}

type DbRepoMock struct {
}

func (m DbRepoMock) Find(transactionId string) (*model.Transaction, error) {
	if transactionId == "tqZi6QapS41zcEHy" {
		return &model.Transaction{
			TransactionId:   "tqZi6QapS41zcEHy",
			TransactionType: model.Sale,
			OrderId:         "c66oxMaisTwJQXjD",
			Amount:          10,
			Currency:        model.Eur,
			Description:     "Test transaction",
			AccountId:       "001",
		}, nil
	}
	return nil, nil
}
func (m DbRepoMock) FetchAll() ([]*model.Transaction, error) {
	return nil, nil
}
func (m DbRepoMock) FetchAllByAccount(accountId string) ([]*model.Transaction, error) {
	if accountId == "001" {
		return []*model.Transaction{
			{TransactionId: "tqZi6QapS41zcEHy1",
				TransactionType: model.Sale,
				OrderId:         "c66oxMaisTwJQXjD",
				Amount:          10,
				Currency:        model.Eur,
				Description:     "Test transaction",
				AccountId:       "001"},
			{TransactionId: "tqZi6QapS41zcEHy2",
				TransactionType: model.Sale,
				OrderId:         "c66oxMaisTwJQXjD",
				Amount:          10,
				Currency:        model.Eur,
				Description:     "Test transaction",
				AccountId:       "001"},
		}, nil
	}
	return nil, nil
}

func (m DbRepoMock) Insert(ent *model.Transaction) error {
	ent.Id = 1
	return nil
}
func (m DbRepoMock) Empty() error {
	return nil
}
func (m DbRepoMock) Init() error {
	return nil
}
func (m DbRepoMock) Close() {

}
