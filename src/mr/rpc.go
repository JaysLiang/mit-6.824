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
	MapType JobType = 0
	ReduceType JobType =1
)

type MapInfo struct {
	FileName string
	ResFileName string
}

type ReduceInfo struct {
	FileName string
	ResFileName string
}

type JobReq struct {
	JobType JobType
}

type JobReply struct {
	JobType JobType
	Reduce int
	TaskId int
	MapInfo *MapInfo
	ReduceInfo *ReduceInfo
}