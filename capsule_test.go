package capsule

import (
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAdd(t *testing.T) {
	cases := map[string]struct {
		input, want *capsule
	}{
		"with query and args": {input: New("SELECT * FROM users WHERE 1 = 1").Add("AND email = $ AND first_name = $", "foo@bar.com", "foo"), want: &capsule{query: "SELECT * FROM users WHERE 1 = 1 AND email = $ AND first_name = $", args: []any{"foo@bar.com", "foo"}}},
		"without args":        {input: New("SELECT * FROM users WHERE 1 = 1").Add("AND email = $"), want: &capsule{query: "SELECT * FROM users WHERE 1 = 1 AND email = $"}},
		"without query":       {input: New("SELECT * FROM users").Add("", "foo@bar.com"), want: &capsule{query: "SELECT * FROM users", args: []any{"foo@bar.com"}}},
		"with query with leading and trailing whitespaces": {input: New(" SELECT * FROM users WHERE 1 = 1 ").Add(" AND email = $ ", "foo@bar.com"), want: &capsule{query: "SELECT * FROM users WHERE 1 = 1 AND email = $", args: []any{"foo@bar.com"}}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			log.Println("input", c.input.query)
			log.Println("want", c.want.query)

			if diffs := cmp.Diff(c.input, c.want, cmp.AllowUnexported(capsule{})); diffs != "" {
				t.Fatal(diffs)
			}
		})
	}

}

func TestRender(t *testing.T) {
	cases := map[string]struct {
		input *capsule
		want  string
	}{
		"without marker": {input: &capsule{query: "SELECT * FROM users"}, want: "SELECT * FROM users"},
		"with marker":    {input: &capsule{query: "SELECT * FROM users WHERE email = $"}, want: "SELECT * FROM users WHERE email = $1"},
		"with markers":   {input: &capsule{query: "SELECT * FROM users WHERE email = $ AND first_name = $"}, want: "SELECT * FROM users WHERE email = $1 AND first_name = $2"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if diffs := cmp.Diff(c.input.Render(), c.want); diffs != "" {
				t.Fatal(diffs)
			}
		})
	}
}
