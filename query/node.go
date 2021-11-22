package query

import (
	"fmt"
	"strconv"

	"github.com/7vars/leikari"
)

type Type int 

const (
	IDENT Type = iota + 18
	INT
	FLOAT
	BOOL
	STRING
)

type Operator Type

const (
	EQ Operator = iota + 7
	NE
	CO
	SW
	EW
	GT
	GE
	LT
	LE
)

type Logical Type 

const (
	AND Logical = iota + 3
	OR
)

type Prefix Type

const (
	NOT Prefix = iota + 16
	PR
	GROUP Prefix = 23
)

type Node interface {
	Type() Type
	String() string
}

type Identifier interface {
	Node
	Name() string
}

type identifier struct {
	name string
}

func newIdentifier(name string) Identifier {
	return &identifier{
		name: name,
	}
}

func (i *identifier) Type() Type { return IDENT }

func (i *identifier) String() string { return i.Name() }

func (i *identifier) Name() string { return i.name }

type Value interface {
	Node
	Value() interface{}

	IntValue() (int64, bool)
	FloatValue() (float64, bool)
	BoolValue() (bool, bool)
	StringValue() (string, bool)
}

type value struct {
	ntype Type
	value interface{}
}

func newValue(v interface{}) (Value, error) {
	switch x := v.(type) {
	case int64:
		return newIntValue(x), nil
	case float64:
		return newFloatValue(x), nil
	case bool:
		return newBoolValue(x), nil
	case string:
		return newStringValue(x), nil
	}
	return nil, leikari.Errorf("", "unsupported type %T", v)
}

func newIntValue(i int64) Value {
	return &value{
		ntype: INT,
		value: i,
	}
}

func newFloatValue(f float64) Value {
	return &value{
		ntype: FLOAT,
		value: f,
	}
}

func newBoolValue(b bool) Value {
	return &value{
		ntype: BOOL,
		value: b,
	}
}

func newStringValue(s string) Value {
	return &value{
		ntype: STRING,
		value: s,
	}
}

func (v *value) Type() Type { return v.ntype }

func (v *value) String() string {
	switch x := v.value.(type) {
	case int64:
		return strconv.FormatInt(x, 10)
	case float64:
		return fmt.Sprintf("%v", x)
	case bool:
		return strconv.FormatBool(x)
	case string:
		return fmt.Sprintf("'%s'", x)
	}
	return ""
}

func (v *value) Value() interface{} {
	return v.value
}

func (v *value) IntValue() (int64, bool) {
	if i, ok := v.value.(int64); ok {
		return i, true
	}
	return 0, false
}

func (v *value) FloatValue() (float64, bool) {
	if f, ok := v.value.(float64); ok {
		return f, true
	}
	return 0, false
}

func (v *value) BoolValue() (bool, bool) {
	if b, ok := v.value.(bool); ok {
		return b, true
	}
	return false, false
}

func (v *value) StringValue() (string, bool) {
	if s, ok := v.value.(string); ok {
		return s, true
	}
	return "", false
}

type Condition interface {
	Node
	AND(Node) Condition
	OR(Node) Condition
}

type Comparsion interface {
	Condition
	Identifier() Identifier
	Operator() Operator
	Value() Value
}

type comparsion struct {
	identifier Identifier
	operator Operator
	value Value
}

func newComparsion(identifier Identifier, operator Operator, value Value) (Comparsion, error) {
	switch operator {
	case EQ,NE,CO,SW,EW,GT,GE,LT,LE:
		return &comparsion{
			identifier: identifier,
			operator: operator,
			value: value,
		}, nil
	default:
		return nil, leikari.Errorf("", "unsupported operator %v", operator)
	}
}

func Compare(attr string, operator Operator, v interface{}) (Comparsion, error) {
	val, err := newValue(v)
	if err != nil {
		return nil, err
	}
	return newComparsion(newIdentifier(attr), operator, val)
}

