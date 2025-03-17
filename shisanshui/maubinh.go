package shisanshui

import (
	"sort"
)

// MauBinh 一种十三水组合
type MauBinh struct {
	cards []WrapCard
	// 牌型大小依次为: One < Two < Three
	One   *Group
	Two   *Group
	Three *Group
}

// Cards 获取排序后的手牌
func (m *MauBinh) Cards() []ICard {
	cards := make([]ICard, 0, 13)
	for _, g := range []*Group{m.Three, m.Two, m.One} {
		for _, val := range g.cards {
			cards = append(cards, val.ori)
		}
	}
	return cards
}

// OriCards 返回原始手牌
func (m *MauBinh) OriCards() []ICard {
	cards := make([]ICard, 0, 13)
	for _, val := range m.cards {
		cards = append(cards, val.ori)
	}
	return cards
}

// Equal 比较 m 和 m1 是否相等，仅比较每墩大小是否相等，不比较具体的花色和点数
func (m *MauBinh) Equal(m1 *MauBinh) bool {
	if m.Three.rank != m1.Three.rank || m.Two.rank != m1.Two.rank ||
		m.One.rank != m.One.rank {
		return false
	}
	return true
}

// sortCards 按maubinh游戏规则大小排序
func sortCards(cards []WrapCard) {
	sort.Slice(cards, func(i, j int) bool {
		// 比较点数大小
		pi, pj := cards[i].Point(), cards[j].Point()
		if pi == PA {
			pi += PK
		}
		if pj == PA {
			pj += PK
		}

		if pi != pj {
			return pi < pj
		}

		// 比较花色大小
		si, sj := cards[i].Suit(), cards[j].Suit()
		switch si {
		case Spade:
			si = 1
		case Heart:
			si = 4
		case Club:
			si = 2
		case Diamond:
			si = 3
		}
		switch sj {
		case Spade:
			sj = 1
		case Heart:
			sj = 4
		case Club:
			sj = 2
		case Diamond:
			sj = 3
		}

		return si < sj
	})
}

func (m *MauBinh) isBetter(mau *MauBinh) bool {
	if mau.Three.rank > m.Three.rank {
		return false
	}

	if mau.Two.rank > m.Two.rank {
		return false
	}

	if mau.One.rank > m.One.rank {
		return false
	}

	return true
}

// NewMauBinh 返回一个十三水组合
// cards: 按照十三水第三墩，第二墩，第一墩墩顺序排列
func NewMauBinh(cards []ICard) *MauBinh {
	if len(cards) != 13 {
		return nil
	}

	cs := make([]WrapCard, 0, len(cards))
	for _, val := range cards {
		cs = append(cs, newCard(val))
	}

	ori := newGroup(cs)
	three := newGroup(cs[0:5])
	two := newGroup(cs[5:10])
	one := newGroup(cs[10:])

	return newMauBinh(ori, three, two, one)
}

// newMauBinh 按所给的13张cards和已经确定的三墩牌，构造一个符合十三水规则的组牌
// 注意事项：
// - 不解析特殊牌型
// - Three 一定不为空，因为非特殊牌型不可能13张牌没有一个对子，那么three至少是一个对子
// - cards 中的牌去除 Three、Two、one中牌后一定是散牌
// todo：后续再优化构建组牌的算法
func newMauBinh(ori *Group, groups ...*Group) *MauBinh {
	m := &MauBinh{
		cards: ori.cards,
	}

	// 获取未使用的牌
	// todo:利用更好的算法快速查找未使用卡牌
	usedCards := make([]WrapCard, 0)
	for _, g := range groups {
		usedCards = append(usedCards, g.cards...)
	}
	useGroup := newGroup(usedCards)
	unuseCards := ori.excludeCard(useGroup)

	// 未使用的牌全都是散牌，不存在一对及以上组合，从第三到到第一道依次从小到大填充踢脚
	// 因此从小到大进行排序
	sortCards(unuseCards)

	gs := make([]*Group, 0, 3)
	// 填充踢脚
	for i, g := range groups {
		n := 5
		if i == 2 {
			n = 3
		}

		n = n - len(g.cards)
		if n > 0 {
			// note:不能如注释写法，因为 groups 中的元素如果被复用，那么会导致原切片内容发生改变
			// g = newGroup(append(g.cards, unuseCards[:n]...))
			cards := append([]WrapCard{}, g.cards...)
			cards = append(cards, unuseCards[:n]...)
			g = newGroup(cards)
			unuseCards = unuseCards[n:]
		}

		gs = append(gs, g)
	}

	n := 3 - len(gs)
	for i := 0; i < n; i++ {
		nc := 5
		if len(gs) == 2 {
			nc = 3
		}

		gs = append(gs, newGroup(unuseCards[:nc]))
		unuseCards = unuseCards[nc:]
	}

	m.Three, m.Two, m.One = gs[0], gs[1], gs[2]

	// 验证是否合法，在添加踢脚的过程中可能形成新的牌型
	if m.Three.rank >= m.Two.rank && m.Two.rank > m.One.rank {
		return m
	}

	return nil
}

