package pbservice

import "viewservice"
import "fmt"
import "testing"
import "time"
import "log"
import "runtime"
import "math/rand"
import "os"
import "net/rpc"
import "strconv"

func check(ck *Clerk, key string, value string) {
  v, err := ck.Get(key)
  if v != value {
    log.Fatalf("Get(%v) -> %v, expected %v, err %v", key, v, value, err)
  }
}

func port(suffix string) string {
  s := "/var/tmp/pb-" 
  s += strconv.Itoa(os.Getuid()) + "-"
  s += strconv.Itoa(os.Getpid()) + "-"
  s += suffix
  return s
}

func TestBasicFail(t *testing.T) {
  runtime.GOMAXPROCS(4)

  vshost := port("v")
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)

  ck := MakeClerk(vshost, "")
  _, err := ck.Get("aaa")
  if err == nil {
    t.Fatal("expected nil")
  }

  // single primary, no backup
  fmt.Printf("Single primary, no backup: ")

  s1 := StartServer(vshost, port("1"))

  deadtime := viewservice.PingInterval * viewservice.DeadPings
  time.Sleep(deadtime * 2)
  if vck.Primary() != s1.me {
    t.Fatal("first primary never formed view")
  }
  
  ck.Put("111", "v1")
  check(ck, "111", "v1")

  ck.Put("2", "v2")
  check(ck, "2", "v2")

  ck.Put("1", "v1a")
  check(ck, "1", "v1a")

  fmt.Printf("OK\n")

  // add a backup

  fmt.Printf("Add a backup: ")

  s2 := StartServer(vshost, port("2"))
  for i := 0; i < viewservice.DeadPings * 2; i++ {
    v, _ := vck.Get()
    if v.Backup == s2.me {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  v, _ := vck.Get()
  if v.Backup != s2.me {
    t.Fatal("backup never came up")
  }

  ck.Put("3", "33")
  check(ck, "3", "33")

  // give the backup time to initialize
  time.Sleep(3 * viewservice.PingInterval)

  ck.Put("4", "44")
  check(ck, "4", "44")

  fmt.Printf("OK\n")

  // kill the primary

  fmt.Printf("Primary failure: ")

  s1.kill()
  for i := 0; i < viewservice.DeadPings * 2; i++ {
    v, _ := vck.Get()
    if v.Primary == s2.me {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  v, _ = vck.Get()
  if v.Primary != s2.me {
    t.Fatal("backup never switched to primary")
  }

  check(ck, "1", "v1a")
  check(ck, "3", "33")
  check(ck, "4", "44")

  fmt.Printf("OK\n")

  // kill solo server, start new server, check that
  // it does not start serving as primary

  fmt.Printf("Kill last server, new one should not be active: ")

  s2.kill()
  s3 := StartServer(vshost, port("3"))
  for i := 0; i < viewservice.DeadPings * 2; i++ {
    v, _ := vck.Get()
    if v.Primary == s3.me {
      t.Fatal("uninitialized backup promoted to primary")
    }
    time.Sleep(viewservice.PingInterval)
  }

  fmt.Printf("OK\n")

  s1.kill()
  s2.kill()
  s3.kill()
  time.Sleep(time.Second)
  vs.Kill()
  time.Sleep(time.Second)
}

// Put right after a backup dies.
func TestFailPut(t *testing.T) {
  runtime.GOMAXPROCS(4)

  vshost := port("v")
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)

  s1 := StartServer(vshost, port("1"))
  time.Sleep(time.Second)
  s2 := StartServer(vshost, port("2"))
  time.Sleep(time.Second)
  s3 := StartServer(vshost, port("3"))

  for i := 0; i < viewservice.DeadPings * 3; i++ {
    v, _ := vck.Get()
    if v.Primary != "" && v.Backup != "" {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  time.Sleep(time.Second) // wait for backup initializion
  v1, _ := vck.Get()
  if v1.Primary != s1.me || v1.Backup != s2.me {
    t.Fatalf("wrong primary or backup")
  }

  ck := MakeClerk(vshost, "")

  ck.Put("a", "aa")
  ck.Put("b", "bb")
  ck.Put("c", "cc")
  check(ck, "a", "aa")
  check(ck, "b", "bb")
  check(ck, "c", "cc")

  // kill backup, then immediate Put
  fmt.Printf("Put() immediately after backup failure: ")
  s2.kill()
  ck.Put("a", "aaa")
  check(ck, "a", "aaa")

  for i := 0; i < viewservice.DeadPings * 3; i++ {
    v, _ := vck.Get()
    if v.Viewnum > v1.Viewnum && v.Primary != ""  && v.Backup != "" {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  time.Sleep(time.Second) // wait for backup initialization
  v2, _ := vck.Get()
  if v2.Primary != s1.me || v2.Backup != s3.me {
    t.Fatal("wrong primary or backup")
  }

  check(ck, "a", "aaa")
  fmt.Printf("OK\n")

  // kill primary, then immediate Put
  fmt.Printf("Put() immediately after primary failure: ")
  s1.kill()
  ck.Put("b", "bbb")
  check(ck, "b", "bbb")

  for i := 0; i < viewservice.DeadPings * 3; i++ {
    v, _ := vck.Get()
    if v.Viewnum > v2.Viewnum && v.Primary != "" {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  time.Sleep(time.Second)

  check(ck, "a", "aaa")
  check(ck, "b", "bbb")
  check(ck, "c", "cc")
  fmt.Printf("OK\n")

  s1.kill()
  s2.kill()
  s3.kill()
  time.Sleep(viewservice.PingInterval * 2)
  vs.Kill()
}

// do a bunch of concurrent Put()s on the same key,
// then check that primary and backup have identical values.
// i.e. that they processed the Put()s in the same order.
func TestConcurrentSame(t *testing.T) {
  runtime.GOMAXPROCS(4)

  vshost := port("v")
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)

  fmt.Printf("Concurrent Put()s to the same key: ")

  const nservers = 2
  var sa [nservers]*PBServer
  for i := 0; i < nservers; i++ {
    sa[i] = StartServer(vshost, port(strconv.Itoa(i+1)))
  }

  for iters := 0; iters < viewservice.DeadPings*2; iters++ {
    view, _ := vck.Get()
    if view.Primary != "" && view.Backup != "" {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  
  // give p+b time to ack, initialize
  time.Sleep(viewservice.PingInterval * viewservice.DeadPings)

  done := false

  view1, _ := vck.Get()
  const nclients = 3
  const nkeys = 2
  for xi := 0; xi < nclients; xi++ {
    go func(i int) {
      ck := MakeClerk(vshost, "")
      rr := rand.New(rand.NewSource(int64(os.Getpid()+i)))
      for done == false {
        k := strconv.Itoa(rr.Int() % nkeys)
        v := strconv.Itoa(rr.Int())
        err := ck.Put(k, v)
        if err != nil {
          t.Fatalf("Put(%v, %v) failed, err=%v", k, v, err)
        }
      }
    }(xi)
  }

  time.Sleep(5 * time.Second)
  done = true
  time.Sleep(time.Second)

  // read from primary
  ck := MakeClerk(vshost, "")
  var vals [nkeys]string
  for i := 0; i < nkeys; i++ {
    vals[i], _ = ck.Get(strconv.Itoa(i))
    if vals[i] == "" {
      t.Fatalf("Get(%v) failed from primary", i)
    }
  }

  // kill the primary
  for i := 0; i < nservers; i++ {
    if view1.Primary == sa[i].me {
      sa[i].kill()
      break
    }
  }
  for iters := 0; iters < viewservice.DeadPings*2; iters++ {
    view, _ := vck.Get()
    if view.Primary == view1.Backup {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  view2, _ := vck.Get()
  if view2.Primary != view1.Backup {
    t.Fatal("wrong Primary")
  }

  // read from old backup
  for i := 0; i < nkeys; i++ {
    z, _ := ck.Get(strconv.Itoa(i))
    if z != vals[i] {
      t.Fatalf("Get(%v) from backup; wanted %v, got %v", i, vals[i], z)
    }
  }

  fmt.Printf("OK\n")

  for i := 0; i < nservers; i++ {
    sa[i].kill()
  }
  time.Sleep(time.Second)
  vs.Kill()
  time.Sleep(time.Second)
}

func TestConcurrentSameUnreliable(t *testing.T) {
  runtime.GOMAXPROCS(4)

  vshost := port("v")
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)

  fmt.Printf("Concurrent Put()s to the same key; unreliable: ")

  const nservers = 2
  var sa [nservers]*PBServer
  for i := 0; i < nservers; i++ {
    sa[i] = StartServer(vshost, port(strconv.Itoa(i+1)))
    sa[i].unreliable = true
  }

  for iters := 0; iters < viewservice.DeadPings*2; iters++ {
    view, _ := vck.Get()
    if view.Primary != "" && view.Backup != "" {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  
  // give p+b time to ack, initialize
  time.Sleep(viewservice.PingInterval * viewservice.DeadPings)

  done := false

  view1, _ := vck.Get()
  const nclients = 3
  const nkeys = 2
  for xi := 0; xi < nclients; xi++ {
    go func(i int) {
      ck := MakeClerk(vshost, "")
      rr := rand.New(rand.NewSource(int64(os.Getpid()+i)))
      for done == false {
        k := strconv.Itoa(rr.Int() % nkeys)
        v := strconv.Itoa(rr.Int())
        err := ck.Put(k, v)
        if err != nil {
          t.Fatalf("Put(%v, %v) failed, err=%v", k, v, err)
        }
      }
    }(xi)
  }

  time.Sleep(5 * time.Second)
  done = true
  time.Sleep(time.Second)

  // read from primary
  ck := MakeClerk(vshost, "")
  var vals [nkeys]string
  for i := 0; i < nkeys; i++ {
    vals[i], _ = ck.Get(strconv.Itoa(i))
    if vals[i] == "" {
      t.Fatalf("Get(%v) failed from primary", i)
    }
  }

  // kill the primary
  for i := 0; i < nservers; i++ {
    if view1.Primary == sa[i].me {
      sa[i].kill()
      break
    }
  }
  for iters := 0; iters < viewservice.DeadPings*2; iters++ {
    view, _ := vck.Get()
    if view.Primary == view1.Backup {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  view2, _ := vck.Get()
  if view2.Primary != view1.Backup {
    t.Fatal("wrong Primary")
  }

  // read from old backup
  for i := 0; i < nkeys; i++ {
    z, _ := ck.Get(strconv.Itoa(i))
    if z != vals[i] {
      t.Fatalf("Get(%v) from backup; wanted %v, got %v", i, vals[i], z)
    }
  }

  fmt.Printf("OK\n")

  for i := 0; i < nservers; i++ {
    sa[i].kill()
  }
  time.Sleep(time.Second)
  vs.Kill()
  time.Sleep(time.Second)
}

// check that a partitioned primary does not serve requests.
func TestPartition(t *testing.T) {
  runtime.GOMAXPROCS(4)

  vshost := port("v")
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)

  fmt.Printf("Partition an old primary: ")

  s1 := StartServer(vshost, port("1"))

  for i := 0; i < viewservice.DeadPings; i++ {
    if vck.Primary() != "" {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }

  s2 := StartServer(vshost, port("2"))
  time.Sleep(time.Second)
  s3 := StartServer(vshost, port("3"))
  time.Sleep(time.Second)

  v1, _ := vck.Get()
  if v1.Primary != s1.me || v1.Backup != s2.me {
    t.Fatal("wrong primary or backup")
  }

  ck := MakeClerk(vshost, "")

  ck.Put("a", "aa")
  ck.Put("b", "bb")
  check(ck, "a", "aa")
  check(ck, "b", "bb")

  s1.partitioned = true

  for iters := 0; iters < viewservice.DeadPings * 3; iters++ {
    v, _ := vck.Get()
    if v.Primary == s2.me && v.Backup == s3.me {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  v2, _ := vck.Get()
  if v2.Primary != s2.me || v2.Backup != s3.me {
    t.Fatal("wrong primary or backup 2")
  }
  check(ck, "a", "aa")
  check(ck, "b", "bb")

  {
    // try to Put to the old (partitioned) primary
    c, err := rpc.Dial("unix", s1.me)
    if err != nil {
      t.Fatal("could not dial partitioned server")
    }
    defer c.Close()
    args := &PutArgs{}
    args.Key = "a"
    args.Value = "bad"
    var reply PutReply
    err = c.Call("PBServer.Put", args, &reply)
    if err == nil && reply.Err == OK {
      t.Fatal("Put should have failed")
    }
  }
  check(ck, "a", "aa")

  s1.partitioned = false
  time.Sleep(time.Second)
  s2.kill()
  time.Sleep(time.Second)
  check(ck, "a", "aa")
  check(ck, "b", "bb")

  fmt.Printf("OK\n")

  s1.kill()
  s2.kill()
  s3.kill()
  vs.Kill()
  time.Sleep(time.Second)
}

// constant put/get while crashing and restarting servers
func TestRepeatedCrash(t *testing.T) {
  runtime.GOMAXPROCS(4)

  vshost := port("v")
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)
  
  fmt.Printf("Repeated failures/restarts: ")

  const nservers = 3
  var sa [nservers]*PBServer
  for i := 0; i < nservers; i++ {
    sa[i] = StartServer(vshost, port(strconv.Itoa(i+1)))
  }

  for i := 0; i < viewservice.DeadPings; i++ {
    if vck.Primary() != "" {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }

  done := false

  go func() {
    // kill and restart servers
    rr := rand.New(rand.NewSource(int64(os.Getpid())))
    for done == false {
      i := rr.Int() % nservers
      // fmt.Printf("%v killing %v\n", ts(), 5001+i)
      sa[i].kill()

      // wait long enough for new view to form, backup to be initialized
      time.Sleep(2 * viewservice.PingInterval * viewservice.DeadPings)

      sa[i] = StartServer(vshost, port(strconv.Itoa(i+1)))

      // wait long enough for new view to form, backup to be initialized
      time.Sleep(2 * viewservice.PingInterval * viewservice.DeadPings)
    }
  } ()

  for xi := 0; xi < 2; xi++ {
    go func(i int) {
      ck := MakeClerk(vshost, "")
      data := map[string]string{}
      rr := rand.New(rand.NewSource(int64(os.Getpid()+i)))
      for done == false {
        k := strconv.Itoa((i * 1000000) + (rr.Int() % 10))
        wanted, ok := data[k]
        if ok {
          v, err := ck.Get(k)
          if v != wanted {
            t.Fatalf("key=%v wanted=%v got=%v err=%v", k, wanted, v, err)
          }
        }
        nv := strconv.Itoa(rr.Int())
        err := ck.Put(k, nv)
        if err != nil {
          t.Fatalf("Put(%v, %v) failed, err=%v", k, nv, err)
        }
        // if no sleep here, then server tick() threads do not get 
        // enough time to Ping the viewserver.
        time.Sleep(10 * time.Millisecond)
      }
    }(xi)
  }

  time.Sleep(20 * time.Second)
  done = true
  time.Sleep(time.Second)

  ck := MakeClerk(vshost, "")
  ck.Put("aaa", "bbb")
  if v, _ := ck.Get("aaa"); v != "bbb" {
    t.Fatalf("final Put/Get failed")
  }

  fmt.Printf("OK\n")

  for i := 0; i < nservers; i++ {
    sa[i].kill()
  }
  time.Sleep(time.Second)
  vs.Kill()
  time.Sleep(time.Second)
}

func TestRepeatedCrashUnreliable(t *testing.T) {
  runtime.GOMAXPROCS(4)

  vshost := port("v")
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)
  
  fmt.Printf("Repeated failures/restarts; unreliable: ")

  const nservers = 3
  var sa [nservers]*PBServer
  for i := 0; i < nservers; i++ {
    sa[i] = StartServer(vshost, port(strconv.Itoa(i+1)))
    sa[i].unreliable = true
  }

  for i := 0; i < viewservice.DeadPings; i++ {
    if vck.Primary() != "" {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }

  done := false

  go func() {
    // kill and restart servers
    rr := rand.New(rand.NewSource(int64(os.Getpid())))
    for done == false {
      i := rr.Int() % nservers
      // fmt.Printf("%v killing %v\n", ts(), 5001+i)
      sa[i].kill()

      // wait long enough for new view to form, backup to be initialized
      time.Sleep(2 * viewservice.PingInterval * viewservice.DeadPings)

      sa[i] = StartServer(vshost, port(strconv.Itoa(i+1)))

      // wait long enough for new view to form, backup to be initialized
      time.Sleep(2 * viewservice.PingInterval * viewservice.DeadPings)
    }
  } ()

  for xi := 0; xi < 2; xi++ {
    go func(i int) {
      ck := MakeClerk(vshost, "")
      data := map[string]string{}
      rr := rand.New(rand.NewSource(int64(os.Getpid()+i)))
      for done == false {
        k := strconv.Itoa((i * 1000000) + (rr.Int() % 10))
        wanted, ok := data[k]
        if ok {
          v, err := ck.Get(k)
          if v != wanted {
            t.Fatalf("key=%v wanted=%v got=%v err=%v", k, wanted, v, err)
          }
        }
        nv := strconv.Itoa(rr.Int())
        err := ck.Put(k, nv)
        if err != nil {
          t.Fatalf("Put(%v, %v) failed, err=%v", k, nv, err)
        }
        // if no sleep here, then server tick() threads do not get 
        // enough time to Ping the viewserver.
        time.Sleep(10 * time.Millisecond)
      }
    }(xi)
  }

  time.Sleep(20 * time.Second)
  done = true
  time.Sleep(time.Second)

  ck := MakeClerk(vshost, "")
  ck.Put("aaa", "bbb")
  if v, _ := ck.Get("aaa"); v != "bbb" {
    t.Fatalf("final Put/Get failed")
  }

  fmt.Printf("OK\n")

  for i := 0; i < nservers; i++ {
    sa[i].kill()
  }
  time.Sleep(time.Second)
  vs.Kill()
  time.Sleep(time.Second)
}
