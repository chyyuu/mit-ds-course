package lockservice

import "net"
import "net/rpc"
import "log"
import "sync"
import "fmt"
import "os"

type LockServer struct {
  mu sync.Mutex
  l net.Listener
  dead bool
  am_primary bool
  primary string
  backup string

  // for each lock name, is it locked?
  locks map[string]bool
}


//
// server Lock RPC handler.
//
// you will have to modify this function
//
func (ls *LockServer) Lock(args *LockArgs, reply *LockReply) error {
  ls.mu.Lock()
  defer ls.mu.Unlock()


  locked, _ := ls.locks[args.Lockname]

  if locked {
    reply.OK = false
  } else {
    reply.OK = true
    ls.locks[args.Lockname] = true
  }

  return nil
}

//
// server Unlock RPC handler.
//
func (ls *LockServer) Unlock(args *UnlockArgs, reply *UnlockReply) error {

  // Your code here.

  return nil
}

//
// tell the server to shut itself down.
// for testing.
// please don't change this.
//
func (ls *LockServer) kill() {
  ls.dead = true
  ls.l.Close()
}

func StartServer(primary string, backup string, am_primary bool) *LockServer {
  ls := new(LockServer)
  ls.primary = primary
  ls.backup = backup
  ls.am_primary = am_primary
  ls.locks = map[string]bool{}

  me := ""
  if am_primary {
    me = primary
  } else {
    me = backup
  }

  // tell net/rpc about our RPC server and handlers.
  rpcs := rpc.NewServer()
  rpcs.Register(ls)

  // prepare to receive connections from clients.
  // change "unix" to "tcp" to use over a network.
  os.Remove(me) // only needed for "unix"
  l, e := net.Listen("unix", me);
  if e != nil {
    log.Fatal("listen error: ", e);
  }
  ls.l = l

  // please don't change any of the following code,
  // or do anything to subvert it.

  // create a thread to accept RPC connections from clients.
  go func() {
    for ls.dead == false {
      conn, err := ls.l.Accept()
      if err == nil && ls.dead == false {
        go rpcs.ServeConn(conn)
      } else if err == nil {
        conn.Close()
      }
      if err != nil && ls.dead == false {
        fmt.Printf("LockServer(%v) accept: %v\n", me, err.Error())
        ls.kill()
      }
    }
  }()

  return ls
}
