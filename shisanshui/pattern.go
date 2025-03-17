package shisanshui

type Pattern int

const (
	_             Pattern = iota
	High                  // 高牌
	Pair                  // 一对
	TwoPair               // 两对
	ThreeKind             // 三条
	Straight              // 顺子
	Flush                 // 同花
	FullHouse             // 葫芦
	FourKind              // 四条
	FlushStraight         // 同花顺
)

// 各种牌型的数量
const (
	highNum          = 1563 // 散牌有 1563 种:即 C(13,5) - 10 + C(13,3)(三张的散牌)
	pairNum          = 3029 // 一对有 3029 种:即 13*C(12,3) + 13(不带踢脚的对子) + 13*12(带一个踢脚的对子)
	twoPairNum       = 936  // 两对有 936 种:即 C(13,2)*11 + 78(不带踢脚的两对)
	threeKindNum     = 871  // 三条有 871 种:即 13*C(12,2) + 13(不带踢脚的三条)
	straightNum      = 10   // 顺子有 10 种
	flushNum         = 1277 // 同花有 1277 种: 即 C(13,5) - 10
	fullHouseNum     = 156  // 葫芦有 156 种: 即 13*12
	fourKindNum      = 169  // 四条有 169 种：即 13*12 + 13(不带踢脚的四条)
	flushStraightNum = 10   // 同花顺有 10 种
)

var GlobalParser = NewPatternParser()

// PatternParser 牌型解析器
// 用13个不同的质数表示牌的点数，将五张牌的点数相乘得到一个数，这个数因式分解后只会得到唯一的一种组合
// 那么可以建立一张牌型和大小的映射表，可以快速确定五张牌是什么牌型
type PatternParser struct {
	straightMap map[uint64]int
	flushMap    map[uint64]int
	setMap      map[uint64]int

	ranks    []int
	patterns []Pattern
}

func NewPatternParser() *PatternParser {
	p := &PatternParser{
		straightMap: highMaps(),
		flushMap:    flushMaps(),
		setMap:      setMaps(),

		ranks:    []int{highNum, pairNum, twoPairNum, threeKindNum, straightNum, flushNum, fullHouseNum, fourKindNum, flushStraightNum},
		patterns: []Pattern{High, Pair, TwoPair, ThreeKind, Straight, Flush, FullHouse, FourKind, FlushStraight},
	}

	return p
}
func (p *PatternParser) parse(cards []WrapCard) (pattern Pattern, rank int) {
	rank = p.rank(cards)

	inx := 0
	r := rank
	for i, n := range p.ranks {
		inx = i
		if r <= n {
			break
		}
		r -= n
	}

	pattern = p.patterns[inx]

	return pattern, rank
}

// rank 计算牌型的大小
func (p *PatternParser) rank(cards []WrapCard) int {
	nums := uint64(1) // 手牌乘级
	isSuit := 0xF     // 判断是否同花
	for _, val := range cards {
		nums *= uint64(resetSuit(val.card))
		isSuit &= 1 << (val.ori.Suit() - 1)
	}

	r := 0 // 牌力值

	// 是否同花
	if isSuit > 0 && len(cards) == 5 {
		r = p.flushMap[nums]
	} else if val, ok := p.straightMap[nums]; ok { // 是否顺子、高牌
		r = val
	} else {
		r = p.setMap[nums]
	}

	return r
}

// highTable 所有的高排组合大小映射表
func highMaps() map[uint64]int {
	maps := make(map[uint64]int, 1287)

	nums := []uint64{pa, p2, p3, p4, p5, p6, p7, p8, p9, pt, pj, pq, pk, pa}

	// 计算顺子映射
	straights := highNum + pairNum + twoPairNum + threeKindNum
	for i := 0; i <= len(nums)-5; i++ {
		num := nums[i]
		for j := i + 1; j < i+5; j++ {
			num *= nums[j]
		}
		straights++
		maps[num] = straights
	}

	// 计算高牌映射
	highs := 0
	nums = []uint64{p2, p3, p4, p5, p6, p7, p8, p9, pt, pj, pq, pk, pa}
	for i := 2; i < len(nums); i++ {
		for j := 1; j < i; j++ {
			for k := 0; k < j; k++ {
				// 添加三张高牌的映射
				val := nums[i] * nums[j] * nums[k]
				// 跳过顺子
				if _, ok := maps[val]; ok {
					continue
				}
				highs++
				maps[val] = highs

				for l := 1; l < k; l++ {
					for m := 0; m < l; m++ {
						val := nums[i] * nums[j] * nums[k] * nums[l] * nums[m]
						// 跳过顺子
						if _, ok := maps[val]; ok {
							continue
						}
						highs++
						maps[val] = highs
					}
				}
			}
		}
	}

	return maps
}

