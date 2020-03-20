package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	pb "github.com/louisjimenez/skaffold-demo-web-server/todo"
)

const (
	address = "localhost:50051"
)

var client pb.TodoClient

func handler(w http.ResponseWriter, r *http.Request) {
	str := listTodoItems(&pb.Category{Name: "work"})
	fmt.Fprint(w, str)
}

func listTodoItems(cat *pb.Category) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := client.ListTodos(ctx, cat)
	if err != nil {
		log.Fatalf("Unable to ListTodos: %v", err)
	}
	var todoList strings.Builder
	for {
		todo, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unable to ListTodos: %v", err)
		}
		todoList.WriteString(todo.Description)
		todoList.WriteString("\n")
	}
	return todoList.String()
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()
	client = pb.NewTodoClient(conn)
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8080", nil)
}
