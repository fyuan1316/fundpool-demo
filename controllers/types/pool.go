package types

import "fmt"

type Pool struct {
	Name    string
	Balance int64
	Minus   int64
}

func NewPool(name string, b, m int64) Pool {
	return Pool{
		Name:    name,
		Balance: b,
		Minus:   m,
	}
}

func (p *Pool) Max() int64 {
	return p.Balance - p.Minus
}
func (p *Pool) StepFetch(step int64) {
	p.Minus = p.Minus + step
}
func (p *Pool) GetShortfalls() int64 {
	return p.Minus
}
func (p Pool) Info() string {
	return fmt.Sprintf("Pool: %s, fetch=%v", p.Name, p.Minus)
}
