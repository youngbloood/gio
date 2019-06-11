package gio_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/youngbloood/gio"
)

func TestMaster(t *testing.T) {

	result := make(gio.Result, 5)
	m := gio.NewMaster(5, result)
	m.Run()
	defer m.Stop()

	m.Push(
		job1{1},
		job2{2},
		job3{3},
		job4{4},
		job5{5},
	)

	time.Sleep(2 * time.Second)
	fmt.Println("reuslt=", result)
}

type job1 struct {
	A int
}

func (j job1) Run(result gio.Result) error {
	result["job1"] = j.A
	return nil
}

type job2 struct {
	A int
}

func (j job2) Run(result gio.Result) error {
	result["job2"] = j.A
	return nil
}

type job3 struct {
	A int
}

func (j job3) Run(result gio.Result) error {
	result["job3"] = j.A
	return nil
}

type job4 struct {
	A int
}

func (j job4) Run(result gio.Result) error {
	result["job4"] = j.A
	return nil
}

type job5 struct {
	A int
}

func (j job5) Run(result gio.Result) error {
	result["job"] = j.A
	return nil
}
