package main

import "fmt"

const MAX_LEN = 5

type Node struct {
	Val   string
	Left  *Node
	Right *Node
}

type Queue struct {
	Head *Node
	Tail *Node
	Size int
}

type Cache struct {
	Queue Queue
	Hash  Hash
}

func NewCache() Cache {
	return Cache{Queue: NewQueue(), Hash: Hash{}}
}

func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}
	head.Right = tail
	tail.Left = head
	return Queue{head, tail, 0}
}

type Hash map[string]*Node

func (c *Cache) Check(word string) {
	node := &Node{}

	if val, ok := c.Hash[word]; ok {
		node = c.Remove(val)
	} else {
		node = &Node{Val: word}
	}
	c.Add(node)
	c.Hash[word] = node
}

func (c *Cache) Remove(node *Node) *Node {
	fmt.Printf("Removing %s\n", node.Val)
	left := node.Left
	right := node.Right
	left.Right = right
	right.Left = left
	c.Queue.Size--
	delete(c.Hash, node.Val)
	return node
}

func (c *Cache) Add(node *Node) {
	fmt.Printf("Adding %s\n", node.Val)
	temp := c.Queue.Head.Right
	c.Queue.Head.Right = node
	node.Left = c.Queue.Head
	node.Right = temp
	temp.Left = node
	c.Queue.Size++
	if c.Queue.Size > MAX_LEN {
		c.Remove(c.Queue.Tail.Left)
	}
}

func (c *Cache) Display() {
	c.Queue.Display()
}

func (q *Queue) Display() {
	node := q.Head.Right
	fmt.Printf("%d - [", q.Size)
	for i := 0; i < q.Size; i++ {
		fmt.Printf("{%s}", node.Val)
		if i < q.Size-1 {
			fmt.Printf("<--->")
		}
		node = node.Right
	}
	fmt.Printf("]\n")
}

func main() {
	fmt.Println("Cache Start")
	cache := NewCache()
	for _, word := range []string{"parrot", "test", "a", "b", "c", "d", "d", "parrot"} {
		cache.Check(word)
		cache.Display()
	}
}
