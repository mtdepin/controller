package cache

import (
	"errors"
	"fmt"
	"sync"
)

type Node struct {
	pre  *Node
	next *Node
	key  interface{}
}

type Cache struct {
	keysMap map[interface{}]*Node
	head    *Node
	last    *Node
	size    int32
	capture int32

	lock *sync.RWMutex
}

func (p *Cache) InitCache(size int32) error {
	if size < 1 {
		return errors.New("size must >  0")
	}
	p.keysMap = make(map[interface{}]*Node, size)
	p.head = &Node{}
	p.last = nil
	p.capture = 0
	p.size = size

	p.lock = new(sync.RWMutex)
	return nil
}

//1.add
func (p *Cache) Add(key interface{}) {
	defer p.lock.RUnlock()
	p.lock.RLock()

	node, ok := p.keysMap[key]
	if ok { //存在
		if p.capture <= 1 { //1个节点不用调整.
			return
		}

		//如果命中末尾节点，末尾节点指向前一个节点.
		if p.last == node { //末尾节点.
			p.last = p.last.pre
			p.last.next = nil
			node.pre = nil
		} else { //中间节点.
			node.pre.next = node.next
			node.next.pre = node.pre
		}

		p.head.next.pre = node
		node.next = p.head.next
		p.head.next = node
		node.pre = p.head
	} else {
		p.addNewNode(key)
	}
}

func (p *Cache) addNewNode(key interface{}) {
	if p.capture < p.size {
		node := &Node{key: key}
		tmp := p.head.next
		p.head.next = node
		node.pre = p.head
		if tmp != nil {
			tmp.pre = node
			node.next = tmp
		}
		//第一个节点
		if p.last == nil {
			p.last = node
		}
		//p.last = p.last.next

		p.keysMap[key] = node

		p.capture++

	} else { //删除末尾节点.
		delete(p.keysMap, p.last.key)
		//复用末尾节点，减少内存分配.
		tmp := p.last
		tmp.key = key

		if p.capture > 1 {
			p.last = p.last.pre //
			p.last.next = nil

			tmp.next = p.head.next
			p.head.next.pre = tmp
			tmp.pre = p.head
			p.head.next = tmp
		}

		p.keysMap[key] = tmp
	}
}

func (p *Cache) Delete(key interface{}) {
	defer p.lock.RUnlock()
	p.lock.RLock()

	node, ok := p.keysMap[key]
	if !ok {
		return
	}

	delete(p.keysMap, node.key)
	p.capture--

	if p.last == node {
		if p.last.pre == p.head {
			p.head.next = nil
			p.last.pre = nil
			p.last = nil
		} else { //删除末尾
			p.last = p.last.pre
			p.last.next = nil
			node.pre = nil
			node = nil
		}
	} else {
		node.pre.next = node.next
		node.next.pre = node.pre
		node = nil
	}
}

func (p *Cache) Search(key interface{}) bool {
	defer p.lock.RUnlock()
	p.lock.RLock()
	_, ok := p.keysMap[key]
	if ok { //放到队列头。
		p.Add(key)
	}
	return ok
}

func (p *Cache) Print() {
	fmt.Printf("\nsize: %d, capture: %d\n", p.size, p.capture)
	tmp := p.head.next
	for tmp != nil {
		fmt.Printf("%v\t", tmp.key)
		tmp = tmp.next
	}
	if p.last != nil {
		fmt.Printf("-----last: %d\n", p.last.key)
	} else {
		fmt.Printf("-----last: nil\n")
	}
}
