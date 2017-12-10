// Package muxer provides functions that helps
// you implement the fan-in of concurrent operations
// by providing a way to coalesce the output of multiple
// channels on just one.
package muxer

import "reflect"

// Do will mux all the provided inputs channels on the given
// output channel. Both are interface{} since this will
// mux any type of channel (a point for generics =/).
//
// A goroutine is created to perform the muxing.
// It is a severe programming error to call this function
// with a parameter that is not a channel or with channels
// of different types.
//
// Nil channels will be ignored.
func Do(output interface{}, inputs ...interface{}) error {
	// TODO type checking
	// TODO channel direction checking
	// TODO empty inputs checking
	outputVal := reflect.ValueOf(output)
	inputVal := reflect.ValueOf(inputs[0])

	go func() {
		// TODO handle close
		for {
			v, ok := inputVal.Recv()
			if !ok {
				outputVal.Close()
				return
			}
			outputVal.Send(v)
		}
	}()

	return nil
}
