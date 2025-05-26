package lemin_test

import (
	"fmt"
	"lemin/pkg/lemin"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Set the default to the maximum value time.Duration can take
const defaultDuration time.Duration = 1<<63 - 1

type result struct {
	t []lemin.Turn
	e error
}

type testCase struct {
	filepath string
	timeout  time.Duration
	turns    int
	err      error
}

func run(file *os.File) result {
	farm, err := lemin.ParseFromFile(file)
	if err != nil {
		return result{t: nil, e: err}
	}

	turns, err := lemin.Lemin(farm)
	if err != nil {
		return result{t: nil, e: err}
	}

	return result{t: turns, e: nil}
}

func runTest(tC testCase) func(t *testing.T) {
	return func(t *testing.T) {
		file, err := os.Open(tC.filepath)
		if err != nil {
			t.Fatal(err)
		}

		if tC.timeout == 0 {
			tC.timeout = defaultDuration
		}

		result := make(chan result, 1)
		go func() {
			result <- run(file)
		}()

		select {
		case <-time.After(tC.timeout):
			t.Fatalf("expected %s to take at most %v", tC.filepath, tC.timeout)
		case result := <-result:
			turns := result.t
			err := result.e

			if turns == nil && tC.err != err {
				t.Fatalf("expected %s to fail with '%s' but was '%s'", tC.filepath, err.Error(), tC.err.Error())
			}

			if tC.turns != 0 && len(turns) > tC.turns {
				t.Fatalf("expected %s to be %d turns at most but was %d", tC.filepath, tC.turns, len(turns))
			}
		}
	}
}

func TestGoodSamples(t *testing.T) {
	fmt.Println(os.Getwd())
	testCases := []testCase{
		{
			filepath: "../../samples/example00",
			turns:    6,
		},
		{
			filepath: "../../samples/example01",
			turns:    8,
		},
		{
			filepath: "../../samples/example02",
			turns:    11,
		},
		{
			filepath: "../../samples/example03",
			turns:    6,
		},
		{
			filepath: "../../samples/example04",
			turns:    6,
		},
		{
			filepath: "../../samples/example05",
			turns:    8,
		},
		{
			filepath: "../../samples/example06",
			timeout:  time.Minute + 30*time.Second,
		},
		{
			filepath: "../../samples/example07",
			timeout:  2*time.Minute + 30*time.Second,
		},
	}

	for _, tC := range testCases {
		t.Run(filepath.Base(tC.filepath), runTest(tC))
	}
}

func TestBadSamples(t *testing.T) {
	testCases := []testCase{
		{
			filepath: "../../samples/badapple00",
			err:      lemin.ErrCommandIsNotAllowed,
		},

		{
			filepath: "../../samples/badexample00",
			err:      lemin.ErrInvalidAntCount,
		},
		{
			filepath: "../../samples/badexample01",
			err:      lemin.ErrInsaneFarm,
		},
	}

	for _, tC := range testCases {
		t.Run(filepath.Base(tC.filepath), runTest(tC))
	}
}
