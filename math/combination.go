package math

// Combination 返回C(len(data), m)的所有组合结果，len(data) <= 64
// - 创建有n个元素数组，数组元素的值为1表示选中，为0则没选中。
// - 初始化，将数组前m个元素置1，表示第一个组合为前m个数。
// - 从左到右扫描数组元素值的“10”组合，找到第一个“10”组合后将其变为“01”组合，同时将其左边的所有“1”全部移动到数组的最左端。
// - 当某次循环没有找到“10“组合时，说明得到了最后一个组合，循环结束。
func Combination[T any](data []T, m int) [][]T {
	n := len(data)
	nums := comb(n, m)
	list := make([][]T, 0, nums)

	tag := uint64(((1<<n - 1) >> (n - m)) << (n - m))
	for i := 0; i < nums; i++ {
		list = append(list, bitsToData(data, tag))
		tag = bitChange(tag, n)
	}

	return list
}

// bitsToData 返回 n 中 bit 位为1的对应的data中的数据
func bitsToData[T any](data []T, n uint64) []T {
	list := make([]T, 0, len(data))

	for i, val := range data {
		if n&uint64(1<<(len(data)-1-i)) != 0 {
			list = append(list, val)
		}
	}

	return list
}

// bitChange 从左到右扫描 n 中“10”组合，找到第一个“10”组合后将其变为“01”组合，同时将其左边的所有“1”最左端，返回变换后的数字
// start：从右往左第几位
func bitChange(n uint64, start int) uint64 {
	b := uint64(1 << (start - 1)) // 从左往右遍历时第一个 10 的大小
	b1 := b | (1 << (start - 2))  // 用来取tags对应的两位值
	leftOneNum := 0               // 第一个 10 左边的 1 的个数
	for j := 0; j < start-1; j++ {
		if b^(n&b1) == 0 {
			// 找到“10”组合后将其变为“01”组合
			n -= b >> 1
			// 将其左边的所有“1”全部移动到数组的最左端
			b2 := uint64((1<<leftOneNum - 1) << (start - leftOneNum))
			n = b2 | (n & (b - 1))
			break
		}
		if b&n == b {
			leftOneNum++
		}
		b >>= 1
		b1 >>= 1
	}

	return n
}

// factorial 计算给定非负整数 n 的阶乘
func factorial(n int) int {
	return multiply(1, n)
}

// multiply 返回从start累乘n个连续数的结果
func multiply(start, n int) int {
	result := start
	for i := 0; i < n-1; i++ {
		start++
		result *= start
	}
	return result
}

// comb 计算组合数 C(n, k)
func comb(n, k int) int {
	if k > n || n < 0 || k < 0 {
		return 0
	}
	// 优化计算，C(n, k) == C(n, n-k)
	if k < n-k {
		k = n - k
	}
	return multiply(k+1, n-k) / factorial(n-k)
}

// CartesianProduct 计算多个数组的笛卡尔积
func CartesianProduct[T any](arrays [][]T) [][]T {
	if len(arrays) == 0 {
		return nil
	}

	// 保存结果的切片
	var result = [][]T{{}}

	for _, array := range arrays {
		var temp [][]T
		for _, res := range result {
			for _, elem := range array {
				// 将当前结果与新数组的元素组合
				// note:如果需要一个新的切片元素，一定一定要拷贝原切片的内容
				item := append([]T{}, res...)
				newCombination := append(item, elem)
				temp = append(temp, newCombination)
			}
		}
		result = temp
	}

	return result
}
