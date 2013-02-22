package pbservice

import "viewservice"
import "net/rpc"
//import "fmt"
// You'll probably need to uncomment this:
// import "time"


type Clerk struct {
  vs *viewservice.Clerk
}

func MakeClerk(vshost string, me string) *Clerk {
  ck := new(Clerk)
  ck.vs = viewservice.MakeClerk(me, vshost)
  return ck
}


//
// call() sends an RPC to the rpcname handler on server srv
// with arguments args, waits for the reply, and leaves the
// reply in reply. the reply argument should be a pointer
// to a reply structure.
//
// the return value is true if the server responded, and false
// if call() was not able to contact the server. in particular,
// the reply's contents are only valid if call() returned true.
//
// you should assume that call() will time out and return an
// error after a while if it doesn't get a reply from the server.
//
// please use call() to send all RPCs, in client.go and server.go.
// please don't change this function.
//
func call(srv string, rpcname string,
          args interface{}, reply interface{}) bool {
  c, errx := rpc.Dial("unix", srv)
  if errx != nil {
    return false
  }
  defer c.Close()
    
  err := c.Call(rpcname, args, reply)
  if err == nil {
    return true
  }
  return false
}

//
// fetch a key's value from the current primary;
// if they key has never been set, return "".
// Get() must keep trying until it either the
// primary replies with the value or the primary
// says the key doesn't exist (has never been Put().
//
func (ck *Clerk) Get(key string) string {

  // Your code here.
  args := &GetArgs{}
  args.Key = key
  var reply GetReply
  //var i uint
  for {
    //fmt.Println("ck:get call vs.pri and pb.get")
    prisrv := ck.vs.Primary()
    if prisrv == "" {
      //i++
      //if( i== 8000) {
        //fmt.Println("ck:get pri is NULL")
        //i=0
      //}
      continue
    } else {
       // fmt.Printf("ck:get prisrv %s\n",prisrv)
    }
    if call(prisrv, "PBServer.Get", args, &reply){
        if reply.Err==OK {
        //fmt.Printf("ck.get key %t, vale %t OK\n",args.Key,reply.Value)
            return reply.Value
        } //else if reply.Err==ErrNoKey{
        //fmt.Printf("ck.get NoKey\n")
       // return ""
    }
  }
  return "???"
}

//
// tell the primary to update key's value.
// must keep trying until it succeeds.
//
func (ck *Clerk) Put(key string, value string) {

  // Your code here.
  args := &PutArgs{}
  args.Key, args.Value = key, value
  var reply PutReply
//  var i uint
  for {
    //fmt.Println("ck:get call vs.pri and pb.put")
    prisrv := ck.vs.Primary();
    if prisrv == "" {
      //i++
      //if( i== 4000) {
        //fmt.Println("ck:put pri is NULL")
        //i=0
      //}
      continue
    }
    if call(prisrv, "PBServer.Put", args, &reply){
        if reply.Err==OK {
        //fmt.Printf("ck.put  %t OK\n",args) 
        return
        } else {
            //fmt.Printf("ERR: ck.put  %t \n",reply) 
        }
    }
  }
}
