package pbservice

import "viewservice"
import "fmt"
import "testing"
import "time"
import "log"
import "runtime"
import "math/rand"
import "os"
import "strconv"

func check(ck *Clerk, key string, value string) {
  v := ck.Get(key)
  if v != value {
    log.Fatalf("Get(%v) -> %v, expected %v", key, v, value)
  }
}

func port(tag string, host int) string {
  s := "/var/tmp/824-"
  s += strconv.Itoa(os.Getuid()) + "/"
  os.Mkdir(s, 0777)
  s += "pb-"
  s += strconv.Itoa(os.Getpid()) + "-"
  s += tag + "-"
  s += strconv.Itoa(host)
  return s
}

func TestBasicFail(t *testing.T) {
  runtime.GOMAXPROCS(4)

  tag := "basic"
  vshost := port(tag+"v", 1)
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)

  ck := MakeClerk(vshost, "")

  fmt.Printf("Test: Single primary, no backup ...\n")

  s1 := StartServer(vshost, port(tag, 1))

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

  fmt.Printf("  ... Passed\n")

  // add a backup

  fmt.Printf("Test: Add a backup ...\n")

  s2 := StartServer(vshost, port(tag, 2))
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

  fmt.Printf("  ... Passed\n")

  // kill the primary

  fmt.Printf("Test: Primary failure ...\n")

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

  fmt.Printf("  ... Passed\n")

  // kill solo server, start new server, check that
  // it does not start serving as primary

  fmt.Printf("Test: Kill last server, new one should not be active ...\n")

  s2.kill()
  s3 := StartServer(vshost, port(tag, 3))
  time.Sleep(1 * time.Second)
  get_done := false
  go func() {
    ck.Get("1")
    get_done = true
  }()
  time.Sleep(2 * time.Second)
  if get_done {
    t.Fatalf("ck.Get() returned even though no initialized primary")
  }

  fmt.Printf("  ... Passed\n")

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

  tag := "failput"
  vshost := port(tag+"v", 1)
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)

  s1 := StartServer(vshost, port(tag, 1))
  time.Sleep(time.Second)
  s2 := StartServer(vshost, port(tag, 2))
  time.Sleep(time.Second)
  s3 := StartServer(vshost, port(tag, 3))

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
  fmt.Printf("Test: Put() immediately after backup failure ...\n")
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
  fmt.Printf("  ... Passed\n")

  // kill primary, then immediate Put
  fmt.Printf("Test: Put() immediately after primary failure ...\n")
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
  fmt.Printf("  ... Passed\n")

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

  tag := "cs"
  vshost := port(tag+"v", 1)
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)

  fmt.Printf("Test: Concurrent Put()s to the same key ...\n")

  const nservers = 2
  var sa [nservers]*PBServer
  for i := 0; i < nservers; i++ {
    sa[i] = StartServer(vshost, port(tag, i+1))
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
        ck.Put(k, v)
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
    vals[i] = ck.Get(strconv.Itoa(i))
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
    z := ck.Get(strconv.Itoa(i))
    if z != vals[i] {
      t.Fatalf("Get(%v) from backup; wanted %v, got %v", i, vals[i], z)
    }
  }

  fmt.Printf("  ... Passed\n")

  for i := 0; i < nservers; i++ {
    sa[i].kill()
  }
  time.Sleep(time.Second)
  vs.Kill()
  time.Sleep(time.Second)
}

func TestConcurrentSameUnreliable(t *testing.T) {
  runtime.GOMAXPROCS(4)

  tag := "csu"
  vshost := port(tag+"v", 1)
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)

  fmt.Printf("Test: Concurrent Put()s to the same key; unreliable ...\n")

  const nservers = 2
  var sa [nservers]*PBServer
  for i := 0; i < nservers; i++ {
    sa[i] = StartServer(vshost, port(tag, i+1))
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
        ck.Put(k, v)
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
    vals[i] = ck.Get(strconv.Itoa(i))
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
    z := ck.Get(strconv.Itoa(i))
    if z != vals[i] {
      t.Fatalf("Get(%v) from backup; wanted %v, got %v", i, vals[i], z)
    }
  }

  fmt.Printf("  ... Passed\n")

  for i := 0; i < nservers; i++ {
    sa[i].kill()
  }
  time.Sleep(time.Second)
  vs.Kill()
  time.Sleep(time.Second)
}

// constant put/get while crashing and restarting servers
func TestRepeatedCrash(t *testing.T) {
  runtime.GOMAXPROCS(4)

  tag := "rc"
  vshost := port(tag+"v", 1)
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)
  
  fmt.Printf("Test: Repeated failures/restarts ...\n")

  const nservers = 3
  var sa [nservers]*PBServer
  for i := 0; i < nservers; i++ {
    sa[i] = StartServer(vshost, port(tag, i+1))
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

      sa[i] = StartServer(vshost, port(tag, i+1))

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
          v := ck.Get(k)
          if v != wanted {
            t.Fatalf("key=%v wanted=%v got=%v", k, wanted, v)
          }
        }
        nv := strconv.Itoa(rr.Int())
        ck.Put(k, nv)
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
  if v := ck.Get("aaa"); v != "bbb" {
    t.Fatalf("final Put/Get failed")
  }

  fmt.Printf("  ... Passed\n")

  for i := 0; i < nservers; i++ {
    sa[i].kill()
  }
  time.Sleep(time.Second)
  vs.Kill()
  time.Sleep(time.Second)
}

