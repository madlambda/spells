// Package muxer provides functions that helps
// you implement the fan-in of concurrent operations
// by providing a way to coalesce the sink of multiple
// channels on just one.
package muxer

import (
	"fmt"
	"reflect"
)

// Do will mux all the provided source channels on the given
// sink channel. Both are interface{} since this will
// mux any type of channel (a point for generics =/).
//
// A goroutine is created to perform the muxing.
// It is a severe programming error to call this function
// with a parameter that is not a channel or with channels
// of different types.
//
// The source channels will be used only for reading.
// While there is an open source channel the sink channel
// will also remain closed.
//
// The provided sink channel will be closed by the muxer
// when all source channels are closed, so
// it must be used ONLY for reading operations (never close it).
//
// The sink and source channels must transport values of the same
// type. No nil channels are allowed on the parameters.
func Do(sink interface{}, sources ...interface{}) error {

	if err := checkParams(sink, sources); err != nil {
		return err
	}

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

func checkParams(sink interface{}, sources []interface{}) error {
	sinktype := reflect.TypeOf(sink)
	if sinktype.Kind() != reflect.Chan {
		return fmt.Errorf("sink has invalid type[%s] kind[%s]", sinktype, sinktype.Kind())
	}

	for i, source := range sources {
		sourcetype := reflect.TypeOf(source)
		if sourcetype.Kind() != reflect.Chan {
			return fmt.Errorf(
				"source[%d] has invalid type[%s] kind[%s]",
				i,
				sourcetype,
				sourcetype.Kind(),
			)
		}
		if sourcetype.Elem() != sinktype.Elem() {
			return fmt.Errorf(
				"source[%d] is [chan %s] but sink is [chan %s]",
				i,
				sourcetype.Elem(),
				sinktype.Elem(),
			)
		}
	}

	return nil
}
