/**
 * @Author: ChenJunJi
 * @Desc:
 * @Date: 2021/8/27 16:40
 */

package queue_list

import "jjyz/base/lock"

// QueueListSt 线程安全队列 使用切片实现
type QueueListSt struct {
	lock.LockSt
	defCap     uint32
	dataList   []interface{}
	appendList []interface{}
}

func NewQueueList(capacity uint32) *QueueListSt {
	if capacity <= 0 {
		capacity = 8
	}
	st := &QueueListSt{}
	st.dataList = make([]interface{}, 0, capacity)
	st.appendList = make([]interface{}, 0, capacity)
	return st
}

func (list *QueueListSt) Append(data interface{}) {
	list.Lock()
	defer list.Unlock()
	list.appendList = append(list.appendList, data)
}

func (list *QueueListSt) Flush() {
	list.Lock()
	defer list.Unlock()
	list.dataList = append(list.dataList, list.appendList[:]...)
	list.appendList = make([]interface{}, 0, list.defCap)
}

func (list *QueueListSt) Traverse(fn func(args interface{})) {
	if nil == fn {
		return
	}
	for _, line := range list.dataList {
		fn(line)
	}
	list.dataList = make([]interface{}, 0, list.defCap)
}
