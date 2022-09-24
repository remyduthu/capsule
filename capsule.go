package capsule

import (
	"database/sql"
	"fmt"
	"strings"
)

type capsule struct {
	query string
	args  []any
}

// New returns a new capsule with the query, query and the arguments, args.
func New(query string, args ...any) *capsule {
	result := &capsule{args: args}
	result.setQuery(query)

	return result
}

// Args returns the arguments of the capsule.
func (c capsule) Args() []any {
	return c.args
}

// Render renders the capsule query with ordinal markers. It suffixes every '$'
// character with the appropriate ordinal marker.
func (c capsule) Render() string {
	marker, result := 0, ""

	for _, char := range c.query {
		if string(char) != "$" {
			result = result + string(char)
			continue
		}

		marker++
		result = result + "$" + fmt.Sprint(marker)
	}

	return result
}

// Add adds the query, query and the arguments, args to the capsule.
func (c *capsule) Add(query string, args ...any) *capsule {
	c.setQuery(query)
	c.args = append(c.args, args...)

	return c
}

func (c *capsule) setQuery(query string) {
	if query == "" {
		return
	}

	c.query = strings.TrimSpace(c.query + " " + strings.TrimSpace(query))
}

type ScanRowHandler[T any] func(handler func(dest ...any) error) (*T, error)

func Scan[T any](rows *sql.Rows, handler ScanRowHandler[T]) ([]*T, error) {
	defer rows.Close()

	var results []*T
	for rows.Next() {
		result, err := handler(rows.Scan)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func ScanRow[T any](row *sql.Row, handler ScanRowHandler[T]) (*T, error) {
	return handler(row.Scan)
}
