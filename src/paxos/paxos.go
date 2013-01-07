package paxos

//
// Paxos library
//
// Manages a sequence of agreed-on values.
// The set of peers is fixed.
// Copes with network failures (partition, msg loss, &c).
// Does not store anything persistently, so cannot really handle crashes.
//
// px = paxos.Make(peers []string, me string)
// px.Max() int -- highest instance seq known, or -1
// px.Min() int -- lowest instance seq known, or -1
// px.Get(seq int) (decided bool, v interface{}) -- get info about an instance
// px.Start(seq int, v interface{}) -- start agreement on new instance
// px.Done(seq int) -- ok to GC all instances <= seq
//

