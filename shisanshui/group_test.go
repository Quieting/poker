package shisanshui

import (
	"reflect"
	"sort"
	"testing"
)

func Test_group_fullHouse(t *testing.T) {
	type fields struct {
		cards []card
	}
	tests := []struct {
		name   string
		fields fields
		want   [][]card
	}{
		{name: "fullHouse", fields: fields{cards: []card{s2, h2, c2, s3, d3, h3, d5, d8, ht, cj}}, want: [][]card{
			{s2, h2, c2, s3, d3},
			{s2, h2, c2, s3, h3},
			{s2, h2, c2, h3, d3},
			{s2, h2, h3, s3, d3},
			{s2, h3, c2, s3, d3},
			{h3, h2, c2, s3, d3},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newGroup(cardsToWarpCards(tt.fields.cards))

			got := g.fullHouse()
			want := cardsListToGroupList(tt.want)

			groupsEqual(got, want, t)
		})
	}
}

func Test_group_flushStraight(t *testing.T) {
	type fields struct {
		cards []card
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [][]card
	}{
		{name: "case1", fields: fields{cards: []card{c2, c3, c4, c5, c6, h7}}, args: args{n: 5}, want: [][]card{
			{c2, c3, c4, c5, c6},
		}},
		{name: "case2:含A的同花顺", fields: fields{cards: []card{sa, s2, s3, s4, s5, st, sj, sq, sk}}, args: args{n: 5}, want: [][]card{
			{sa, s2, s3, s4, s5},
			{sa, st, sj, sq, sk},
		}},
		{name: "case3:不成同花顺子", fields: fields{cards: []card{sa, s2, d3, c4, ht, hj, hq, hk}}, args: args{n: 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cards := make([]WrapCard, 0, len(tt.fields.cards))
			for _, val := range tt.fields.cards {
				cards = append(cards, newCard(val))
			}
			g := newGroup(cards)
			got := g.flushStraight(tt.args.n)

			want := cardsListToGroupList(tt.want)

			groupsEqual(got, want, t)
		})
	}
}

func Test_group_straight(t *testing.T) {
	type fields struct {
		cards []card
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [][]card
	}{
		{name: "case1", fields: fields{cards: []card{s2, d3, c4, c5, h6, h7}}, args: args{n: 5}, want: [][]card{
			{s2, d3, c4, c5, h6},
			{d3, c4, c5, h6, h7},
		}},
		{name: "case2:含A的两头顺子", fields: fields{cards: []card{sa, s2, d3, c4, c5, ht, hj, cq, dk}}, args: args{n: 5}, want: [][]card{
			{sa, s2, d3, c4, c5},
			{sa, ht, hj, cq, dk},
		}},
		{name: "case3:不成顺子", fields: fields{cards: []card{sa, s2, d3, c4, ht, hj, cq}}, args: args{n: 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newGroup(cardsToWarpCards(tt.fields.cards))

			got := g.straight(tt.args.n)
			want := cardsListToGroupList(tt.want)

			groupsEqual(got, want, t)
		})
	}
}

func sortGroups(list []*Group) {
	sort.Slice(list, func(i, j int) bool {
		for k := 0; k < len(list[i].cardsBySuit); k++ {
			if k >= len(list[j].cardsBySuit) {
				return false
			}
			if list[i].cardsBySuit[k].card == list[j].cardsBySuit[k].card {
				continue
			}
			return list[i].cardsBySuit[k].card < list[j].cardsBySuit[k].card
		}
		return true
	})
}

