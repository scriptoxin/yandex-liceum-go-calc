package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/scriptoxin/yandex-liceum-go-calc/internal/evaluator"
	"github.com/scriptoxin/yandex-liceum-go-calc/pkg/errors"
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, errors.ErrBadRequest)
		return
	}

	result, err := evaluator.Calc(request.Expression)
	if err != nil {
		if err == evaluator.ErrInvalidExpression {
			respondWithError(w, errors.ErrInvalidExpression)
		} else {
			respondWithError(w, errors.ErrInternalServer)
		}
		return
	}

	response := struct {
		Result string `json:"result"`
	}{
		Result: fmt.Sprintf("%.0f", result),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func respondWithError(w http.ResponseWriter, appErr *errors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	errorResponse := map[string]string{"error": appErr.Message}
	json.NewEncoder(w).Encode(errorResponse)
}
