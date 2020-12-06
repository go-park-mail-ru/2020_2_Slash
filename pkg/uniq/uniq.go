package uniq

func RemoveDuplicates(elements []uint64) []uint64 {
	encountered := map[uint64]bool{}
	result := []uint64{}

	for _, elem := range elements {
		if !encountered[elem] {
			encountered[elem] = true
			result = append(result, elem)
		}
	}
	return result
}
