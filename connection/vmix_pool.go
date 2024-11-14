package connection

import "github.com/puzpuzpuz/xsync/v3"

type VMixConnectionPool interface {
	Add(addr string, vc VMixConnection)
	Remove(vc VMixConnection)
	// RemoveByAddr?
	AddNew(addr string) VMixConnection
}

type vMixConnectionPool struct {
	pool *xsync.MapOf[string, VMixConnection]
}

func NewvMixConnectionPool() VMixConnectionPool {
	return &vMixConnectionPool{
		pool: xsync.NewMapOf[string, VMixConnection](),
	}
}

// Add implements VMixConnectionPool.
func (v *vMixConnectionPool) Add(addr string, vc VMixConnection) {
	v.pool.Store(addr, vc)
}

// AddNew implements VMixConnectionPool.
func (v *vMixConnectionPool) AddNew(addr string) VMixConnection {
	vc := NewVMixConnection(addr)
	v.pool.Store(addr, vc)
	return vc
}

// Remove implements VMixConnectionPool.
func (v *vMixConnectionPool) Remove(vc VMixConnection) {
	panic("unimplemented")
}
