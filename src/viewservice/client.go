package viewservice

import "net/rpc"
import "fmt"

//
// the viewservice Clerk lives in the client
// and maintains a little state.
//
type Clerk struct {
  me string      // client's name (host:port)
  server string  // viewservice's host:port
}

func MakeClerk(me string, server string) *Clerk {
  ck := new(Clerk)
  ck.me = me
  ck.server = server
  return ck
}

func (ck *Clerk) Ping(viewnum uint) (View, error) {
  // create a connection to the server.
  c, err := rpc.Dial("unix", ck.server)
  if err != nil {
    return View{}, fmt.Errorf("Ping(%v) Dial(%v) failed: %v",
      viewnum, ck.server, err)
  }
  defer c.Close()

  // prepare the arguments.
  args := &PingArgs{}
  args.Me = ck.me
  args.Viewnum = viewnum
  var reply PingReply

  // send an RPC request, wait for the reply.
  err = c.Call("ViewServer.Ping", args, &reply)
  if err != nil {
    return View{}, fmt.Errorf("Ping(%v) failed: %v", viewnum, err)
  }

  return reply.V, nil
}

func (ck *Clerk) Get() (View, bool) {
  c, err := rpc.Dial("unix", ck.server)
  if err != nil {
    return View{}, false
  }
  defer c.Close()
  
  args := &GetArgs{}
  var reply GetReply
  err = c.Call("ViewServer.Get", args, &reply)
  if err != nil {
    fmt.Printf("Get() failed: %v", err)
    return View{}, false
  }
  return reply.View, true
}

func (ck *Clerk) Primary() string {
  v, ok := ck.Get()
  if ok {
    return v.Primary
  }
  return ""
}
