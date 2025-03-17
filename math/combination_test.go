package math

import (
	"reflect"
	"testing"
)

func Test_combination(t *testing.T) {
	type args[T any] struct {
		data []T
		m    int
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want [][]T
	}
	tests := []testCase[int]{
		{name: "c(3,2)", args: args[int]{data: []int{1, 2, 3}, m: 2}, want: [][]int{{1, 2}, {1, 3}, {2, 3}}},
		{name: "c(2,1)", args: args[int]{data: []int{1, 2}, m: 1}, want: [][]int{{1}, {2}}},
		{name: "c(4,2)", args: args[int]{data: []int{1, 2, 3, 4}, m: 2}, want: [][]int{{1, 2}, {1, 3}, {2, 3}, {1, 4}, {2, 4}, {3, 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Combination(tt.args.data, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Combination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bitChange(t *testing.T) {
	type args struct {
		n     uint64
		start int
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{name: "case: 1111-1011", args: args{n: uint64(0xFB), start: 8}, want: uint64(0xF7)},
		{name: "case: 0011-1011", args: args{n: uint64(0x3B), start: 8}, want: uint64(0xC7)},
		{name: "case: 0011-1111", args: args{n: uint64(0x3F), start: 8}, want: uint64(0x3F)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bitChange(tt.args.n, tt.args.start); got != tt.want {
				t.Errorf("bitChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cartesianProduct(t *testing.T) {
	type args[T any] struct {
		arrays [][]T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want [][]T
	}
	tests := []testCase[int]{
		{name: "case:两个数组", args: args[int]{arrays: [][]int{{1, 2}, {3, 4}}}, want: [][]int{{1, 3}, {1, 4}, {2, 3}, {2, 4}}},
		{name: "case:多个数组", args: args[int]{arrays: [][]int{{1, 2, 3}, {4}, {5, 6}}}, want: [][]int{{1, 4, 5}, {1, 4, 6}, {2, 4, 5}, {2, 4, 6}, {3, 4, 5}, {3, 4, 6}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CartesianProduct(tt.args.arrays); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CartesianProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_comb(t *testing.T) {
	type args struct {
		n int
		k int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "C(6,4)", args: args{n: 6, k: 4}, want: 15},
		{name: "C(8,3)", args: args{n: 8, k: 3}, want: 56},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := comb(tt.args.n, tt.args.k); got != tt.want {
				t.Errorf("comb() = %v, want %v", got, tt.want)
			}
		})
	}
}
