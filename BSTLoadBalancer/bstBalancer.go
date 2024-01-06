package loadBalancer

import (
	"fmt"
	"net/http/httputil"
	"sync"
)

// BSTBalancer
// We do not need to find min/max values,
// because there is no max/min balancer
// We can use BST, since it has O(log n) time complexity for searching
// all values, which is exactly what we need
type BSTBalancer struct {
	Key     int
	LastKey int
	Val     *Backend
	Left    *BSTBalancer
	Right   *BSTBalancer
}

func InitializeBSTBalancers(servers []string, proxy *httputil.ReverseProxy) *BSTBalancer {
	if len(servers) == 0 {
		return nil
	}

	bstBalancer := &BSTBalancer{
		Key: 0,
		Val: &Backend{
			URL:          servers[0],
			Alive:        true,
			mux:          sync.RWMutex{},
			ReverseProxy: proxy,
		},
	}

	for i := 1; i <= len(servers)-1; i++ {
		backend := &Backend{
			Root:         bstBalancer,
			URL:          servers[i],
			Alive:        true,
			mux:          sync.RWMutex{},
			ReverseProxy: bstBalancer.Val.ReverseProxy,
		}
		bstBalancer.Insert(bstBalancer.Key, backend)
	}

	return bstBalancer
}

func (b *BSTBalancer) Insert(key int, backend *Backend) *BSTBalancer {
	if b == nil {
		backend.Root = nil

		b = &BSTBalancer{
			Key:   key,
			Val:   backend,
			Left:  &BSTBalancer{},
			Right: &BSTBalancer{},
		}

		return b
	}

	if key < b.Key {
		return b.Left.Insert(key, backend)
	} else {
		return b.Right.Insert(key, backend)
	}
}

func (b *BSTBalancer) Search(key int) *BSTBalancer {
	if b == nil {
		return nil
	}

	if b.Key == key {
		return b
	}

	if b.Key > key {
		return b.Right.Search(key)
	} else {
		return b.Left.Search(key)
	}
}

func (b *BSTBalancer) Delete(key int) *BSTBalancer {
	if b == nil {
		return nil
	}

	if key < b.Key {
		b.Left = b.Left.Delete(key)
	} else {
		b.Right = b.Right.Delete(key)
	}

	// b.Key == key
	if b.Right == nil && b.Left == nil {
		b = nil
		return nil
	}

	if b.Left == nil {
		b = b.Right
		return b
	}

	if b.Right == nil {
		b = b.Left
		return b
	}

	smallestOnTheRight := b.Right
	for {
		if smallestOnTheRight != nil && smallestOnTheRight.Left != nil {
			smallestOnTheRight = smallestOnTheRight.Left
		} else {
			break
		}
	}

	b.Key, b.Val = smallestOnTheRight.Key, smallestOnTheRight.Val
	b.Right = b.Right.Delete(b.Key)
	return b
}

// String prints a visual representation of the tree
func (b *BSTBalancer) String() {
	b.Val.mux.Lock()
	defer b.Val.mux.Unlock()
	fmt.Println("------------------------------------------------")
	stringify(b.Val.Root, 0)
	fmt.Println("------------------------------------------------")
}

// Max returns maximal value in the tree
func (b *BSTBalancer) Max() int {
	b.Val.mux.RLock()
	defer b.Val.mux.RUnlock()

	n := b.Val.Root
	if n == nil {
		return n.Key
	}

	for {
		if n.Right == nil {
			return n.Key
		}

		n = n.Right
	}
}

// Min returns the lowest value in the tree
func (b *BSTBalancer) Min() int {
	b.Val.mux.RLock()
	defer b.Val.mux.RUnlock()

	n := b.Val.Root
	if n == nil {
		return n.Key
	}

	for {
		if n.Left == nil {
			return n.Key
		}

		n = n.Left
	}
}

// stringify internal recursive function to print a tree
func stringify(b *BSTBalancer, level int) {
	if b != nil {
		format := ""
		for i := 0; i < level; i++ {
			format += "       "
		}
		format += "---[ "
		level++
		stringify(b.Left, level)
		fmt.Printf(format+"%d\n", b.Key)
		stringify(b.Right, level)
	}
}
