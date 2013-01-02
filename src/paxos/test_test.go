package paxos

import "testing"
import "runtime"
import "strconv"
import "os"
import "time"
import "fmt"

func port(tst string, host int) string {
  s := "/var/tmp/px-"
  s += strconv.Itoa(os.Getuid()) + "-"
  s += strconv.Itoa(os.Getpid()) + "-"
  s += tst + "-"
  s += strconv.Itoa(host)
  return s
}

func ndecided(t *testing.T, pxa []*Paxos, seq int) int {
  count := 0
  var v interface{}
  for i := 0; i < len(pxa); i++ {
    if pxa[i] != nil {
      decided, v1 := pxa[i].Get(seq)
      if decided {
        if count > 0 && v != v1 {
          t.Fatalf("decided values do not match; seq=%v i=%v v=%v v1=%v",
            seq, i, v, v1)
        }
        count++
        v = v1
      }
    }
  }
  return count
}

func waitn(t *testing.T, pxa[]*Paxos, seq int, wanted int) {
  for iters := 0; iters < 50; iters++ {
    if ndecided(t, pxa, seq) >= wanted {
      break
    }
    time.Sleep(100 * time.Millisecond)
  }
  nd := ndecided(t, pxa, seq)
  if nd < wanted {
    t.Fatalf("too few decided; seq=%v ndecided=%v wanted=%v", seq, nd, wanted)
  }
}

func waitmajority(t *testing.T, pxa[]*Paxos, seq int) {
  waitn(t, pxa, seq, len(pxa) / 2)
}

func checkmax(t *testing.T, pxa[]*Paxos, seq int, max int) {
  time.Sleep(3 * time.Second)
  nd := ndecided(t, pxa, seq)
  if nd > max {
    t.Fatalf("too many decided; seq=%v ndecided=%v max=%v", seq, nd, max)
  }
}

func cleanup(pxa []*Paxos) {
  for i := 0; i < len(pxa); i++ {
    if pxa[i] != nil {
      pxa[i].kill()
    }
  }
}

func TestBasic(t *testing.T) {
  runtime.GOMAXPROCS(4)

  const npaxos = 3
  var pxa []*Paxos = make([]*Paxos, npaxos)
  var pxh []string = make([]string, npaxos)
  defer cleanup(pxa)

  v0 := "hello"
  
  for i := 0; i < npaxos; i++ {
    pxh[i] = port("basic", i)
  }
  for i := 0; i < npaxos; i++ {
    pxa[i] = Make(pxh, i)
    pxa[i].Start(0, v0)
  }

  // very first agreement

  waitn(t, pxa, 0, npaxos)

  // single proposer

  pxa[1].Start(1, 77)
  waitn(t, pxa, 1, npaxos)

  // multiple proposers, different values

  pxa[0].Start(2, 100)
  pxa[1].Start(2, 101)
  pxa[2].Start(2, 102)
  waitn(t, pxa, 2, npaxos)
}

//
// many agreements (without failures)
//
func TestMany(t *testing.T) {
  runtime.GOMAXPROCS(4)

  const npaxos = 4
  var pxa []*Paxos = make([]*Paxos, npaxos)
  var pxh []string = make([]string, npaxos)
  defer cleanup(pxa)

  v0 := 0
  
  for i := 0; i < npaxos; i++ {
    pxh[i] = port("many", i)
  }
  for i := 0; i < npaxos; i++ {
    pxa[i] = Make(pxh, i)
    pxa[i].Start(0, v0)
  }

  const ninst = 50
  for seq := 1; seq < ninst; seq++ {
    for i := 0; i < npaxos; i++ {
      pxa[i].Start(seq, (seq * 10) + i)
    }
  }

  for {
    done := true
    for seq := 1; seq < ninst; seq++ {
      if ndecided(t, pxa, seq) < npaxos {
        done = false
      }
    }
    if done {
      break
    }
    time.Sleep(100 * time.Millisecond)
  }
}

//
// basic failure / restart cases
//
func TestBasicCrash(t *testing.T) {
  runtime.GOMAXPROCS(4)

  const npaxos = 5
  var pxa []*Paxos = make([]*Paxos, npaxos)
  var pxh []string = make([]string, npaxos)
  defer cleanup(pxa)

  v0 := 0
  
  for i := 0; i < npaxos; i++ {
    pxh[i] = port("basiccrash", i)
  }
  for i := 0; i < npaxos; i++ {
    pxa[i] = Make(pxh, i)
    pxa[i].Start(0, v0)
  }

  // minority crashed
  pxa[0].kill()
  pxa[1].kill()
  pxa[2].Start(1, 101)
  waitmajority(t, pxa, 1)

  // majority crashed
  pxa[2].kill()
  pxa[3].Start(2, 102)
  checkmax(t, pxa, 2, 0)

  // restart two crashed peers
  pxa[1] = Make(pxh, 1)
  pxa[2] = Make(pxh, 2)
  pxa[3].Start(3, 103)

  // do agreements now complete?
  // XXX requiring completion is a little dubious, since peers that
  // have lost on-disk Paxos state really should not
  // let themselves restart.
  // on the other hand, in this test no crashed peer had
  // made any prepare_ok promise.
  waitmajority(t, pxa, 2)
  waitmajority(t, pxa, 3)

  // can restarted peers see past agreements?
  ok := false
  for iters := 0; iters < 50; iters++ {
    ok = true
    for seq := 0; seq <= 3; seq++ {
      if ndecided(t, pxa, seq) != 4 {
        ok = false
      }
    }
    if ok {
      break
    }
    time.Sleep(100 * time.Millisecond)
  }
  if ok == false {
    t.Fatalf("restarted peer doesn't see old agreements")
  }
}

