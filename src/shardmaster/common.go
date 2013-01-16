package shardmaster

//
// master configuration server.
// 256 shards, by first char of key.
// master's main job is to assign shards to replication groups.
//
// RPC interface:
// Join(gid int64, servers []string)
// Leave(gid)
// Query() -> Config
//
// client responsible for ensuring GID is unique (e.g. rand number)
// never call ck.Join() twice with same GID!
//


