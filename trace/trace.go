package trace

import (
	"github.com/995933447/simpletrace"
)

func init() {
	Ctx = &Context{
		gidToTraceIdMap: map[int64]string{},
	}
}

var Ctx *Context

type Context struct {
	gidToTraceIdMap map[int64]string
}

func (c *Context) GetCurGTrace(gid int64) (string, bool) {
	traceId, ok := c.gidToTraceIdMap[gid]
	return traceId, ok
}

func (c *Context) SetCurGTrace(gid int64, traceId string) {
	c.gidToTraceIdMap[gid] = traceId
}

func GenTraceId() string {
	return simpletrace.NewTraceId()
}
