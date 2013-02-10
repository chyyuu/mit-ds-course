package lockservice

import "net"
import "net/rpc"
import "log"
import "sync"
import "fmt"
import "os"
import "io"
import "time"

type LockServer struct {
	mu    sync.Mutex
	l     net.Listener
	dead  bool // for test_test.go
	dying bool // for test_test.go

	am_primary bool   // am I the primary?
	backup     string // backup's port

	// for each lock name, is it locked?
	locks map[string]bool
	xids  map[int64]bool //key: seq num, value: reply of the require with seq num
}

//
// server Lock RPC handler.
// you will have to modify this function
//
func (ls *LockServer) Lock(args *LockArgs, reply *LockReply) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	_, isused := ls.xids[args.Xid]
	if isused == true { //used xid
		//fmt.Printf("ls(%v):Lock: used xid(%v)\n", ls.am_primary, args.Xid)
		reply.OK = ls.xids[args.Xid]
	} else {
		locked, _ := ls.locks[args.Lockname]
		//fmt.Printf("ls(%v): Lock before: Lock(%v), locked(%v), reply(%v)\n", ls.am_primary, args.Lockname, ls.locks[args.Lockname], reply.OK)
		if locked {
			reply.OK = false
		} else {
			reply.OK = true
			ls.locks[args.Lockname] = true
		}
		ls.xids[args.Xid] = reply.OK
		//fmt.Printf("ls(%v): Lock after: Lock(%v), xid(%v), locked(%v)\n", ls.am_primary, args.Lockname, ls.xids[args.Xid], ls.locks[args.Lockname], reply.OK)
		// call B srv
		if ls.am_primary {
			var breply LockReply
			call(ls.backup, "LockServer.Lock", args, &breply)
		}
	}
	return nil
}

//
// server Unlock RPC handler.
//
func (ls *LockServer) Unlock(args *UnlockArgs, reply *UnlockReply) error {

	// Your code here.
	ls.mu.Lock()
	defer ls.mu.Unlock()

	_, isused := ls.xids[args.Xid]
	if isused == true { //used xid
		//fmt.Printf("ls(%v):Unlock: used xid(%v)\n", ls.am_primary, args.Xid)
		reply.OK = ls.xids[args.Xid]
	} else {
		locked, _ := ls.locks[args.Lockname]
		//fmt.Printf("ls(%v): Unlock before: Lock(%v), locked(%v), reply(%v)\n", ls.am_primary, args.Lockname, ls.locks[args.Lockname], reply.OK)
		if locked {
			reply.OK = true
			ls.locks[args.Lockname] = false
		} else {
			reply.OK = false
		}
		ls.xids[args.Xid] = reply.OK
		//fmt.Printf("ls(%v): Unlock after: Lock(%v), xid(%v), locked(%v), reply(%v)\n", ls.am_primary, args.Lockname, ls.xids[args.Xid], ls.locks[args.Lockname], reply.OK)
		// call B srv
		if ls.am_primary {
			var breply LockReply
			call(ls.backup, "LockServer.Unlock", args, &breply)
		}
	}
	return nil
}

//
// tell the server to shut itself down.
// for testing.
// please don't change this.
//
func (ls *LockServer) kill() {
	ls.dead = true
	ls.l.Close()
	//fmt.Printf("ls(%v): KILLed\n", ls.am_primary)
}

//
// hack to allow test_test.go to have primary process
// an RPC but not send a reply. can't use the shutdown()
// trick b/c that causes client to immediately get an
// error and send to backup before primary does.
// please don't change anything to do with DeafConn.
//
type DeafConn struct {
	c io.ReadWriteCloser
}

func (dc DeafConn) Write(p []byte) (n int, err error) {
	return len(p), nil
}
func (dc DeafConn) Close() error {
	return dc.c.Close()
}
func (dc DeafConn) Read(p []byte) (n int, err error) {
	return dc.c.Read(p)
}

func StartServer(primary string, backup string, am_primary bool) *LockServer {
	ls := new(LockServer)
	ls.backup = backup
	ls.am_primary = am_primary
	ls.locks = map[string]bool{}
	// Your initialization code here.
	ls.xids = map[int64]bool{}

	me := ""
	if am_primary {
		me = primary
	} else {
		me = backup
	}

	// tell net/rpc about our RPC server and handlers.
	rpcs := rpc.NewServer()
	rpcs.Register(ls)

	// prepare to receive connections from clients.
	// change "unix" to "tcp" to use over a network.
	os.Remove(me) // only needed for "unix"
	l, e := net.Listen("unix", me)
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	ls.l = l

	// please don't change any of the following code,
	// or do anything to subvert it.

	// create a thread to accept RPC connections from clients.
	go func() {
		for ls.dead == false {
			conn, err := ls.l.Accept()
			if err == nil && ls.dead == false {
				if ls.dying {
					// process the request but force discard of reply.

					// without this the connection is never closed,
					// b/c ServeConn() is waiting for more requests.
					go func() {
						time.Sleep(2 * time.Second)
						conn.Close()
					}()
					ls.l.Close()

					// this object has the type ServeConn expects,
					// but discards writes (i.e. discards the RPC reply).
					deaf_conn := DeafConn{c: conn}

					rpcs.ServeConn(deaf_conn)
					//fmt.Printf("ls(%v): DYING\n", ls.am_primary)
					ls.dead = true
				} else {
					go rpcs.ServeConn(conn)
				}
			} else if err == nil {
				conn.Close()
			}
			if err != nil && ls.dead == false {
				fmt.Printf("LockServer(%v) accept: %v\n", me, err.Error())
				ls.kill()
			}
		}
	}()

	return ls
}
