package mr

//
// RPC definitions.
//

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.

type JobType int

const (
	MapType    JobType = 0
	ReduceType JobType = 1
	IdleType   JobType = 9999999
)

type MapInfo struct {
	FileName    string
	ResFileName string
}

type ReduceInfo struct {
	ReduceIdx int
	Task      []int
}

type JobReq struct {
	JobType JobType
}

type JobReply struct {
	JobType    JobType
	Reduce     int
	TaskId     int
	MapInfo    *MapInfo
	ReduceInfo *ReduceInfo
}

type MapDoneReq struct {
	JobType JobType
	TaskId  int
}

type EmtpyReply struct {
}

type ReduceDoneReq struct {
	JobType   JobType
	ReduceIdx int
}
