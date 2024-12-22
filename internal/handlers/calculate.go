package handlers

import (
	"encoding/json"
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
		Result float64 `json:"result"`
	}{
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// respondWithError - удобный метод для отправки ошибок.
func respondWithError(w http.ResponseWriter, appErr *errors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	json.NewEncoder(w).Encode(appErr)
}
