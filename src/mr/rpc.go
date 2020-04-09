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
}

type ReduceInfo struct {
	FileName string
}

type JobReq struct {
	JobType JobType
	MapInfo *MapInfo
	ReduceInfo* ReduceInfo
}

type JobReply struct {

}