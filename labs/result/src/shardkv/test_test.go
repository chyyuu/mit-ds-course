package shardkv

import "testing"
import "shardmaster"
import "runtime"
import "strconv"
import "os"
import "time"
import "fmt"
import "sync"
import "math/rand"

func port(tag string, host int) string {
  s := "/var/tmp/824-"
  s += strconv.Itoa(os.Getuid()) + "/"
  os.Mkdir(s, 0777)
  s += "skv-"
  s += strconv.Itoa(os.Getpid()) + "-"
  s += tag + "-"
  s += strconv.Itoa(host)
  return s
}

func mcleanup(sma []*shardmaster.ShardMaster) {
  for i := 0; i < len(sma); i++ {
    if sma[i] != nil {
      sma[i].Kill()
    }
  }
}

func cleanup(sa [][]*ShardKV) {
  for i := 0; i < len(sa); i++ {
    for j := 0; j < len(sa[i]); j++ {
      sa[i][j].kill()
    }
  }
}

func setup(tag string, unreliable bool) ([]string, []int64, [][]string, [][]*ShardKV, func()) {
  runtime.GOMAXPROCS(4)
  
  const nmasters = 3
  var sma []*shardmaster.ShardMaster = make([]*shardmaster.ShardMaster, nmasters)
  var smh []string = make([]string, nmasters)
  // defer mcleanup(sma)
  for i := 0; i < nmasters; i++ {
    smh[i] = port(tag+"m", i)
  }
  for i := 0; i < nmasters; i++ {
    sma[i] = shardmaster.StartServer(smh, i)
  }

  const ngroups = 3   // replica groups
  const nreplicas = 3 // servers per group
  gids := make([]int64, ngroups)    // each group ID
  ha := make([][]string, ngroups)   // ShardKV ports, [group][replica]
  sa := make([][]*ShardKV, ngroups) // ShardKVs
  // defer cleanup(sa)
  for i := 0; i < ngroups; i++ {
    gids[i] = int64(i + 100)
    sa[i] = make([]*ShardKV, nreplicas)
    ha[i] = make([]string, nreplicas)
    for j := 0; j < nreplicas; j++ {
      ha[i][j] = port(tag+"s", (i*nreplicas)+j)
    }
    for j := 0; j < nreplicas; j++ {
      sa[i][j] = StartServer(gids[i], smh, ha[i], j)
      sa[i][j].unreliable = unreliable
    }
  }

  clean := func() { mcleanup(sma) ; cleanup(sa) }
  return smh, gids, ha, sa, clean
}

func TestBasic(t *testing.T) {
  smh, gids, ha, _, clean := setup("basic", false)
  defer clean()

  fmt.Printf("Test: Basic Join/Leave ...\n")

  mck := shardmaster.MakeClerk(smh)
  mck.Join(gids[0], ha[0])

  ck := MakeClerk(smh)

  ck.Put("a", "b")
  if ck.Get("a") != "b" {
    t.Fatalf("got wrong value")
  }

  keys := make([]string, 10)
  vals := make([]string, len(keys))
  for i := 0; i < len(keys); i++ {
    keys[i] = strconv.Itoa(rand.Int())
    vals[i] = strconv.Itoa(rand.Int())
    ck.Put(keys[i], vals[i])
  }

  // are keys still there after joins?
  for g := 1; g < len(gids); g++ {
    mck.Join(gids[g], ha[g])
    time.Sleep(1 * time.Second)
    for i := 0; i < len(keys); i++ {
      v := ck.Get(keys[i])
      if v != vals[i] {
        t.Fatalf("joining; wrong value; g=%v k=%v wanted=%v got=%v",
          g, keys[i], vals[i], v)
      }
      vals[i] = strconv.Itoa(rand.Int())
      ck.Put(keys[i], vals[i])
    }
  }
  
  // are keys still there after leaves?
  for g := 0; g < len(gids)-1; g++ {
    mck.Leave(gids[g])
    time.Sleep(1 * time.Second)
    for i := 0; i < len(keys); i++ {
      v := ck.Get(keys[i])
      if v != vals[i] {
        t.Fatalf("leaving; wrong value; g=%v k=%v wanted=%v got=%v",
          g, keys[i], vals[i], v)
      }
      vals[i] = strconv.Itoa(rand.Int())
      ck.Put(keys[i], vals[i])
    }
  }

  fmt.Printf("  ... Passed\n")
}

