/**
 * @Author: ChenJunJi
 * @Desc:
 * @Date: 2021/8/27 17:08
 */

package lock_list

import "sync"

type Lock struct {
	sync.RWMutex
}
