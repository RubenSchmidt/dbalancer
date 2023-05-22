package dbalancer

import (
	"context"
	"database/sql"
	"sync/atomic"
)

// DBalancer is a database balancer that can be used to balance read queries
type DBalancer struct {
	master *sql.DB
	rrs    []*sql.DB
	next   uint32
}

// NewDBalancer creates a new DBalancer with the master DB
// The master DB is the DB that will be used for write queries
// The master DB will also be used for read queries
func NewDBalancer(master *sql.DB, rrs ...*sql.DB) *DBalancer {
	rrs = append(rrs, master)
	return &DBalancer{
		master: master,
		rrs:    rrs,
	}
}

// AddReadReplica adds a read replica to the DBalancer
func (b *DBalancer) AddReadReplica(r *sql.DB) {
	b.rrs = append(b.rrs, r)
}

// Master returns the master DB
func (b *DBalancer) Master() *sql.DB {
	return b.master
}

// GetReadConn returns a read connection from the DBalancer using round robin load balancing
func (b *DBalancer) ReadConn(ctx context.Context) (*sql.Conn, error) {
	n := atomic.AddUint32(&b.next, 1)
	return b.rrs[(int(n)-1)%len(b.rrs)].Conn(ctx)
}

// GetWriteConn returns a write connection from the DBalancer
func (b *DBalancer) WriteConn(ctx context.Context) (*sql.Conn, error) {
	return b.master.Conn(ctx)
}

// SetMaxOpenConns sets the maximum number of open connections to the master and read replicas
func (b *DBalancer) SetMaxOpenConns(n int) {
	b.master.SetMaxOpenConns(n)
	for _, r := range b.rrs {
		r.SetMaxOpenConns(n)
	}
}

// SetMaxIdleConns sets the maximum number of idle connections to the master and read replicas
func (b *DBalancer) SetMaxIdleConns(n int) {
	b.master.SetMaxIdleConns(n)
	for _, r := range b.rrs {
		r.SetMaxIdleConns(n)
	}
}

func (b *DBalancer) Close() error {
	for _, r := range b.rrs {
		err := r.Close()
		if err != nil {
			return err
		}
	}
	return b.master.Close()
}
