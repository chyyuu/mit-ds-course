package kvpaxos

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
  CID int64  // client ID
  Cseq uint64 // client seq #
  Key string
  Value string
}

type PutReply struct {
  Err Err
}

type GetArgs struct {
  CID int64  // client ID
  Cseq uint64 // client seq #
  Key string
}

type GetReply struct {
  Err Err
  Value string
}
