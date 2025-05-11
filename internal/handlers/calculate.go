package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/scriptoxin/yandex-liceum-go-calc/pkg/db"
	"github.com/scriptoxin/yandex-liceum-go-calc/pkg/jwt"
)

type calcRequest struct {
	Expression string `json:"expression"`
}

// AuthMiddleware проверяет JWT и кладёт user_id в контекст
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if len(header) < 8 || header[:7] != "Bearer " {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := header[7:]
		uid, err := jwt.Parse(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Calculate — POST /api/v1/calculate
// Генерируем UUID, сохраняем в SQLite, статус = pending
func Calculate(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("user_id").(int)

	var req calcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Генерируем уникальный ID задачи
	id := uuid.NewString()

	// Сохраняем в БД
	_, err := db.Conn.Exec(
		"INSERT INTO expressions(id, user_id, expression, status) VALUES(?, ?, ?, ?)",
		id, uid, req.Expression, "pending",
	)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// TODO: отправить task по gRPC агентам

	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// GetExpressions — GET /api/v1/expressions
func GetExpressions(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("user_id").(int)

	rows, err := db.Conn.Query(
		"SELECT id, status, result FROM expressions WHERE user_id = ?",
		uid,
	)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var (
			id     string
			status string
			res    sql.NullFloat64
		)
		rows.Scan(&id, &status, &res)
		item := map[string]interface{}{
			"id":     id,
			"status": status,
		}
		if res.Valid {
			item["result"] = res.Float64
		}
		list = append(list, item)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": list})
}

// GetExpression — GET /api/v1/expressions/{id}
func GetExpression(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("user_id").(int)
	id := mux.Vars(r)["id"]

	var (
		expr   string
		status string
		res    sql.NullFloat64
	)
	err := db.Conn.QueryRow(
		"SELECT expression, status, result FROM expressions WHERE id = ? AND user_id = ?",
		id, uid,
	).Scan(&expr, &status, &res)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	out := map[string]interface{}{
		"id":     id,
		"status": status,
	}
	if res.Valid {
		out["result"] = res.Float64
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": out})
}
