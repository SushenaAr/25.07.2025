package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Constraint struct {
	Op  string
	Ver []int
}

func ParseConstraint(s string) (Constraint, error) {
	ops := []string{"<=", ">=", "<", ">", "="}
	for _, op := range ops {
		if strings.HasPrefix(s, op) {
			ver, err := parseVersion(strings.TrimPrefix(s, op))
			return Constraint{Op: op, Ver: ver}, err
		}
	}
	ver, err := parseVersion(s)
	return Constraint{Op: "=", Ver: ver}, err
}

func parseVersion(s string) ([]int, error) {
	parts := strings.Split(s, ".")
	out := make([]int, len(parts))
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid version %s", s)
		}
		out[i] = n
	}
	return out, nil
}

func Compare(a, b []int) int {
	for i := 0; i < max(len(a), len(b)); i++ {
		x, y := 0, 0
		if i < len(a) {
			x = a[i]
		}
		if i < len(b) {
			y = b[i]
		}
		if x < y {
			return -1
		} else if x > y {
			return 1
		}
	}
	return 0
}

func (c Constraint) Match(v string) bool {
	parsed, err := parseVersion(v)
	if err != nil {
		return false
	}
	cmp := Compare(parsed, c.Ver)
	switch c.Op {
	case "=":
		return cmp == 0
	case "<":
		return cmp == -1
	case "<=":
		return cmp <= 0
	case ">":
		return cmp == 1
	case ">=":
		return cmp >= 0
	default:
		return false
	}
}

func FindBestMatch(files []string, pkg, constraintStr string) (string, string, error) {
	re := regexp.MustCompile(pkg + `-(\d+(?:\.\d+)*).tar.gz`)
	var matched []string
	var versions [][]int
	var names []string

	c, err := ParseConstraint(constraintStr)
	if err != nil {
		return "", "", err
	}
	for _, f := range files {
		m := re.FindStringSubmatch(f)
		if len(m) == 2 && c.Match(m[1]) {
			v, _ := parseVersion(m[1])
			matched = append(matched, f)
			versions = append(versions, v)
			names = append(names, strings.TrimSuffix(f, ".tar.gz"))
		}
	}
	if len(matched) == 0 {
		return "", "", fmt.Errorf("no match for %s %s", pkg, constraintStr)
	}
	i := maxVersionIndex(versions)
	return matched[i], names[i], nil
}

func maxVersionIndex(vers [][]int) int {
	max := 0
	for i := 1; i < len(vers); i++ {
		if Compare(vers[i], vers[max]) == 1 {
			max = i
		}
	}
	return max
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
