# Dbalancer

Golang connection load balancer for master slave replication. It assumes a single master that can be written to and zero-to-many read replicas that can be load balanced to.

```go

// Create a new DBalancer with the master DB
bl := dbalancer.NewDBalancer(db, rep, rep2)
defer bl.Close()

// Optional configuration for all databases
bl.SetMaxOpenConns(100)
bl.SetMaxIdleConns(50)

db := bl.ReadDB() // will balance between master, and replicas

db = bl.WriteDB() // returns master db

```
