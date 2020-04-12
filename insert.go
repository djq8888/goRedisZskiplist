package goRedisZskiplist

import (
	"math/rand"
	"time"
)

//跳跃表插入对象
func (zsl *Zskiplist) Insert(score float64, obj interface{}) {
	//update数组用于记录新节点在每一层的插入位置（即在每一层，新节点插入在哪个节点后面）
	update := make([]*zskiplistNode, ZSKIPLIST_MAXLEVEL)
	//rank数组用于记录每一层update节点和头节点的节点个数（底层节点个数），该数据用于在插入节点之后更新span
	rank := make([]uint, ZSKIPLIST_MAXLEVEL)

	x := zsl.header

	//1.生成update数组和rank数组
	for i := zsl.level - 1; i >= 0; i-- {
		//初始化该层的rank值（不用每一层都从头节点开始计数，除顶层之外，每一层都可以从上一层的update节点位置开始）
		if i == zsl.level - 1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		//如果在该层的下一个节点的评分小于新节点的评分
		//或者评分相同，但是下一个节点存储的对象小于新节点的对象（TODO：暂未实现）
		//则指针右移
		for x.level[i].forward != nil && x.level[i].forward.score <= score {
			//更新该层的rank，增加该层当前节点的span
			rank[i] += x.level[i].span
			//将该层当前节点指向该层下一个节点
			x = x.level[i].forward
		}

		//记录该层的update节点
		update[i] = x
	}

	//2.生成新节点的层数（随机）
	level := randomLevel()
	//如果生成的level比跳跃表当前最大层数大，则生成超过跳跃表当前最大层数的这些层对应的update和rank
	if level > zsl.level {
		for i := zsl.level; i < level; i++ {
			rank[i] = 0
			update[i] = zsl.header
			update[i].level[i].span = uint(zsl.length)
		}
		zsl.level = level
	}

	//3.插入新节点
	//创建新节点
	x = createNode(level, score, obj)
	//在每一层插入新节点
	for i := 0; i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		//调整新节点上一个节点的span，计算新节点的span
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = rank[0] - rank[1] + 1
	}
	//如果level小于原跳跃表的level，还需要把level上面那些层的update节点的span+1
	for i := level; i < zsl.level; i++ {
		update[i].level[i].span++
	}
	//设置新节点的上一个节点
	if update[0] == zsl.header {
		x.backward = nil
	} else {
		x.backward = update[0]
	}
	//设置新节点的下一个节点的上一个节点
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		zsl.tail = x
	}
	//更新跳跃表节点数量
	zsl.length++
}

//生成一个1-ZSKIPLIST_MAXLEVEL之间的随机数，作为节点层数
func randomLevel() int {
	level := 1
	for float32(random() & 0xFFFF) < ZSKIPLIST_P * 0xFFFF {
		level++
	}
	if level < ZSKIPLIST_MAXLEVEL {
		return level
	}
	return ZSKIPLIST_MAXLEVEL
}

//生成随机int
func random() int {
	rand.Seed(time.Now().UnixNano())
	var limit int = 0xFFFF
	return rand.Intn(limit)
}