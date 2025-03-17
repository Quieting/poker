package shisanshui

import (
	"sort"
	"testing"
)

// func TestMauBinhList(t *testing.T) {
// 	type args struct {
// 		cards []card
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want [][]*Group
// 	}{
//
// 		{name: "case:多牌型", args: args{cards: []card{c2, dq, c9, h2, h8, h7, ha, hj, s4, c4, sk, ck, dk}}, want: [][]*Group{
// 			// 第三道葫芦
// 			{cardsToGroup([]card{s4, c4, sk, ck, dk}), cardsToGroup([]card{h2, h8, h7, ha, hj}), cardsToGroup([]card{c2, dq, c9})},
// 			{cardsToGroup([]card{h2, c2, sk, ck, dk}), cardsToGroup([]card{s4, c4})},
// 			{cardsToGroup([]card{s4, c4, sk, ck, dk}), cardsToGroup([]card{h2, c2})},
//
// 			// 第三道同花
// 			{cardsToGroup([]card{h2, h8, h7, hj, ha}), cardsToGroup([]card{sk, ck, dk}), cardsToGroup([]card{s4, c4})},
// 			{cardsToGroup([]card{h2, h8, h7, hj, ha}), cardsToGroup([]card{sk, ck, s4, c4})},
// 			{cardsToGroup([]card{h2, h8, h7, hj, ha}), cardsToGroup([]card{sk, dk, s4, c4})},
// 			{cardsToGroup([]card{h2, h8, h7, hj, ha}), cardsToGroup([]card{ck, dk, s4, c4})},
// 			{cardsToGroup([]card{h2, h8, h7, hj, ha}), cardsToGroup([]card{ck, dk}), cardsToGroup([]card{s4, c4})},
// 			{cardsToGroup([]card{h2, h8, h7, hj, ha}), cardsToGroup([]card{sk, ck}), cardsToGroup([]card{s4, c4})},
// 			{cardsToGroup([]card{h2, h8, h7, hj, ha}), cardsToGroup([]card{sk, dk}), cardsToGroup([]card{s4, c4})},
// 			// 第三道三条
// 			{cardsToGroup([]card{sk, ck, dk}), cardsToGroup([]card{h2, c2, s4, c4})},
// 			{cardsToGroup([]card{sk, ck, dk}), cardsToGroup([]card{s4, c4}), cardsToGroup([]card{h2, c2})},
// 			// 第三道两对
// 			{cardsToGroup([]card{ck, dk, s4, c4}), cardsToGroup([]card{h2, c2})},
// 			{cardsToGroup([]card{sk, ck, s4, c4}), cardsToGroup([]card{h2, c2})},
// 			{cardsToGroup([]card{sk, dk, s4, c4}), cardsToGroup([]card{h2, c2})},
// 			{cardsToGroup([]card{ck, dk, h2, c2}), cardsToGroup([]card{s4, c4})},
// 			{cardsToGroup([]card{sk, ck, h2, c2}), cardsToGroup([]card{s4, c4})},
// 			{cardsToGroup([]card{sk, dk, h2, c2}), cardsToGroup([]card{s4, c4})},
// 			{cardsToGroup([]card{h2, c2, s4, c4}), cardsToGroup([]card{ck, dk})},
// 			{cardsToGroup([]card{h2, c2, s4, c4}), cardsToGroup([]card{sk, ck})},
// 			{cardsToGroup([]card{h2, c2, s4, c4}), cardsToGroup([]card{sk, dk})},
// 			// 第三道一对
// 			{cardsToGroup([]card{sk, dk}), cardsToGroup([]card{s4, c4}), cardsToGroup([]card{h2, c2})},
// 			{cardsToGroup([]card{sk, ck}), cardsToGroup([]card{s4, c4}), cardsToGroup([]card{h2, c2})},
// 			{cardsToGroup([]card{ck, dk}), cardsToGroup([]card{s4, c4}), cardsToGroup([]card{h2, c2})},
// 		},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := MauBinhList(cardsToWarpCards(tt.args.cards))
//
// 			want := make([]*MauBinh, 0, len(tt.want))
// 			for _, val := range tt.want {
// 				want = append(want, newMauBinh(cardsToGroup(tt.args.cards), val...))
// 			}
//
// 			mauBinhsEqual(got, want, t)
// 		})
// 	}
// }

