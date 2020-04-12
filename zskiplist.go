package goRedisZskiplist

const (
	ZSKIPLIST_MAXLEVEL = 32
)

type Zskiplist struct {
	//头、尾节点
	header, tail *zskiplistNode
	//节点数量
	length uint64
	//当前最大层数
	level int
}

//创建跳跃表
func Create() *Zskiplist {
	zsl := new(Zskiplist)
	//当前层数置为1
	zsl.level = 1
	//初始化头节点（头节点不用于存储数据）
	zsl.header = new(zskiplistNode)
	zsl.header.level = make([]zskiplistLevel, ZSKIPLIST_MAXLEVEL)
	for i := 0; i < ZSKIPLIST_MAXLEVEL; i++ {
		zsl.header.level[i].forward = nil
		zsl.header.level[i].span = 0
	}
	zsl.header.backward = nil
	zsl.tail = nil
	return zsl
}
