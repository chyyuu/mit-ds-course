package pbservice

import "net"
import "fmt"
import "net/rpc"
import "log"
import "time"
import "viewservice"
import "sync"
import "os"


type PBServer struct {
  mu sync.Mutex
  l net.Listener
  dead bool // for testing
  partitioned bool // for testing
  me string
  vs *viewservice.Clerk
  // Your declarations here.
}

func (pb *PBServer) Get(args *GetArgs, reply *GetReply) error {

  // Your code here.

  return nil
}

func (pb *PBServer) Put(args *PutArgs, reply *PutReply) error {
  reply.Err = OK


  // Your code here.

  return nil
}


//
// ping the viewserver periodically.
// if view changed:
//   transition to new view.
//   manage transfer of state from primary to new backup.
//
func (pb *PBServer) tick() {

  // Your code here.
}

// tell the server to shut itself down.
// please do not change this function.
func (pb *PBServer) kill() {
  pb.dead = true
  pb.l.Close()
}


func StartServer(vshost string, me string) *PBServer {
  pb := new(PBServer)
  pb.me = me
  pb.vs = viewservice.MakeClerk(me, vshost)
  // Your pb.* initializations here.

  rpcs := rpc.NewServer()
  rpcs.Register(pb)

  os.Remove(pb.me)
  l, e := net.Listen("unix", pb.me);
  if e != nil {
    log.Fatal("listen error: ", e);
  }
  pb.l = l

  // please do not change any of the following code.

  go func() {
    for pb.dead == false {
      conn, err := pb.l.Accept()
      if err == nil && pb.dead == false {
        go rpcs.ServeConn(conn)
      } else if err == nil {
        conn.Close()
      }
      if err != nil && pb.dead == false {
        fmt.Printf("PBServer(%v) accept: %v\n", me, err.Error())
        pb.kill()
      }
    }
  }()

  go func() {
    for pb.dead == false {
      if pb.partitioned == false {
        pb.tick()
      }
      time.Sleep(viewservice.PingInterval)
    }
  }()

  return pb
}
