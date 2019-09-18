// Copyright (c) 2018 Ashley Jeffs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, sub to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package metrics

//------------------------------------------------------------------------------

type fCounterVec struct {
	f func([]string) StatCounter
}

func (f *fCounterVec) With(labels ...string) StatCounter {
	return f.f(labels)
}

func fakeCounterVec(f func([]string) StatCounter) StatCounterVec {
	return &fCounterVec{
		f: f,
	}
}

//------------------------------------------------------------------------------

type fTimerVec struct {
	f func([]string) StatTimer
}

func (f *fTimerVec) With(labels ...string) StatTimer {
	return f.f(labels)
}

func fakeTimerVec(f func([]string) StatTimer) StatTimerVec {
	return &fTimerVec{
		f: f,
	}
}

//------------------------------------------------------------------------------

type fGaugeVec struct {
	f func([]string) StatGauge
}

func (f *fGaugeVec) With(labels ...string) StatGauge {
	return f.f(labels)
}

func fakeGaugeVec(f func([]string) StatGauge) StatGaugeVec {
	return &fGaugeVec{
		f: f,
	}
}

//------------------------------------------------------------------------------
