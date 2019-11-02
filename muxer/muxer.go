// Package muxer provides functions that helps
// you implement the fan-in of concurrent operations
// by providing a way to coalesce the sink of multiple
// channels on just one.
package muxer

import (
	"errors"
	"fmt"
	"reflect"
)

// Do will mux all the provided source channels on the given
// sink channel. All channels are interface{} since this will
// mux any type of channel (a point for generics =/).
//
// A goroutine is created to perform the muxing.
// It is a programming error to call this function
// with a parameter that is not a channel or with channels
// of different types.
//
// The source channels will be used only for reading.
// While there is an open source channel the sink channel
// will also remain open.
//
// The provided sink channel will be closed by the muxer
// when all source channels are closed, so
// it must be used ONLY for reading operations,
// never write on it or close it.
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
	if sink == nil {
		return errors.New("sink channel is a nil empty interface")
	}

	sinktype := reflect.TypeOf(sink)
	if sinktype.Kind() != reflect.Chan {
		return fmt.Errorf("sink has invalid type[%s] kind[%s]", sinktype, sinktype.Kind())
	}

	if sinktype.ChanDir() == reflect.RecvDir {
		return errors.New("sink channel is receive only, sink channels MUST be able to send")
	}

	if reflect.ValueOf(sink).IsNil() {
		return errors.New("sink channel is nil")
	}

	for i, source := range sources {
		if source == nil {
			return fmt.Errorf("invalid nil source channel at position[%d]", i)
		}

		sourcetype := reflect.TypeOf(source)
		if sourcetype.Kind() != reflect.Chan {
			return fmt.Errorf("source[%d] has invalid type[%s] kind[%s]", i, sourcetype, sourcetype.Kind())
		}

		if sourcetype.ChanDir() == reflect.SendDir {
			return errors.New("source channel is send only, source channels MUST be able to receive")
		}

		if reflect.ValueOf(source).IsNil() {
			return errors.New("source channel is nil")
		}

		if sourcetype.Elem() != sinktype.Elem() {
			return fmt.Errorf("source[%d] is [chan %s] but sink is [chan %s]", i, sourcetype.Elem(), sinktype.Elem())
		}
	}

	return nil
}
