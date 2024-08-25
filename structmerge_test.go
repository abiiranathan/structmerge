package structmerge

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

type Address struct {
	Street  string
	City    string
	Country string
}

type TestStruct struct {
	Name    string
	Age     int
	Address Address
	Active  bool
	Count   uint
	hidden  string
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name     string
		dst      TestStruct
		src      TestStruct
		cfg      Config
		expected TestStruct
		wantErr  bool
	}{
		{
			name: "IncludeAll",
			dst: TestStruct{
				Name: "Alice",
				Age:  30,
				Address: Address{
					Street: "123 Old St",
					City:   "Old City",
				},
			},
			src: TestStruct{
				Name: "Bob",
				Age:  25,
				Address: Address{
					Street:  "456 New St",
					City:    "New City",
					Country: "New Country",
				},
				Active: true,
				Count:  5,
			},
			cfg: Config{Option: IncludeAll},
			expected: TestStruct{
				Name: "Bob",
				Age:  25,
				Address: Address{
					Street:  "456 New St",
					City:    "New City",
					Country: "New Country",
				},
				Active: true,
				Count:  5,
			},
		},
		{
			name: "ExcludeEmpty",
			dst: TestStruct{
				Name: "Alice",
				Age:  30,
				Address: Address{
					Street: "123 Old St",
					City:   "Old City",
				},
				Count: 10,
			},
			src: TestStruct{
				Name: "",
				Age:  0,
				Address: Address{
					Street:  "",
					City:    "New City",
					Country: "New Country",
				},
				Active: true,
				Count:  0,
			},
			cfg: Config{Option: ExcludeEmpty},
			expected: TestStruct{
				Name: "Alice",
				Age:  30,
				Address: Address{
					Street:  "123 Old St",
					City:    "New City",
					Country: "New Country",
				},
				Active: true,
				Count:  10,
			},
		},
		{
			name: "OverwriteEmpty",
			dst: TestStruct{
				Name: "Alice",
				Age:  0,
				Address: Address{
					Street: "",
					City:   "Old City",
				},
				Count: 0,
			},
			src: TestStruct{
				Name: "Bob",
				Age:  25,
				Address: Address{
					Street:  "456 New St",
					City:    "New City",
					Country: "New Country",
				},
				Active: true,
				Count:  5,
			},
			cfg: Config{Option: OverwriteEmpty},
			expected: TestStruct{
				Name: "Alice",
				Age:  25,
				Address: Address{
					Street:  "456 New St",
					City:    "Old City",
					Country: "New Country",
				},
				Active: true,
				Count:  5,
			},
		},
		{
			name: "Include",
			dst: TestStruct{
				Name: "Alice",
				Age:  30,
				Address: Address{
					Street: "123 Old St",
					City:   "Old City",
				},
				Count: 10,
			},
			src: TestStruct{
				Name: "Bob",
				Age:  25,
				Address: Address{
					Street:  "456 New St",
					City:    "New City",
					Country: "New Country",
				},
				Active: true,
				Count:  5,
			},
			cfg: Config{Option: IncludeAll, Include: []string{"Name", "Address.Street", "Address.Country", "Count"}},
			expected: TestStruct{
				Name: "Bob",
				Age:  30,
				Address: Address{
					Street:  "456 New St",
					City:    "Old City",
					Country: "New Country",
				},
				Active: false,
				Count:  5,
			},
		},
		{
			name: "Exclude",
			dst: TestStruct{
				Name: "Alice",
				Age:  30,
				Address: Address{
					Street: "123 Old St",
					City:   "Old City",
				},
				Count: 10,
			},
			src: TestStruct{
				Name: "Bob",
				Age:  25,
				Address: Address{
					Street:  "456 New St",
					City:    "New City",
					Country: "New Country",
				},
				Active: true,
				Count:  5,
			},
			cfg: Config{Option: IncludeAll, Exclude: []string{"Name", "Age", "Address.City"}},
			expected: TestStruct{
				Name: "Alice",
				Age:  30,
				Address: Address{
					Street:  "456 New St",
					City:    "Old City",
					Country: "New Country",
				},
				Active: true,
				Count:  5,
			},
		},
		{
			name: "UintHandling",
			dst: TestStruct{
				Name:  "Alice",
				Age:   30,
				Count: 0,
			},
			src: TestStruct{
				Name:  "Bob",
				Age:   25,
				Count: 5,
			},
			cfg: Config{Option: ExcludeEmpty},
			expected: TestStruct{
				Name:  "Bob",
				Age:   25,
				Count: 5,
			},
		},
		{
			name: "HiddenFields",
			dst: TestStruct{
				Name:   "Alice",
				Age:    30,
				hidden: "old",
			},
			src: TestStruct{
				Name:   "Bob",
				Age:    25,
				hidden: "new",
			},
			cfg: Config{Option: IncludeAll},
			expected: TestStruct{
				Name:   "Bob",
				Age:    25,
				hidden: "old", // This should not change
			},
		},
		{
			name: "NestedStructMerge",
			dst: TestStruct{
				Name: "Alice",
				Age:  30,
				Address: Address{
					Street: "123 Old St",
					City:   "Old City",
				},
			},
			src: TestStruct{
				Name: "Bob",
				Age:  25,
				Address: Address{
					Street:  "456 New St",
					Country: "New Country",
				},
			},
			cfg: Config{Option: IncludeAll},
			expected: TestStruct{
				Name: "Bob",
				Age:  25,
				Address: Address{
					Street:  "456 New St",
					City:    "",
					Country: "New Country",
				},
			},
		},
		{
			name: "NestedStructExcludeEmpty",
			dst: TestStruct{
				Name: "Alice",
				Age:  30,
				Address: Address{
					Street: "123 Old St",
					City:   "Old City",
				},
			},
			src: TestStruct{
				Name: "Bob",
				Age:  25,
				Address: Address{
					Street:  "",
					Country: "New Country",
				},
			},
			cfg: Config{Option: ExcludeEmpty},
			expected: TestStruct{
				Name: "Bob",
				Age:  25,
				Address: Address{
					Street:  "123 Old St",
					City:    "Old City",
					Country: "New Country",
				},
			},
		},
		{
			name: "NestedStructInclude",
			dst: TestStruct{
				Name: "Alice",
				Age:  30,
				Address: Address{
					Street: "Old Street",
					City:   "Old City",
				},
			},
			src: TestStruct{
				Name: "Bob",
				Age:  25,
				Address: Address{
					Street:  "New Street",
					City:    "New City",
					Country: "New Country",
				},
			},
			cfg: Config{Option: IncludeAll, Include: []string{"Name", "Address.Street", "Address.Country"}},
			expected: TestStruct{
				Name: "Bob",
				Age:  30,
				Address: Address{
					Street:  "New Street",
					City:    "Old City",
					Country: "New Country",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Merge(&tt.dst, tt.src, tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Merge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.dst, tt.expected) {
				t.Errorf("Merge() = %#v, want %#v", tt.dst, tt.expected)
			}
		})
	}
}

