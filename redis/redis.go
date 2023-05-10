package redis

import "github.com/995933447/redisgroup"

type NodeConf struct {
	Host     string
	Port     int
	Password string
	DB       int
}

var redisGrp *redisgroup.Group

func Init(nodeCfgs []*NodeConf) {
	var nodes []*redisgroup.Node
	for _, cfg := range nodeCfgs {
		nodes = append(nodes, redisgroup.NewNodeV2(cfg.Host, cfg.Port, cfg.Password, cfg.DB))
	}
	redisGrp = redisgroup.NewGroup(nodes, Logger)
}

func MustRedisGroup() *redisgroup.Group {
	if redisGrp == nil {
		panic("redisGroup not init")
	}
	return redisGrp
}
