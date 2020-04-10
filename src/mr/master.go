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

type metaData struct{
	fileName string
	state int
}

type Master struct {
	// Your definitions here.
	mapBuckets [][]metaData
	reduceBuckets [][]metaData
	mu sync.Mutex
	reduce int
}

// Your code here -- RPC handlers for the worker to call.
func(m* Master) GetJob(req *JobReq, reply* JobReply) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for idx:=range m.mapBuckets {
		bucket:=m.mapBuckets[idx]
		if len(bucket) ==0 {
			continue
		}
		for i:=range bucket {
			if bucket[i].state == 0 {
				reply.JobType = MapType
				reply.TaskId = 1
				reply.Reduce = m.reduce
				reply.MapInfo = &MapInfo{FileName:bucket[i].fileName}
				bucket[i].state = 1
			}
		}
	}

}

func (m* Master) ReportMapDone(req , resp) {

}

func (m *Master) ReportReduceDone(req, resp) {

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
	timer:= time.NewTimer(time.Second* 10)
	for {
		select {
		case <-timer.C:
			m.mu.Lock()

			for idx := range m.buckets {
				bucket := m.buckets[idx]
				if len(bucket) == 0 {
					continue
				}
				ret = false
				break
			}
			m.mu.Unlock()
			if ret {
				break
			}
		}
	}


	return ret
}

//
// create a Master.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	// Your code here.
	m.mapBuckets = make([][]metaData, nReduce)
	for idx:=range files {
		bucket:=m.mapBuckets[idx%nReduce]
		bucket = append(bucket, metaData{
			fileName: files[idx],
			state:    0,
		})
	}

	m.server()

	return &m
}