func TestMove(t *testing.T) {
  smh, gids, ha, _, clean := setup("basic", false)
  defer clean()

  fmt.Printf("Test: Shards really move ...\n")

  mck := shardmaster.MakeClerk(smh)
  mck.Join(gids[0], ha[0])

  ck := MakeClerk(smh)

  // insert one key per shard
  for i := 0; i < shardmaster.NShards; i++ {
    ck.Put(string('0'+i), string('0'+i))
  }

  // add group 1.
  mck.Join(gids[1], ha[1])
  time.Sleep(5 * time.Second)
  
  // check that keys are still there.
  for i := 0; i < shardmaster.NShards; i++ {
    if ck.Get(string('0'+i)) != string('0'+i) {
      t.Fatalf("missing key/value")
    }
  }

  // remove sockets from group 0.
  for i := 0; i < len(ha[0]); i++ {
    os.Remove(ha[0][i])
  }

  count := 0
  var mu sync.Mutex
  for i := 0; i < shardmaster.NShards; i++ {
    go func(me int) {
      myck := MakeClerk(smh)
      v := myck.Get(string('0'+me))
      if v == string('0'+me) {
        mu.Lock()
        count++
        mu.Unlock()
      } else {
        t.Fatalf("Get(%v) yielded %v\n", i, v)
      }
    }(i)
  }

  time.Sleep(10 * time.Second)

  if count > shardmaster.NShards / 3 && count < 2*(shardmaster.NShards/3) {
    fmt.Printf("  ... Passed\n")
  } else {
    t.Fatalf("%v keys worked after killing 1/2 of groups; wanted %v",
      count, shardmaster.NShards / 2)
  }
}

func TestLimp(t *testing.T) {
  smh, gids, ha, sa, clean := setup("basic", false)
  defer clean()

  fmt.Printf("Test: Reconfiguration with some dead replicas ...\n")

  mck := shardmaster.MakeClerk(smh)
  mck.Join(gids[0], ha[0])

  ck := MakeClerk(smh)

  ck.Put("a", "b")
  if ck.Get("a") != "b" {
    t.Fatalf("got wrong value")
  }

  for g := 0; g < len(sa); g++ {
    sa[g][rand.Int() % len(sa[g])].kill()
  }

  keys := make([]string, 10)
  vals := make([]string, len(keys))
  for i := 0; i < len(keys); i++ {
    keys[i] = strconv.Itoa(rand.Int())
    vals[i] = strconv.Itoa(rand.Int())
    ck.Put(keys[i], vals[i])
  }

  // are keys still there after joins?
  for g := 1; g < len(gids); g++ {
    mck.Join(gids[g], ha[g])
    time.Sleep(1 * time.Second)
    for i := 0; i < len(keys); i++ {
      v := ck.Get(keys[i])
      if v != vals[i] {
        t.Fatalf("joining; wrong value; g=%v k=%v wanted=%v got=%v",
          g, keys[i], vals[i], v)
      }
      vals[i] = strconv.Itoa(rand.Int())
      ck.Put(keys[i], vals[i])
    }
  }
  
  // are keys still there after leaves?
  for g := 0; g < len(gids)-1; g++ {
    mck.Leave(gids[g])
    time.Sleep(2 * time.Second)
    for i := 0; i < len(sa[g]); i++ {
      sa[g][i].kill()
    }
    for i := 0; i < len(keys); i++ {
      v := ck.Get(keys[i])
      if v != vals[i] {
        t.Fatalf("leaving; wrong value; g=%v k=%v wanted=%v got=%v",
          g, keys[i], vals[i], v)
      }
      vals[i] = strconv.Itoa(rand.Int())
      ck.Put(keys[i], vals[i])
    }
  }

  fmt.Printf("  ... Passed\n")
}

func doConcurrent(t *testing.T, unreliable bool) {
  smh, gids, ha, _, clean := setup("conc", unreliable)
  defer clean()

  mck := shardmaster.MakeClerk(smh)
  for i := 0; i < len(gids); i++ {
    mck.Join(gids[i], ha[i])
  }

  const npara = 11
  var ca [npara]chan bool
  for i := 0; i < npara; i++ {
    ca[i] = make(chan bool)
    go func(me int) {
      ok := true
      defer func() { ca[me] <- ok }()
      ck := MakeClerk(smh)
      mymck := shardmaster.MakeClerk(smh)
      key := strconv.Itoa(me)
      last := ""
      for iters := 0; iters < 3; iters++ {
        v := ck.Get(key)
        if v != last {
          ok = false
          t.Fatalf("Get(%v) expected %v got %v\n", key, last, v)
        }
        last = strconv.Itoa(rand.Int())
        ck.Put(key, last)

        mymck.Move(rand.Int() % shardmaster.NShards,
          gids[rand.Int() % len(gids)])

        time.Sleep(time.Duration(rand.Int() % 30) * time.Millisecond)
      }
    }(i)
  }

  for i := 0; i < npara; i++ {
    x := <- ca[i]
    if x == false {
      t.Fatalf("something is wrong")
    }
  }
}

func TestConcurrent(t *testing.T) {
  fmt.Printf("Test: Concurrent Put/Get/Move ...\n")
  doConcurrent(t, false)
  fmt.Printf("  ... Passed\n")
}

func TestConcurrentUnreliable(t *testing.T) {
  fmt.Printf("Test: Concurrent Put/Get/Move (unreliable) ...\n")
  doConcurrent(t, true)
  fmt.Printf("  ... Passed\n")
}
