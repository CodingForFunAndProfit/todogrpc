package json

import (
	"encoding/json"
	"io"

	"github.com/CodingForFunAndProfit/todogrpc/internal/jobs"
)

func Load(r io.Reader) ([]*jobs.Todo, error) {
	var todos []*jobs.Todo
	return todos, json.NewDecoder(r).Decode(&todos)
}
func Flush(w io.Writer, todos []*jobs.Todo) error {
	return json.NewEncoder(w).Encode(todos)
}
