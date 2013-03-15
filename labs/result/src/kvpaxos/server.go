package kvpaxos

import "net"
import "fmt"
import "net/rpc"
import "log"
import "paxos"
import "sync"
import "os"
import "syscall"
import "encoding/gob"
import "math/rand"
import "time"

const (
  GetOp = "Get"
  PutOp = "Put"
)

type Op struct {
  // Your definitions here.
  // Field names must start with capital letters,
  // otherwise RPC will break.
  OpTyp string
  ReqNum int64
  Key string
  Value string
}

type KVPaxos struct {
  mu sync.Mutex
  l net.Listener
  me int
  dead bool // for testing
  unreliable bool // for testing
  px *paxos.Paxos

  // Your definitions here.
  // the key/value db
  mydb map[string]string 
  //historys of reqnum mapping to the value 
  history map[int64]string
  //The instance of done
  DoneInst int
}

//This function puts/gets value in mydb, and is protectd by mutex of caller
func (kv *KVPaxos) putgetmydb(arg Op) {

switch arg.OpTyp {
  case PutOp:
      kv.mydb[arg.Key] = arg.Value
      kv.history[arg.ReqNum] = ""
  case GetOp:
      val, ok := kv.mydb[arg.Key]
      if ok {
        kv.history[arg.ReqNum] = val
      } else {
        kv.history[arg.ReqNum] = ""
      }
  }
  kv.px.Done(kv.DoneInst)
  kv.DoneInst++
}

func (kv *KVPaxos) processmydb(arg *Op) (Err, string) {
  kv.mu.Lock()
  defer kv.mu.Unlock()

  for {
    //return the duplicated arg's decided value
    if val, ok := kv.history[arg.ReqNum]; ok {
      return OK, val
    }

    //If the instance has been decided
    var val Op
    decided, ret := kv.px.Status(kv.DoneInst)
    if decided {
      val = ret.(Op)
      kv.putgetmydb(val)
      continue
    }
    //The instance isn't been decided, try start it up
    kv.px.Start(kv.DoneInst, *arg)
    to := 10 * time.Millisecond
    decided, ret = kv.px.Status(kv.DoneInst)

    for !decided {
      time.Sleep(to)
      if to < 10 * time.Second {
        to *= 2
      }
      decided, ret = kv.px.Status(kv.DoneInst)
    }
    //Apply whatever instance we got
    val = ret.(Op)
    kv.putgetmydb(val)
    //Find the instance return the val, or may do it again.
    if val == *arg {
      val, ok := kv.mydb[arg.Key]
      if ok {
        return OK, val
      } else {
        return ErrNoKey, ""
      }
    }
  }
  return OK, ""
}


func (kv *KVPaxos) Get(args *GetArgs, reply *GetReply) error {
  // Your code here.
  op := new(Op)
  op.ReqNum = args.ReqNum
  op.OpTyp = GetOp
  op.Key = args.Key
  op.Value = ""

  err, value := kv.processmydb(op)
  reply.Err = err
  reply.Value = value

  return nil
}


func (kv *KVPaxos) Put(args *PutArgs, reply *PutReply) error {
  // Your code here.
  op := new(Op)
  op.ReqNum = args.ReqNum
  op.OpTyp = PutOp
  op.Key = args.Key
  op.Value = args.Value

  err, _ := kv.processmydb(op)
  reply.Err = err

  return nil
}

// tell the server to shut itself down.
// please do not change this function.
func (kv *KVPaxos) kill() {
  kv.dead = true
  kv.l.Close()
  kv.px.Kill()
}

//
// servers[] contains the ports of the set of
// servers that will cooperate via Paxos to
// form the fault-tolerant key/value service.
// me is the index of the current server in servers[].
// 
func StartServer(servers []string, me int) *KVPaxos {
  // this call is all that's needed to persuade
  // Go's RPC library to marshall/unmarshall
  // struct Op.
  gob.Register(Op{})

  kv := new(KVPaxos)
  kv.me = me

  // Your initialization code here.
  kv.mydb = map[string]string{}
  kv.history = map[int64]string{}
  kv.DoneInst = 0

  rpcs := rpc.NewServer()
  rpcs.Register(kv)

  kv.px = paxos.Make(servers, me, rpcs)

  os.Remove(servers[me])
  l, e := net.Listen("unix", servers[me]);
  if e != nil {
    log.Fatal("listen error: ", e);
  }
  kv.l = l

  // please do not change any of the following code,
  // or do anything to subvert it.

  go func() {
    for kv.dead == false {
      conn, err := kv.l.Accept()
      if err == nil && kv.dead == false {
        if kv.unreliable && (rand.Int63() % 1000) < 100 {
          // discard the request.
          conn.Close()
        } else if kv.unreliable && (rand.Int63() % 1000) < 200 {
          // process the request but force discard of reply.
          c1 := conn.(*net.UnixConn)
          f, _ := c1.File()
          err := syscall.Shutdown(int(f.Fd()), syscall.SHUT_WR)
          if err != nil {
            fmt.Printf("shutdown: %v\n", err)
          }
          go rpcs.ServeConn(conn)
        } else {
          go rpcs.ServeConn(conn)
        }
      } else if err == nil {
        conn.Close()
      }
      if err != nil && kv.dead == false {
        fmt.Printf("KVPaxos(%v) accept: %v\n", me, err.Error())
        kv.kill()
      }
    }
  }()

  return kv
}

