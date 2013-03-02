package paxos

//
// Paxos library, to be included in an application.
// Multiple applications will run, each including
// a Paxos peer.
//
// Manages a sequence of agreed-on values.
// The set of peers is fixed.
// Copes with network failures (partition, msg loss, &c).
// Does not store anything persistently, so cannot handle crash+restart.
//
// The application interface:
//
// px = paxos.Make(peers []string, me string)
// px.Start(seq int, v interface{}) -- start agreement on new instance
// px.Status(seq int) (decided bool, v interface{}) -- get info about an instance
// px.Done(seq int) -- ok to forget all instances <= seq
// px.Max() int -- highest instance seq known, or -1
// px.Min() int -- instances before this seq have been forgotten
//

import "net"
import "net/rpc"
import "log"
import "os"
import "syscall"
import "sync"
import "fmt"
import "math"
import "math/rand"
import "time"


type Paxos struct {
  mu sync.Mutex
  l net.Listener
  dead bool
  unreliable bool
  rpcCount int
  peers []string
  me int // index into peers[]


  // Your data here.
  seq_inst map[int]InstanceStatus   //key:seq value:instance
  maxseq int
  peers_doneseq map[string]int      //key: peers value:doneseq
  doneseq int                       //done seq num
}

type InstanceStatus struct {
  nPrepare int        // highest Prepare number seen
  nAccept int         // highest Accept number seen
  vAccept interface{} // highest Accept value seen

  completed bool
  value interface{}   // desided value
}

type PrepareArgs struct {
  Seq int
  N int
}

type PrepareReply struct {
  Ok bool
  N int
  V interface{}
  Done int
}

type AcceptArgs struct {
  Seq int
  N int
  V interface{}
}

type AcceptReply struct {
  Ok bool
  N int
  Done int
}

type DecidedArgs struct {
  Seq int
  V interface{}
}

type DecidedReply struct {
  Done int
}

//
// call() sends an RPC to the rpcname handler on server srv
// with arguments args, waits for the reply, and leaves the
// reply in reply. the reply argument should be a pointer
// to a reply structure.
//
// the return value is true if the server responded, and false
// if call() was not able to contact the server. in particular,
// the replys contents are only valid if call() returned true.
//
// you should assume that call() will time out and return an
// error after a while if it does not get a reply from the server.
//
// please use call() to send all RPCs, in client.go and server.go.
// please do not change this function.
//
func call(srv string, name string, args interface{}, reply interface{}) bool {
  c, err := rpc.Dial("unix", srv)
  if err != nil {
    err1 := err.(*net.OpError)
    if err1.Err != syscall.ENOENT && err1.Err != syscall.ECONNREFUSED {
      fmt.Printf("paxos Dial() failed: %v\n", err1)
    }
    return false
  }
  defer c.Close()
    
  err = c.Call(name, args, reply)
  if err == nil {
    return true
  }
  return false
}


//
// the application wants paxos to start agreement on
// instance seq, with proposed value v.
// Start() returns right away; the application will
// call Status() to find out if/when agreement
// is reached.
//
func (px *Paxos) Start(seq int, v interface{}) {
  // Your code here.
  if seq < px.Min() {
    //fmt.Printf("px::Start seq  %d < px.Min %d!!!!\n",seq, px.Min())
    return
  }
  px.mu.Lock()
  if _, ok := px.seq_inst[seq]; !ok {
      //fmt.Printf("px::Start create new instance[%d]\n",seq) 
    px.createInstance(seq)
  }
  px.mu.Unlock()

  go px.go_proposer(seq,v)
}

// create new instancestatus
func (px *Paxos) createInstance(seq int) {
  is := new(InstanceStatus)
  is.nPrepare = -1
  is.nAccept = -1
  is.vAccept = nil
  is.completed = false
  is.value = nil
  px.seq_inst[seq] = *is
  if seq > px.maxseq {
      //fmt.Printf("px::createInstance set new max instance %d\n", seq)
    px.maxseq = seq
  }
}

// the go routine for proposer
func (px *Paxos) go_proposer(seq int, v interface{}) {
  ok := false
  for !ok {
    ok = px.Propose(seq, v)
    time.Sleep(time.Duration(rand.Int() % 50) * time.Millisecond)
  }
}

//should uesed in locks
func (px *Paxos) uniqueHigherNumber(seq int) int {
  var maxnum int
  if (px.seq_inst[seq].nPrepare >px.seq_inst[seq].nAccept){
     maxnum = px.seq_inst[seq].nPrepare
 }else {
     maxnum = px.seq_inst[seq].nAccept
 }
 return (maxnum/len(px.peers)+1)*len(px.peers) + px.me
}

//update peers_downseq, forget the old instance
func (px *Paxos) updateDone(server string, done int) {
  px.mu.Lock()
  defer px.mu.Unlock()

  if (done > px.peers_doneseq[server]) {
    px.peers_doneseq[server] = done
    min := math.MaxInt32
    for _,v := range px.peers_doneseq {
      if v < min {
        min = v
      }
    }
    for k, _ := range px.seq_inst {
      if (k < min) {
        delete(px.seq_inst, k)
      }
    }
  }
}

