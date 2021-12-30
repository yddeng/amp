// MIT License
//
// Copyright (c) 2021 ydd
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package heap

import (
	"container/heap"
)

// Element is an element of a heap.
type Element interface {
	// Less reports whether the element sort before the other.
	Less(other interface{}) bool
}

// inlineHeap inline of Heap
type inlineHeap struct {
	elements []Element
	elemIdx  map[Element]int
}

// Heap
type Heap struct {
	inlineHeap *inlineHeap
}

// New returns an initialized heap.
func New() *Heap {
	h := &Heap{
		inlineHeap: &inlineHeap{
			elements: make([]Element, 0, 32),
			elemIdx:  make(map[Element]int, 32),
		},
	}
	return h
}

// Less returns the opposite of the embedded implementation's Less method.
func (h *inlineHeap) Less(i, j int) bool {
	return h.elements[i].Less(h.elements[j])
}

// Swap swaps the elements with indexes i and j.
func (h *inlineHeap) Swap(i, j int) {
	h.elemIdx[h.elements[i]] = j
	h.elemIdx[h.elements[j]] = i
	h.elements[i], h.elements[j] = h.elements[j], h.elements[i]
}

// Len is the number of elements in the heap.
func (h *inlineHeap) Len() int {
	return len(h.elements)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
func (h *inlineHeap) Pop() (v interface{}) {
	h.elements, v = h.elements[:h.Len()-1], h.elements[h.Len()-1]
	item := v.(Element)
	delete(h.elemIdx, item)
	return item
}

// Push pushes the element x onto the heap.
func (h *inlineHeap) Push(v interface{}) {
	h.elemIdx[v.(Element)] = h.Len()
	h.elements = append(h.elements, v.(Element))
}

// Len is the number of elements in the heap.
func (h *Heap) Len() int {
	return h.inlineHeap.Len()
}

// Push pushes the element e onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap) Push(e Element) {
	heap.Push(h.inlineHeap, e)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func (h *Heap) Pop() Element {
	if h.Len() > 0 {
		return heap.Pop(h.inlineHeap).(Element)
	}
	return nil
}

// Top returns the minimum element (according to Less) from the heap e or nil if the heap is empty.
func (h *Heap) Top() Element {
	if h.Len() > 0 {
		return h.inlineHeap.elements[0]
	}
	return nil
}

// IsExist returns whether the element in the heap.
func (h *Heap) IsExist(ele Element) bool {
	_, ok := h.inlineHeap.elemIdx[ele]
	return ok
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap) Remove(ele Element) {
	i, ok := h.inlineHeap.elemIdx[ele]
	if ok {
		heap.Remove(h.inlineHeap, i)
	}
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func (h *Heap) Fix(ele Element) {
	i, ok := h.inlineHeap.elemIdx[ele]
	if ok {
		heap.Fix(h.inlineHeap, i)
	}
}

// Reset resets the heap to be empty,
func (h *Heap) Reset() {
	h.inlineHeap.elements = h.inlineHeap.elements[0:0]
	h.inlineHeap.elemIdx = make(map[Element]int, cap(h.inlineHeap.elements))
}

// PushHeap inserts a copy of another heap at the heap h.
// The heap l and other may be the same. They must not be nil.
func (h *Heap) PushHeap(other *Heap) {
	for _, e := range other.inlineHeap.elements {
		h.Push(e)
	}
}

// PushList inserts a slice element elements in the heap h.
func (h *Heap) PushList(elements ...Element) {
	for _, e := range elements {
		h.Push(e)
	}
}
