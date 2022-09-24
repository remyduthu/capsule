# 🚀 Capsule

[![GoReference](https://pkg.go.dev/badge/github.com/remyduthu/capsule.svg)](https://pkg.go.dev/github.com/remyduthu/capsule)

Capsule is a Go package with a minimalist API to wrap the official
[database/sql](https://pkg.go.dev/database/sql) package. It provides idiomatic
functions to build SQL [queries with ordinal
markers](https://pkg.go.dev/github.com/lib/pq#hdr-Queries) and to scan SQL
row(s). Capsule is not (and will never be) an ORM.

## Examples

### Query

The following example shows how to query rows using capsule and the official
[database/sql](https://pkg.go.dev/database/sql) packages.

```go
var db *sql.DB

type user struct {
	ID int
}

func main() {
	// The ordinal markers will be added after every '$' character by using the
	// `Render` method.
	c := capsule.New("SELECT * FROM users WHERE email = $", "foo@bar.com")

	// Dynamically extend the capsule.
	if ... {
		c.Add("WHERE first_name = $", "foo")
	}

	rows, err := db.Query(c.Render(), c.Args()...)
	if err != nil {
		panic(err)
	}

	users, err := capsule.Scan(rows, scanUser)
	if err != nil {
		panic(err)
	}
}

func scanUser(handler func(dest ...any) error) (*user, error) {
	var result user
	return &result, handler(&result.ID)
}
```

### Query Row

The following example shows how to query a single row using capsule and the
official [database/sql](https://pkg.go.dev/database/sql) packages.

```go
var db *sql.DB

type user struct {
	ID int
}

func main() {
	c := capsule.New("SELECT * FROM users WHERE id = $", 1)

	user, err := capsule.ScanRow(db.QueryRow(c.Render(), c.Args()...), scanUser)
	if err != nil {
		panic(err)
	}
}

func scanUser(handler func(dest ...any) error) (*user, error) {
	var result user
	return &result, handler(&result.ID)
}
```