//
// a peer starts up, with proposal, after others decide.
// then another peer starts, without a proposal.
// 
func TestOld(t *testing.T) {
  runtime.GOMAXPROCS(4)

  const npaxos = 5
  var pxa []*Paxos = make([]*Paxos, npaxos)
  var pxh []string = make([]string, npaxos)
  defer cleanup(pxa)

  for i := 0; i < npaxos; i++ {
    pxh[i] = port("old", i)
  }

  pxa[1] = Make(pxh, 1)
  pxa[2] = Make(pxh, 2)
  pxa[3] = Make(pxh, 3)
  pxa[1].Start(1, 111)

  waitmajority(t, pxa, 1)

  pxa[0] = Make(pxh, 0)
  pxa[0].Start(1, 222)

  waitn(t, pxa, 1, 4)

  if false {
    pxa[4] = Make(pxh, 4)
    waitn(t, pxa, 1, npaxos)
  }
}

//
// many agreements, with unreliable RPC
//
func TestManyUnreliable(t *testing.T) {
  runtime.GOMAXPROCS(4)

  const npaxos = 4
  var pxa []*Paxos = make([]*Paxos, npaxos)
  var pxh []string = make([]string, npaxos)
  defer cleanup(pxa)

  v0 := 0
  
  for i := 0; i < npaxos; i++ {
    pxh[i] = port("many", i)
  }
  for i := 0; i < npaxos; i++ {
    pxa[i] = Make(pxh, i)
    pxa[i].unreliable = true
    pxa[i].Start(0, v0)
  }

  const ninst = 50
  for seq := 1; seq < ninst; seq++ {
    for i := 0; i < npaxos; i++ {
      pxa[i].Start(seq, (seq * 10) + i)
    }
  }

  for {
    done := true
    for seq := 1; seq < ninst; seq++ {
      if ndecided(t, pxa, seq) < npaxos {
        done = false
      }
    }
    if done {
      break
    }
    time.Sleep(100 * time.Millisecond)
  }
}

func pp(src int, dst int) string {
  s := "/var/tmp/px-" 
  s += strconv.Itoa(os.Getuid()) + "-"
  s += strconv.Itoa(os.Getpid()) + "-"
  s += strconv.Itoa(src) + "-"
  s += strconv.Itoa(dst)
  return s
}

func part(t *testing.T, npaxos int, p1 []int, p2 []int, p3 []int) {
  for i := 0; i < npaxos; i++ {
    for j := 0; j < npaxos; j++ {
      ij := pp(i, j)
      os.Remove(ij)
    }
  }

  pa := [][]int{p1, p2, p3}
  for pi := 0; pi < len(pa); pi++ {
    p := pa[pi]
    for i := 0; i < len(p); i++ {
      for j := 0; j < len(p); j++ {
        ij := pp(p[i], p[j])
        pj := port("partition", p[j])
        err := os.Link(pj, ij)
        if err != nil {
          t.Fatalf("os.Link(%v, %v): %v\n", pj, ij, err)
        }
      }
    }
  }
}

func TestPartition(t *testing.T) {
  runtime.GOMAXPROCS(4)

  const npaxos = 5
  var pxa []*Paxos = make([]*Paxos, npaxos)
  defer cleanup(pxa)

  for i := 0; i < npaxos; i++ {
    var pxh []string = make([]string, npaxos)
    for j := 0; j < npaxos; j++ {
      if j == i {
        pxh[j] = port("partition", i)
      } else {
        pxh[j] = pp(i, j)
      }
    }
    pxa[i] = Make(pxh, i)
  }
  defer part(t, npaxos, []int{}, []int{}, []int{})

  seq := 0

  fmt.Printf("No decision if partitioned: ")

  pxa[1].Start(seq, 111)
  checkmax(t, pxa, seq, 0)
  
  fmt.Printf("OK\n")

  fmt.Printf("Decision in majority partition: ")

  part(t, npaxos, []int{0}, []int{1,2,3}, []int{4})
  time.Sleep(2 * time.Second)
  waitmajority(t, pxa, seq)

  fmt.Printf("OK\n")

  fmt.Printf("All agree after full heal: ")

  pxa[0].Start(seq, 1000) // poke them
  pxa[4].Start(seq, 1004)
  part(t, npaxos, []int{0,1,2,3,4}, []int{}, []int{})

  waitn(t, pxa, seq, npaxos)

  fmt.Printf("OK\n")

  fmt.Printf("One peer switches partitions: ")

  for iters := 0; iters < 20; iters++ {
    seq++

    part(t, npaxos, []int{0,1,2}, []int{3,4}, []int{})
    pxa[0].Start(seq, seq * 10)
    pxa[3].Start(seq, (seq * 10) + 1)
    waitmajority(t, pxa, seq)
    if ndecided(t, pxa, seq) > 3 {
      t.Fatalf("too many decided")
    }
    
    part(t, npaxos, []int{0,1}, []int{2,3,4}, []int{})
    waitn(t, pxa, seq, npaxos)
  }

  fmt.Printf("OK\n")

  fmt.Printf("One peer switches partitions, unreliable: ")

  for i := 0; i < npaxos; i++ {
    pxa[i].unreliable = true
  }

  for iters := 0; iters < 20; iters++ {
    seq++

    part(t, npaxos, []int{0,1,2}, []int{3,4}, []int{})
    pxa[0].Start(seq, seq * 10)
    pxa[3].Start(seq, (seq * 10) + 1)
    waitmajority(t, pxa, seq)
    if ndecided(t, pxa, seq) > 3 {
      t.Fatalf("too many decided")
    }
    
    part(t, npaxos, []int{0,1}, []int{2,3,4}, []int{})
    waitn(t, pxa, seq, 4)
  }

  fmt.Printf("OK\n")
}