func TestMergeErrors(t *testing.T) {
	tests := []struct {
		name    string
		dst     interface{}
		src     interface{}
		cfg     Config
		wantErr error
	}{
		{
			name:    "Invalid destination",
			dst:     TestStruct{},
			src:     TestStruct{},
			cfg:     Config{},
			wantErr: ErrInvalidDestination,
		},
		{
			name:    "Invalid source",
			dst:     &TestStruct{},
			src:     &TestStruct{},
			cfg:     Config{},
			wantErr: ErrInvalidSource,
		},
		{
			name:    "Type mismatch",
			dst:     &TestStruct{},
			src:     struct{ Foo string }{},
			cfg:     Config{},
			wantErr: ErrTypeMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Merge(tt.dst, tt.src, tt.cfg)
			if err != tt.wantErr {
				t.Errorf("Merge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type Person struct {
	Name    string
	Age     int
	Address *Address
	Active  bool
}

func TestMergeWithPointer(t *testing.T) {
	address1 := &Address{
		Street:  "123 Old St",
		City:    "Old City",
		Country: "Old Country",
	}

	person1 := Person{
		Name:    "Alice",
		Age:     25,
		Address: address1,
		Active:  true,
	}

	address2 := &Address{
		Street:  "456 New St",
		City:    "New City",
		Country: "",
	}

	person2 := Person{
		Name:    "Bob",
		Age:     30,
		Address: address2,
		Active:  false,
	}

	cfg := Config{
		Option: IncludeAll,
	}

	err := Merge(&person1, person2, cfg)
	if err != nil {
		t.Errorf("merge failed: %v\n", err)
	}

	expected := Person{
		Name:    "Bob",
		Age:     30,
		Address: address2,
		Active:  false,
	}

	if !reflect.DeepEqual(expected, person1) {
		t.Errorf("expected %#v, got %#v\n", expected, person1)
	}
}

type Floats struct {
	Value  float32
	Value2 float64
}

func TestMergeWithFloats(t *testing.T) {
	f1 := Floats{
		Value:  4.6,
		Value2: 100.0,
	}

	var f2 Floats

	err := Merge(&f2, f1, Config{Option: IncludeAll})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(f1, f2) {
		t.Errorf("expected f2 to equal %v, got %v\n", f1, f2)
	}
}

type Date time.Time

func (d *Date) Merge(src reflect.Value) error {
	srcDate, ok := src.Interface().(Date)
	if !ok {
		return ErrInvalidSource
	}
	*d = srcDate
	return nil
}

type Plan struct {
	Date Date
	Data json.RawMessage
}

func TestCopyTime(t *testing.T) {
	var t1, t2 time.Time
	t1 = time.Now()

	err := Merge(&t2, t1, Config{Option: IncludeAll})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(t1, t2) {
		t.Errorf("expected t2 to equal %v, got %v\n", t1, t2)
	}

	var d1, d2 Date

	d1 = Date(time.Now())
	Merge(&d2, d1, Config{Option: IncludeAll})

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("expected d2 to equal %v, got %v\n", d1, d2)
	}

	// field implementing Merger
	var p1, p2 Plan
	p1.Date = Date(time.Now())
	p1.Data = json.RawMessage([]byte("Hello world"))

	err = Merge(&p2, p1, Config{Option: IncludeAll})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(p1, p2) {
		t.Fatalf("p1 and p2 are not equal")
	}
}
