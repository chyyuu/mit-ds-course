package kvpaxos

import "net/rpc"
import "time"

type Clerk struct {
  servers []string
  // You will have to modify this struct.
}

func MakeClerk(servers []string) *Clerk {
  ck := new(Clerk)
  ck.servers = servers
  // You'll have to add code here.
  return ck
}

//
// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
func (ck *Clerk) Get(key string) string {
  // You will have to modify this function.

  for {
    // try each known server.
    for _, srv := range ck.servers {
      c, err := rpc.Dial("unix", srv)
      if err == nil {
        defer c.Close()
        
        args := &GetArgs{}
        args.Key = key
        var reply GetReply
        err := c.Call("KVPaxos.Get", args, &reply)
        c.Close()
        if err == nil && (reply.Err == OK || reply.Err == ErrNoKey) {
          return reply.Value
        }
      }
      
    }
    time.Sleep(100 * time.Millisecond)
  }
  return ""
}

//
// set the value for a key.
// keeps trying until it succeeds.
//
func (ck *Clerk) Put(key string, value string) {
  // You will have to modify this function.

  for {
    for _, srv := range ck.servers {
      c, err := rpc.Dial("unix", srv)
      if err == nil {
        defer c.Close()
        
        args := &PutArgs{}
        args.Key = key
        args.Value = value
        var reply PutReply
        err := c.Call("KVPaxos.Put", args, &reply)
        c.Close()
        if err == nil && reply.Err == OK {
          return 
        }
      }
    }
    time.Sleep(100 * time.Millisecond)
  }
}
