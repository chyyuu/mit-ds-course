package shardkv

//
// sharded key/value server.
// lots of replica groups, each running op-at-a-time paxos.
// shardmaster decides which group serves each shard.
// shardmaster may change shard assignment from time to time.
//

