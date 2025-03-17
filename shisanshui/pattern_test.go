package shisanshui

import (
	"testing"
)

var testParser = NewPatternParser()

func TestPatternParser_parse(t *testing.T) {
	type args struct {
		cards []card
	}

	patterSumNum := []int{0, 0}
	n := 0
	for _, val := range []int{highNum, pairNum, twoPairNum, threeKindNum, straightNum, flushNum, fullHouseNum, fourKindNum} {
		n += val
		patterSumNum = append(patterSumNum, n)
	}
	tests := []struct {
		name        string
		args        args
		wantPattern Pattern
		wantRank    int
	}{
		{name: "case:最小的散牌(五张)", args: args{cards: []card{c2, c3, c4, c5, h7}}, wantPattern: High, wantRank: patterSumNum[High] + 17},
		{name: "case:最小的散牌(三张)", args: args{cards: []card{c2, c3, c4}}, wantPattern: High, wantRank: patterSumNum[High] + 1},
		{name: "case:最大的一对(不带踢脚)", args: args{cards: []card{ca, ha}}, wantPattern: Pair, wantRank: patterSumNum[Pair] + 12*233 + 1},
		{name: "case:最大的一对(带踢脚)", args: args{cards: []card{ca, ha, hk, cq, cj}}, wantPattern: Pair, wantRank: patterSumNum[TwoPair]},
		{name: "case:最小的两对(不带踢脚)", args: args{cards: []card{c2, c3, h2, h3}}, wantPattern: TwoPair, wantRank: patterSumNum[TwoPair] + 1},
		{name: "case:最小的两对(带踢脚)", args: args{cards: []card{c2, c3, h4, h2, h3}}, wantPattern: TwoPair, wantRank: patterSumNum[TwoPair] + 2},
		{name: "case:最大的三条(不带踢脚)", args: args{cards: []card{ca, ha, da}}, wantPattern: ThreeKind, wantRank: patterSumNum[ThreeKind] + 12*67 + 1},
		{name: "case:最大的三条(带踢脚)", args: args{cards: []card{ca, ha, da, dk, cq}}, wantPattern: ThreeKind, wantRank: patterSumNum[Straight]},
		{name: "case:最小的顺子", args: args{cards: []card{ca, c2, c3, h4, d5}}, wantPattern: Straight, wantRank: patterSumNum[Straight] + 1},
		{name: "case:最大的顺子", args: args{cards: []card{ca, ck, cq, hj, ht}}, wantPattern: Straight, wantRank: patterSumNum[Flush]},
		{name: "case:最小的同花", args: args{cards: []card{c2, c3, c4, c5, c7}}, wantPattern: Flush, wantRank: patterSumNum[Flush] + 1},
		{name: "case:最小的葫芦", args: args{cards: []card{c2, c3, s2, h2, h3}}, wantPattern: FullHouse, wantRank: patterSumNum[FullHouse] + 1},
		{name: "case:四条q(不带踢脚)", args: args{cards: []card{cq, dq, sq, hq}}, wantPattern: FourKind, wantRank: patterSumNum[FourKind] + 131},
		{name: "case:四条q(带踢脚)", args: args{cards: []card{cq, dq, sq, hq, sa}}, wantPattern: FourKind, wantRank: patterSumNum[FourKind] + 143},
		{name: "case:同花顺", args: args{cards: []card{cq, cj, ct, c9, c8}}, wantPattern: FlushStraight, wantRank: patterSumNum[FlushStraight] + 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPattern, gotRank := testParser.parse(cardsToWarpCards(tt.args.cards))
			if gotPattern != tt.wantPattern {
				t.Errorf("parse() gotPattern = %v, want %v", gotPattern, tt.wantPattern)
			}
			if gotRank != tt.wantRank {
				t.Errorf("parse() gotRank = %v, want %v", gotRank, tt.wantRank)
			}
		})
	}
}
