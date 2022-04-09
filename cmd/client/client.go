package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CodingForFunAndProfit/todogrpc/internal/jobs/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr, caCertFn string

func init() {
	flag.StringVar(&addr, "address", "localhost:34443", "server address")
	flag.StringVar(&caCertFn, "ca-cert", "certs/cert.pem", "CA certificate")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			`Usage: %s [flags] [add todo, ...|complete #]
 add add comma-separated todos
 completed complete designated todo
Flags:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func list(ctx context.Context, client jobs.TodoSvcClient) error {
	todos, err := client.List(ctx, new(jobs.Empty))
	if err != nil {
		return err
	}
	if len(todos.Todos) == 0 {
		fmt.Println("You have nothing to do!")
		return nil
	}
	fmt.Println("#\t[X]\tDescription")
	for i, todo := range todos.Todos {
		c := " "
		if todo.Completed {
			c = "X"
		}
		fmt.Printf("%d\t[%s]\t%s\n", i+1, c, todo.Description)
	}
	return nil
}

func add(ctx context.Context, client jobs.TodoSvcClient,
	s string) error {
	todos := new(jobs.Todos)
	for _, todo := range strings.Split(s, ",") {
		if desc := strings.TrimSpace(todo); desc != "" {
			todos.Todos = append(todos.Todos, &jobs.Todo{
				Description: desc,
			})
		}
	}
	var err error
	if len(todos.Todos) > 0 {
		_, err = client.Add(ctx, todos)
	}
	return err
}

func completed(ctx context.Context, client jobs.TodoSvcClient,
	s string) error {
	i, err := strconv.Atoi(s)
	if err == nil {
		_, err = client.Completed(ctx,
			&jobs.CompletedRequest{TodoNumber: int32(i)})
	}
	return err
}

func main() {
	flag.Parse()
	caCert, err := ioutil.ReadFile(caCertFn)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("failed to add certificate to pool")
	}
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(
			credentials.NewTLS(
				&tls.Config{
					CurvePreferences: []tls.CurveID{tls.CurveP256},
					MinVersion:       tls.VersionTLS12,
					RootCAs:          certPool,
				},
			),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	todosvc := jobs.NewTodoSvcClient(conn)
	ctx := context.Background()
	switch strings.ToLower(flag.Arg(0)) {
	case "add":
		err = add(ctx, todosvc, strings.Join(flag.Args()[1:], " "))
	case "completed":
		err = completed(ctx, todosvc, flag.Arg(1))
	}
	if err != nil {
		log.Fatal(err)
	}
	err = list(ctx, todosvc)
	if err != nil {
		log.Fatal(err)
	}
}
