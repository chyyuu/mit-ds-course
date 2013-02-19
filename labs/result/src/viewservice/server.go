package viewservice

import "net"
import "net/rpc"
import "log"
import "time"
import "sync"
import "fmt"
import "os"

type ViewServer struct {
  mu sync.Mutex
  l net.Listener
  dead bool
  me string


  // Your declarations here.
  pingtime      map[string]time.Time   // server id: ping time
  viewnum       map[string]uint        // server id: viewnum
  curview       View                   // the current view from view server
  psack         bool                   // the pri srv acked?  according the pingargs.viewnum 
  nextviewnum   uint                   // the next server's viewnum
  nextprisrv    string                 // the next primary server
  nextbaksrv    string                 // the next backup server
}

//
// server Ping RPC handler.
//
func (vs *ViewServer) Ping(args *PingArgs, reply *PingReply) error {

  // Your code here.
  vs.mu.Lock()
  defer vs.mu.Unlock()

  //fmt.Printf("vs:ping: curview: %t ||\n",vs.curview)
  if  vs.psack == false {
       if (vs.curview.Primary == args.Me || vs.nextprisrv == args.Me) &&
           vs.viewnum[vs.curview.Primary] == vs.curview.Viewnum  {
                //fmt.Printf("  vs: cond 1: psack = true\n")
                vs.psack = true
       } else {
          if vs.viewnum[vs.curview.Primary] == 0 && vs.curview.Viewnum == args.Viewnum { // pri srv first run
                //fmt.Printf("  vs: cond 2: psack = true\n")
		            vs.psack = true
		      }
	     }
  }
  
  vs.pingtime[args.Me] = time.Now()
  vs.viewnum[args.Me]  = args.Viewnum
  reply.View = vs.curview
  //fmt.Printf("vs: psack: %t end\n",vs.psack)
  return nil
}

// 
// server Get() RPC handler.
//
func (vs *ViewServer) Get(args *GetArgs, reply *GetReply) error {

  // Your code here.
  //fmt.Printf("vs:get arg (%t)\n",*args)
  reply.View.Primary = vs.curview.Primary
  reply.View.Backup  = vs.curview.Backup
  reply.View.Viewnum = vs.curview.Viewnum
  //fmt.Printf("vs:get reply (%t)\n",*reply)
  return nil
}


