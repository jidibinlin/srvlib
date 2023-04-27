/**
 * @Author: ChenJunJi
 * @Desc:
 * @Date: 2021/8/27 17:18
 */

package queue_list

import (
	"log"
	"testing"
)

func Que() {
	st := NewQueueList(8)
	st.Append(1)
	st.Append(2)
	st.Append(3)
	st.Append(4)
	st.Append(5)
	st.Flush()
	st.Append(6)
	st.Append(7)
	st.Append(8)
	st.Append(9)
	st.Append(10)

	log.Printf("%v ||| %v", st.appendList, st.dataList)
}

func TestQue(t *testing.T) {
	Que()
}
