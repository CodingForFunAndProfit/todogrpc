package protobuf

import (
	"io"
	"io/ioutil"

	"github.com/CodingForFunAndProfit/todogrpc/internal/jobs/v1"
	"google.golang.org/protobuf/proto"
)

func Load(r io.Reader) ([]*jobs.Todo, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var todos jobs.Todos
	return todos.Todos, proto.Unmarshal(b, &todos)
}

func Flush(w io.Writer, todos []*jobs.Todo) error {
	b, err := proto.Marshal(&jobs.Todos{Todos: todos})
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
