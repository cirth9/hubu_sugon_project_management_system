package t

import "testing"

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if ans := Add(1, 2); ans != 3 {
			b.Errorf("1 + 2 expected be 3, but %d got", ans)
		}
	}
}
