package gio

import (
	"context"
	"sync"
)

type pool struct {
	ctx     context.Context
	jobPool chan chan Jober
	result  Saver

	erw  sync.RWMutex
	errs []error
}

func (p *pool) addErr(err error) {
	p.erw.Lock()
	p.errs = append(p.errs, err)
	p.erw.Unlock()
}

func (p *pool) getErrs() []error {
	p.erw.RLock()
	defer p.erw.RUnlock()
	return p.errs
}
