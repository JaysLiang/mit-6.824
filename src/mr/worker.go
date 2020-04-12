package mr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)
import "log"
import "net/rpc"
import "hash/fnv"

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

func CallGetJob() JobReply {
	req := JobReq{}
	reply := JobReply{}
	s := call("Master.GetJob", &req, &reply)
	if !s {
		log.Fatal("connect master failed! exit")
	}
	return reply
}

func CallMapDone(taskId int) {
	req := MapDoneReq{
		JobType: MapType,
		TaskId:  taskId,
	}
	reply := EmtpyReply{}
	s := call("Master.ReportMapDone", &req, &reply)
	if !s {
		log.Fatal("connect master failed! exit")
	}
}

func CallReduceDone(reduceIdx int) {
	req := ReduceDoneReq{
		JobType:   MapType,
		ReduceIdx: reduceIdx,
	}
	reply := EmtpyReply{}
	s := call("Master.ReportReduceDone", &req, &reply)
	if !s {
		log.Fatal("connect master failed! exit")
	}
}

func WriteBucket(kva []KeyValue, taksId int, reduce int) {
	reduceBucket := make([][]KeyValue, reduce)
	for idx := range kva {
		bucketIdx := ihash(kva[idx].Key) % reduce
		reduceBucket[bucketIdx] = append(reduceBucket[bucketIdx], kva[idx])
	}

	for i := 0; i < reduce; i++ {
		fileName := fmt.Sprintf("mr-%d-%d", taksId, i)
		resFile, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("cannot open %v", fileName)
		}
		enc := json.NewEncoder(resFile)
		for _, kv := range reduceBucket[i] {
			enc.Encode(&kv)
		}
		resFile.Sync()
		resFile.Close()
	}
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
	log.SetFlags(log.Llongfile)
	// Your worker implementation here.
	// uncomment to send the Example RPC to the master.
	// CallExample()
	for {
		reply := CallGetJob()
		switch reply.JobType {
		case MapType:
			log.Println("get Job MapType:----:", reply.TaskId, reply.MapInfo.FileName)
			file, err := os.Open(reply.MapInfo.FileName)
			if err != nil {
				log.Fatalf("cannot open %v", reply.MapInfo.FileName)
			}
			content, err := ioutil.ReadAll(file)
			if err != nil {
				log.Fatalf("cannot read %v", reply.MapInfo.FileName)
			}
			file.Close()
			kva := mapf(reply.MapInfo.FileName, string(content))
			WriteBucket(kva, reply.TaskId, reply.Reduce)
			CallMapDone(reply.TaskId)
		case ReduceType:
			log.Println("get Job ReduceType:----:", reply.TaskId, reply.ReduceInfo.ReduceIdx)
			kva := make([]KeyValue, 0)
			oname := fmt.Sprintf("mr-out-%d", reply.ReduceInfo.ReduceIdx)
			ofile, _ := os.Create(oname)
			for _, v := range reply.ReduceInfo.Task {
				fileName := fmt.Sprintf("mr-%d-%d", v, reply.ReduceInfo.ReduceIdx)
				resFile, err := os.Open(fileName)
				if err != nil {
					log.Fatalf("cannot open %v", fileName)
				}
				dec := json.NewDecoder(resFile)
				for {
					var kv KeyValue
					if err := dec.Decode(&kv); err != nil {
						break
					}
					kva = append(kva, kv)
				}
			}
			sort.Sort(ByKey(kva))
			i := 0
			for i < len(kva) {
				j := i + 1
				for j < len(kva) && kva[j].Key == kva[i].Key {
					j++
				}
				values := []string{}
				for k := i; k < j; k++ {
					values = append(values, kva[k].Value)
				}
				output := reducef(kva[i].Key, values)

				// this is the correct format for each line of Reduce output.
				fmt.Fprintf(ofile, "%v %v\n", kva[i].Key, output)

				i = j
			}
			ofile.Close()
			CallReduceDone(reply.ReduceInfo.ReduceIdx)
		case IdleType:
			//time.Sleep(time.Millisecond * 100)
			// no do nothing
		}
	}

}

//
// example function to show how to make an RPC call to the master.
//
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	call("Master.Example", &args, &reply)

	// reply.Y should be 100.
	fmt.Printf("reply.Y %v\n", reply.Y)
}

//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	c, err := rpc.DialHTTP("unix", "mr-socket")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
