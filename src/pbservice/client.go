package pbservice

import "viewservice"
import "fmt"
// You'll probably need to uncomment these:
// import "net/rpc"
// import "time"


type Clerk struct {
  vs *viewservice.Clerk
}

func MakeClerk(vshost string, me string) *Clerk {
  ck := new(Clerk)
  ck.vs = viewservice.MakeClerk(me, vshost)
  return ck
}


func (ck *Clerk) Get(key string) (string, error) {
  err := fmt.Errorf("pb Get(%v): unknown error", key)


  // Your code here.

  return "", err
}

func (ck *Clerk) Put(key string, value string) error {
  err := fmt.Errorf("pb Put(%v): unknown error", key)


  // Your code here.

  return err
}
