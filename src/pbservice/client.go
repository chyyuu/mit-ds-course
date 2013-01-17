package pbservice

import "viewservice"
import "net/rpc"
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
// please use call() to send all RPCs, both in client.go and in server.go.
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


//
// fetch a key's value from the current primary;
// if they key has never been set, return "".
// Get() must keep trying until it either the
// primary replies with the value or the primary
// says the key doesn't exist (has never been Put().
//
func (ck *Clerk) Get(key string) string {

  // Your code here.

  return "???"
}

//
// tell the primary to update key's value.
// must keep trying until it succeeds.
//
func (ck *Clerk) Put(key string, value string) {

  // Your code here.
}
