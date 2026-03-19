package ds

type Node[T any] struct {
	Value       T
	Left, Right *Node[T]
}
