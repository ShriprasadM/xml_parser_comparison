package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringBased(t *testing.T) {
	vast := stringBased(vast)
	assert.True(t, strings.Contains(vast, SampleTrackingEvent))
	// fmt.Println(vast)
}

func TestEtreeBased(t *testing.T) {
	vast, _ := etreeBased(vast)
	assert.True(t, strings.Contains(vast, SampleTrackingEvent))
	fmt.Println(vast)
}

func TestXmlEncodingBased(t *testing.T) {
	vast, _ := xmlEncodingBased(vast)
	assert.True(t, strings.Contains(vast, SampleTrackingEvent))
	// fmt.Println(vast)
}

func BenchmarkStringBased(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stringBased(vast)
	}
}

func BenchmarkEtreeBased(b *testing.B) {
	for i := 0; i < b.N; i++ {
		etreeBased(vast)
	}
}

func BenchmarkXmlEncodingBased(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xmlEncodingBased(vast)
	}
}
