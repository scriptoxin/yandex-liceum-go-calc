package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Task представляет задачу, которую агент получает от оркестратора.
type Task struct {
	ID            int    `json:"id"`
	ExpressionID  int    `json:"expression_id"`
	Arg1          string `json:"arg1"`
	Arg2          string `json:"arg2"`
	Operation     string `json:"operation"`
	OperationTime int    `json:"operation_time"`
}

// getOrchestratorURL возвращает базовый URL оркестратора из переменной окружения ORCHESTRATOR_URL.
// Если переменная не установлена, используется значение по умолчанию "http://localhost:8080".
func getOrchestratorURL() string {
	url := os.Getenv("ORCHESTRATOR_URL")
	if url == "" {
		url = "http://localhost:8080"
	}
	return url
}

// fetchTask запрашивает задачу у оркестратора по эндпоинту GET /internal/task.
func fetchTask() (*Task, error) {
	resp, err := http.Get(getOrchestratorURL() + "/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Если задач нет, возвращаем nil, не ошибку.
		return nil, nil
	}

	var res struct {
		Task Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res.Task, nil
}

// sendTaskResult отправляет результат задачи обратно оркестратору через POST /internal/task.
func sendTaskResult(taskID int, result float64) error {
	payload := map[string]interface{}{
		"id":     taskID,
		"result": result,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(getOrchestratorURL()+"/internal/task", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

// performOperation выполняет указанную в задаче операцию с задержкой, заданной в миллисекундах.
func performOperation(task *Task) (float64, error) {
	arg1, err := strconv.ParseFloat(task.Arg1, 64)
	if err != nil {
		return 0, err
	}
	arg2, err := strconv.ParseFloat(task.Arg2, 64)
	if err != nil {
		return 0, err
	}
	// Симулируем длительное выполнение операции.
	time.Sleep(time.Millisecond * time.Duration(task.OperationTime))
	switch task.Operation {
	case "+":
		return arg1 + arg2, nil
	case "-":
		return arg1 - arg2, nil
	case "*":
		return arg1 * arg2, nil
	case "/":
		if arg2 == 0 {
			return 0, nil // Можно вернуть ошибку, если требуется
		}
		return arg1 / arg2, nil
	default:
		return 0, nil
	}
}

// agentWorker – функция-воркер, которая постоянно запрашивает задачу, выполняет её и отправляет результат.
func agentWorker() {
	for {
		task, err := fetchTask()
		if err != nil {
			log.Println("Error fetching task:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		if task == nil {
			// Если задач нет, ждём секунду и пробуем снова.
			time.Sleep(1 * time.Second)
			continue
		}
		log.Printf("Fetched task: %+v", task)
		result, err := performOperation(task)
		if err != nil {
			log.Println("Error performing operation:", err)
			continue
		}
		if err := sendTaskResult(task.ID, result); err != nil {
			log.Println("Error sending task result:", err)
		} else {
			log.Printf("Task %d completed with result %f", task.ID, result)
		}
	}
}

func main() {
	// Определяем количество воркеров из переменной окружения COMPUTING_POWER (по умолчанию 1).
	numWorkers := 1
	if cp := os.Getenv("COMPUTING_POWER"); cp != "" {
		if v, err := strconv.Atoi(cp); err == nil && v > 0 {
			numWorkers = v
		}
	}
	log.Printf("Starting agent with %d worker(s)...", numWorkers)
	for i := 0; i < numWorkers; i++ {
		go agentWorker()
	}
	// Блокировка main, чтобы программа не завершалась.
	select {}
}
