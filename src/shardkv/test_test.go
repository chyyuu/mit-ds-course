package shardkv

import "testing"
import "shardmaster"
import "runtime"
import "strconv"
import "os"
import "time"
import "fmt"
// import "math/rand"

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

func TestBasic(t *testing.T) {
  runtime.GOMAXPROCS(4)

  const nmasters = 3
  var sma []*shardmaster.ShardMaster = make([]*shardmaster.ShardMaster, nmasters)
  var smh []string = make([]string, nmasters)
  defer mcleanup(sma)
  for i := 0; i < nmasters; i++ {
    smh[i] = port("basicm", i)
  }
  for i := 0; i < nmasters; i++ {
    sma[i] = shardmaster.StartServer(smh, i)
  }

  const ngroups = 5   // replica groups
  const nreplicas = 3 // servers per group
  gids := make([]int64, ngroups)
  sa := make([][]*ShardKV, ngroups)
  ha := make([][]string, ngroups) // ShardKV ports
  for i := 0; i < ngroups; i++ {
    gids[i] = int64(i + 100)
    sa[i] = make([]*ShardKV, nreplicas)
    ha[i] = make([]string, nreplicas)
    for j := 0; j < nreplicas; j++ {
      ha[i][j] = port("basics", j)
    }
    for j := 0; j < nreplicas; j++ {
      sa[i][j] = StartServer(gids[i], smh, ha[i], j)
    }
  }

  mck := shardmaster.MakeClerk(smh)
  mck.Join(gids[0], ha[0])
  mck.Join(gids[1], ha[1])
  mck.Join(gids[2], ha[2])

  ck := MakeClerk(smh)

  ck.Put("a", "b")
  fmt.Printf("%v\n", ck.Get("a"))

  time.Sleep(5 * time.Second)
}
