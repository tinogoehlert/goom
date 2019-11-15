package utils

func Clamp(num, min, max float32) float32 {
	if num < min {
		return min
	}
	if num > max {
		return max
	}
	return num
}