func flushMaps() map[uint64]int {
	maps := make(map[uint64]int, 1287)
	nums := []uint64{pa, p2, p3, p4, p5, p6, p7, p8, p9, pt, pj, pq, pk, pa}

	// 计算同花顺子映射
	flushStraigh := highNum + pairNum + twoPairNum + threeKindNum + straightNum + flushNum + fullHouseNum + fourKindNum
	for i := 0; i <= len(nums)-5; i++ {
		num := nums[i]
		for j := i + 1; j < i+5; j++ {
			num *= nums[j]
		}
		flushStraigh++
		maps[num] = flushStraigh
	}

	// 计算同花映射
	nums = []uint64{p2, p3, p4, p5, p6, p7, p8, p9, pt, pj, pq, pk, pa}
	flush := highNum + pairNum + twoPairNum + threeKindNum + straightNum
	for i := 5; i < len(nums); i++ {
		for j := 3; j < i; j++ {
			for k := 2; k < j; k++ {
				for l := 1; l < k; l++ {
					for m := 0; m < l; m++ {
						val := nums[i] * nums[j] * nums[k] * nums[l] * nums[m]
						if _, ok := maps[val]; ok {
							continue
						}
						flush++
						maps[val] = flush
					}
				}
			}
		}
	}
	return maps
}

func setMaps() map[uint64]int {
	nums := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41}
	sets := make(map[uint64]int, 4888)

	// 四条
	four := highNum + pairNum + twoPairNum + threeKindNum + straightNum + flushNum + fullHouseNum
	for i := 0; i < len(nums); i++ {
		// 添加不带踢脚的四条
		val := nums[i] * nums[i] * nums[i] * nums[i]
		four++
		sets[val] = four
		for j := 0; j < len(nums); j++ {
			if j == i {
				continue
			}

			val := nums[i] * nums[i] * nums[i] * nums[i] * nums[j]
			four++
			sets[val] = four
		}
	}

	// 葫芦
	fullHouse := highNum + pairNum + twoPairNum + threeKindNum + straightNum + flushNum
	for i := 0; i < len(nums); i++ {
		for j := 0; j < len(nums); j++ {
			if j == i {
				continue
			}
			val := nums[i] * nums[i] * nums[i] * nums[j] * nums[j]
			fullHouse++
			sets[val] = fullHouse
		}
	}

	// 三条
	three := highNum + pairNum + twoPairNum
	for i := 0; i < len(nums); i++ { // 三条牌
		val := nums[i] * nums[i] * nums[i]
		three++
		sets[val] = three
		for j := 1; j < len(nums); j++ { // 散牌
			if j == i {
				continue
			}
			for k := 0; k < j; k++ { // 散牌
				if k == i {
					continue
				}

				val := nums[i] * nums[i] * nums[i] * nums[j] * nums[k]
				three++
				sets[val] = three
			}
		}
	}
	// 两对
	twoPair := highNum + pairNum
	for i := 1; i < len(nums); i++ { // 大对子
		for j := 0; j < i; j++ { // 小对子
			val := nums[i] * nums[i] * nums[j] * nums[j]
			twoPair++
			sets[val] = twoPair
			for k := 0; k < len(nums); k++ { // 踢脚
				if k == i || k == j {
					continue
				}

				val := nums[i] * nums[i] * nums[j] * nums[j] * nums[k]
				twoPair++
				sets[val] = twoPair
			}
		}
	}
	// 一对
	onePair := highNum
	for i := 0; i < len(nums); i++ { // 对子
		val := nums[i] * nums[i]
		onePair++
		sets[val] = onePair
		for j := 0; j < len(nums); j++ { // 大踢脚
			if j == i {
				continue
			}
			val := nums[i] * nums[i] * nums[j] // 一对带一个踢脚
			onePair++
			sets[val] = onePair

			for k := 1; k < j; k++ { // 中踢脚
				if k == i {
					continue
				}

				for l := 0; l < k; l++ { // 小踢脚
					if l == i {
						continue
					}

					val := nums[i] * nums[i] * nums[j] * nums[k] * nums[l]
					onePair++
					sets[val] = onePair
				}
			}
		}
	}

	return sets
}