// BestMaubinh 返回最优组牌
// better: 更优判定，moreBetter 是否比 better 更好
func BestMaubinh(cards []ICard, better func(better, moreBetter *MauBinh) bool) *MauBinh {
	cs := make([]WrapCard, 0, len(cards))
	for _, val := range cards {
		cs = append(cs, newCard(val))
	}

	list := mauBinhBestList(cs)

	var best *MauBinh
	for _, val := range list {
		if best == nil || better(best, val) {
			best = val
		}
	}

	return best
}

// mauBinhBestList 返回的十三水组合相互之间互有输赢，不会存在一个组合全输另一个组合
func mauBinhBestList(cards []WrapCard) []*MauBinh {
	if len(cards) != 13 {
		return nil
	}

	// 获取所有可能的牌型
	g := newGroup(cards)
	threes := g.all()
	sort.Slice(threes, func(i, j int) bool {
		return threes[i].rank > threes[j].rank
	})

	// note：当前组牌算法不会将对子分成两张散牌的组合，例如：
	// 手牌：ck hk dk cq dq sa ha s2 c3 d4 d6 h7 h8
	// 不会出现如下组合：
	// sa h7 h8
	// ha s2 c3 d4 d6
	// ck hk dk cq dq
	list := make([]*MauBinh, 0)

	// list 中是否有比 groups 更好的组合，每一道都大于 m
	hasBetter := func(list []*MauBinh, m *MauBinh) bool {
		if len(list) == 0 {
			return false
		}

		for _, val := range list {
			if !val.isBetter(m) {
				continue
			}
			return true
		}
		return false
	}
	for i := 0; i < len(threes); i++ {
		three := threes[i]
		twos := remainGroups(threes[i:], three)
		if len(twos) == 0 {
			maubinh := newMauBinh(g, three)
			if maubinh == nil || hasBetter(list, maubinh) {
				continue
			}

			list = append(list, maubinh)
			continue
		}
		for j := 0; j < len(twos); j++ {
			two := twos[j]
			ones := remainGroups(twos[j:], two)
			if len(ones) == 0 {
				maubinh := newMauBinh(g, three, two)
				if maubinh == nil || hasBetter(list, maubinh) {
					continue
				}

				list = append(list, maubinh)
				continue
			}
			for k := 0; k < len(ones); k++ {
				one := ones[k]
				if len(one.cards) > 3 {
					continue
				}
				maubinh := newMauBinh(g, three, two, one)
				if maubinh == nil || hasBetter(list, maubinh) {
					continue
				}
				list = append(list, maubinh)
			}
		}
	}
	return list
}

// remainGroups 从 list 中删除 exclude 包含的 cards 的 Group，返回剩余的 list
func remainGroups(list []*Group, exclude ...*Group) []*Group {
	res := make([]*Group, 0)
	i := 0
loop:
	for i < len(list) {
		val := list[i]
		i++
		for _, v := range exclude {
			// todo：优化算法，更快的判断是否包含
			if len(val.intersectionCards(v)) > 0 {
				goto loop
			}
		}

		res = append(res, val)
	}

	return res
}
