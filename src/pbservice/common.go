package pbservice

const (
  OK = "OK"
  ErrNoKey = "ErrNoKey"
  ErrWrongServer = "ErrWrongServer"
  ErrWrongView = "ErrWrongView"
  ErrWrongState = "ErrWrongState"
  ErrUnknown = "ErrUnknown"
)
type Err string

type PutArgs struct {
  Key string
  Value string
}

type PutReply struct {
  Err Err
}

type GetArgs struct {
  Key string
}

type GetReply struct {
  Err Err
  Value string
}


// Your RPC definitions here.
