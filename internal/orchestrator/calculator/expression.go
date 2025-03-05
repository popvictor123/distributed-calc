package calculator
import "fmt"

type ASTNode interface {
	String() string
}

type NumberNode struct {
	Value float64
}

func (n *NumberNode) String() string {
	return fmt.Sprintf("%g", n.Value)
}

type BinaryOpNode struct {
	Left  ASTNode
	Op    string
	Right ASTNode
}

func (n *BinaryOpNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left.String(), n.Op, n.Right.String())
}