func Test_group_all(t *testing.T) {
	type fields struct {
		cards []card
	}
	tests := []struct {
		name   string
		fields fields
		want   [][]card
	}{
		{name: "case1", fields: fields{cards: []card{c2, dq, c9, h2, h8, h7, ha, hj, s4, c4, sk, ck, hk}}, want: [][]card{
			// 葫芦
			{s4, c4, sk, ck, hk},
			{h2, c2, sk, ck, hk},
			// 同花
			{h2, h8, h7, hj, ha},
			{h2, h8, h7, hj, hk},
			{h2, h8, h7, hk, ha},
			{h2, h8, hk, hj, ha},
			{h2, hk, h7, hj, ha},
			{hk, h8, h7, hj, ha},
			// 三条
			{sk, ck, hk},
			// 两对
			{c2, h2, s4, c4},
			{c2, h2, sk, ck},
			{c2, h2, sk, hk},
			{c2, h2, ck, hk},
			{s4, c4, sk, ck},
			{s4, c4, sk, hk},
			{s4, c4, ck, hk},

			// 一对
			{c2, h2},
			{sk, ck},
			{sk, hk},
			{ck, hk},
			{s4, c4},
		},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newGroup(cardsToWarpCards(tt.fields.cards))

			got := g.all()
			want := cardsListToGroupList(tt.want)

			groupsEqual(got, want, t)
		})
	}
}

func groupsEqual(got, want []*Group, t *testing.T) {
	sortGroups(got)
	sortGroups(want)
	if len(got) != len(want) {
		t.Errorf("len(got) = %d, len(want) = %d", len(got), len(want))
	}
	for i := 0; i < len(got); i++ {
		if !reflect.DeepEqual(got[i].cardsBySuit, want[i].cardsBySuit) {
			t.Errorf("got = %v, want %v", got[i].cardsBySuit, want[i].cardsBySuit)
		}
	}
}

func Test_group_intersectionCards(t *testing.T) {
	type fields struct {
		cards []card
	}
	type args struct {
		a []card
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []card
	}{
		{name: "case:返回空", fields: fields{cards: []card{c2, d2, sq, sa, ck}}, args: args{a: []card{h2, d3, dj, dt}}},
		{name: "case:a是b的子集", fields: fields{cards: []card{c2, d2, sq, sa, ck}}, args: args{a: []card{c2, d2, sa, sq}}, want: []card{c2, d2, sa, sq}},
		{name: "case：一般情况", fields: fields{cards: []card{c2, d2, sq, sa, ck}}, args: args{a: []card{c2, d3, sq, dj, dt}}, want: []card{c2, sq}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newGroup(cardsToWarpCards(tt.fields.cards))
			got := g.intersectionCards(newGroup(cardsToWarpCards(tt.args.a)))

			groupsEqual([]*Group{newGroup(got)}, []*Group{newGroup(cardsToWarpCards(tt.want))}, t)
		})
	}
}

func Test_group_excludeCard(t *testing.T) {
	type fields struct {
		cards []card
	}
	type args struct {
		a []card
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []card
	}{
		{name: "case:没有交集", fields: fields{cards: []card{c2, d2, sq, sa, ck}}, args: args{a: []card{h2, d3, dj, dt}}, want: []card{c2, d2, sq, sa, ck}},
		{name: "case:a是b的子集", fields: fields{cards: []card{c2, d2, sq, sa, ck}}, args: args{a: []card{c2, d2, sa, sq}}, want: []card{ck}},
		{name: "case：一般情况", fields: fields{cards: []card{c2, d2, sq, sa, ck}}, args: args{a: []card{c2, d3, sq, dj, dt}}, want: []card{d2, sa, ck}},
		{name: "case：一般情况2", fields: fields{cards: []card{c2, d2, sq, dj, dt}}, args: args{a: []card{d2, d3, sq, sa, ck}}, want: []card{c2, dj, dt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newGroup(cardsToWarpCards(tt.fields.cards))

			got := g.excludeCard(newGroup(cardsToWarpCards(tt.args.a)))
			groupsEqual([]*Group{newGroup(got)}, []*Group{newGroup(cardsToWarpCards(tt.want))}, t)
		})
	}
}

func cardsToWarpCards(cards []card) []WrapCard {
	cs := make([]WrapCard, 0, len(cards))
	for _, val := range cards {
		cs = append(cs, newCard(val))
	}

	return cs
}

func cardsListToGroupList(cards [][]card) []*Group {
	list := make([]*Group, 0, len(cards))
	for _, val := range cards {
		list = append(list, newGroup(cardsToWarpCards(val)))
	}

	return list
}
