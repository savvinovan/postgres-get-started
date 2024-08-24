### Best practices postgres in golang

best lib - https://github.com/jackc/pgx

1. Configure your connection pool size
```golang
db.SetMaxIdleConns(MaxIdleConns)
db.SetMaxOpenConns(MaxOpenConns)
```

2. Monitor your connection pool
```golang
db.Stats()

// pgx
connPool.Stat()
```

3. Log your pgx
```golang
pgx.Logger // interface
```

4. Use pg-bouncer for connection pooling
- https://www.pgbouncer.org/features.html
- Recommendation: Choose transaction pooling mode
- A server connection is assigned to a client only during a transaction.
- Needs to note that prepared statements would work only inside a transaction.
- Needs to note that every query in golang uses prepared statements.
- Use the Simple Query mode to avoid problems with prepared statements in the transactional mode PgBouncer.
- https://pkg.go.dev/github.com/jackc/pgx/v5#hdr-PgBouncer
  By default pgx automatically uses prepared statements. Prepared statements are incompatible with PgBouncer. This can be disabled by setting a different QueryExecMode in ConnConfig.DefaultQueryExecMode.

5. Use HaProxy for load balancing across multiple pg-bouncers

