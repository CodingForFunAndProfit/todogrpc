package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CodingForFunAndProfit/todogrpc/internal/jobs/v1"
	// storage "github.com/CodingForFunAndProfit/todogrpc/internal/json"
	// storage "github.com/awoodbeck/gnp/ch12/gob"
	storage "github.com/CodingForFunAndProfit/todogrpc/internal/protobuf"
)

var dataFile string

func init() {
	flag.StringVar(&dataFile, "file", "db/todos.db", "data file")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			`Usage: %s [flags] [add todo, ...|complete #]
 add add comma-separated todos
 complete complete designated todo
Flags:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func load() ([]*jobs.Todo, error) {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return make([]*jobs.Todo, 0), nil
	}
	df, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := df.Close(); err != nil {
			fmt.Printf("closing data file: %v", err)
		}
	}()
	return storage.Load(df)
}

func flush(todos []*jobs.Todo) error {
	df, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := df.Close(); err != nil {
			fmt.Printf("closing data file: %v", err)
		}
	}()
	return storage.Flush(df, todos)
}

func list() error {
	todos, err := load()
	if err != nil {
		return err
	}
	if len(todos) == 0 {
		fmt.Println("You're all caught up!")
		return nil
	}
	fmt.Println("#\t[X]\tDescription")
	for i, todo := range todos {
		c := " "
		if todo.Completed {
			c = "X"
		}
		fmt.Printf("%d\t[%s]\t%s\n", i+1, c, todo.Description)
	}
	return nil
}

func add(s string) error {
	todos, err := load()
	if err != nil {
		return err
	}
	for _, todo := range strings.Split(s, ",") {
		if desc := strings.TrimSpace(todo); desc != "" {
			todos = append(todos, &jobs.Todo{
				Description: desc,
			})
		}
	}
	return flush(todos)
}

func completed(s string) error {
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	todos, err := load()
	if err != nil {
		return err
	}
	if i < 1 || i > len(todos) {
		return fmt.Errorf("todo %d not found", i)
	}
	todos[i-1].Completed = true
	return flush(todos)
}

func main() {
	flag.Parse()
	var err error
	switch strings.ToLower(flag.Arg(0)) {
	case "add":
		err = add(strings.Join(flag.Args()[1:], " "))
	case "complete":
		err = completed(flag.Arg(1))
	}
	if err != nil {
		log.Fatal(err)
	}
	err = list()
	if err != nil {
		log.Fatal(err)
	}
}