func TestRepeatedCrashUnreliable(t *testing.T) {
  runtime.GOMAXPROCS(4)

  tag := "rcu"
  vshost := port(tag+"v", 1)
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)
  
  fmt.Printf("Test: Repeated failures/restarts; unreliable ...\n")

  const nservers = 3
  var sa [nservers]*PBServer
  for i := 0; i < nservers; i++ {
    sa[i] = StartServer(vshost, port(tag, i+1))
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

      sa[i] = StartServer(vshost, port(tag, i+1))

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
          v := ck.Get(k)
          if v != wanted {
            t.Fatalf("key=%v wanted=%v got=%v", k, wanted, v)
          }
        }
        nv := strconv.Itoa(rr.Int())
        ck.Put(k, nv)
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
  if v := ck.Get("aaa"); v != "bbb" {
    t.Fatalf("final Put/Get failed")
  }

  fmt.Printf("  ... Passed\n")

  for i := 0; i < nservers; i++ {
    sa[i].kill()
  }
  time.Sleep(time.Second)
  vs.Kill()
  time.Sleep(time.Second)
}

func TestPartition1(t *testing.T) {
  runtime.GOMAXPROCS(4)

  tag := "part1"
  vshost := port(tag+"v", 1)
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)

  ck := MakeClerk(vshost, "")

  fmt.Printf("Test: Old primary does not serve Gets ...\n")

  vshosta := vshost + "a"
  os.Link(vshost, vshosta)

  s1 := StartServer(vshosta, port(tag, 1))

  deadtime := viewservice.PingInterval * viewservice.DeadPings
  time.Sleep(deadtime * 2)
  if vck.Primary() != s1.me {
    t.Fatal("primary never formed initial view")
  }

  s2 := StartServer(vshost, port(tag, 2))
  time.Sleep(deadtime * 2)
  v1, _ := vck.Get()
  if v1.Primary != s1.me || v1.Backup != s2.me {
    t.Fatal("backup did not join view")
  }
  
  ck.Put("a", "1")
  check(ck, "a", "1")

  os.Remove(vshosta)

  // now s1 cannot talk to viewserver, so view will change.

  for iter := 0; iter < viewservice.DeadPings * 3; iter++ {
    if vck.Primary() == s2.me {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  if vck.Primary() != s2.me {
    t.Fatalf("primary never changed")
  }

  // wait long enough that s2 is guaranteed to have Pinged
  // the viewservice, and thus that s2 must know about
  // the new view.
  time.Sleep(2 * viewservice.PingInterval)

  // s1 can talk to s2, so s1 should learn that it
  // should not act as primary.

  get_succeeded := false

  go func(){
    args := &GetArgs{}
    args.Key = "a"
    var reply GetReply
    ok := call(s1.me, "PBServer.Get", args, &reply)
    if ok && reply.Err == OK {
      get_succeeded = true
    }
  }()

  time.Sleep(3 * time.Second)
  if get_succeeded {
    t.Fatalf("Get to old server succeeded, but should not have")
  }

  check(ck, "a", "1")

  fmt.Printf("  ... Passed\n")

  s1.kill()
  s2.kill()
  vs.Kill()
}

func TestPartition2(t *testing.T) {
  runtime.GOMAXPROCS(4)

  tag := "part2"
  vshost := port(tag+"v", 1)
  vs := viewservice.StartServer(vshost)
  time.Sleep(time.Second)
  vck := viewservice.MakeClerk("", vshost)

  ck := MakeClerk(vshost, "")

  vshosta := vshost + "a"
  os.Link(vshost, vshosta)

  s1 := StartServer(vshosta, port(tag, 1))

  fmt.Printf("Test: Partitioned old primary does not complete Gets ...\n")

  deadtime := viewservice.PingInterval * viewservice.DeadPings
  time.Sleep(deadtime * 2)
  if vck.Primary() != s1.me {
    t.Fatal("primary never formed initial view")
  }

  s2 := StartServer(vshost, port(tag, 2))
  time.Sleep(deadtime * 2)
  v1, _ := vck.Get()
  if v1.Primary != s1.me || v1.Backup != s2.me {
    t.Fatal("backup did not join view")
  }
  
  ck.Put("a", "1")
  check(ck, "a", "1")

  os.Remove(vshosta)

  // now s1 cannot talk to viewserver, so view will change.

  for iter := 0; iter < viewservice.DeadPings * 3; iter++ {
    if vck.Primary() == s2.me {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  if vck.Primary() != s2.me {
    t.Fatalf("primary never changed")
  }

  s3 := StartServer(vshost, port(tag, 3))
  for iter := 0; iter < viewservice.DeadPings * 3; iter++ {
    v, _ := vck.Get()
    if v.Backup == s3.me && v.Primary == s2.me {
      break
    }
    time.Sleep(viewservice.PingInterval)
  }
  v2, _ := vck.Get()
  if v2.Primary != s2.me || v2.Backup != s3.me {
    t.Fatalf("new backup never joined")
  }
  time.Sleep(2 * time.Second)
  ck.Put("a", "2")

  s2.kill()
  time.Sleep(1 * time.Second)

  get_finished := false
  go func(){
    args := &GetArgs{}
    args.Key = "a"
    var reply GetReply
    call(s1.me, "PBServer.Get", args, &reply)
    get_finished = true
  }()

  time.Sleep(2 * time.Second)
  if get_finished == true {
    t.Fatalf("partitioned primary replied to a Get")
  }

  check(ck, "a", "2")

  fmt.Printf("  ... Passed\n")

  s1.kill()
  s2.kill()
  s3.kill()
  vs.Kill()
}