func IsEqual(attr string, v interface{}) (Comparsion, error) { return Compare(attr, EQ, v) }
func IsNotEqual(attr string, v interface{}) (Comparsion, error) { return Compare(attr, NE, v) }
func Contains(attr string, v interface{}) (Comparsion, error) { return Compare(attr, CO, v) }
func StartWith(attr string, v interface{}) (Comparsion, error) { return Compare(attr, SW, v) }
func EndWith(attr string, v interface{}) (Comparsion, error) { return Compare(attr, EW, v) }
func IsGreater(attr string, v interface{}) (Comparsion, error) { return Compare(attr, GT, v) }
func IsGreaterOrEqual(attr string, v interface{}) (Comparsion, error) { return Compare(attr, GE, v) }
func IsLess(attr string, v interface{}) (Comparsion, error) { return Compare(attr, LT, v) }
func IsLessOrEqual(attr string, v interface{}) (Comparsion, error) { return Compare(attr, LE, v) }

func (c *comparsion) Type() Type {
	return Type(c.Operator())
}

func (c *comparsion) String() string {
	op := ""
	switch c.operator {
	case EQ:
		op = "EQ"
	case NE:
		op = "NE"
	case CO:
		op = "CO"
	case SW:
		op = "SW"
	case EW:
		op = "EW"
	case GT:
		op = "GT"
	case GE:
		op = "GE"
	case LT:
		op = "LT"
	case LE:
		op = "LE"
	}
	return fmt.Sprintf("%v %v %v", c.identifier, op, c.value)
}

func (c *comparsion) AND(co Node) Condition {
	return And(c, co)
}

func (c *comparsion) OR(co Node) Condition {
	return Or(c, co)
}

func (c *comparsion) Identifier() Identifier {
	return c.identifier
}

func (c *comparsion) Operator() Operator {
	return c.operator
}

func (c *comparsion) Value() Value {
	return c.value
}

type PrefixCondition interface {
	Condition
	Prefix() Prefix
	Right() Node
}

type prefixCondition struct {
	prefix Prefix
	right Node
}

func Pr(i Identifier) PrefixCondition {
	return &prefixCondition{
		prefix: PR,
		right: i,
	}
}

func Not(n Node) PrefixCondition {
	return &prefixCondition{
		prefix: NOT,
		right: n,
	}
}

func (p *prefixCondition) Type() Type { return Type(p.Prefix()) }

func (p *prefixCondition) String() string {
	pr := ""
	switch p.Prefix() {
	case PR:
		pr = "PR"
	case NOT:
		pr = "NOT"
	}
	return fmt.Sprintf("%v %v", pr, p.right)
}

func (p *prefixCondition) AND(c Node) Condition {
	return And(p, c)
}

func (p *prefixCondition) OR(c Node) Condition {
	return Or(p, c)
}

func (p *prefixCondition) Prefix() Prefix {
	return p. prefix
}

func (p *prefixCondition) Right() Node {
	return p.right
}

type group struct {
	node Node
}

func Group(n Node) PrefixCondition {
	return &group{
		node: n,
	}
}

func (g *group) Type() Type { return Type(g.Prefix()) }

func (g *group) String() string { return fmt.Sprintf("(%v)", g.node) }

func (g *group) AND(c Node) Condition {
	return And(g, c)
}

func (g *group) OR(c Node) Condition {
	return Or(g, c)
}

func (g *group) Prefix() Prefix { return GROUP }

func (g *group) Right() Node { return g.node }

type LogicalCondition interface {
	Condition
	Logical() Logical
	Left() Node
	Right() Node
}

type logical struct {
	logical Logical
	left Node
	right Node
}

func And(left, right Node) LogicalCondition {
	return &logical{
		logical: AND,
		left: left,
		right: right,
	}
}

func Or(left, right Node) LogicalCondition {
	return &logical{
		logical: OR,
		left: left,
		right: right,
	}
}

func (l *logical) Type() Type { return Type(l.Logical()) }

func (l *logical) String() string {
	lo := "AND"
	if l.logical == OR {
		lo = "OR"
	}
	return fmt.Sprintf("%v %v %v", l.left, lo, l.right)
}

func (l *logical) AND(c Node) Condition {
	return And(l, c)
}

func (l *logical) OR(c Node) Condition {
	return Or(l, c)
}

func (l *logical) Logical() Logical {
	return l.logical
}

func (l *logical) Left() Node {
	return l.left
}

func (l *logical) Right() Node {
	return l.right
}