package hashring

import "testing"

func BenchmarkNew(b *testing.B) {
	nodes := []string{"a", "b", "c", "d", "e", "f", "g"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(nodes)
	}
}

func BenchmarkHashes(b *testing.B) {
	nodes := []string{"a", "b", "c", "d", "e", "f", "g"}
	ring := New(nodes)
	tt := []struct {
		key   string
		nodes []string
	}{
		{"test", []string{"a", "b"}},
		{"test", []string{"a", "b"}},
		{"test1", []string{"b", "d"}},
		{"test2", []string{"f", "b"}},
		{"test3", []string{"f", "c"}},
		{"test4", []string{"c", "b"}},
		{"test5", []string{"f", "a"}},
		{"aaaa", []string{"b", "a"}},
		{"bbbb", []string{"f", "a"}},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o := tt[i%len(tt)]
		ring.GetNodes(o.key, 2)
	}
}

func BenchmarkHashesSingle(b *testing.B) {
	nodes := []string{"a", "b", "c", "d", "e", "f", "g"}
	ring := New(nodes)
	tt := []struct {
		key   string
		nodes []string
	}{
		{"test", []string{"a", "b"}},
		{"test", []string{"a", "b"}},
		{"test1", []string{"b", "d"}},
		{"test2", []string{"f", "b"}},
		{"test3", []string{"f", "c"}},
		{"test4", []string{"c", "b"}},
		{"test5", []string{"f", "a"}},
		{"aaaa", []string{"b", "a"}},
		{"bbbb", []string{"f", "a"}},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o := tt[i%len(tt)]
		ring.GetNode(o.key)
	}
}
