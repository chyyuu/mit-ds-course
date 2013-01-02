package lockservice

import "net/rpc"
import "log"

//
// the lockservice Clerk lives in the client
// and maintains a little state.
//
type Clerk struct {
  servers [2]string // primary, backup
}

func MakeClerk(primary string, backup string) *Clerk {
  ck := new(Clerk)
  ck.servers[0] = primary
  ck.servers[1] = backup
  return ck
}

//
// ask the lock service for a lock.
// returns true if the lock service
// granted the lock, false otherwise.
//
// you will have to modify this function.
//
func (ck *Clerk) Lock(lockname string) bool {
  // create a connection to the server.
  c, err := rpc.Dial("unix", ck.servers[0])
  if err != nil {
    // log.Printf("Lock Dial(%v) failed: %v\n", ck.servers[0], err)
    return false
  }
  defer c.Close()
  
  // prepare the arguments.
  args := &LockArgs{}
  args.Lockname = lockname
  var reply LockReply
  
  // send an RPC request, wait for the reply.
  err = c.Call("LockServer.Lock", args, &reply)
  if err != nil {
    // RPC-level failure
    log.Printf("Lock(%v) RPC failed: %v\n", ck.servers[0], err)
    return false
  }
  
  return reply.OK
}


//
// ask the lock service to unlock a lock.
// returns true if the lock service
// unlocked it, false otherwise.
//

func (ck *Clerk) Unlock(lockname string) bool {

  // Your code here.

  return false
}
