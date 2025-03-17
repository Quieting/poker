package shisanshui

import (
	"github.com/Quieting/poker/math"
)

type Group struct {
	cards []WrapCard
	// 按花色，大小排序的牌组，注意：A会出现在同一个花色的两端方便查找顺子，
	// 排序后结果类似以下示例:
	// S2 S3 S4 HA H2 H3 H4 HT HA
	cardsBySuit []WrapCard
	// 按大小，花色排序
	// 排序后结果类似以下示例:
	// [[SA HA],[S2 H2], [H5 C5 D5]]
	cardsByPoint [][]WrapCard

	// 如果 len(cards) == 5，则以下字段有值
	pattern Pattern
	// 牌型绝对大小
	rank int
}

func (g *Group) Pattern() Pattern {
	return g.pattern
}

func (g *Group) Cards() []ICard {
	cards := make([]ICard, 0, len(g.cards))
	for _, val := range g.cards {
		cards = append(cards, val.ori)
	}
	return cards
}

func newGroup(cards []WrapCard) *Group {
	g := &Group{
		cards:        cards,
		cardsBySuit:  make([]WrapCard, 0, len(cards)+4),
		cardsByPoint: make([][]WrapCard, 0),
	}

	list := [56]WrapCard{}
	for _, val := range cards {
		s, p := val.ori.Suit(), val.ori.Point()
		list[(s-1)*14+p-1] = val
		if p == PA {
			list[(s-1)*14+PK] = val
		}
	}

	// 按花色大小排序
	for _, val := range list {
		if val.card == 0 {
			continue
		}
		g.cardsBySuit = append(g.cardsBySuit, val)
	}

	// 按大小，花色排序
	for p := PA; p < PK+1; p++ {
		cs := make([]WrapCard, 0)
		for _, s := range []Suit{Spade, Heart, Club, Diamond} {
			i := (s-1)*14 + p - 1
			if list[i].card == 0 {
				continue
			}
			cs = append(cs, list[i])
		}
		if len(cs) == 0 {
			continue
		}

		g.cardsByPoint = append(g.cardsByPoint, cs)
	}

	if len(g.cards) <= 5 {
		g.pattern, g.rank = GlobalParser.parse(g.cards)
	}

	return g
}

// all 返回可能的高牌以上的牌型
func (g *Group) all() []*Group {
	funcs := []func() []*Group{
		func() []*Group {
			return g.flushStraight(5)
		},
		g.fourKind,
		g.fullHouse,
		func() []*Group {
			return g.flush(5)
		},
		func() []*Group {
			return g.straight(5)
		},
		g.three,
		g.twoPairs,
		g.pair,
	}

	list := make([]*Group, 0)
	for _, f := range funcs {
		groups := f()
		list = append(list, groups...)
	}
	return list
}

// intersectionCards 返回g和a相同的牌(花色和点数相同)
func (g *Group) intersectionCards(a *Group) []WrapCard {
	cards := make([]WrapCard, 0)

	ai, gi := 0, 0
	for {
		if ai >= len(a.cardsByPoint) || gi >= len(g.cardsByPoint) {
			break
		}
		ap, gp := a.cardsByPoint[ai][0].Point(), g.cardsByPoint[gi][0].Point()
		if ap < gp {
			ai++
			continue
		}

		if ap > gp {
			gi++
			continue
		}

		aj, gj := 0, 0
		for {
			if aj >= len(a.cardsByPoint[ai]) || gj >= len(g.cardsByPoint[gi]) {
				break
			}

			aCard, gCard := a.cardsByPoint[ai][aj], g.cardsByPoint[gi][gj]

			if aCard.card < gCard.card {
				aj++
				continue
			}

			if aCard.card > gCard.card {
				gj++
				continue
			}

			cards = append(cards, aCard)
			aj++
			gj++
		}

		ai++
		gi++
	}

	return cards
}

// excludeCard 返回g中不包含a的牌(花色和点数相同)
func (g *Group) excludeCard(a *Group) []WrapCard {
	cards := make([]WrapCard, 0)

	ai, gi := 0, 0
	for {
		if ai >= len(a.cardsByPoint) || gi >= len(g.cardsByPoint) {
			break
		}
		ap, gp := a.cardsByPoint[ai][0].Point(), g.cardsByPoint[gi][0].Point()
		if ap < gp {
			ai++
			continue
		}

		if ap > gp {
			cards = append(cards, g.cardsByPoint[gi]...)
			gi++
			continue
		}

		aj, gj := 0, 0
		for {
			if aj >= len(a.cardsByPoint[ai]) || gj >= len(g.cardsByPoint[gi]) {
				break
			}

			aCard, gCard := a.cardsByPoint[ai][aj], g.cardsByPoint[gi][gj]

			if aCard.card < gCard.card {
				aj++
				continue
			}

			if aCard.card > gCard.card {
				cards = append(cards, gCard)
				gj++
				continue
			}

			aj++
			gj++
		}

		cards = append(cards, g.cardsByPoint[gi][gj:]...)

		ai++
		gi++
	}

	for _, v := range g.cardsByPoint[gi:] {
		cards = append(cards, v...)
	}

	return cards
}

