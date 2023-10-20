package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"plutus/internal/model"

	"github.com/gorilla/mux"
)

const (
	GetRepo = "getRepo"
)

func notFound(w http.ResponseWriter, _ *http.Request) {
	responseError(w, "Not Found", http.StatusNotFound)
}
func writeTransaction(w http.ResponseWriter, r *http.Request) {
	p := model.Transaction{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate data
	if !p.Validate() {
		responseError(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	// Get DBRepository
	repo := r.Context().Value(GetRepo)
	repoInst, err := repo.(func(f context.Context) (model.DbRepository, error))(r.Context())
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer repoInst.Close()

	if err = repoInst.Insert(&p); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJson(w, map[string]string{"transaction_id": p.TransactionId}, http.StatusOK)
}

func accountBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		responseError(w, "Missing transaction_id", http.StatusBadRequest)
		return
	}

	// Get DBRepository
	repo := r.Context().Value(GetRepo)
	repoInst, err := repo.(func(f context.Context) (model.DbRepository, error))(r.Context())
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer repoInst.Close()

	trs, err := repoInst.FetchAllByAccount(id)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(trs) == 0 {
		responseError(w, "No transactions found", http.StatusNotFound)
		return
	}

	balance := model.NewBalance(trs)
	responseJson(w, map[string]string{"balance": fmt.Sprintf("%.2f", balance.GetAmount())}, http.StatusOK)
}

func fetchAll(w http.ResponseWriter, r *http.Request) {
	// Get DBRepository
	repo := r.Context().Value(GetRepo)
	repoInst, err := repo.(func(f context.Context) (model.DbRepository, error))(r.Context())
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer repoInst.Close()

	pr, err := repoInst.FetchAll()
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if pr == nil {
		responseError(w, "Transaction not found", http.StatusNotFound)
		return
	}

	responseJson(w, pr, http.StatusOK)
}

func SetRoutes(router *mux.Router) {
	// Setup api
	router.HandleFunc("/webhooks/transaction", writeTransaction).Methods(http.MethodPost)
	router.HandleFunc("/account/{id:[a-zA-Z0-9]+}", accountBalance).Methods(http.MethodGet)
	router.HandleFunc("/all", fetchAll).Methods(http.MethodGet)

	// Default handler
	router.NotFoundHandler = http.HandlerFunc(notFound)
}

func logError(err error) {
	if err == nil {
		return
	}
	slog.Error(err.Error())
}

type ResponseError struct {
	Timestamp string `json:"timestamp"`
	Error     any    `json:"error"`
}

func responseError(w http.ResponseWriter, errString string, status int) {
	responseJson(w, ResponseError{Timestamp: time.Now().Format(time.RFC3339), Error: errString}, status)
}

func responseJson(w http.ResponseWriter, resp any, statusCode int) {
	code := http.StatusOK
	if statusCode > 0 {
		code = statusCode
	}
	w.WriteHeader(code)
	logError(json.NewEncoder(w).Encode(resp))
}
