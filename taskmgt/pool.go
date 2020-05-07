package taskmgt

import "sync"

type Pool struct {
	tq *TaskSimpleQueue
	EntryChannel chan *TaskSimpleQueue
	worker_num int
	JobsChannel chan *TaskSimpleQueue

	JobsCntRwLock sync.RWMutex
	JobsCnt int
}

func NewPool(cap int) *Pool {
	p := Pool {
		tq: NewTaskSimpleQueue(),
		EntryChannel: make(chan *TaskSimpleQueue),
		worker_num: cap,
		JobsChannel: make(chan *TaskSimpleQueue),
	}
	return &p
}

func (p *Pool) worker(workId int) {
	for taskqueue := range p.JobsChannel {
		//task.Execute()
		task := taskqueue.PopFirstTask(DeadlineFirstScheduler)
		if task != nil {
			task.RunSimulation()
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
	for tq := range p.EntryChannel {
		p.JobsChannel <- tq
	}
}

/*
func (p *Pool) Submit(t *TaskEntity) {
	p.EntryChannel <- t
}*/

func (p *Pool) SafeSubmit(task *TaskEntity) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()

	p.tq.EnqueueTask(task)
	p.JobsCntRwLock.Lock()
	p.JobsCnt ++
	p.JobsCntRwLock.Unlock()
	p.EntryChannel <- p.tq
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
	p.tq.Clear()
	return true
}

func (p *Pool) GetJobsCnt() int {
	len := 0
	p.JobsCntRwLock.RLock()
	len = p.JobsCnt
	p.JobsCntRwLock.RUnlock()
	return len
}