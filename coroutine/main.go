package coroutine

import _ "unsafe"

type coro struct{}

//go:linkname newcoro runtime.newcoro
func newcoro(f func(*coro)) *coro

//go:linkname coroswitch runtime.coroswitch
func coroswitch(c *coro)

type Coro struct {
	handle *coro
}

func Create(f func(self *Coro, parent *coro)) *Coro {
	res := &Coro{}
	res.handle = newcoro(func(internalCoro *coro) {
		f(res, internalCoro)
	})
	return res
}

func SwitchRaw(c *coro) {
	if c != nil {
		coroswitch(c)
	}
}

func (c *Coro) SwitchTo() {
	if c.handle != nil {
		coroswitch(c.handle)
	}
}
