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
)

// Copy is a copy of a tree which can be used to apply changes to
// the radix tree. All changes are applied atomically and a new tree
// is returned when committed. A Copy is not thread safe.
type Copy struct {
	size int
	root *Node
}

// Size is used to return the total number of elements in the tree.
func (c *Copy) Size() int {
	return c.size
}

// Root returns the root node withing this copy of the radix tree.
func (c *Copy) Root() *Node {
	return c.root
}

// Tree returns a new tree with the changes committed in memory.
func (c *Copy) Tree() *Tree {
	return &Tree{c.size, c.root}
}

// Cursor returns a new cursor for iterating through the radix tree.
func (c *Copy) Cursor() *Cursor {
	return &Cursor{tree: c}
}

// Get is used to retrieve a specific key, returning the current value.
func (c *Copy) Get(key []byte) interface{} {
	return c.root.get(key)
}

// Del is used to delete a given key, returning the previous value.
func (c *Copy) Del(key []byte) interface{} {
	root, leaf := c.del(nil, c.root, key)
	if root != nil {
		c.root = root
	}
	if leaf != nil {
		c.size--
		return leaf.val
	}
	return nil
}

// Put is used to insert a specific key, returning the previous value.
func (c *Copy) Put(key []byte, val interface{}) interface{} {
	root, leaf := c.put(nil, c.root, key, key, val)
	if root != nil {
		c.root = root
	}
	if leaf == nil {
		c.size++
		return nil
	}
	return leaf.val
}

// ---------------------------------------------------------------------------

func prefix(a, b []byte) (i int) {
	for i = 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			break
		}
	}
	return
}

func concat(a, b []byte) (c []byte) {
	c = make([]byte, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return
}

func (c *Copy) del(p, n *Node, s []byte) (*Node, *leaf) {

	if len(s) == 0 {

		if !n.isLeaf() {
			return nil, nil
		}

		d := n.dup()

		// Remove the leaf node
		d.leaf = nil

		// Check if the node should be merged
		if n != c.root && len(d.edges) == 1 {
			d.mergeChild()
		}

		// Return the found node and leaf node
		return d, n.leaf

	}

	// Look for an edge
	l := s[0]
	i, e := n.getEdge(l)
	if e == nil || !bytes.HasPrefix(s, e.prefix) {
		return nil, nil
	}

	// Consume the search prefix
	s = s[len(e.prefix):]

	node, leaf := c.del(n, e, s)
	if node == nil {
		return nil, nil
	}

	// Copy this node
	d := n.dup()

	// Delete the edge if the node has no edges
	if node.leaf == nil && len(node.edges) == 0 {
		d.delEdge(l)
		if n != c.root && len(d.edges) == 1 && !d.isLeaf() {
			d.mergeChild()
		}
	} else {
		d.edges[i].node = node
	}

	return d, leaf

}

func (c *Copy) put(p, n *Node, s, k []byte, v interface{}) (*Node, *leaf) {

	if len(s) == 0 {

		d := n.dup()

		if n.isLeaf() {
			d.leaf.val = v
			return d, n.leaf
		} else {
			d.leaf = &leaf{key: k, val: v}
			return d, n.leaf
		}

	}

	// Look for the edge
	i, e := n.getEdge(s[0])

	// No edge, create one
	if e == nil {
		e := edge{
			label: s[0],
			node: &Node{
				leaf: &leaf{
					key: k,
					val: v,
				},
				prefix: s,
			},
		}
		d := n.dup()
		d.addEdge(e)
		return d, nil
	}

	// Determine longest prefix of the search key on match
	cl := prefix(s, e.prefix)

	if cl == len(e.prefix) {
		s = s[cl:]
		node, leaf := c.put(n, e, s, k, v)
		if node != nil {
			nc := n.dup()
			nc.edges[i].node = node
			return nc, leaf
		}
		return nil, leaf
	}

	// Split the node
	nc := n.dup()
	splitNode := &Node{
		prefix: s[:cl],
	}
	nc.replaceEdge(edge{
		label: s[0],
		node:  splitNode,
	})

	// Restore the existing child node
	modChild := e.dup()
	splitNode.addEdge(edge{
		label: modChild.prefix[cl],
		node:  modChild,
	})
	modChild.prefix = modChild.prefix[cl:]

	// Create a new leaf node
	leaf := &leaf{
		key: k,
		val: v,
	}

	// If the new key is a subset, add to to this node
	s = s[cl:]
	if len(s) == 0 {
		splitNode.leaf = leaf
		return nc, nil
	}

	// Create a new edge for the node
	splitNode.addEdge(edge{
		label: s[0],
		node: &Node{
			leaf:   leaf,
			prefix: s,
		},
	})

	return nc, nil

}
