package xnet

import (
	"fmt"
	"tcpx/xiface"
)

type Job struct {
	Task xiface.ITask
}

//工人
type Worker struct {
	//工人自己的任务通道
	jobQueue chan Job
	//通知停止通道
	quit chan bool
}

func (w *Worker) run(wpq chan chan Job) {
	go func() {
		for {
			wpq <- w.jobQueue
			select {
			case job := <-w.jobQueue:
				err := job.Task.DoTask()
				if err != nil {
					fmt.Printf("job.Task.DoTask error : %v \n", err)
				}
			case <-w.quit:
				return
			}
		}
	}()
}
func newWorker() *Worker {
	w := &Worker{
		jobQueue:make(chan Job),
		quit:make(chan bool),
	}
	return w
}

//工作池
type WorkerPool struct {
	workerLen int
	workerPoolQueue chan chan Job
	//工作池接收任务通道
	jobQueue chan Job
	//通知停止通道
	quit chan bool
}

func (wp *WorkerPool)SendTask () {
	go func() {
		for  {
			select {
			case job := <-wp.jobQueue: //接收到任务
				go func(job Job) {
					//监听空闲的工人 把任务传给工人任务通道
					jobQueue := <- wp.workerPoolQueue
					jobQueue <- job
				}(job)
			case <-wp.quit:
				return
			}
		}
	}()
}
func (wp *WorkerPool)Run()  {
	for i:=0;i<wp.workerLen;i++{
		worker := newWorker()
		worker.run(wp.workerPoolQueue)
	}
	//分发任务
	wp.SendTask()
}

func (wp *WorkerPool)Submit(task xiface.ITask)  {
	job := Job{
		Task:task,
	}
	wp.jobQueue <- job
}

func NewWorkPool( maxWorkLen int)  *WorkerPool{
	wp := &WorkerPool{
		workerLen: maxWorkLen,
		workerPoolQueue:make(chan chan Job, maxWorkLen),
		jobQueue:make(chan Job),
		quit:make(chan bool),
	}
	return wp
}
