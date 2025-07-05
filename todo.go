package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task       string
	Done       bool
	CreatAt    time.Time
	CompleteAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:       task,
		Done:       false,
		CreatAt:    time.Now(),
		CompleteAt: time.Time{},
	}

	*t = append(*t, todo)
}

func (t *Todos) Complete(index int) error {
	ls := *t
	if index < 0 || index > len(ls) {
		return errors.New("invalid index")
	}
	ls[index-1].CompleteAt = time.Now()
	ls[index-1].Done = true

	return nil
}

func (t *Todos) Delete(index int) error {
	ls := *t
	if index < 0 || index > len(ls) {
		return errors.New("invalid index")
	}

	*t = append(ls[:index-1], ls[index:]...)

	return nil

}

func (t *Todos) Load(filename string) error {

	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (t *Todos) Print() {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done"},
			{Align: simpletable.AlignCenter, Text: "Created at"},
			{Align: simpletable.AlignCenter, Text: "Completed at"},
		},
	}

	var cells [][]*simpletable.Cell

	for idx, item := range *t {
		idx++

		task := blue(item.Task)
		done := red("No")
		createAt := blue(item.CreatAt.Format(time.RFC822))
		completeAt := blue(item.CompleteAt.Format(time.RFC822))

		if item.Done {
			task = green(fmt.Sprintf("\u2705 %s", item.Task))
			done = green("Yes")
			createAt = green(item.CreatAt.Format(time.RFC822))
			completeAt = green(item.CompleteAt.Format(time.RFC822))
		}
		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", idx)},
			{Text: task},
			{Text: done},
			{Text: createAt},
			{Text: completeAt},
		})
	}

	table.Body = &simpletable.Body{
		Cells: cells,
	}

	// table.Footer = &simpletable.Footer{
	// 	Cells: []*simpletable.Cell{
	// 		{Align: simpletable.AlignCenter, Span: 5, Text: "Your TODOS are here"},
	// 	},
	// }

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}