//
// tick() is called once per PingInterval; it should notice
// if servers have died or recovered, and change the view
// accordingly.
//
func (vs *ViewServer) tick() {

  // Your code here.
  vs.mu.Lock()
  defer vs.mu.Unlock()

  //fmt.Printf("vs:tick:BEG curview: %t ||\n",vs.curview)
  //fmt.Print("vs:tick BEG-- ")
  //fmt.Print("psvnum: ", vs.viewnum[vs.curview.Primary]," || ");
  //fmt.Print("psack: ", vs.psack, " || ");
	//fmt.Print("nvnum :", vs.nextviewnum, " || ")
  //fmt.Print("nextps: ",vs.nextprisrv," || ")
  //fmt.Print("nextbs: ",vs.nextbaksrv," || ")
  //fmt.Print("cvnum: ", vs.curview.Viewnum, " || ");
	//fmt.Print("prisrv: ", vs.curview.Primary, " || ");
	//fmt.Print("baksrv: ", vs.curview.Backup, "\n");  

  // no pri srv, prepare next pri srv
  if	vs.curview.Primary == "" &&	vs.curview.Viewnum == vs.nextviewnum {
				//fmt.Println("   vs:tick: search next primary server");
				nextprisrv := ""
				for server, pingtime := range vs.pingtime {
						if time.Now().Sub(pingtime) < DeadLine && server != "" {
								 nextprisrv = server
                            break
						}
				}
				if nextprisrv != ""{
						vs.nextprisrv = nextprisrv
						vs.nextviewnum = vs.curview.Viewnum + 1;
						vs.psack = false
				}
	}

  // no bak srv, prepare next bak srv   
	if  vs.curview.Backup == "" && vs.curview.Viewnum == vs.nextviewnum {
			  //fmt.Println("   vs:tick: search next backup server");
				nextbaksrv := ""
				for server, pingtime := range vs.pingtime {
						if	time.Now().Sub(pingtime) < DeadLine && 	server != "" &&
							  server != vs.curview.Primary  {
									nextbaksrv = server
                            break
						}
				}
				if nextbaksrv != "" {
						vs.nextbaksrv = nextbaksrv
						vs.nextviewnum = vs.curview.Viewnum + 1
						vs.psack = false
				}
		}

  // check pri/bak srv dead or restart
  pridead, bakdead :=false,false
  // pri srv dead or restart?    
  if  vs.curview.Primary != "" && 
     (time.Now().Sub(vs.pingtime[vs.curview.Primary]) > DeadLine ||  // pri srv dead
		  vs.viewnum[vs.curview.Primary] == 0) {                         // pri srv restart
				pridead =	true
	}

  // bak srv dead or restart?
	if  vs.curview.Backup != "" &&
      (time.Now().Sub(vs.pingtime[vs.curview.Backup]) > DeadLine ||  // bak srv dead
						vs.viewnum[vs.curview.Backup] == 0) {                    // bak srv restart
				bakdead =	true
	}
  
  //if pridead { 
     //if vs.viewnum[vs.curview.Primary] == 0{
         //fmt.Println("   vs:tick: pri srv restart") 
     //} else {
         //fmt.Println("   vs:tick: pri srv dead") 
     //}
  //}
  
  //if bakdead {
     //if vs.viewnum[vs.curview.Backup] == 0{
         //fmt.Println("   vs:tick: bak srv restart") 
     //} else {
         //fmt.Println("   vs:tick: bak srv dead") 
     //}
  //}
  
  // pri srv need change    
  if pridead && vs.curview.Viewnum == vs.nextviewnum {
			if vs.curview.Backup != ""{
				//fmt.Println("   vs:tick: pri srv dead or restart, change to backup")
				vs.nextprisrv = vs.curview.Backup;
				vs.nextviewnum = vs.curview.Viewnum + 1
				vs.psack = false
			}
	}
 
  // bak srv need change    
  if bakdead && vs.curview.Viewnum == vs.nextviewnum {
		 //fmt.Println("   vs:tick: bak srv dead or restart")
		 vs.curview.Backup = ""
		 vs.nextviewnum = vs.curview.Viewnum + 1
		 vs.psack = false
	}

  // pri srv ack the last view, and nexeviewnum updated, so change new view now
	if vs.psack == true && vs.nextviewnum != vs.curview.Viewnum {
     //fmt.Println("   vs:tick: the curview num need change!")
		 if vs.nextprisrv != vs.curview.Primary && vs.nextprisrv != "" {
            //fmt.Println("       vs:tick: curview change pri srv!")
						vs.curview.Primary = vs.nextprisrv
						vs.nextprisrv = ""
				}
				if vs.nextbaksrv != vs.curview.Backup && vs.nextbaksrv != "" {
            //fmt.Println("       vs:tick: curview change bak srv!")
						vs.curview.Backup = vs.nextbaksrv;
						vs.nextbaksrv = "";
				}
        
				if vs.curview.Primary == vs.curview.Backup {
           //fmt.Println("       vs:tick: curview pri == srv!")
					 vs.curview.Backup = ""
				}
				vs.curview.Viewnum = vs.nextviewnum
	}

  //fmt.Printf("vs:tick:END curview: %t ||\n",vs.curview)
  //fmt.Print("vs:tick END--")
  //fmt.Print("psvnum: ", vs.viewnum[vs.curview.Primary]," || ");
  //fmt.Print("psack: ", vs.psack, " || ");
  //fmt.Print("nvnum :", vs.nextviewnum, " || ")
  //fmt.Print("nextps: ",vs.nextprisrv," || ")
  //fmt.Print("nextbs: ",vs.nextbaksrv," || ")
  //fmt.Print("cvnum: ", vs.curview.Viewnum, " || ");
	//fmt.Print("prisrv: ", vs.curview.Primary, " || ");
	//fmt.Print("baksrv: ", vs.curview.Backup, "|| ");  
	//fmt.Println("");  
}

//
// tell the server to shut itself down.
// for testing.
// please don't change this function.
//
func (vs *ViewServer) Kill() {
  vs.dead = true
  vs.l.Close()
}

func StartServer(me string) *ViewServer {
  vs := new(ViewServer)
  vs.me = me
  // Your vs.* initializations here.
  vs.pingtime = map[string]time.Time{}
  vs.viewnum = map[string]uint{}
  vs.curview.Viewnum = 0
  vs.curview.Primary = ""
  vs.curview.Backup = ""
  vs.psack = false
  vs.nextviewnum=0
  
  // tell net/rpc about our RPC server and handlers.
  rpcs := rpc.NewServer()
  rpcs.Register(vs)

  // prepare to receive connections from clients.
  // change "unix" to "tcp" to use over a network.
  os.Remove(vs.me) // only needed for "unix"
  l, e := net.Listen("unix", vs.me);
  if e != nil {
    log.Fatal("listen error: ", e);
  }
  vs.l = l

  // please don't change any of the following code,
  // or do anything to subvert it.

  // create a thread to accept RPC connections from clients.
  go func() {
    for vs.dead == false {
      conn, err := vs.l.Accept()
      if err == nil && vs.dead == false {
        go rpcs.ServeConn(conn)
      } else if err == nil {
        conn.Close()
      }
      if err != nil && vs.dead == false {
        fmt.Printf("ViewServer(%v) accept: %v\n", me, err.Error())
        vs.Kill()
      }
    }
  }()

  // create a thread to call tick() periodically.
  go func() {
    for vs.dead == false {
      vs.tick()
      time.Sleep(PingInterval)
    }
  }()

  return vs
}
