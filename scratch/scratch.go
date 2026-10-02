package main

import (
	"github.com/goforj/godump"
	"github.com/mugtree/magpie/lib"
)

type Log struct {
	ID              int64
	RunType         string
	ArticlesCreated int64
}

func main() {

	logs := []Log{
		{ID: 1, RunType: "Today", ArticlesCreated: 15},
		{ID: 2, RunType: "Today", ArticlesCreated: 12},
		{ID: 3, RunType: "Tomorrow", ArticlesCreated: 1},
		{ID: 4, RunType: "Tomorrow", ArticlesCreated: 3},
		{ID: 5, RunType: "Never", ArticlesCreated: 4},
	}

	logsMap := lib.GroupBy(logs, func(l Log) string {
		return l.RunType
	})

	godump.Dump(logsMap)

}
