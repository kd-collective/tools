// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains tests for the useless-assignment checker.

package testdata

import "math/rand"

type ST[T interface{ ~int }] struct {
	x T
	l []T
}

func (s *ST[T]) SetX(x T, ch chan T) {
	// Accidental self-assignment; it should be "s.x = x"
	x = x // want "self-assignment of x"
	// Another mistake
	s.x = s.x // want "self-assignment of s.x"

	s.l[0] = s.l[0] // want "self-assignment of s.l.0."

	// Bail on any potential side effects to avoid false positives
	s.l[num()] = s.l[num()]
	rng := rand.New(rand.NewSource(0))
	s.l[rng.Intn(len(s.l))] = s.l[rng.Intn(len(s.l))]
	s.l[<-ch] = s.l[<-ch]
}

func num() int { return 2 }
