package lockservice

import "testing"
import "runtime"
import "math/rand"
import "os"
import "strconv"
import "time"
import "fmt"

func tl(t *testing.T, ck *Clerk, lockname string, expected bool) {
  x := ck.Lock(lockname)
  if x != expected {
    t.Fatalf("Lock(%v) returned %v; expected %v", lockname, x, expected)
  }
}

func tu(t *testing.T, ck *Clerk, lockname string, expected bool) {
  x := ck.Unlock(lockname)
  if x != expected {
    t.Fatalf("Unlock(%v) returned %v; expected %v", lockname, x, expected)
  }
}

//
// cook up a unique-ish UNIX-domain socket name
// in /var/tmp. can't use current directory since
// AFS doesn't support UNIX-domain sockets.
//
func port(suffix string) string {
  s := "/var/tmp/lockserver-" 
  s += strconv.Itoa(os.Getuid()) + "-"
  s += strconv.Itoa(os.Getpid()) + "-"
  s += suffix
  return s
}

func TestBasic(t *testing.T) {
  fmt.Printf("Basic lock/unlock: ")

  runtime.GOMAXPROCS(4)

  phost := port("p")
  bhost := port("b")
  p := StartServer(phost, bhost, true)  // primary
  b := StartServer(phost, bhost, false) // backup

  ck := MakeClerk(phost, bhost)

  tl(t, ck, "a", true)
  tu(t, ck, "a", true)

  tl(t, ck, "a", true)
  tl(t, ck, "b", true)
  tu(t, ck, "a", true)
  tu(t, ck, "b", true)

  tl(t, ck, "a", true)
  tl(t, ck, "a", false)
  tu(t, ck, "a", true)
  tu(t, ck, "a", false)

  p.kill()
  b.kill()

  fmt.Printf("OK\n")
}

func TestPrimaryFail(t *testing.T) {
  fmt.Printf("Primary failure: ")
  runtime.GOMAXPROCS(4)

  phost := port("p")
  bhost := port("b")
  p := StartServer(phost, bhost, true)  // primary
  b := StartServer(phost, bhost, false) // backup

  ck := MakeClerk(phost, bhost)

  tl(t, ck, "a", true)

  tl(t, ck, "b", true)
  tu(t, ck, "b", true)

  tl(t, ck, "c", true)
  tl(t, ck, "c", false)

  tl(t, ck, "d", true)
  tu(t, ck, "d", true)
  tl(t, ck, "d", true)

  p.kill()
  
  tl(t, ck, "a", false)
  tu(t, ck, "a", true)

  tu(t, ck, "b", false)
  tl(t, ck, "b", true)

  tu(t, ck, "c", true)

  tu(t, ck, "d", true)

  b.kill()
  fmt.Printf("OK\n")
}

func TestBackupFail(t *testing.T) {
  fmt.Printf("Backup failure: ")
  runtime.GOMAXPROCS(4)

  phost := port("p")
  bhost := port("b")
  p := StartServer(phost, bhost, true)  // primary
  b := StartServer(phost, bhost, false) // backup

  ck := MakeClerk(phost, bhost)

  tl(t, ck, "a", true)

  tl(t, ck, "b", true)
  tu(t, ck, "b", true)

  tl(t, ck, "c", true)
  tl(t, ck, "c", false)

  tl(t, ck, "d", true)
  tu(t, ck, "d", true)
  tl(t, ck, "d", true)

  b.kill()
  
  tl(t, ck, "a", false)
  tu(t, ck, "a", true)

  tu(t, ck, "b", false)
  tl(t, ck, "b", true)

  tu(t, ck, "c", true)

  tu(t, ck, "d", true)

  p.kill()
  fmt.Printf("OK\n")
}

func TestMany(t *testing.T) {
  fmt.Printf("Multiple clients with primary failure: ")
  runtime.GOMAXPROCS(4)

  phost := port("p")
  bhost := port("b")
  p := StartServer(phost, bhost, true)  // primary
  b := StartServer(phost, bhost, false) // backup

  const nclients = 2
  const nlocks = 10
  done := false
  var state [nclients][nlocks]bool
  var acks [nclients]bool

  for xi := 0; xi < nclients; xi++ {
    go func(i int){
      ck := MakeClerk(phost, bhost)
      rr := rand.New(rand.NewSource(int64(os.Getpid()+i)))
      for done == false {
        locknum := (rr.Int() % nlocks)
        lockname := strconv.Itoa(locknum + (i * 1000))
        what := rr.Int() % 2
        if what == 0 {
          ck.Lock(lockname)
          state[i][locknum] = true
        } else {
          ck.Unlock(lockname)
          state[i][locknum] = false
        }
      }
      acks[i] = true
    }(xi)
  }

  time.Sleep(2 * time.Second)
  p.kill()
  time.Sleep(2 * time.Second)
  done = true
  time.Sleep(time.Second)
  ck := MakeClerk(phost, bhost)
  for xi := 0; xi < nclients; xi++ {
    if acks[xi] == false {
      t.Fatal("one client didn't complete")
    }
    for locknum := 0; locknum < nlocks; locknum++ {
      lockname := strconv.Itoa(locknum + (xi * 1000))
      locked := ck.Unlock(lockname)
      if locked != state[xi][locknum] {
        t.Fatal("bad final state")
      }
    }
  }

  b.kill()
  fmt.Printf("OK\n")
}

func TestConcurrentCounts(t *testing.T) {
  fmt.Printf("Multiple clients, single lock, primary failure: ")
  runtime.GOMAXPROCS(4)

  phost := port("p")
  bhost := port("b")
  p := StartServer(phost, bhost, true)  // primary
  b := StartServer(phost, bhost, false) // backup

  const nclients = 2
  const nlocks = 1
  done := false
  var acks [nclients]bool
  var locks [nclients][nlocks] int
  var unlocks [nclients][nlocks] int

  for xi := 0; xi < nclients; xi++ {
    go func(i int){
      ck := MakeClerk(phost, bhost)
      rr := rand.New(rand.NewSource(int64(os.Getpid()+i)))
      for done == false {
        locknum := rr.Int() % nlocks
        lockname := strconv.Itoa(locknum)
        what := rr.Int() % 2
        if what == 0 {
          if ck.Lock(lockname) {
            locks[i][locknum]++
          }
        } else {
          if ck.Unlock(lockname) {
            unlocks[i][locknum]++
          }
        }
      }
      acks[i] = true
    }(xi)
  }

  time.Sleep(2 * time.Second)
  p.kill()
  time.Sleep(2 * time.Second)
  done = true
  time.Sleep(time.Second)
  for xi := 0; xi < nclients; xi++ {
    if acks[xi] == false {
      t.Fatal("one client didn't complete")
    }
  }
  ck := MakeClerk(phost, bhost)
  for locknum := 0; locknum < nlocks; locknum++ {
    nl := 0
    nu := 0
    for xi := 0; xi < nclients; xi++ {
      nl += locks[xi][locknum]
      nu += unlocks[xi][locknum]
    }
    locked := ck.Unlock(strconv.Itoa(locknum))
    // fmt.Printf("lock=%d nl=%d nu=%d locked=%v\n",
    //   locknum, nl, nu, locked)
    if nl < nu || nl > nu + 1 {
      t.Fatal("lock race 1")
    }
    if nl == nu && locked != false {
      t.Fatal("lock race 2")
    }
    if nl != nu && locked != true {
      t.Fatal("lock race 3")
    }
  }

  b.kill()
  fmt.Printf("OK\n")
}
