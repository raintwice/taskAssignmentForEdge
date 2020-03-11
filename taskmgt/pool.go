package taskmgt

import (
	"sync"
	"taskAssignmentForEdge/common"
)

type Pool struct {
	EntryChannel chan *TaskEntity
	worker_num int
	JobsChannel chan *TaskEntity

	JobsCntRwLock sync.RWMutex
	JobsCnt int
}

func NewPool(cap int) *Pool {
	p := Pool {
		EntryChannel: make(chan *TaskEntity),
		worker_num: cap,
		JobsChannel: make(chan *TaskEntity),
	}
	return &p
}

func (p *Pool) worker(workId int) {
	for task := range p.JobsChannel {
		//task.Execute()
		if common.EvalType == common.SIMULATION {
			task.RunSimulation()
		} else {
			task.RunEvaluation()
		}
		//log.Printf("Task %d completed in worker ID %d", task.TaskId, workId)
		p.JobsCntRwLock.Lock()
		p.JobsCnt--
		p.JobsCntRwLock.Unlock()
	}
}

/*
func (p *Pool) Run() {
	for i := 0; i < p.worker_num; i++ {
		go p.worker(i)
	}

	for task := range p.EntryChannel {
		p.JobsChannel <- task
	}

	close(p.JobsChannel)
	close(p.EntryChannel)
}*/


func (p *Pool) SafeRun() {
	for i := 0; i < p.worker_num; i++ {
		go p.worker(i)
	}

	defer func() {
		if recover() != nil {

		}
	}()
	for task := range p.EntryChannel {
		p.JobsChannel <- task
	}
}

/*
func (p *Pool) Submit(t *TaskEntity) {
	p.EntryChannel <- t
}*/

func (p *Pool) SafeSubmit(t *TaskEntity) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()
	p.JobsCntRwLock.Lock()
	p.JobsCnt ++
	p.JobsCntRwLock.Unlock()
	p.EntryChannel <- t
	return  false
}

func (p *Pool) SafeStop() (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()
	close(p.EntryChannel)
	close(p.JobsChannel)
	return true
}

func (p *Pool) GetJobsCnt() int {
	len := 0
	p.JobsCntRwLock.RLock()
	len = p.JobsCnt
	p.JobsCntRwLock.RUnlock()
	return len
}