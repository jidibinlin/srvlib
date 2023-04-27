package v2

import (
	"github.com/gzjjyz/srvlib/logger"
	"github.com/gzjjyz/srvlib/utils"
	work "github.com/gzjjyz/srvlib/worker"
	"log"
	"sync/atomic"
	"time"
)

const (
	revBatchMsgMaxWait    = time.Millisecond * 10
	loopEventProcInterval = time.Millisecond * 10
)

type Worker struct {
	stopped             atomic.Bool
	exitCh              chan bool
	loopFunc            func()
	mHdl                map[uint32]work.MsgHdlType
	msgCh               chan *work.MsgSt
	procBatchMsgMaxSize uint32
	firstMsgInLoop      *work.MsgSt
	stoppedRecvFromGate atomic.Bool
	exitGateCh          chan bool
}

func NewWorker(msgCapacity uint32, loopFunc func()) *Worker {
	if nil == loopFunc {
		log.Fatalf("Worker Loop Func IsNil")
	}
	worker := &Worker{}
	worker.exitCh = make(chan bool)
	worker.exitGateCh = make(chan bool)
	worker.loopFunc = loopFunc
	worker.msgCh = make(chan *work.MsgSt, msgCapacity)
	worker.procBatchMsgMaxSize = msgCapacity
	worker.mHdl = make(map[uint32]work.MsgHdlType)
	return worker
}

func (w *Worker) RegisterMsgHandler(msgId uint32, hdl work.MsgHdlType) {
	if nil == hdl {
		log.Fatalf("注册消息处理函数为空, 消息id=%d", msgId)
		return
	}
	_, repeat := w.mHdl[msgId]
	if repeat {
		log.Fatalf("注册消息重复, 消息id=%d", msgId)
		return
	}
	w.mHdl[msgId] = hdl
}

func (w *Worker) SendMsg(id uint32, params ...interface{}) {
	if w.stopped.Load() {
		return
	}
	st := &work.MsgSt{
		MsgId: id,
		Param: params,
	}
	w.msgCh <- st
}

func (w *Worker) SendMsgFromGate(id uint32, params ...interface{}) {
	if w.stoppedRecvFromGate.Load() {
		return
	}
	st := &work.MsgSt{
		MsgId: id,
		Param: params,
	}
	w.msgCh <- st
}

func (w *Worker) GoStart() bool {
	utils.ProtectGo(func() {
		doLoopFuncTk := time.NewTicker(loopEventProcInterval)
		defer doLoopFuncTk.Stop()
	out:
		for {
			select {
			case <-w.exitCh:
				w.stopped.Store(true)
				break out
			case <-w.exitGateCh:
				w.stoppedRecvFromGate.Store(true)
				break out
			case w.firstMsgInLoop = <-w.msgCh:
				w.loop()
				// let tick restart to calc interval
				doLoopFuncTk.Stop()
				doLoopFuncTk = time.NewTicker(loopEventProcInterval)
			case <-doLoopFuncTk.C:
				w.loop()
			}
		}
	})
	return true
}

func (w *Worker) Stop() {
	w.exitCh <- true
	w.ProcessMsg(w.FetchAndMergeBatch(nil))
}

func (w *Worker) StopGate() {
	w.exitGateCh <- true
	w.ProcessMsg(w.FetchAndMergeBatch(nil))
}

func (w *Worker) loop() {
	defer func() {
		if err := recover(); err != nil {
			logger.Stack("循环中出现错误:%v", err)
		}
	}()

	var msgList []*work.MsgSt
	if w.firstMsgInLoop != nil {
		msgList = w.FetchAndMergeBatch([]*work.MsgSt{w.firstMsgInLoop})
		w.firstMsgInLoop = nil
	} else {
		msgList = w.FetchAndMergeBatch(nil)
	}
	utils.ProtectRun(w.loopFunc)
	w.ProcessMsg(msgList)
}

func (w *Worker) ProcessMsg(msgList []*work.MsgSt) {
	for _, msg := range msgList {
		t := time.Now()
		if fn, ok := w.mHdl[msg.MsgId]; ok {
			utils.ProtectRun(func() {
				fn(msg.Param[:]...)
			})
		}
		if since := time.Since(t); since > 20*time.Millisecond {
			logger.Debug("process msg end! id:%d, cost:%v", msg.MsgId, since)
		}
	}
}

func (w *Worker) FetchAndMergeBatch(msgList []*work.MsgSt) []*work.MsgSt {
	t := time.Now()
	for {
		select {
		case msg := <-w.msgCh:
			msgList = append(msgList, msg)
			if uint32(len(msgList)) >= w.procBatchMsgMaxSize {
				goto out
			}
			if since := time.Since(t); since > revBatchMsgMaxWait {
				goto out
			}
		default:
			goto out
		}
	}
out:
	return msgList
}
