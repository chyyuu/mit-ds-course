package shardmaster

//
// master configuration server.
// 256 shards, by first char of key.
// the master assigns shards to replication groups.
//
// RPC interface:
// Join(gid int64, servers []string)
// Leave(gid)
// Query(num) -> fetch Config # num, or latest config if num too big or small.
//
// client responsible for ensuring GID is unique (e.g. rand number).
// a GID must be greater than zero.
// never call ck.Join() twice with same GID!
//
// clients can use Query(num) to ask for old configs,
// so the server must keep a history back to the beginning of time.
//
// config #0 is the initial configuration, with no
// groups and all shards assigned to group 0 (the invalid group).
//


