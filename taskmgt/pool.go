package taskmgt

type Pool struct {
	EntryChannel chan *TaskEntity
	worker_num int
	JobsChannel chan *TaskEntity
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
		task.RunSimulation()
		//log.Printf("Task %d completed in worker ID %d", task.TaskId, workId)
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.worker_num; i++ {
		go p.worker(i)
	}

	for task := range p.EntryChannel {
		p.JobsChannel <- task
	}

	close(p.JobsChannel)
	close(p.EntryChannel)
}

func (p *Pool) Submit(t *TaskEntity) {
	p.EntryChannel <- t
}