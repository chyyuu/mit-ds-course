package shardmaster

// 
// Shardmaster clerk.
// Please don't change this file.
//

import "net/rpc"
import "time"

type Clerk struct {
  servers []string // shardmaster replicas
}

func MakeClerk(servers []string) *Clerk {
  ck := new(Clerk)
  ck.servers = servers
  return ck
}

//
// please use call() to send all RPCs.
// please don't change this function.
//
func call(srv string, name string, args interface{}, reply interface{}) bool {
  c, errx := rpc.Dial("unix", srv)
  if errx != nil {
    return false
  }
  defer c.Close()
    
  err := c.Call(name, args, reply)
  if err == nil {
    return true
  }
  return false
}

func (ck *Clerk) Query(num int) Config {
  for {
    // try each known server.
    for _, srv := range ck.servers {
      args := &QueryArgs{}
      args.Num = num
      var reply QueryReply
      ok := call(srv, "ShardMaster.Query", args, &reply)
      if ok {
        return reply.Config
      }
    }
    time.Sleep(100 * time.Millisecond)
  }
  return Config{}
}

func (ck *Clerk) Join(gid int64, servers []string) {
  for {
    // try each known server.
    for _, srv := range ck.servers {
      args := &JoinArgs{}
      args.GID = gid
      args.Servers = servers
      var reply JoinReply
      ok := call(srv, "ShardMaster.Join", args, &reply)
      if ok {
        return
      }
    }
    time.Sleep(100 * time.Millisecond)
  }
}

func (ck *Clerk) Leave(gid int64) {
  for {
    // try each known server.
    for _, srv := range ck.servers {
      args := &LeaveArgs{}
      args.GID = gid
      var reply LeaveReply
      ok := call(srv, "ShardMaster.Leave", args, &reply)
      if ok {
        return
      }
    }
    time.Sleep(100 * time.Millisecond)
  }
}

func (ck *Clerk) Move(shard int, gid int64) {
  for {
    // try each known server.
    for _, srv := range ck.servers {
      args := &MoveArgs{}
      args.Shard = shard
      args.GID = gid
      var reply LeaveReply
      ok := call(srv, "ShardMaster.Move", args, &reply)
      if ok {
        return
      }
    }
    time.Sleep(100 * time.Millisecond)
  }
}
