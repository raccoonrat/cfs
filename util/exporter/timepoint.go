// Copyright 2018 The Chubao Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package exporter

import (
	"fmt"
	"sync"
	"time"
)

var (
	TPPool = &sync.Pool{New: func() interface{} {
		return new(TimePoint)
	}}
	TPCh chan *TimePoint
)

func collectTP() {
	TPCh = make(chan *TimePoint, ChSize)
	for {
		m := <-TPCh
		metric := m.Metric()
		metric.Set(float64(m.val))
		TPPool.Put(m)
	}
}

type TimePoint struct {
	Gauge
	startTime time.Time
}

type TimePointCount struct {
	tp  *TimePoint
	cnt *Counter
}

func NewTP(name string) (tp *TimePoint) {
	if !enabled {
		return
	}
	tp = TPPool.Get().(*TimePoint)
	tp.name = metricsName(name)
	tp.startTime = time.Now()
	return
}

func (tp *TimePoint) Set() {
	if !enabled {
		return
	}
	val := time.Since(tp.startTime).Nanoseconds()
	tp.val = val
	tp.publish()
}

func NewTPCnt(name string) (tpc *TimePointCount) {
	if !enabled {
		return
	}
	tpc = new(TimePointCount)
	tpc.tp = NewTP(name)
	tpc.cnt = NewCounter(fmt.Sprintf("%s_count", name))
	return
}

func (tpc *TimePointCount) Set() {
	if !enabled {
		return
	}
	tpc.tp.Set()
	tpc.cnt.Add(1)
}

func (tp *TimePoint) publish() {
	select {
	case TPCh <- tp:
	default:
	}
}
