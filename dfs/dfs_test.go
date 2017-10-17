package dfs

import (
	"testing"
)

func GenerateJsonPayload(b *testing.B, maxWidth, maxDepth int) []byte {

	data, err := CreateTestJson(maxWidth, maxDepth)
	if err != nil {
		b.Fatal(err)
	}

	//fmt.Printf("mw: %d\tmd: %d\t size: %0.02fkb\n", maxWidth, maxDepth, KB(data))

	return data
}

func ValidateArgs(b *testing.B, d *dfs, jsonPayload []byte, maxDepth int) {
	// fmt.Println(d.Validate(jsonPayload))
	//d.Validate(jsonPayload)
	CalculateJsonDepth(jsonPayload, maxDepth)
}

func Benchmark_pw1_pd1_maxw1_maxd1(b *testing.B) {

	jsonPayload := GenerateJsonPayload(b, 1, 1)
	d := New(1, 1)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ValidateArgs(b, d, jsonPayload, 1)
	}
}

func Benchmark_pw2_pd2_maxw1_maxd1(b *testing.B) {
	jsonPayload := GenerateJsonPayload(b, 2, 2)
	d := New(1, 1)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ValidateArgs(b, d, jsonPayload, 1)
	}
}

func Benchmark_pw2_pd2_maxw2_maxd2(b *testing.B) {
	jsonPayload := GenerateJsonPayload(b, 2, 2)
	d := New(2, 2)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ValidateArgs(b, d, jsonPayload, 2)
	}
}

func Benchmark_pw4_pd4_maxw2_maxd2(b *testing.B) {
	jsonPayload := GenerateJsonPayload(b, 4, 4)
	d := New(2, 2)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ValidateArgs(b, d, jsonPayload, 2)
	}
}

func Benchmark_pw4_pd4_maxw4_maxd4(b *testing.B) {
	jsonPayload := GenerateJsonPayload(b, 4, 4)
	d := New(4, 4)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ValidateArgs(b, d, jsonPayload, 4)
	}
}

func Benchmark_pw8_pd8_maxw4_maxd4(b *testing.B) {
	jsonPayload := GenerateJsonPayload(b, 8, 8)
	d := New(4, 4)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ValidateArgs(b, d, jsonPayload, 4)
	}
}

func Benchmark_pw8_pd8_maxw8_maxd8(b *testing.B) {
	jsonPayload := GenerateJsonPayload(b, 8, 8)
	//d := New(8, 8)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ValidateArgs(b, nil, jsonPayload, 8)
	}
}

//////////////////////////////////////////////////////////////////////.
func Benchmark_pw10_pd4_maxw10_maxd3(b *testing.B) {
	jsonPayload := GenerateJsonPayload(b, 10, 4)
	d := New(10, 3)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ValidateArgs(b, d, jsonPayload, 3)
	}
}

func Benchmark_pw10_pd5_maxw10_maxd3(b *testing.B) {
	jsonPayload := GenerateJsonPayload(b, 10, 5)
	d := New(10, 3)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ValidateArgs(b, d, jsonPayload, 3)
	}
}

func Benchmark_pw140_pd3_maxw10_maxd3(b *testing.B) {
	jsonPayload := GenerateJsonPayload(b, 140, 3)
	d := New(10, 3)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ValidateArgs(b, d, jsonPayload, 3)
	}
}
