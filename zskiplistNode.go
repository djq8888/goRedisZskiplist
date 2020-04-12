package goRedisZskiplist

type zskiplistLevel struct {
	//节点在该层的下一个节点
	forward *zskiplistNode
	//节点距离该层下一个节点的距离
	span uint
}

type zskiplistNode struct {
	//节点内容
	obj interface{}
	//节点分数（链表按照分数从下到大排序）
	score float64
	//上一个节点
	backward *zskiplistNode
	//该节点在各层的信息
	level []zskiplistLevel
}

func createNode(level int, score float64, obj interface{}) *zskiplistNode {
	zn := new(zskiplistNode)
	zn.level = make([]zskiplistLevel, level)
	zn.score = score
	zn.obj = obj
	return zn
}
