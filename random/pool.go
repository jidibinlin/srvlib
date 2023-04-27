package random

// Item 权重项
type Item struct {
	weight uint32
	item   interface{}
}

// Pool 权重池
type Pool struct {
	totalWeight uint32
	list        []*Item
}

func (pool *Pool) Clear() {
	pool.totalWeight = 0
	pool.list = nil
}

func (pool *Pool) Size() int {
	return len(pool.list)
}

// AddItem 添加到权重池
func (pool *Pool) AddItem(item interface{}, weight uint32) {
	pool.totalWeight += weight
	pool.list = append(pool.list, &Item{weight: weight, item: item})
}

// RandomOne 从权重池随机一项
func (pool *Pool) RandomOne() interface{} {
	if pool.totalWeight <= 0 {
		return nil
	}

	value := UintU(pool.totalWeight)
	total := pool.totalWeight
	for _, line := range pool.list {
		total -= line.weight
		if total < value {
			return line.item
		}
	}
	return nil
}

// RandomMany 从权重池中随机多项
func (pool *Pool) RandomMany(count uint32) (ret []interface{}) {
	if pool.totalWeight <= 0 {
		return nil
	}

	tail := len(pool.list)
	if count >= uint32(tail) {
		for _, line := range pool.list {
			ret = append(ret, line.item)
		}
		return ret
	}

	total := pool.totalWeight
	for i := uint32(1); i <= count; i++ {
		temp := total
		value := UintU(total)
		for j := 0; j < tail; j++ {
			one := pool.list[j]
			temp -= one.weight
			if temp < value {
				ret = append(ret, one.item)
				total -= one.weight
				pool.list[j], pool.list[tail-1] = pool.list[tail-1], pool.list[j]
				tail -= 1
				break
			}
		}
	}
	return
}