// the real part of propose
func (px *Paxos) Propose(seq int, v interface{}) bool {
  if ok, _ := px.Status(seq); ok {
    //fmt.Printf("px::Propose this propose %d is done!!!\n",seq)
    return true
  }
  px.mu.Lock()
  N := px.uniqueHigherNumber(seq)
  px.mu.Unlock()
  //proposer rpc Prepare to peers 
  preArgs := &PrepareArgs{}
  preArgs.Seq = seq
  preArgs.N = N
  preAckN := 0
  preoknum := 0
  var preAckV interface{}
  preAckV = nil
  for _, val  := range px.peers {
    var preReply PrepareReply
    var ok bool
    if val == px.peers[px.me] {
       result := px.Prepare(preArgs, &preReply)
       if result == nil {
           ok = true
       }else {
           ok = false
       }
    } else {
       ok = call(val, "Paxos.Prepare", preArgs, &preReply)
    }
    if ok {
      if preReply.Ok {
        preoknum++
        if preAckN <= preReply.N && preReply.V != nil {
          preAckN = preReply.N
          preAckV = preReply.V
        }
      } else {
        px.mu.Lock()

        if preReply.N/len(px.peers) > px.seq_inst[seq].nPrepare{
          is := px.seq_inst[seq]
          is.nPrepare = preReply.N/len(px.peers)
          px.seq_inst[seq] = is
        }
        px.mu.Unlock()
      }
      px.updateDone(val, preReply.Done)
    }
  }
  if (2*preoknum < len(px.peers)) {
    //fmt.Printf("prposer NO get majority prepare %d, return false!!!!\n",preoknum)
    return false
  }
  //proposer rpc Accept to peers 
  accoknum := 0
  accArgs := &AcceptArgs{}
  accArgs.Seq = seq
  accArgs.N = N

  if (preAckV == nil) {
    accArgs.V = v
  } else {
    accArgs.V = preAckV
  }

  for _, val := range px.peers {
    var accReply AcceptReply
    var ok bool
    if val == px.peers[px.me] {
        result := px.Accept(accArgs, &accReply)
        if result == nil {
           ok = true
        }else {
           ok = false
        }
    } else {
      ok = call(val, "Paxos.Accept", accArgs, &accReply)
    }
    if ok {
      if accReply.Ok {
        accoknum++
      } else {
        px.mu.Lock()
        if accReply.N/len(px.peers) > px.seq_inst[seq].nAccept{
          is := px.seq_inst[seq]
          is.nAccept = accReply.N/len(px.peers)
          px.seq_inst[seq] = is
        }
        px.mu.Unlock()
      }
      px.updateDone(val, accReply.Done)
    }
  }

  if (2*accoknum < len(px.peers)) {
    //fmt.Printf("prposer NO get majority accept %d, return false!!!!\n",accoknum)
    return false
  }
  //proposer rpc Decided to peers 
  decArgs := &DecidedArgs{}
  decArgs.Seq = seq
  decArgs.V = accArgs.V
  for _, val := range px.peers {
    var decReply DecidedReply
    ok := true
    if val == px.peers[px.me] {
        result := px.Decided(decArgs, &decReply)
        if result == nil {
           ok = true
        }else {
           ok = false
        }
    } else {
      ok = call(val, "Paxos.Decided", decArgs, &decReply)
    }
    if ok {
      px.updateDone(val, decReply.Done)
    }
  }
  return true
}

// rpc Prepare
func (px *Paxos) Prepare(args *PrepareArgs, reply *PrepareReply) error {
  px.mu.Lock()
  defer px.mu.Unlock()
  is, ok := px.seq_inst[args.Seq]
  if !ok {
    px.createInstance(args.Seq)
    is = px.seq_inst[args.Seq]
  }
  if (args.N > is.nPrepare) {
    is.nPrepare = args.N
    reply.Ok = true
    reply.N = is.nAccept
    reply.V = is.vAccept
    px.seq_inst[args.Seq] = is
  } else {
    reply.Ok = false
    reply.N = is.nPrepare
    reply.V = nil
  }

  reply.Done = px.doneseq
  return nil
}

// rpc Accept
func (px *Paxos) Accept(args *AcceptArgs, reply *AcceptReply) error {
  px.mu.Lock()
  defer px.mu.Unlock()
  is, ok := px.seq_inst[args.Seq]
  if !ok {
    px.createInstance(args.Seq)
    is = px.seq_inst[args.Seq]
  }
  if (args.N >= is.nPrepare) {
    is.nAccept = args.N
    is.vAccept = args.V
    reply.Ok = true
    reply.N = args.N
    px.seq_inst[args.Seq] = is
  } else {
    reply.Ok = false
    reply.N = is.nPrepare
  }

  reply.Done = px.doneseq
  return nil
}