func TestMauMauBinhBestList(t *testing.T) {
	type args struct {
		cards []card
	}
	tests := []struct {
		name string
		args args
		want [][]*Group
	}{
		{name: "case1", args: args{cards: []card{c2, dq, c9, h2, h8, h7, ha, hj, s4, c4, sk, ck, dk}}, want: [][]*Group{
			// 第三道葫芦
			{cardsToGroup([]card{s4, c4, sk, ck, dk}), cardsToGroup([]card{h2, h8, h7, ha, hj})},
			// 第三道同花
			{cardsToGroup([]card{h2, h8, h7, hj, ha}), cardsToGroup([]card{sk, ck, dk}), cardsToGroup([]card{s4, c4})},
		}},
		{name: "case2", args: args{cards: []card{h2, s6, ca, dt, d3, dq, d6, dk, s8, h8, sj, cj, hj}}, want: [][]*Group{
			// 第三道葫芦
			{cardsToGroup([]card{s8, h8, sj, cj, hj}), cardsToGroup([]card{dt, d3, dq, d6, dk})},
			// 第三道同花
			{cardsToGroup([]card{dt, d3, dq, d6, dk}), cardsToGroup([]card{sj, cj, hj}), cardsToGroup([]card{s8, h8})},
			// 第三道顺子
			{cardsToGroup([]card{dt, hj, dq, dk, ca}), cardsToGroup([]card{s8, h8, d6, s6}), cardsToGroup([]card{sj, cj})},
		}},
		{name: "case3", args: args{cards: []card{c5, h5, d5, h4, s4, c6, h7, c9, dt, hj, dq, dk, ca}}, want: [][]*Group{
			{cardsToGroup([]card{c5, h5, d5, h4, s4}), cardsToGroup([]card{dt, hj, dq, dk, ca})},
			{cardsToGroup([]card{c5, h5, d5, h4, s4}), cardsToGroup([]card{c9, dt, hj, dq, dk})},
			{cardsToGroup([]card{c9, dt, hj, dq, dk}), cardsToGroup([]card{c6, h7, c5, h5, d5}), cardsToGroup([]card{ca, h4, s4})},
			{cardsToGroup([]card{ca, dt, hj, dq, dk}), cardsToGroup([]card{c6, h7, c5, h5, d5}), cardsToGroup([]card{c9, h4, s4})},
		}},
		{name: "case4", args: args{cards: []card{dq, d8, d6, d7, d2, h3, hj, sj, c5, s5, ha, h4, s4}}, want: [][]*Group{
			{cardsToGroup([]card{dq, d8, d6, d7, d2}), cardsToGroup([]card{hj, sj, c5, s5}), cardsToGroup([]card{h4, s4})},
			{cardsToGroup([]card{dq, d8, d6, d7, d2}), cardsToGroup([]card{hj, sj, h4, s4}), cardsToGroup([]card{c5, s5})},
			{cardsToGroup([]card{dq, d8, d6, d7, d2}), cardsToGroup([]card{h4, s4, c5, s5}), cardsToGroup([]card{hj, sj})},
			{cardsToGroup([]card{s4, s5, d6, d7, d8}), cardsToGroup([]card{ha, d2, h3, h4, c5}), cardsToGroup([]card{hj, sj})},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mauBinhBestList(cardsToWarpCards(tt.args.cards))

			want := make([]*MauBinh, 0, len(tt.want))
			for _, val := range tt.want {
				want = append(want, newMauBinh(cardsToGroup(tt.args.cards), val...))
			}

			mauBinhsEqual(got, want, t)
		})
	}
}

func mauBinhsEqual(got, want []*MauBinh, t *testing.T) {
	if len(got) != len(want) {
		t.Errorf("len(got) = %d, len(want) = %d", len(got), len(want))
	}

	sortMaubinhs(got)
	sortMaubinhs(want)

	for i := 0; i < len(got); i++ {
		if got[i].One.rank != want[i].One.rank || got[i].Two.rank != want[i].Two.rank || got[i].Three.rank != want[i].Three.rank {
			t.Errorf("got = %v, want %v", got[i], want[i])
		}
	}
}

func sortMaubinhs(list []*MauBinh) {
	sort.Slice(list, func(i, j int) bool {
		if list[i].Three.rank == list[j].Three.rank {
			if list[i].Two.rank == list[j].Two.rank {
				if list[i].One.rank == list[j].One.rank {
					return false
				}
				return list[i].One.rank > list[j].One.rank
			}
			return list[i].Two.rank > list[j].Two.rank
		}
		return list[i].Three.rank > list[j].Three.rank

	})
}

func cardsToGroup(cards []card) *Group {
	return newGroup(cardsToWarpCards(cards))
}
