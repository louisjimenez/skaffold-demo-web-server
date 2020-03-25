package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/louisjimenez/skaffold-demo-config"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	host string
	port string
)

const (
	defGrpcServerAddr = "localhost:8090"
	GrpcHostEnv       = "GRPC_SERVER_HOST"
	GrpcPortEnv       = "GRPC_SERVER_PORT"
)

var client pb.TodoClient

func handler(w http.ResponseWriter, r *http.Request) {
	str := listTodoItems(&pb.Category{Name: "work"})
	fmt.Fprint(w, str)
}

func listTodoItems(cat *pb.Category) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := client.List(ctx, cat)
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

func connectGRPC(host string) (*grpc.ClientConn, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client = pb.NewTodoClient(conn)
	http.HandleFunc("/", handler)
	return conn, nil
}

func healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "OK")
}

func main() {
	flag.StringVar(&host, "host", "", "host of grpc server")
	flag.StringVar(&port, "port", "", "port of grpc server")
	flag.Parse()
	if host == "" {
		if host, port = os.Getenv(GrpcHostEnv), os.Getenv(GrpcPortEnv); host == "" || port == "" {
			log.Printf("host %v port %v \n", host, port)
			host = defGrpcServerAddr
			log.Printf("addr %v \n", host)

		} else {
			host = strings.Join([]string{host, port}, ":")
		}
	}
	log.Printf("GRPC server address is %v \n", host)


	http.HandleFunc("/health", healthcheckHandler)
	conn, err := connectGRPC(host)
	if err != nil {
		log.Printf("Unable to connect: %v", err)
	}
	defer conn.Close()
	http.ListenAndServe(":9000", nil)
}
