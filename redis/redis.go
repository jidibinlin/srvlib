package redis

import "github.com/995933447/redisgroup"

type NodeConf struct {
	host     string
	port     int
	password string
}

var redisGrp *redisgroup.Group

func Init(nodeCfgs []*NodeConf) {
	var nodes []*redisgroup.Node
	for _, cfg := range nodeCfgs {
		nodes = append(nodes, redisgroup.NewNode(cfg.host, cfg.port, cfg.password))
	}
	redisGrp = redisgroup.NewGroup(nodes, Logger)
}

func MustRedisGroup() *redisgroup.Group {
	if redisGrp == nil {
		panic("redisGroup not init")
	}
	return redisGrp
}
