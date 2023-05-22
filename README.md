# Dbalancer

Golang connection load balancer for master slave replication. It assumes a single master that can be written to and zero-to-many read replicas that can be load balanced to.

```go

// Create a new DBalancer with the master DB
bl := dbalancer.NewDBalancer(db, rep, rep2)
defer bl.Close()

// Optional configuration
bl.SetMaxOpenConns(100)
bl.SetMaxIdleConns(50)

c, _ := bl.ReadConn(ctx)
// read with conn, will balance between master, and replicas
defer c.Close()

c, _ = bl.WriteConn(ctx)
// write with conn
defer c.Close()

```
