// Copyright © 2016 Abcum Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ptree

import (
	"bytes"
	"sort"
)

// Node represents an immutable node in the radix tree which
// can be either an edge node or a leaf node.
type Node struct {
	leaf   *leaf
	edges  []edge
	prefix []byte
}

type leaf struct {
	key []byte
	val interface{}
}

type edge struct {
	label byte
	node  *Node
}

// Min returns the key and value of the minimum item in the
// subtree of the current node.
func (n *Node) Min() ([]byte, interface{}) {

	for {

		if n.isLeaf() {
			return n.leaf.key, n.leaf.val
		}

		if len(n.edges) > 0 {
			n = n.edges[0].node
		} else {
			break
		}

	}

	return nil, nil

}

// Max returns the key and value of the maximum item in the
// subtree of the current node.
func (n *Node) Max() ([]byte, interface{}) {

	for {

		if num := len(n.edges); num > 0 {
			n = n.edges[num-1].node
			continue
		}

		if n.isLeaf() {
			return n.leaf.key, n.leaf.val
		}

		break

	}

	return nil, nil

}

// Path is used to recurse over the tree only visiting nodes
// which are above this node in the tree.
func (n *Node) Path(k []byte, f Walker) {

	s := k

	for {

		if n.leaf != nil && f(n.leaf.key, n.leaf.val) {
			return
		}

		if len(s) == 0 {
			return
		}

		if _, n = n.getEdge(s[0]); n == nil {
			return
		}

		if bytes.HasPrefix(s, n.prefix) {
			s = s[len(n.prefix):]
		} else {
			break
		}

	}

}

// Subs is used to recurse over the tree only visiting nodes
// which are directly under this node in the tree.
func (n *Node) Subs(k []byte, f Walker) {

	s := k

	for {

		// Check for key exhaution
		if len(s) == 0 {
			subs(n, f, false)
			return
		}

		// Look for an edge
		if _, n = n.getEdge(s[0]); n == nil {
			break
		}

		// Consume the search prefix
		if bytes.HasPrefix(s, n.prefix) {
			s = s[len(n.prefix):]
		} else if bytes.HasPrefix(n.prefix, s) {
			subs(n, f, false)
			return
		} else {
			break
		}

	}

}

// Walk is used to recurse over the tree only visiting nodes
// which are under this node in the tree.
func (n *Node) Walk(k []byte, f Walker) {

	s := k

	for {

		// Check for key exhaution
		if len(s) == 0 {
			walk(n, f, false)
			return
		}

		// Look for an edge
		if _, n = n.getEdge(s[0]); n == nil {
			break
		}

		// Consume the search prefix
		if bytes.HasPrefix(s, n.prefix) {
			s = s[len(n.prefix):]
		} else if bytes.HasPrefix(n.prefix, s) {
			walk(n, f, false)
			return
		} else {
			break
		}

	}

}

// ------------------------------
// ------------------------------
// ------------------------------
// ------------------------------
// ------------------------------

func (n *Node) isLeaf() bool {
	return n.leaf != nil
}

func (n *Node) dup() *Node {
	d := &Node{}
	if n.leaf != nil {
		d.leaf = &leaf{}
		*d.leaf = *n.leaf
	}
	if n.prefix != nil {
		d.prefix = make([]byte, len(n.prefix))
		copy(d.prefix, n.prefix)
	}
	if len(n.edges) != 0 {
		d.edges = make([]edge, len(n.edges))
		copy(d.edges, n.edges)
	}
	return d
}

func (n *Node) addEdge(e edge) {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= e.label
	})
	n.edges = append(n.edges, e)
	if idx != num {
		copy(n.edges[idx+1:], n.edges[idx:num])
		n.edges[idx] = e
	}
}

func (n *Node) replaceEdge(e edge) {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= e.label
	})
	if idx < num && n.edges[idx].label == e.label {
		n.edges[idx].node = e.node
		return
	}
	panic("replacing missing edge")
}

func (n *Node) getEdge(label byte) (int, *Node) {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= label
	})
	if idx < num && n.edges[idx].label == label {
		return idx, n.edges[idx].node
	}
	return -1, nil
}

func (n *Node) delEdge(label byte) {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= label
	})
	if idx < num && n.edges[idx].label == label {
		copy(n.edges[idx:], n.edges[idx+1:])
		n.edges[len(n.edges)-1] = edge{}
		n.edges = n.edges[:len(n.edges)-1]
	}
}

func (n *Node) mergeChild() {
	e := n.edges[0]
	child := e.node
	n.prefix = concat(n.prefix, child.prefix)
	if child.leaf != nil {
		n.leaf = new(leaf)
		*n.leaf = *child.leaf
	} else {
		n.leaf = nil
	}
	if len(child.edges) != 0 {
		n.edges = make([]edge, len(child.edges))
		copy(n.edges, child.edges)
	} else {
		n.edges = nil
	}
}

func subs(n *Node, f Walker, sub bool) bool {

	// Visit the leaf values if any
	if sub && n.leaf != nil {
		if f(n.leaf.key, n.leaf.val) {
			return true
		}
		return false
	}

	// Recurse on the children
	for _, e := range n.edges {
		if subs(e.node, f, true) {
			return true
		}
	}

	return false

}

func walk(n *Node, f Walker, sub bool) bool {

	// Visit the leaf values if any
	if n.leaf != nil {
		if f(n.leaf.key, n.leaf.val) {
			return true
		}
	}

	// Recurse on the children
	for _, e := range n.edges {
		if walk(e.node, f, true) {
			return true
		}
	}

	return false

}

func (n *Node) get(key []byte) interface{} {

	s := key

	for {

		// Check for key exhaution
		if len(s) == 0 {
			if n.isLeaf() {
				return n.leaf.val
			}
			break
		}

		// Look for an edge
		_, n = n.getEdge(s[0])
		if n == nil {
			break
		}

		// Consume the search prefix
		if bytes.HasPrefix(s, n.prefix) {
			s = s[len(n.prefix):]
		} else {
			break
		}

	}

	return nil

}
