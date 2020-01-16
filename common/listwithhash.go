package common

import (
	"container/list"
	"sync"
)

type ListWithHash struct {
	lst *list.List
	hashtb map[interface{}](*list.Element)
	rwlock sync.RWMutex
}

func NewListWithHash() *ListWithHash {
	return &ListWithHash{
		lst:    list.New(),
		hashtb: make(map[interface{}](*list.Element)),
	}
}

func (lh *ListWithHash) Find(key interface{})  (res interface{}) {
	lh.rwlock.RLock()
	if e,ok := lh.hashtb[key]; ok {
		res = e.Value
	} else {
		res = nil
	}
	lh.rwlock.RUnlock()
	return
}

func (lh *ListWithHash) Add(key interface{}, value interface{}) {
	lh.rwlock.Lock()
	e := lh.lst.PushBack(value)
	lh.hashtb[key] = e
	lh.rwlock.Unlock()
}

//false for existing key
func (lh *ListWithHash) AddUnique(key interface{}, value interface{})  (res bool) {
	lh.rwlock.Lock()
	if _,ok := lh.hashtb[key]; ok {
		res = false
	} else {
		e := lh.lst.PushBack(value)
		lh.hashtb[key] = e
		res = true
	}
	lh.rwlock.Unlock()
	return  res
}

func (lh *ListWithHash) Remove(key interface{})  (v interface{}) {
	lh.rwlock.Lock()
	if e, ok := lh.hashtb[key]; ok{
		v = e.Value
		lh.lst.Remove(e)
		delete(lh.hashtb, key)
	} else {
		v = nil
	}
	lh.rwlock.Unlock()
	return v
}

func (lh *ListWithHash) GetListNum()  (len int) {
	lh.rwlock.RLock()
	len = lh.lst.Len()
	lh.rwlock.RUnlock()
	return  len
}

func (lh *ListWithHash) GetList() *list.List{
	return lh.lst
}
