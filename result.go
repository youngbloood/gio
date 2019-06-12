package gio

// Saver interface
type Saver interface {
	Save(...interface{}) error
}
