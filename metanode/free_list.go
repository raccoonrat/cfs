// Copyright 2018 The Chubao Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package metanode

import (
	"container/list"
	"sync"
)

type freeList struct {
	sync.RWMutex
	list *list.List
}

func newFreeList() *freeList {
	return &freeList{
		list: list.New(),
	}
}

// Pop removes the first item on the list and returns it.
func (i *freeList) Pop() (ino *Inode) {
	i.Lock()
	defer i.Unlock()
	item := i.list.Front()
	if item == nil {
		return
	}
	val := i.list.Remove(item)
	ino = val.(*Inode)
	return
}

// Push inserts a new item at the back of the list.
func (i *freeList) Push(ino *Inode) {
	i.Lock()
	defer i.Unlock()
	i.list.PushBack(ino)
}

// GetFront returns the first item on the list.
func (i *freeList) GetFront() (ino *Inode) {
	i.Lock()
	defer i.Unlock()
	item := i.list.Front()
	if item == nil {
		return
	}
	ino = item.Value.(*Inode)
	return
}

// FrontMoveToBack moves the front item to the back of the list.
func (i *freeList) FrontMoveToBack() {
	i.Lock()
	defer i.Unlock()
	item := i.list.Front()
	if item == nil {
		return
	}
	i.list.MoveToBack(item)
}
