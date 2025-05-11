package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/scriptoxin/yandex-liceum-go-calc/internal/evaluator"
	pb "github.com/scriptoxin/yandex-liceum-go-calc/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	client := pb.NewCalculatorClient(conn)

	for {
		// Запрашиваем задачу у оркестратора
		task, err := client.GetTask(context.Background(), &pb.Empty{})
		if err != nil {
			log.Printf("GetTask error: %v", err)
			time.Sleep(time.Second)
			continue
		}

		// Вычисляем с помощью evaluator.Calc
		result, err := evaluator.Calc(task.Expression)
		if err != nil {
			log.Printf("Calc error for %q: %v", task.Expression, err)
			// можно отправить статус «error», но пока отправляем 0
			result = 0
		}

		// Отправляем результат обратно
		_, err = client.SubmitResult(context.Background(), &pb.Result{
			Id:    task.Id,
			Value: result,
		})
		if err != nil {
			log.Printf("SubmitResult error: %v", err)
		}

		time.Sleep(time.Second)
	}
}
