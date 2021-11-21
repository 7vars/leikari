package mapper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/7vars/leikari/query"
)

func Check(node query.Node, v interface{}, mapname ...MapName) bool {
	switch n := node.(type) {
	case query.Comparsion:
		switch n.Value().Type() {
		case query.INT:
			if va, ok := Int64Value(n.Identifier().Name(), v, mapname...); ok {
				if vb, ok := n.Value().IntValue(); ok {
					return compareInt(va, n.Operator(), vb)
				}
			}
		case query.FLOAT:
			if va, ok := Float64Value(n.Identifier().Name(), v, mapname...); ok {
				if vb, ok := n.Value().FloatValue(); ok {
					return compareFloat(va, n.Operator(), vb)
				}
			}
		case query.BOOL:
			if va, ok := BoolValue(n.Identifier().Name(), v, mapname...); ok {
				if vb, ok := n.Value().BoolValue(); ok {
					return compareBool(va, n.Operator(), vb)
				}
			}
		case  query.STRING:
			if va, ok := StringValue(n.Identifier().Name(), v, mapname...); ok {
				if vb, ok := n.Value().StringValue(); ok {
					return compareString(va, n.Operator(), vb)
				}
			}
		}
	case query.PrefixCondition:
		switch n.Prefix() {
		case query.GROUP:
			return Check(n.Right(), v, mapname...)
		case query.NOT:
			return !Check(n.Right(), v, mapname...)
		case query.PR:
			if ident, ok := n.Right().(query.Identifier); ok {
				_, ok := Value(ident.Name(), v, mapname...)
				return ok
			}
		}
	case query.LogicalCondition:
		switch n.Logical() {
		case query.AND:
			return Check(n.Left(), v, mapname...) && Check(n.Right(), v, mapname...)
		case query.OR:
			return Check(n.Left(), v, mapname...) || Check(n.Right(), v, mapname...)
		}
	}
	return false
}

func compareInt(a int64, op query.Operator, b int64) bool {
	switch op {
	case query.EQ:
		return a == b
	case query.NE:
		return a != b
	case query.CO:
		return strings.Contains(strconv.Itoa(int(a)), strconv.Itoa(int(b)))
	case query.SW:
		return strings.HasPrefix(strconv.Itoa(int(a)), strconv.Itoa(int(b)))
	case query.EW:
		return strings.HasSuffix(strconv.Itoa(int(a)), strconv.Itoa(int(b)))
	case query.GT:
		return a > b
	case query.GE:
		return a >= b
	case query.LT:
		return a < b
	case query.LE:
		return a <= b
	}
	return false
}

func compareFloat(a float64, op query.Operator, b float64) bool {
	switch op {
	case query.EQ:
		return a == b
	case query.NE:
		return a != b
	case query.CO:
		return strings.Contains(fmt.Sprintf("%f", a), fmt.Sprintf("%f", b))
	case query.SW:
		return strings.HasPrefix(fmt.Sprintf("%f", a), fmt.Sprintf("%f", b))
	case query.EW:
		return strings.HasSuffix(fmt.Sprintf("%f", a), fmt.Sprintf("%f", b))
	case query.GT:
		return a > b
	case query.GE:
		return a >= b
	case query.LT:
		return a < b
	case query.LE:
		return a >= b
	}
	return false
}

func compareBool(a bool, op query.Operator, b bool) bool {
	switch op {
	case query.EQ:
		return a == b
	case query.NE:
		return a != b
	case query.CO:
		return a == b
	case query.SW:
		return a == b
	case query.EW:
		return a == b
	case query.GT:
		return a != b
	case query.GE:
		return true
	case query.LT:
		return a != b
	case query.LE:
		return true
	}
	return false
}

func compareString(a string, op query.Operator, b string) bool {
	switch op {
	case query.EQ:
		return a == b
	case query.NE:
		return a != b
	case query.CO:
		return strings.Contains(a, b)
	case query.SW:
		return strings.HasPrefix(a, b)
	case query.EW:
		return strings.HasSuffix(a, b)
	case query.GT:
		return a > b
	case query.GE:
		return a >= b
	case query.LT:
		return a < b
	case query.LE:
		return a <= b
		
	}
	return false
}