// flushStraight 获取可能的同花顺组合
func (g *Group) flushStraight(n int) []*Group {
	if g.cardsNum() < n {
		return nil
	}

	list := make([]*Group, 0)
	for i := 0; i < len(g.cardsBySuit)-n+1; i++ {
		minCard, maxCard := g.cardsBySuit[i], g.cardsBySuit[i+n-1]

		if minCard.Suit() != maxCard.Suit() {
			continue
		}

		maxPoint, minPoint := maxCard.Point(), minCard.Point()
		if (maxPoint-minPoint != n-1) && !(minPoint == PT && maxPoint == PA) {
			continue
		}

		list = append(list, newGroup(g.cardsBySuit[i:i+n]))
	}
	return list
}

// fourKind 获取可能的四条组合
func (g *Group) fourKind() []*Group {
	return g.sets(4, 4)
}

// fullHouse 获取可能的葫芦组合
func (g *Group) fullHouse() []*Group {
	// 获取有三张及以上和两张及以上相同点数的牌集合
	setGroups, pairGroups := g.sets(3, 4), g.sets(2, 4)

	if len(setGroups)+len(pairGroups) <= 1 || len(setGroups) == 0 {
		return nil
	}

	list := make([]*Group, 0)
	for _, sg := range setGroups {
		for _, pg := range pairGroups {
			if sg.maxCard().Point() == pg.maxCard().Point() {
				continue
			}

			sgs, pgs := sg.combination(3), pg.combination(2)
			for _, s := range sgs {
				for _, p := range pgs {
					list = append(list, newGroup(append(s.cards, p.cards...)))
				}
			}
		}
	}
	return list
}

// flush 获取可能的同花组合
func (g *Group) flush(n int) []*Group {
	if g.cardsNum() < n {
		return nil
	}

	list := make([]*Group, 0)
	cards := make([]WrapCard, 0, 13)
	for _, c := range g.cardsBySuit {
		if len(cards) == 0 || (cards[0].Suit() == c.Suit() && cards[0].Point() != c.Point()) {
			cards = append(cards, c)
			continue
		}

		if len(cards) >= n {
			for _, val := range math.Combination[WrapCard](cards, n) {
				list = append(list, newGroup(val))
			}
		}

		cards = []WrapCard{c}
	}

	if len(cards) >= n {
		for _, val := range math.Combination[WrapCard](cards, n) {
			list = append(list, newGroup(val))
		}
	}
	return list
}

// straight 获取可能的顺子组合
func (g *Group) straight(n int) []*Group {
	if g.cardsNum() < n {
		return nil
	}

	list := make([]*Group, 0)

	cards := g.cardsByPoint
	if cards[0][0].Point() == PA {
		cards = append(cards, cards[0])
	}

	for i := 0; i < len(cards)-n+1; i++ {
		minPoint, maxPoint := cards[i][0].Point(), cards[i+n-1][0].Point()
		if (maxPoint-minPoint) != n-1 && !(maxPoint == PA && minPoint == PT) {
			continue
		}

		for _, val := range math.CartesianProduct(cards[i : i+n]) {
			list = append(list, newGroup(val))
		}
	}

	return list
}

// three 获取可能的三条，仅返回三张相同的牌，踢脚暂不返回
func (g *Group) three() []*Group {
	list := make([]*Group, 0)
	for _, sg := range g.sets(3, 4) {
		sgs := sg.combination(3)
		list = append(list, sgs...)
	}

	return list
}

// three 获取可能的两队，仅返回两队相同的牌，踢脚暂不返回
func (g *Group) twoPairs() []*Group {
	list := make([]*Group, 0)
	pairs := g.sets(2, 4)
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			pairA, pairB := pairs[i], pairs[j]
			sgs, pgs := pairA.combination(2), pairB.combination(2)
			for _, s := range sgs {
				for _, p := range pgs {
					list = append(list, newGroup(append(s.cards, p.cards...)))
				}
			}
		}
	}
	return list
}

// pair 获取可能的一对，仅返回对子，踢脚暂不返回
func (g *Group) pair() []*Group {
	list := make([]*Group, 0)
	for _, sg := range g.sets(2, 4) {
		sgs := sg.combination(2)
		list = append(list, sgs...)
	}

	return list
}

// sets 返回同一点数张数在 [min, max] 内组合
func (g *Group) sets(min, max int) []*Group {
	list := make([]*Group, 0)
	for _, cards := range g.cardsByPoint {
		if len(cards) < min || len(cards) > max {
			continue
		}

		g := &Group{
			cards:        cards,
			cardsBySuit:  cards,
			cardsByPoint: [][]WrapCard{cards},
		}
		g.pattern, g.rank = GlobalParser.parse(cards)
		list = append(list, g)
	}

	return list
}

// maxCard 返回点数最大的牌，A最小，K最大（不分花色）
func (g *Group) maxCard() WrapCard {
	return g.cardsByPoint[len(g.cardsByPoint)-1][0]
}

func (g *Group) cardsNum() int {
	return len(g.cards)
}

// combination 返回返回随机 c(len(g.cars), m) 的组合
func (g *Group) combination(m int) []*Group {
	if len(g.cards) < m {
		return nil
	}

	res := math.Combination[WrapCard](g.cards, m)
	list := make([]*Group, 0, len(res))
	for _, val := range res {
		list = append(list, newGroup(val))
	}
	return list
}
