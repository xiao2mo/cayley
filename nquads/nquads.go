// Copyright 2014 The Cayley Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nquads

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/google/cayley/graph"
)

var (
	ErrAbsentSubject   = errors.New("nqauds: absent subject")
	ErrAbsentPredicate = errors.New("nqauds: absent predicate")
	ErrAbsentObject    = errors.New("nqauds: absent object")
	ErrUnterminated    = errors.New("nqauds: unterminated quad")
)

func Parse(str string) (*graph.Triple, error) {
	t, err := parse([]rune(str))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

type Decoder struct {
	r    *bufio.Reader
	line []byte
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r)}
}

func (dec *Decoder) Unmarshal() (*graph.Triple, error) {
	dec.line = dec.line[:0]
	for {
		l, pre, err := dec.r.ReadLine()
		if err != nil {
			return nil, err
		}
		dec.line = append(dec.line, l...)
		if !pre {
			break
		}
	}
	triple, err := Parse(string(dec.line))
	if err != nil {
		return nil, fmt.Errorf("failed to parse %q: %v", dec.line, err)
	}
	if triple == nil {
		return dec.Unmarshal()
	}
	return triple, nil
}
