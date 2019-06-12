package gio_test

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"runtime"

	"github.com/youngbloood/gio"
)

func TestMaster(t *testing.T) {

	result := make(Result, 5)
	m := gio.NewMaster(5, result)
	m.Run()
	var jobs = newJobs(100)
	go m.Push(jobs...)
	m.Stop()
	runtime.Gosched()
	log.Println("result = ", result)
}

func newJobs(length int) (jobs []gio.Jober) {
	for i := 0; i < length; i++ {
		jobs = append(jobs, job{i})
	}
	return
}

// Next implement by yourself.

// Job is custom job.
type job struct {
	A int
}

// implement the gio.Jober interface
func (j job) Run(result gio.Saver) error {
	// ...

	// save the result
	return result.Save(fmt.Sprintf("job%d", j.A), j.A)
}

var rwmux sync.RWMutex

// Result is custom result
type Result map[string]interface{}

// implement the gio.Saver interface.(you can handle it with sync or async)
func (r Result) Save(v ...interface{}) error {
	rwmux.Lock()
	r[v[0].(string)] = v[1]
	rwmux.Unlock()
	return nil
}