// rpc Decided
func (px *Paxos) Decided(args *DecidedArgs, reply *DecidedReply) error {
  px.mu.Lock()
  defer px.mu.Unlock()
  is, ok := px.seq_inst[args.Seq]
  if !ok {
    px.createInstance(args.Seq)
    is = px.seq_inst[args.Seq]
  }

  is.completed = true
  is.value = args.V
  px.seq_inst[args.Seq] = is

  reply.Done = px.doneseq
  return nil
}

//
// the application on this machine is done with
// all instances <= seq.
//
// see the comments for Min() for more explanation.
//
func (px *Paxos) Done(seq int) {
  // Your code here.
  px.mu.Lock()
  defer px.mu.Unlock()
  px.peers_doneseq[px.peers[px.me]] = seq
  if (seq > px.doneseq) {
    px.doneseq = seq
  }
}

//
// the application wants to know the
// highest instance sequence known to
// this peer.
//
func (px *Paxos) Max() int {
  // Your code here.
  px.mu.Lock()
  defer px.mu.Unlock()
  return px.maxseq
}

//
// Min() should return one more than the minimum among z_i,
// where z_i is the highest number ever passed
// to Done() on peer i. A peers z_i is -1 if it has
// never called Done().
//
// Paxos is required to have forgotten all information
// about any instances it knows that are < Min().
// The point is to free up memory in long-running
// Paxos-based servers.
//
// It is illegal to call Done(i) on a peer and
// then call Start(j) on that peer for any j <= i.
//
// Paxos peers need to exchange their highest Done()
// arguments in order to implement Min(). These
// exchanges can be piggybacked on ordinary Paxos
// agreement protocol messages, so it is OK if one
// peers Min does not reflect another Peers Done()
// until after the next instance is agreed to.
//
// The fact that Min() is defined as a minimum over
// *all* Paxos peers means that Min() cannot increase until
// all peers have been heard from. So if a peer is dead
// or unreachable, other peers Min()s will not increase
// even if all reachable peers call Done. The reason for
// this is that when the unreachable peer comes back to
// life, it will need to catch up on instances that it
// missed -- the other peers therefor cannot forget these
// instances.
// 
func (px *Paxos) Min() int {
  // You code here.
  px.mu.Lock()
  defer px.mu.Unlock()
  min := math.MaxInt32
  for _,v := range px.peers_doneseq {
    if v < min {
      min = v
    }
  }
  return min+1
}

//
// the application wants to know whether this
// peer thinks an instance has been decided,
// and if so what the agreed value is. Status()
// should just inspect the local peers state;
// it should not contact other Paxos peers.
//
func (px *Paxos) Status(seq int) (bool, interface{}) {
  // Your code here.
  if seq < px.Min() {
    return false, nil
  }
  px.mu.Lock()
  defer px.mu.Unlock()
  if is,ok := px.seq_inst[seq]; ok{
      return is.completed, is.value
  }
  return false, nil
}


//
// tell the peer to shut itself down.
// for testing.
// please do not change this function.
//
func (px *Paxos) Kill() {
  px.dead = true
  if px.l != nil {
    px.l.Close()
  }
}

//
// the application wants to create a paxos peer.
// the ports of all the paxos peers (including this one)
// are in peers[]. this servers port is peers[me].
//
func Make(peers []string, me int, rpcs *rpc.Server) *Paxos {
  px := &Paxos{}
  px.peers = peers
  px.me = me


  // Your initialization code here.
  px.seq_inst = map[int]InstanceStatus{}
  px.maxseq = -1
  px.peers_doneseq = map[string]int{}
  for _, val  := range px.peers {
    px.peers_doneseq[val] = -1
  }
  px.doneseq = -1

  if rpcs != nil {
    // caller will create socket &c
    rpcs.Register(px)
  } else {
    rpcs = rpc.NewServer()
    rpcs.Register(px)

    // prepare to receive connections from clients.
    // change "unix" to "tcp" to use over a network.
    os.Remove(peers[me]) // only needed for "unix"
    l, e := net.Listen("unix", peers[me]);
    if e != nil {
      log.Fatal("listen error: ", e);
    }
    px.l = l
    
    // please do not change any of the following code,
    // or do anything to subvert it.
    
    // create a thread to accept RPC connections
    go func() {
      for px.dead == false {
        conn, err := px.l.Accept()
        if err == nil && px.dead == false {
          if px.unreliable && (rand.Int63() % 1000) < 100 {
            // discard the request.
            conn.Close()
          } else if px.unreliable && (rand.Int63() % 1000) < 200 {
            // process the request but force discard of reply.
            c1 := conn.(*net.UnixConn)
            f, _ := c1.File()
            err := syscall.Shutdown(int(f.Fd()), syscall.SHUT_WR)
            if err != nil {
              fmt.Printf("shutdown: %v\n", err)
            }
            px.rpcCount++
            go rpcs.ServeConn(conn)
          } else {
            px.rpcCount++
            go rpcs.ServeConn(conn)
          }
        } else if err == nil {
          conn.Close()
        }
        if err != nil && px.dead == false {
          fmt.Printf("Paxos(%v) accept: %v\n", me, err.Error())
        }
      }
    }()
  }


  return px
}
