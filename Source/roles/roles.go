package roles

import "strconv"

func IsMasterAlive(peers []string) bool {
	for _, v := range peers {
		if v[0] == 'M' {
			return true
		}
	}
	return false
}
func HowManyMasters(peers []string) int {
	i := 0
	for _, v := range peers {
		if v[0] == 'M' {
			i++
		}
	}
	return i
}
func MastersID(peers []string) string {

	for _, v := range peers {
		if v[0] == 'M' {
			return v[1:]
		}
	}
	return ""
}

func MaxIdAlive(peers []string) int {
	max := 0
	for _, v := range peers {
		x, _ := strconv.Atoi(v[1:])
		if x > max {
			max = x

		}
	}
	return max
}
func AmiAlone(peers []string) bool {
	if len(peers) == 1 {
		return true
	} else {
		return false
	}
}
