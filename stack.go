package arriba

// Mutable stack

type Item struct {
	value interface{}
	next *Item // Next stack item
}

type Stack struct {
	top *Item // Top item of the stack
	size int  // item count of the stack
}

// Put the item on top of the stack
func (stack *Stack) Push(value interface{}) {
	stack.top = &Item { value, stack.top }
	stack.size++
}

// Put the items on top of the stack
func (stack *Stack) PushArray(value []interface{}) {
	for item := range value {
		stack.top = &Item { item, stack.top }
		stack.size++
	}
}

// If the stack is not empty, remove the top element and return the value
// If the stack is empty, return nil
func (stack *Stack) Pop() (value interface{}) {
	if stack.size > 0 {
		value, stack.top = stack.top.value, stack.top.next
		stack.size--
		return
	}
	return nil
}

// item count of the stack
func (stack *Stack) Size() int {
	return stack.size
}
