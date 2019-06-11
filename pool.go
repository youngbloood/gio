package gio

import (
	"sync"
)

type pool struct {
	jobPool chan chan Jober
	result  Result

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
