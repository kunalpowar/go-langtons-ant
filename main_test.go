package main

import (
	"testing"

	log "github.com/Sirupsen/logrus"
)

func BenchmarkIterate(b *testing.B) {
	size := 2000
	g := make(grid, size)
	for i, _ := range g {
		g[i] = make([]bool, size)
	}

	a := ant{dir: dirUp, r: size / 2, c: size / 2}

	log.SetLevel(log.InfoLevel)

	for i := 0; i < b.N; i++ {
		g.iterate(&a)
	}
}
