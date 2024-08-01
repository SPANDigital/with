//go:build samples

package catfacts

import (
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
	"reflect"
	"testing"
	"time"
)

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func Test_catFacts_GetRandomFacts(t *testing.T) {
	// Start our recorder
	r, err := recorder.New("fixtures/cat-facts")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	if r.Mode() != recorder.ModeRecordOnce {
		t.Fatal("Recorder should be in ModeRecordOnce")
	}

	catFacts, err := NewCatFacts(
		WithHttpClient(r.GetDefaultClient()),
	)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		limit *int
	}
	tests := []struct {
		name         string
		args         args
		wantCatFacts []*CatFact
		wantErr      bool
	}{
		{
			name: "get random facts default limit",
			args: args{
				limit: nil,
			},
			wantCatFacts: []*CatFact{
				{
					Id:        "63cc0ffeeec42e60ec323297",
					V:         0,
					Text:      "Cats love sleeping.",
					UpdatedAt: must(time.Parse(time.RFC3339, "2023-01-21T16:17:02.873Z")),
					Deleted:   false,
					Source:    "",
					SentCount: 0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCatFacts, err := catFacts.GetRandomFacts(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRandomFacts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCatFacts, tt.wantCatFacts) {
				t.Errorf("GetRandomFacts() gotCatFacts = %v, want %v", gotCatFacts, tt.wantCatFacts)
			}
		})
	}
}
