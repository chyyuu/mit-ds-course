package shardmaster

import "testing"
import "runtime"
import "strconv"
import "os"
// import "time"
import "fmt"
// import "math/rand"

func port(tag string, host int) string {
  s := "/var/tmp/824-"
  s += strconv.Itoa(os.Getuid()) + "/"
  os.Mkdir(s, 0777)
  s += "sm-"
  s += strconv.Itoa(os.Getpid()) + "-"
  s += tag + "-"
  s += strconv.Itoa(host)
  return s
}

func cleanup(sma []*ShardMaster) {
  for i := 0; i < len(sma); i++ {
    if sma[i] != nil {
      sma[i].Kill()
    }
  }
}

//
// maybe should take a cka[] and find the server with
// the highest Num.
//
func check(t *testing.T, groups []int64, ck *Clerk) {
  c := ck.Query()
  if len(c.Groups) != len(groups) {
    t.Fatalf("wanted %v groups, got %v", len(groups), len(c.Groups))
  }

  // are the groups as expected?
  for _, g := range groups {
    _, ok := c.Groups[g]
    if ok != true {
      t.Fatalf("missing group %v", g)
    }
  }

  // any un-allocated shards?
  if len(groups) > 0 {
    for s, g := range c.Shards {
      _, ok := c.Groups[g]
      if ok == false {
        t.Fatalf("shard %v -> invalid group %v", s, g)
      }
    }
  }

  // more or less balanced sharding?
  counts := map[int64]int{}
  for _, g := range c.Shards {
    counts[g] += 1
  }
  min := 257
  max := 0
  for g, _ := range c.Groups {
    if counts[g] > max {
      max = counts[g]
    }
    if counts[g] < min {
      min = counts[g]
    }
  }
  if max > min + 1 {
    t.Fatalf("max %v too much larger than min %v", max, min)
  }
}

func TestBasic(t *testing.T) {
  runtime.GOMAXPROCS(4)

  const nservers = 3
  var sma []*ShardMaster = make([]*ShardMaster, nservers)
  var kvh []string = make([]string, nservers)
  defer cleanup(sma)

  for i := 0; i < nservers; i++ {
    kvh[i] = port("basic", i)
  }
  for i := 0; i < nservers; i++ {
    sma[i] = StartServer(kvh, i)
  }

  ck := MakeClerk(kvh)
  var cka [nservers]*Clerk
  for i := 0; i < nservers; i++ {
    cka[i] = MakeClerk([]string{kvh[i]})
  }

  fmt.Printf("Basic leave/join: ")

  check(t, []int64{}, ck)

  var gid1 int64 = 1
  ck.Join(gid1, []string{"tweedledum", "tweedledee", "crow"})
  check(t, []int64{gid1}, ck)

  var gid2 int64 = 2
  ck.Join(gid2, []string{"a", "b", "c"})
  check(t, []int64{gid1,gid2}, ck)

  ck.Join(gid2, []string{"a", "b", "c"})
  check(t, []int64{gid1,gid2}, ck)

  ck.Leave(gid1)
  check(t, []int64{gid2}, ck)

  ck.Leave(gid1)
  check(t, []int64{gid2}, ck)

  fmt.Printf("OK\n")

  fmt.Printf("Concurrent leave/join: ")

  const npara = 10
  gids := make([]int64, npara)
  var ca [npara]chan bool
  for xi := 0; xi < npara; xi++ {
    gids[xi] = int64(xi+1)
    ca[xi] = make(chan bool)
    go func(i int) {
      defer func() { ca[i] <- true }()
      var gid int64 = gids[i]
      cka[(i+0)%nservers].Join(gid+1000, []string{"a", "b", "c"})
      cka[(i+0)%nservers].Join(gid, []string{"a", "b", "c"})
      cka[(i+1)%nservers].Leave(gid+1000)
    }(xi)
  }
  for i := 0; i < npara; i++ {
    <- ca[i]
  }
  check(t, gids, ck)

  fmt.Printf("OK\n")

  fmt.Printf("Minimal transfers after joins: ")

  c1 := ck.Query()
  for i := 0; i < 5; i++ {
    ck.Join(int64(npara+1+i), []string{"a","b","c"})
  }
  c2 := ck.Query()
  for i := int64(1); i <= npara; i++ {
    for j := 0; j < len(c1.Shards); j++ {
      if c2.Shards[j] == i {
        if c1.Shards[j] != i {
          t.Fatalf("non-minimal transfer after Join()s")
        }
      }
    }
  }

  fmt.Printf("OK\n")

  fmt.Printf("Minimal transfers after leaves: ")

  for i := 0; i < 5; i++ {
    ck.Leave(int64(npara+1+i))
  }
  c3 := ck.Query()
  for i := int64(1); i <= npara; i++ {
    for j := 0; j < len(c1.Shards); j++ {
      if c2.Shards[j] == i {
        if c3.Shards[j] != i {
          t.Fatalf("non-minimal transfer after Leave()s")
        }
      }
    }
  }

  fmt.Printf("OK\n")
}

func TestUnreliable(t *testing.T) {
  runtime.GOMAXPROCS(4)

  const nservers = 3
  var sma []*ShardMaster = make([]*ShardMaster, nservers)
  var kvh []string = make([]string, nservers)
  defer cleanup(sma)

  for i := 0; i < nservers; i++ {
    kvh[i] = port("basic", i)
  }
  for i := 0; i < nservers; i++ {
    sma[i] = StartServer(kvh, i)
    sma[i].unreliable = true
  }

  ck := MakeClerk(kvh)
  var cka [nservers]*Clerk
  for i := 0; i < nservers; i++ {
    cka[i] = MakeClerk([]string{kvh[i]})
  }

  fmt.Printf("Concurrent leave/join, unreliable: ")

  const npara = 20
  gids := make([]int64, npara)
  var ca [npara]chan bool
  for xi := 0; xi < npara; xi++ {
    gids[xi] = int64(xi+1)
    ca[xi] = make(chan bool)
    go func(i int) {
      defer func() { ca[i] <- true }()
      var gid int64 = gids[i]
      cka[(i+0)%nservers].Join(gid+1000, []string{"a", "b", "c"})
      cka[(i+0)%nservers].Join(gid, []string{"a", "b", "c"})
      cka[(i+1)%nservers].Leave(gid+1000)
    }(xi)
  }
  for i := 0; i < npara; i++ {
    <- ca[i]
  }
  check(t, gids, ck)

  fmt.Printf("OK\n")
}
