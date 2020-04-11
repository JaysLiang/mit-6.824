package mr

import (
	"log"
	"sync"
	"time"
)
import "net"
import "os"
import "net/rpc"
import "net/http"

type metaData struct {
	fileName string
	state    int
	start    time.Time
}

type reduceMetaData struct {
	state int
	start time.Time
}
type Master struct {
	// Your definitions here.
	mapBuckets   map[int]metaData
	allMapTaskId []int
	reducePool   map[int]reduceMetaData
	mu           sync.Mutex
	nReduce      int
}

// Your code here -- RPC handlers for the worker to call.
func (m *Master) GetJob(req *JobReq, reply *JobReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range m.mapBuckets {
		if v.state == 1 && time.Since(v.start) < time.Second*10 {
			continue
		}
		reply.JobType = MapType
		reply.TaskId = k
		reply.Reduce = m.nReduce
		reply.MapInfo = &MapInfo{FileName: v.fileName}
		v.start = time.Now()
		v.state = 1
		m.mapBuckets[k] = v
		return nil
	}

	for k, v := range m.reducePool {
		if v.state == 1 && time.Since(v.start) < time.Second*10 {
			continue
		}
		reply.JobType = ReduceType
		reply.Reduce = m.nReduce
		reply.ReduceInfo = &ReduceInfo{
			ReduceIdx: k,
			Task:      m.allMapTaskId,
		}
		v.start = time.Now()
		v.state = 1
		m.reducePool[k] = v
		return nil
	}

	reply.JobType = IdleType
	return nil
}

func (m *Master) ReportMapDone(req *MapDoneReq, reply *EmtpyReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	md, ok := m.mapBuckets[req.TaskId]
	if !ok {
		return nil
	}
	if time.Since(md.start) > time.Second*10 {
		return nil
	}

	delete(m.mapBuckets, req.TaskId)
	m.allMapTaskId = append(m.allMapTaskId, req.TaskId)

	log.Println("ReportMapDone: ", req.TaskId)
	return nil
}

func (m *Master) ReportReduceDone(req *ReduceDoneReq, reply *EmtpyReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	md, ok := m.reducePool[req.ReduceIdx]
	if !ok {
		return nil
	}
	if time.Since(md.start) > time.Second*10 {
		return nil
	}
	delete(m.reducePool, req.ReduceIdx)
	log.Println("ReportReduceDone: ", req.ReduceIdx)
	return nil
}

//
// an example RPC handler.
//
func (m *Master) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	os.Remove("mr-socket")
	l, e := net.Listen("unix", "mr-socket")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := false

	// Your code here.
	timer := time.NewTimer(time.Second * 10)
Loop:
	for {
		select {
		case <-timer.C:
			m.mu.Lock()

			if len(m.mapBuckets) == 0 && len(m.reducePool) == 0 {
				ret = true
				m.mu.Unlock()
				break Loop
			}
			for k, v := range m.mapBuckets {
				if time.Since(v.start) > time.Second*10 {
					v.state = 0
					m.mapBuckets[k] = v
				}
			}
			for k, v := range m.reducePool {
				if time.Since(v.start) > time.Second*10 {
					v.state = 0
					m.reducePool[k] = v
				}
			}
			m.mu.Unlock()
		}
	}

	return ret
}

//
// create a Master.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}
	log.SetFlags(log.Llongfile)

	// Your code here.
	m.nReduce = nReduce
	m.mapBuckets = make(map[int]metaData)
	for idx := range files {
		m.mapBuckets[idx] = metaData{
			fileName: files[idx],
			state:    0,
		}
	}
	m.allMapTaskId = make([]int, 0)
	m.reducePool = make(map[int]reduceMetaData)
	for idx := 0; idx < nReduce; idx++ {
		m.reducePool[idx] = reduceMetaData{
			state: 0,
		}
	}

	m.server()

	return &m
}
