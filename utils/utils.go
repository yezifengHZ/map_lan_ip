package utils

// ConvertStrSliceToMap 将字符串 slice 转为 map[string]struct{}
func ConvertStrSliceToMap(sl []string) map[string]struct{} {
	set := make(map[string]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// ContainsInMap 判断字符串是否在 map 中
func ContainsInMap(m map[string]struct{}, s string) bool {
	_, ok := m[s]
	return ok
}
