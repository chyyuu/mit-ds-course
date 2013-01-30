package shardmaster

import "testing"
import "runtime"
import "strconv"
import "os"
// import "time"
import "fmt"
import "math/rand"

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
  c := ck.Query(-1)
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

  fmt.Printf("Test: Basic leave/join ...\n")

  cfa := make([]Config, 6)
  cfa[0] = ck.Query(-1)

  check(t, []int64{}, ck)

  var gid1 int64 = 1
  ck.Join(gid1, []string{"x", "y", "z"})
  check(t, []int64{gid1}, ck)
  cfa[1] = ck.Query(-1)

  var gid2 int64 = 2
  ck.Join(gid2, []string{"a", "b", "c"})
  check(t, []int64{gid1,gid2}, ck)
  cfa[2] = ck.Query(-1)

  ck.Join(gid2, []string{"a", "b", "c"})
  check(t, []int64{gid1,gid2}, ck)
  cfa[3] = ck.Query(-1)

  cfx := ck.Query(-1)
  sa1 := cfx.Groups[gid1]
  if len(sa1) != 3 || sa1[0] != "x" || sa1[1] != "y" || sa1[2] != "z" {
    t.Fatal("wrong servers for gid %v: %v\n", gid1, sa1)
  }
  sa2 := cfx.Groups[gid2]
  if len(sa2) != 3 || sa2[0] != "a" || sa2[1] != "b" || sa2[2] != "c" {
    t.Fatal("wrong servers for gid %v: %v\n", gid2, sa2)
  }

  ck.Leave(gid1)
  check(t, []int64{gid2}, ck)
  cfa[4] = ck.Query(-1)

  ck.Leave(gid1)
  check(t, []int64{gid2}, ck)
  cfa[5] = ck.Query(-1)

  fmt.Printf("  ... Passed\n")

  fmt.Printf("Test: Historical queries ...\n")

  for i := 0; i < len(cfa); i++ {
    c := ck.Query(cfa[i].Num)
    if c.Num != cfa[i].Num {
      t.Fatalf("historical Num wrong")
    }
    if c.Shards != cfa[i].Shards {
      t.Fatalf("historical Shards wrong")
    }
    if len(c.Groups) != len(cfa[i].Groups) {
      t.Fatalf("number of historical Groups is wrong")
    }
    for gid, sa := range c.Groups {
      sa1, ok := cfa[i].Groups[gid]
      if ok == false || len(sa1) != len(sa) {
        t.Fatalf("historical len(Groups) wrong")
      }
      if ok && len(sa1) == len(sa) {
        for j := 0; j < len(sa); j++ {
          if sa[j] != sa1[j] {
            t.Fatalf("historical Groups wrong")
          }
        }
      }
    }
  }

  fmt.Printf("  ... Passed\n")

  fmt.Printf("Test: Move ...\n")
  {
    var gid3 int64 = 503
    ck.Join(gid3, []string{"3a", "3b", "3c"})
    var gid4 int64 = 504
    ck.Join(gid4, []string{"4a", "4b", "4c"})
    for i := 0; i < NShards; i++ {
      if i < NShards / 2 {
        ck.Move(i, gid3)
      } else {
        ck.Move(i, gid4)
      }
    }
    cf2 := ck.Query(-1)
    for i := 0; i < NShards; i++ {
      if i < NShards / 2 {
        if cf2.Shards[i] != gid3 {
          t.Fatalf("expected shard %v on gid %v actually %v",
                   i, gid3, cf2.Shards[i])
        }
      } else {
        if cf2.Shards[i] != gid4 {
          t.Fatalf("expected shard %v on gid %v actually %v",
                   i, gid4, cf2.Shards[i])
        }
      }
    }
    ck.Leave(gid3)
    ck.Leave(gid4)
  }
  fmt.Printf("  ... Passed\n")

  fmt.Printf("Test: Concurrent leave/join ...\n")

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

  fmt.Printf("  ... Passed\n")

  fmt.Printf("Test: Minimal transfers after joins ...\n")

  c1 := ck.Query(-1)
  for i := 0; i < 5; i++ {
    ck.Join(int64(npara+1+i), []string{"a","b","c"})
  }
  c2 := ck.Query(-1)
  for i := int64(1); i <= npara; i++ {
    for j := 0; j < len(c1.Shards); j++ {
      if c2.Shards[j] == i {
        if c1.Shards[j] != i {
          t.Fatalf("non-minimal transfer after Join()s")
        }
      }
    }
  }

  fmt.Printf("  ... Passed\n")

  fmt.Printf("Test: Minimal transfers after leaves ...\n")

  for i := 0; i < 5; i++ {
    ck.Leave(int64(npara+1+i))
  }
  c3 := ck.Query(-1)
  for i := int64(1); i <= npara; i++ {
    for j := 0; j < len(c1.Shards); j++ {
      if c2.Shards[j] == i {
        if c3.Shards[j] != i {
          t.Fatalf("non-minimal transfer after Leave()s")
        }
      }
    }
  }

  fmt.Printf("  ... Passed\n")
}

func TestUnreliable(t *testing.T) {
  runtime.GOMAXPROCS(4)

  const nservers = 3
  var sma []*ShardMaster = make([]*ShardMaster, nservers)
  var kvh []string = make([]string, nservers)
  defer cleanup(sma)

  for i := 0; i < nservers; i++ {
    kvh[i] = port("unrel", i)
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

  fmt.Printf("Test: Concurrent leave/join, unreliable, failure ...\n")

  const npara = 20
  gids := make([]int64, npara)
  var ca [npara]chan bool
  for xi := 0; xi < npara; xi++ {
    gids[xi] = int64(xi+1)
    ca[xi] = make(chan bool)
    go func(i int) {
      defer func() { ca[i] <- true }()
      var gid int64 = gids[i]
      cka[1 + (rand.Int() % 2)].Join(gid+1000, []string{"a", "b", "c"})
      cka[1 + (rand.Int() % 2)].Join(gid, []string{"a", "b", "c"})
      cka[1 + (rand.Int() % 2)].Leave(gid+1000)
      // server 0 won't be able to hear any RPCs.
      os.Remove(kvh[0])
    }(xi)
  }
  for i := 0; i < npara; i++ {
    <- ca[i]
  }
  check(t, gids, ck)

  fmt.Printf("  ... Passed\n")
}
