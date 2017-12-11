package muxer_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/madlambda/spells/muxer"
)

func BenchmarkMux10(b *testing.B) {
	const jobscount = 10
	benchMux(b, jobscount)
}

func BenchmarkMux100(b *testing.B) {
	const jobscount = 100
	benchMux(b, jobscount)
}

func BenchmarkMux1000(b *testing.B) {
	const jobscount = 1000
	benchMux(b, jobscount)
}

func BenchmarkMux2500(b *testing.B) {
	const jobscount = 2500
	benchMux(b, jobscount)
}

func BenchmarkMux5000(b *testing.B) {
	const jobscount = 5000
	benchMux(b, jobscount)
}

func BenchmarkMux10000(b *testing.B) {
	const jobscount = 10000
	benchMux(b, jobscount)
}

const minTimeMilli = 100
const maxTimeMilli = 1000

func benchMux(b *testing.B, jobscount int) {
	for i := 0; i < b.N; i++ {
		muxJobs(jobscount)
	}
}

func muxJobs(jobscount int) {
	sources := []interface{}{}
	for i := 0; i < jobscount-1; i++ {
		sources = append(sources, newWorker())
	}
	sources = append(sources, newWorstCaseWorker())

	sink := make(chan time.Duration)
	if err := muxer.Do(sink, sources...); err != nil {
		panic(err)
	}

	for range sink {
	}
}

func newWorker() <-chan time.Duration {
	res := make(chan time.Duration)
	go func() {
		sleep := time.Duration(rand.Intn(maxTimeMilli-minTimeMilli) + minTimeMilli)
		time.Sleep(sleep * time.Millisecond)
		res <- sleep
		close(res)
	}()

	return res
}

func newWorstCaseWorker() <-chan time.Duration {
	res := make(chan time.Duration)
	go func() {
		sleep := time.Duration(maxTimeMilli * time.Millisecond)
		time.Sleep(sleep)
		res <- sleep
		close(res)
	}()
	return res
}
