package shisanshui

import (
	"cmp"
)

type Point = int

const (
	_ Point = iota
	PA
	P2
	P3
	P4
	P5
	P6
	P7
	P8
	P9
	PT
	PJ
	PQ
	PK
)

type Suit = int

const (
	_ Suit = iota
	Spade
	Heart
	Club
	Diamond
)

// ICard 卡牌表示
type ICard interface {
	// Point 点数
	Point() Point
	// Suit 花色
	Suit() Suit
}

const (
	suitSpade   = 1 << 28
	suitHeart   = 1 << 29
	suitClub    = 1 << 30
	suitDiamond = 1 << 31
)
const (
	p2 = 2
	p3 = 3
	p4 = 5
	p5 = 7
	p6 = 11
	p7 = 13
	p8 = 17
	p9 = 19
	pt = 23
	pj = 29
	pq = 31
	pk = 37
	pa = 41
)
const (
	sa card = suitSpade | pa
	s2 card = suitSpade | p2
	s3 card = suitSpade | p3
	s4 card = suitSpade | p4
	s5 card = suitSpade | p5
	s6 card = suitSpade | p6
	s7 card = suitSpade | p7
	s8 card = suitSpade | p8
	s9 card = suitSpade | p9
	st card = suitSpade | pt
	sj card = suitSpade | pj
	sq card = suitSpade | pq
	sk card = suitSpade | pk

	ha card = suitHeart | pa
	h2 card = suitHeart | p2
	h3 card = suitHeart | p3
	h4 card = suitHeart | p4
	h5 card = suitHeart | p5
	h6 card = suitHeart | p6
	h7 card = suitHeart | p7
	h8 card = suitHeart | p8
	h9 card = suitHeart | p9
	ht card = suitHeart | pt
	hj card = suitHeart | pj
	hq card = suitHeart | pq
	hk card = suitHeart | pk

	ca card = suitClub | pa
	c2 card = suitClub | p2
	c3 card = suitClub | p3
	c4 card = suitClub | p4
	c5 card = suitClub | p5
	c6 card = suitClub | p6
	c7 card = suitClub | p7
	c8 card = suitClub | p8
	c9 card = suitClub | p9
	ct card = suitClub | pt
	cj card = suitClub | pj
	cq card = suitClub | pq
	ck card = suitClub | pk

	da card = suitDiamond | pa
	d2 card = suitDiamond | p2
	d3 card = suitDiamond | p3
	d4 card = suitDiamond | p4
	d5 card = suitDiamond | p5
	d6 card = suitDiamond | p6
	d7 card = suitDiamond | p7
	d8 card = suitDiamond | p8
	d9 card = suitDiamond | p9
	dt card = suitDiamond | pt
	dj card = suitDiamond | pj
	dq card = suitDiamond | pq
	dk card = suitDiamond | pk
)

// PA-PK 和 pa-pK的映射，下标表示 Point-1
var points = []uint32{pa, p2, p3, p4, p5, p6, p7, p8, p9, pt, pj, pq, pk}

// Spade-Heart-Club-Diamond 和 suitSpade-suitHeart-suitClub-suitDiamond的映射, 下表表示 Suit -1
var suits = []uint32{suitSpade, suitHeart, suitClub, suitDiamond}

// card 前四位表示花色，后28位表示点数，点数使用素数表示从 2-KA 依此为 [2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41]
type card uint32

func newCard(c ICard) WrapCard {
	point := points[c.Point()-1]
	suit := suits[c.Suit()-1]
	return WrapCard{
		card: card(point | suit),
		ori:  c,
	}
}

func (c card) Point() Point {
	if c&pa == pa {
		return PA
	}
	return binarySearch(points[1:], resetSuit(c)) + 2
}

func (c card) Suit() Suit {
	for i, val := range suits {
		if uint32(c)&val > 0 {
			return i + 1
		}
	}
	return 0
}

func resetSuit(c card) uint32 {
	return uint32(c & 0xFFFFFFF)
}

var cardNames = []string{
	"SA", "S2", "S3", "S4", "S5", "S6", "S7", "S8", "S9", "ST", "SJ", "SQ", "SK",
	"HA", "H2", "H3", "H4", "H5", "H6", "H7", "H8", "H9", "HT", "HJ", "HQ", "HK",
	"CA", "C2", "C3", "C4", "C5", "C6", "C7", "C8", "C9", "CT", "CJ", "CQ", "CK",
	"DA", "D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9", "DT", "DJ", "DQ", "DK",
}

func (c card) String() string {
	s, p := c.Suit(), c.Point()
	return cardNames[(s-1)*13+p-1]
}

type WrapCard struct {
	card card
	ori  ICard
}

func (w WrapCard) Point() Point {
	return w.ori.Point()
}

func (w WrapCard) Suit() Suit {
	return w.ori.Suit()
}

// binarySearch 在一个有序的切片中查找目标值
// 如果找到目标值，返回其索引；否则，返回 -1
func binarySearch[T cmp.Ordered](arr []T, target T) int {
	low := 0
	high := len(arr) - 1

	for low <= high {
		mid := low + (high-low)/2 // 防止 (low + high) 溢出
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1 // 没有找到目标值
}
