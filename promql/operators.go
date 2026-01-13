package promql

import (
	"strings"
)

// BinaryOp represents a binary operation between two expressions.
type BinaryOp struct {
	left     Expr
	op       string
	right    Expr
	on       []string
	ignoring []string
}

// String returns the binary operation as a PromQL string.
func (b *BinaryOp) String() string {
	var sb strings.Builder
	sb.WriteByte('(')
	sb.WriteString(b.left.String())
	sb.WriteByte(' ')
	sb.WriteString(b.op)

	if len(b.on) > 0 {
		sb.WriteString(" on (")
		sb.WriteString(strings.Join(b.on, ","))
		sb.WriteByte(')')
	} else if len(b.ignoring) > 0 {
		sb.WriteString(" ignoring (")
		sb.WriteString(strings.Join(b.ignoring, ","))
		sb.WriteByte(')')
	}

	sb.WriteByte(' ')
	sb.WriteString(b.right.String())
	sb.WriteByte(')')
	return sb.String()
}

// On adds vector matching on the specified labels.
func (b *BinaryOp) On(labels ...string) *BinaryOp {
	b.on = labels
	b.ignoring = nil
	return b
}

// Ignoring adds vector matching ignoring the specified labels.
func (b *BinaryOp) Ignoring(labels ...string) *BinaryOp {
	b.ignoring = labels
	b.on = nil
	return b
}

// Arithmetic operators

// Add creates an addition operation.
func Add(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "+", right: right}
}

// Sub creates a subtraction operation.
func Sub(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "-", right: right}
}

// Mul creates a multiplication operation.
func Mul(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "*", right: right}
}

// Div creates a division operation.
func Div(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "/", right: right}
}

// Mod creates a modulo operation.
func Mod(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "%", right: right}
}

// Pow creates an exponentiation operation.
func Pow(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "^", right: right}
}

// Comparison operators

// GT creates a greater-than comparison.
func GT(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: ">", right: right}
}

// LT creates a less-than comparison.
func LT(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "<", right: right}
}

// GTE creates a greater-than-or-equal comparison.
func GTE(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: ">=", right: right}
}

// LTE creates a less-than-or-equal comparison.
func LTE(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "<=", right: right}
}

// Eq creates an equality comparison.
func Eq(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "==", right: right}
}

// Neq creates a non-equality comparison.
func Neq(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "!=", right: right}
}

// Logical operators

// And creates a logical AND operation.
func And(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "and", right: right}
}

// Or creates a logical OR operation.
func Or(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "or", right: right}
}

// Unless creates a logical UNLESS operation.
func Unless(left, right Expr) *BinaryOp {
	return &BinaryOp{left: left, op: "unless", right: right}
}
