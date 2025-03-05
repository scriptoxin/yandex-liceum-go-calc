package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/scriptoxin/yandex-liceum-go-calc/internal/evaluator"
)

// ExpressionStatus представляет статус вычисления.
type ExpressionStatus string

const (
	StatusPending  ExpressionStatus = "pending"
	StatusComplete ExpressionStatus = "complete"
)

// Expression хранит информацию об арифметическом выражении.
type Expression struct {
	ID        int              `json:"id"`
	Expr      string           `json:"expression"`
	Status    ExpressionStatus `json:"status"`
	Result    *float64         `json:"result,omitempty"`
	Error     string           `json:"error,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// Task представляет задачу для агента.
type Task struct {
	ID            int    `json:"id"`
	ExpressionID  int    `json:"expression_id"`
	Arg1          string `json:"arg1"`
	Arg2          string `json:"arg2"`
	Operation     string `json:"operation"`
	OperationTime int    `json:"operation_time"`
}

var (
	expressions      = make(map[int]*Expression)
	expressionsMutex = newSyncMutex()
	nextExpressionID = 1

	tasks      = make(map[int]*Task)
	tasksMutex = newSyncMutex()
	nextTaskID = 1
)

// --- syncMutex из evaluator оставляем такой же ---
type syncMutex struct {
	ch chan struct{}
}

func newSyncMutex() *syncMutex {
	m := &syncMutex{ch: make(chan struct{}, 1)}
	m.ch <- struct{}{}
	return m
}

func (m *syncMutex) Lock() {
	<-m.ch
}

func (m *syncMutex) Unlock() {
	m.ch <- struct{}{}
}

// HandleCalculate обрабатывает POST-запрос на добавление выражения.
func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, `{"error": "Invalid request data"}`, http.StatusUnprocessableEntity)
		return
	}

	exprObj := addExpression(req.Expression)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": exprObj.ID})
}

// HandleExpressions возвращает список всех выражений.
func HandleExpressions(w http.ResponseWriter, r *http.Request) {
	expressionsMutex.Lock()
	defer expressionsMutex.Unlock()
	var list []*Expression
	for _, expr := range expressions {
		list = append(list, expr)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": list})
}

// HandleExpressionByID возвращает выражение по ID.
func HandleExpressionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}
	expressionsMutex.Lock()
	exprObj, exists := expressions[id]
	expressionsMutex.Unlock()
	if !exists {
		http.Error(w, `{"error": "Expression not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": exprObj})
}

// HandleGetTask возвращает задачу для вычисления агенту.
func HandleGetTask(w http.ResponseWriter, r *http.Request) {
	tasksMutex.Lock()
	defer tasksMutex.Unlock()
	for id, task := range tasks {
		delete(tasks, id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
		return
	}
	http.Error(w, `{"error": "No task available"}`, http.StatusNotFound)
}

// HandlePostTask принимает результат выполнения задачи от агента.
func HandlePostTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request data"}`, http.StatusUnprocessableEntity)
		return
	}
	expressionsMutex.Lock()
	exprObj, exists := expressions[req.ID]
	if !exists {
		expressionsMutex.Unlock()
		http.Error(w, `{"error": "Expression not found"}`, http.StatusNotFound)
		return
	}
	// Если выражение уже помечено с ошибкой, не обновляем его.
	if exprObj.Error != "" {
		expressionsMutex.Unlock()
		http.Error(w, `{"error": "Expression already marked as error"}`, http.StatusUnprocessableEntity)
		return
	}
	exprObj.Status = StatusComplete
	exprObj.Result = &req.Result
	exprObj.UpdatedAt = time.Now()
	expressionsMutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Result recorded"})
}

// addExpression создает новое выражение и планирует задачу.
func addExpression(expr string) *Expression {
	expressionsMutex.Lock()
	id := nextExpressionID
	nextExpressionID++
	now := time.Now()
	exprObj := &Expression{
		ID:        id,
		Expr:      expr,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	expressions[id] = exprObj
	expressionsMutex.Unlock()

	// Разбиваем выражение на задачу
	createTask(exprObj.ID, expr)
	return exprObj
}

// markExpressionError обновляет статус выражения как завершенного с ошибкой.
func markExpressionError(expressionID int, errMsg string) {
	expressionsMutex.Lock()
	defer expressionsMutex.Unlock()
	if expr, exists := expressions[expressionID]; exists {
		expr.Status = StatusComplete
		expr.Result = nil
		expr.Error = errMsg
		expr.UpdatedAt = time.Now()
		fmt.Printf("Expression %d marked with error: %s\n", expressionID, errMsg)
	}
}

// createTask разбивает выражение на задачу для агента.
// Если аргументы некорректны или обнаружено деление на ноль, помечает выражение как ошибочное.
func createTask(expressionID int, expr string) {
	var opPos int = -1
	var op byte
	for i := 0; i < len(expr); i++ {
		if expr[i] == '+' || expr[i] == '-' || expr[i] == '*' || expr[i] == '/' {
			opPos = i
			op = expr[i]
			break
		}
	}
	if opPos == -1 {
		return
	}
	arg1 := expr[:opPos]
	arg2 := expr[opPos+1:]

	// Если аргументы не являются числом, пытаемся вычислить их
	if _, err := strconv.ParseFloat(arg1, 64); err != nil && arg1 != "" {
		if val, err := evaluator.Calc(arg1); err == nil {
			arg1 = fmt.Sprintf("%.0f", val)
		} else {
			markExpressionError(expressionID, "invalid arguments")
			return
		}
	}
	if _, err := strconv.ParseFloat(arg2, 64); err != nil && arg2 != "" {
		if val, err := evaluator.Calc(arg2); err == nil {
			arg2 = fmt.Sprintf("%.0f", val)
		} else {
			markExpressionError(expressionID, "invalid arguments")
			return
		}
	}

	// Повторная проверка, чтобы убедиться, что оба аргумента теперь корректны
	_, err1 := strconv.ParseFloat(arg1, 64)
	val2, err2 := strconv.ParseFloat(arg2, 64)
	if err1 != nil || err2 != nil {
		markExpressionError(expressionID, "invalid arguments")
		return
	}

	// Если операция деления и второй аргумент равен 0, помечаем выражение как ошибочное.
	if op == '/' && val2 == 0 {
		markExpressionError(expressionID, "division by zero")
		return
	}

	// Получаем время выполнения операции из переменных среды
	var opTime int
	switch op {
	case '+':
		opTime = getEnvAsInt("TIME_ADDITION_MS", 1000)
	case '-':
		opTime = getEnvAsInt("TIME_SUBTRACTION_MS", 1000)
	case '*':
		opTime = getEnvAsInt("TIME_MULTIPLICATIONS_MS", 1000)
	case '/':
		opTime = getEnvAsInt("TIME_DIVISIONS_MS", 1000)
	}

	// Создаем и регистрируем задачу
	tasksMutex.Lock()
	taskID := nextTaskID
	nextTaskID++
	task := &Task{
		ID:            taskID,
		ExpressionID:  expressionID,
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     string(op),
		OperationTime: opTime,
	}
	tasks[taskID] = task
	tasksMutex.Unlock()
}

// getEnvAsInt читает переменную окружения и возвращает её значение как int.
func getEnvAsInt(name string, defaultVal int) int {
	valStr := os.Getenv(name)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
