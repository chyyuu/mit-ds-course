package paxos

//
// Paxos library
//
// Manages a sequence of agreed-on values.
// The set of peers is fixed.
//
// px = paxos.Make(peers []string, me string)
// px.Max() int -- highest instance seq known
// px.Get(seq int) (decided bool, v interface{}) -- most recent decided instance
// px.Start(seq int, v interface{}) -- start agreement on new instance
//

