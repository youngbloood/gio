package gio_test

import (
	"testing"

	"github.com/youngbloood/gio"
)

func TestMaster(t *testing.T) {

	result := make(gio.Result, 10)
	m := gio.NewMaster(10, result)
	m.Run()
	defer m.Stop()

	m.Push()

}

type job1 struct {
	A int
}
