package tests

import (
	"lemin/pkg/lemin"
	"os"
	"testing"
	"time"
)

type testCase struct {
	filepath string
	timeout  time.Duration
	turns    int
	err      error
}

func run(file *os.File) ([]lemin.Turn, error) {
	farm, err := lemin.ParseFromFile(file)
	if err != nil {
		return nil, err
	}

	turns, err := lemin.Lemin(farm)
	if err != nil {
		return nil, err
	}

	return turns, nil
}

func runTest(tC testCase) func(t *testing.T) {
	return func(t *testing.T) {
		file, err := os.Open(tC.filepath)
		if err != nil {
			t.Fatal(err)
		}

		start := time.Now()
		turns, err := run(file)
		taken := time.Since(start)

		if tC.timeout != 0 && taken > tC.timeout {
			t.Fatalf("expected %s to take at most %v but took %v", tC.filepath, tC.timeout, taken)
		}

		if turns == nil && tC.err != err {
			t.Fatalf("expected %s to fail with '%s' but was '%s'", tC.filepath, err.Error(), tC.err.Error())
		}

		if tC.turns != 0 && len(turns) > tC.turns {
			t.Fatalf("expected %s to be %d turns at most but was %d", tC.filepath, tC.turns, len(turns))
		}
	}
}

func TestBadSamples(t *testing.T) {
	testCases := []testCase{
		{
			filepath: "../samples/badapple00",
			err:      lemin.ErrCommandIsNotAllowed,
		},
		{
			filepath: "../samples/badexample00",
			err:      lemin.ErrInvalidAntCount,
		},
		{
			filepath: "../samples/badexample01",
			err:      lemin.ErrInsaneFarm,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.filepath, runTest(tC))
	}
}

func TestGoodSamples(t *testing.T) {
	testCases := []testCase{
		{
			filepath: "../samples/example00",
			turns:    6,
		},
		{
			filepath: "../samples/example01",
			turns:    8,
		},
		{
			filepath: "../samples/example02",
			turns:    11,
		},
		{
			filepath: "../samples/example03",
			turns:    6,
		},
		{
			filepath: "../samples/example04",
			turns:    6,
		},
		{
			filepath: "../samples/example05",
			turns:    8,
		},
		{
			filepath: "../samples/example06",
			timeout:  time.Minute + 30*time.Second,
		},
		{
			filepath: "../samples/example07",
			timeout:  2*time.Minute + 30*time.Second,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.filepath, runTest(tC))
	}
}
