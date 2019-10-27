package connect

import (
	"google.golang.org/grpc"
	"sync"
	"taskAssignmentForEdge/taskmgt"
)

type Node struct {
	Maddr string
	Mport int
	Saddr string
	Sport int

	conn *grpc.ClientConn  //与master的连接
	//Tq *taskmgt.TaskQueue  //等待队列
	pool *taskmgt.Pool
}

func NewNode(addr string, port int, sport int) (* Node) {
	return &Node {
		Maddr:addr,
		Mport:port,
		Sport:sport,
	}
}

func  (no *Node) StartPool(cap int, wg *sync.WaitGroup) {
	no.pool = taskmgt.NewPool(cap)
	no.pool.Run()
	wg.Done()
}