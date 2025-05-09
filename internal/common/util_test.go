package common

import (
	"testing"
)

func TestCheckLuhn(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"Wikipedia example", args{"4561261212345467"}, true, false},
		{"Amex example", args{"374245455400126"}, true, false},
		{"Visa example", args{"4263982640269299"}, true, false},
		{"Mastercard example", args{"5425233430109903"}, true, false},
		{"Error example 1", args{"5425233430109902"}, false, false},
		{"Error example 2", args{"123"}, false, false},
		{"Error example incorrect numeric string", args{"a"}, false, true},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckLuhn(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckLuhn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckLuhn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckForAllDigits(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"ok", args{"123"}, true},
		{"error1", args{"a123"}, false},
		{"error2", args{"1a23"}, false},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckForAllDigits(tt.args.s); got != tt.want {
				t.Errorf("CheckForAllDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckOrderNumberFormat(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"Ok1", args{"4561261212345467"}, true, false},
		{"Ok2", args{"374245455400126"}, true, false},
		{"Letters in number", args{"abc123"}, false, false},
		{"Wrong Luhn check", args{"123456123456"}, false, false},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckOrderNumberFormat(tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckOrderNumberFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckOrderNumberFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

type args[T any] struct {
	m         map[string]T
	predicate func(T) bool
}

func TestFilterMap(t *testing.T) {
	mapStringInt := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	tests := []struct {
		name string
		args args[int]
		want []int
	}{
		{"Test1", args[int]{mapStringInt, func(x int) bool { return x == 2 }}, []int{2}},
		{"Test2", args[int]{mapStringInt, func(x int) bool { return x > 1 }}, []int{2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterMap(tt.args.m, tt.args.predicate)
			if !equalIgnoreOrder(got, tt.want) {
				t.Errorf("FilterMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func equalIgnoreOrder[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[T]int)
	for _, x := range a {
		counts[x]++
	}
	for _, x := range b {
		if counts[x] == 0 {
			return false
		}
		counts[x]--
	}
	return true
}
