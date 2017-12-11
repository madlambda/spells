// Package muxer provides functions that helps
// you implement the fan-in of concurrent operations
// by providing a way to coalesce the sink of multiple
// channels on just one.
package muxer

import "reflect"

// Do will mux all the provided source channels on the given
// sink channel. Both are interface{} since this will
// mux any type of channel (a point for generics =/).
//
// A goroutine is created to perform the muxing.
// It is a severe programming error to call this function
// with a parameter that is not a channel or with channels
// of different types.
//
// Nil channels will be ignored.
func Do(sink interface{}, sources ...interface{}) error {
	// TODO type checking
	// TODO channel direction checking
	// TODO empty sources checking

	go func() {
		sinkVal := reflect.ValueOf(sink)
		receiveCases := newCases(sources)

		for len(receiveCases) > 0 {
			chosen, recv, recvOK := reflect.Select(receiveCases)
			if recvOK {
				sinkVal.Send(recv)
				continue
			}

			receiveCases = removeClosedCase(receiveCases, chosen)
		}
		sinkVal.Close()
	}()

	return nil
}

func newCases(sources []interface{}) []reflect.SelectCase {
	cases := []reflect.SelectCase{}
	for _, source := range sources {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(source),
		})
	}
	return cases
}

func removeClosedCase(cases []reflect.SelectCase, i int) []reflect.SelectCase {
	return append(cases[:i], cases[i+1:]...)
}
