package workpool

import (
	"testing"
)

func BenchmarkSprintfCommon(b *testing.B) {
	b.ResetTimer()
	numbers := []int{1, 2, 3, 3, 4, 5, 6, 2}
	Process(numbers)
}

func BenchmarkSprintfPool(b *testing.B) {
	b.ResetTimer()
	numbers := []int{1, 2, 3, 3, 4, 5, 6, 2}
	PoolProcess(numbers)
}
