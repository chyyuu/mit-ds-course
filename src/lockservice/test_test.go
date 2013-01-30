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

func tu(t *testing.T, ck *Clerk, lockname string) {
  ck.Unlock(lockname)
}

//
// cook up a unique-ish UNIX-domain socket name
// in /var/tmp. can't use current directory since
// AFS doesn't support UNIX-domain sockets.
//
func port(suffix string) string {
  s := "/var/tmp/824-"
  s += strconv.Itoa(os.Getuid()) + "/"
  os.Mkdir(s, 0777)
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
  tu(t, ck, "a")

  tl(t, ck, "a", true)
  tl(t, ck, "b", true)
  tu(t, ck, "a")
  tu(t, ck, "b")

  tl(t, ck, "a", true)
  tl(t, ck, "a", false)
  tu(t, ck, "a")
  tu(t, ck, "a")

  p.kill()
  b.kill()

  fmt.Printf("OK\n")
}

func TestPrimaryFail1(t *testing.T) {
  fmt.Printf("Primary failure: ")
  runtime.GOMAXPROCS(4)

  phost := port("p")
  bhost := port("b")
  p := StartServer(phost, bhost, true)  // primary
  b := StartServer(phost, bhost, false) // backup

  ck := MakeClerk(phost, bhost)

  tl(t, ck, "a", true)

  tl(t, ck, "b", true)
  tu(t, ck, "b")

  tl(t, ck, "c", true)
  tl(t, ck, "c", false)

  tl(t, ck, "d", true)
  tu(t, ck, "d")
  tl(t, ck, "d", true)

  p.kill()
  
  tl(t, ck, "a", false)
  tu(t, ck, "a")

  tu(t, ck, "b")
  tl(t, ck, "b", true)

  tu(t, ck, "c")

  tu(t, ck, "d")

  b.kill()
  fmt.Printf("OK\n")
}

func TestPrimaryFail2(t *testing.T) {
  fmt.Printf("Primary failure just before successful reply: ")
  runtime.GOMAXPROCS(4)

  phost := port("p")
  bhost := port("b")
  p := StartServer(phost, bhost, true)  // primary
  b := StartServer(phost, bhost, false) // backup

  ck1 := MakeClerk(phost, bhost)
  ck2 := MakeClerk(phost, bhost)

  tl(t, ck1, "a", true)
  tl(t, ck1, "b", true)

  p.dying = true

  tl(t, ck2, "c", true)
  tl(t, ck1, "c", false)
  tu(t, ck2, "c")
  tl(t, ck1, "c", true)

  b.kill()
  fmt.Printf("OK\n")
}

func TestPrimaryFail3(t *testing.T) {
  fmt.Printf("Primary failure just before unsuccessful reply: ")
  runtime.GOMAXPROCS(4)

  phost := port("p")
  bhost := port("b")
  p := StartServer(phost, bhost, true)  // primary
  b := StartServer(phost, bhost, false) // backup

  ck1 := MakeClerk(phost, bhost)

  tl(t, ck1, "a", true)
  tl(t, ck1, "b", true)

  p.dying = true

  tl(t, ck1, "b", false)

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
  tu(t, ck, "b")

  tl(t, ck, "c", true)
  tl(t, ck, "c", false)

  tl(t, ck, "d", true)
  tu(t, ck, "d")
  tl(t, ck, "d", true)

  b.kill()
  
  tl(t, ck, "a", false)
  tu(t, ck, "a")

  tu(t, ck, "b")
  tl(t, ck, "b", true)

  tu(t, ck, "c")

  tu(t, ck, "d")

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
      locked := ! ck.Lock(lockname)
      if locked != state[xi][locknum] {
        t.Fatal("bad final state")
      }
    }
  }

  b.kill()
  fmt.Printf("OK\n")
}
