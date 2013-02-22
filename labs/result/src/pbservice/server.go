package pbservice

import "net"
import "fmt"
import "net/rpc"
import "log"
import "time"
import "viewservice"
import "sync"
import "os"
import "syscall"
import "math/rand"


type PBServer struct {
  mu sync.Mutex
  l net.Listener
  dead bool // for testing
  unreliable bool // for testing
  me string
  vs *viewservice.Clerk
  // Your declarations here.
  myview  viewservice.View
  kv      map[string]string
}

func (pb *PBServer) Get(args *GetArgs, reply *GetReply) error {

  // Your code here.
  pb.mu.Lock()
  defer pb.mu.Unlock()
  vsview, ok := pb.vs.Get()
  for !ok {
    vsview, ok = pb.vs.Get()
    //fmt.Printf("pb.get: get vsview\n")
  }
  if pb.me == vsview.Primary {
    val, ok := pb.kv[args.Key]
    if ok {
      reply.Value = val
      reply.Err = OK
    } else {
      reply.Value = ""
      reply.Err = ErrNoKey
    }
  } else {
    //reply.Value=""
    reply.Err = ErrWrongServer
  }

  return nil

}

func (pb *PBServer) BackupPut(args *PutArgs, reply *PutReply) error {
  pb.mu.Lock()
  defer pb.mu.Unlock()
  reply.Err = OK
  pb.kv[args.Key] = args.Value
  return nil
}

func (pb *PBServer) Put(args *PutArgs, reply *PutReply) error {
  pb.mu.Lock()
  defer pb.mu.Unlock()

  reply.Err = OK
  // Your code here.
  vsview, ok := pb.vs.Get()
  for !ok {
    vsview, ok = pb.vs.Get()
    //fmt.Printf("pb.put: get vsview\n")
  }
  if pb.me == vsview.Primary {
        pb.kv[args.Key] = args.Value
        if vsview.Backup != "" {
            var bakReply PutReply
            //fmt.Printf("prisrv %t RPC backuput to baksrv %t\n",pb.myview.Primary, pb.myview.Backup) 
            // rpc call baksrv BackupPut may fail
            for !call(vsview.Backup, "PBServer.BackupPut", args, &bakReply) {
                 //fmt.Printf("pb.tick:  call bakcupput\n")
                 newview, newerr := pb.vs.Ping(pb.myview.Viewnum)
                 if newerr != nil || newview != vsview{
                    break
                 }
             }
        }
  } else {
    reply.Err = ErrWrongServer
  }
  return nil
}
func (pb *PBServer) KV(args *KVArgs, reply *KVReply) error {
  pb.mu.Lock()
  defer pb.mu.Unlock()
  pb.kv= args.KV
  //fmt.Printf("Srv %s KV %t\n",pb.me, pb.kv)
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
  pb.mu.Lock()
  defer pb.mu.Unlock()
  //var vsview viewservice.View
  //var err error
  vsview, err := pb.vs.Ping(pb.myview.Viewnum) 
  if  err != nil {
     //fmt.Printf("ERROR: pb.tick.ping: myview%t, vs view %t, err%t\n",pb.myview, vsview, err)
    return
  }

  if pb.myview != vsview {
    if pb.me == vsview.Primary && vsview.Backup != "" {
      // send pri's key/value to bak
      //fmt.Println("prisrv send KV to bak")
      //fmt.Printf("ps:tick FISRT: prisrv %s send KV to bak %s \n", vsview.Primary, vsview.Backup)
      args := &KVArgs{}
      args.KV = pb.kv
      var reply KVReply
      // rpc call backsrv KV may fail
      for !call(vsview.Backup, "PBServer.KV", args, &reply) {
          //fmt.Printf("ps:tick TRY AGAIN: prisrv %s send KV to bak %s \n", vsview.Primary, vsview.Backup)
          newview, newerr := pb.vs.Ping(pb.myview.Viewnum)
          if newerr != nil || vsview != newview {
            break
          }
      } 
    }

    pb.myview = vsview
  } 
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
  pb.myview.Viewnum = 0
  pb.kv = map[string]string{}
  //-----------------------------
  rpcs := rpc.NewServer()
  rpcs.Register(pb)

  os.Remove(pb.me)
  l, e := net.Listen("unix", pb.me);
  if e != nil {
    log.Fatal("listen error: ", e);
  }
  pb.l = l

  // please do not change any of the following code,
  // or do anything to subvert it.

  go func() {
    for pb.dead == false {
      conn, err := pb.l.Accept()
      if err == nil && pb.dead == false {
        if pb.unreliable && (rand.Int63() % 1000) < 100 {
          // discard the request.
          conn.Close()
        } else if pb.unreliable && (rand.Int63() % 1000) < 200 {
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
      if err != nil && pb.dead == false {
        fmt.Printf("PBServer(%v) accept: %v\n", me, err.Error())
        pb.kill()
      }
    }
  }()

  go func() {
    for pb.dead == false {
      pb.tick()
      time.Sleep(viewservice.PingInterval)
    }
  }()

  return pb
}
