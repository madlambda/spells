// Package muxer provides functions that helps
// you implement the fan-in of concurrent operations
// by providing a way to coalesce the output of multiple
// channels on just one.
package muxer

// Do will mux all the provided inputs channels on the given
// output channel. Both are interface{} since this will
// mux any type of channel (a point for generics =/).
//
// It is a severe programming error to call this function
// with a parameter that is not a channel or with channels
// of different types.
//
// Nil channels will be ignored.
func Do(output interface{}, inputs ...interface{}) error {
	return nil
}
