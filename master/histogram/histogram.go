package histogram

import (
	"taskAssignmentForEdge/common"
	"math"
	"sort"
	"sync"
)

type HistoPoint struct{
	Value float64
	Freq int
}

type HistoFrac struct {
	Value float64
	Fraction float64
}

type Histogram struct {
	maxNum int
	//totalNum int
	histo  [] HistoPoint
	rwlock sync.RWMutex
}

//method of Histogram
func NewHistogram(maxPointNum int) *Histogram {
	histog := new(Histogram)
	histog.maxNum = maxPointNum
	histog.histo = make([]HistoPoint, 0)//[]HistoPoint{}
	return histog
}

func (histog *Histogram) GetHistoPoint() []HistoPoint {
	return histog.histo
}

func NewCopyHistogram(maxPointNum int, histog *Histogram) *Histogram {
	newhistog := new(Histogram)
	newhistog.maxNum = maxPointNum
	histog.rwlock.RLock()
	newhistog.histo = make([]HistoPoint, len(histog.GetHistoPoint()))
	copy(newhistog.histo, histog.GetHistoPoint())
	histog.rwlock.RUnlock()
	return newhistog
}

func (histog *Histogram) GetMaxNum() int {
	histog.rwlock.RLock()
	defer histog.rwlock.RUnlock()
	return histog.maxNum
}

func (histog *Histogram) ResetMaxNum(newMaxNum int) {
	histog.rwlock.Lock()
	histog.maxNum = newMaxNum
	histog.rwlock.Unlock()
}

func (histog *Histogram) Update(newValue float64, freq int) {
	histog.rwlock.Lock()
	defer histog.rwlock.Unlock()
	for i, v := range histog.histo {
		if common.IsValueEqual(v.Value, newValue) {
			histog.histo[i].Freq += freq
			return
		}
	}
	//add a new point
	histog.histo = append(histog.histo, HistoPoint{newValue, freq})

	//increase
	sort.Slice(histog.histo, func(i,j int) bool {
		return histog.histo[i].Value < histog.histo[j].Value
	})

	if len(histog.histo) <= histog.maxNum {
		return
	}

	// need to merge

	for len(histog.histo) > histog.maxNum {
		//find a point i that minimaizes value(i+1) minus value(i)
		pos := 0
		diffMin := histog.histo[1].Value-histog.histo[0].Value
		for i := 1; i < len(histog.histo) - 1; i++ {
			diff := histog.histo[i+1].Value - histog.histo[i].Value
			if  diff < diffMin {
				diffMin = diff
				pos = i
			}
		}
		newFreq := histog.histo[pos].Freq + histog.histo[pos+1].Freq
		newValue := (histog.histo[pos].Value * float64(histog.histo[pos].Freq) + histog.histo[pos+1].Value * float64(histog.histo[pos+1].Freq)) / float64(newFreq)
		histog.histo[pos].Value = newValue
		histog.histo[pos].Freq = newFreq
		histog.histo = append(histog.histo[:pos+1], histog.histo[pos+2:]...)
	}
}

func NewHistogramWithQueue(q []int, maxPointNum int) *Histogram {
	histog := NewHistogram(maxPointNum)
	for _, v := range q {
		histog.Update(float64(v),1)
	}
	return histog
}

func (histog *Histogram) Merge(otherHistog *Histogram, maxNum int) {
	if otherHistog == nil {
		return
	}

	if maxNum > histog.maxNum {
		histog.maxNum = maxNum
	}
	for _, v := range otherHistog.histo {
		histog.Update(v.Value, v.Freq)
	}
}

func (histog *Histogram) SumFromNInf(b float64)  float64 {
	if b < histog.histo[0].Value {
		return 0.0
	}

	i := 0
	for ; i < len(histog.histo) - 1; i++ {
		if (histog.histo[i].Value < b || common.IsValueEqual(histog.histo[i].Value, b) ) && histog.histo[i+1].Value > b {
			break
		}
	}

	sum := 0.0
	for j:=0; j < i; j++ {
		sum += float64(histog.histo[j].Freq)
	}

	if i >= len(histog.histo) - 1 {
	    if histog.histo[i].Value < b { //b large than any value, add up all
			sum += float64(histog.histo[i].Freq)
		} else {
			sum += float64(histog.histo[i].Freq)/2  // b equals the last value, add up half of it
		}
	} else {
		sum += float64(histog.histo[i].Freq)/2
		var Freqb float64 = float64(histog.histo[i].Freq) + float64(histog.histo[i+1].Freq - histog.histo[i].Freq)*(b - histog.histo[i].Value)/(histog.histo[i+1].Value - histog.histo[i].Value)
		sum += (float64(histog.histo[i].Freq) + Freqb)*(b - histog.histo[i].Value)/float64(histog.histo[i+1].Value - histog.histo[i].Value)/2
	}

	return sum
}

func (histog *Histogram) SumBetweenRange(begin float64, end float64)  float64 {
	if begin > end {
		return 0
	}
	return histog.SumFromNInf(end) - histog.SumFromNInf(begin)
}

func (histog *Histogram) Uniform(B int)  []float64 {
	if B >= histog.maxNum {
		return nil
	}

	uniform := []float64{}
	totalFreq := 0
	for _, v := range histog.histo {
		totalFreq += v.Freq
	}

	i := 0
	for j := 1; j < B; j++ {
		s := float64(j)/float64(B)*float64(totalFreq)
		for ; i < len(histog.histo)-1; i++ {
			s1 := histog.SumFromNInf(histog.histo[i].Value)
			s2 := histog.SumFromNInf(histog.histo[i+1].Value)
			if s1 <= s &&  s2 > s {
				break
			}
		}
		if i == len(histog.histo) - 1 {
			break
		}
		d := s - histog.SumFromNInf(histog.histo[i].Value)
		a := float64(histog.histo[i+1].Freq - histog.histo[i].Freq) + 0.0000001
		b := 2 * float64(histog.histo[i].Freq)
		c := -2 * d
		z := (-b + math.Sqrt(b*b-4*a*c))/(2*a)
		u := histog.histo[i].Value + (histog.histo[i+1].Value - histog.histo[i].Value)*z
		uniform = append(uniform, u)
	}
	return uniform
}

//distribution function
type HistogramFrac struct {
	Frac []HistoFrac
}

func NewHistogramFrac() *HistogramFrac{
	histogfrac := new(HistogramFrac)
	histogfrac.Frac = make([]HistoFrac, 0)
	return  histogfrac
}

func NewHistogramFracWithLen(len int) *HistogramFrac{
	histogfrac := new(HistogramFrac)
	histogfrac.Frac = make([]HistoFrac, len)
	return  histogfrac
}

//shift the histogram to the right with [offset], then calculate the fraction of points
//仅平移, 用于任务时间
func (histog *Histogram) CalFractionAfterShifting(offset float64)  *HistogramFrac {
	histog.rwlock.RLock()
	defer histog.rwlock.RUnlock()
	histogfrac := NewHistogramFracWithLen(len(histog.histo))

	totalNum := 0;
	for _,v := range histog.histo {
		totalNum += v.Freq
	}

	for i:=0;i<len(histog.histo);i++ {
		histogfrac.Frac[i].Value = histog.histo[i].Value + offset
		histogfrac.Frac[i].Fraction = float64(histog.histo[i].Freq)/float64(totalNum)
	}
	return histogfrac
}

//calculate the fraction from [offset]
//按比例更改概率, 用于存在时间,
func (histog *Histogram) CalFractionFromOffset(offset float64) *HistogramFrac{
	if histog.histo == nil {
		return nil
	}

	if len(histog.histo) == 0 { //empty
		return NewHistogramFrac()
	}

	histog.rwlock.RLock()
	defer histog.rwlock.RUnlock()

	if histog.histo[0].Value > offset {
		histog.rwlock.RUnlock()
		res := histog.CalFractionAfterShifting(-offset)
		histog.rwlock.RLock()
		return res
	}

	histogfrac := NewHistogramFrac()

	if histog.histo[len(histog.histo) - 1].Value < offset {
		return histogfrac
	}

	totalFreq := 0.0
	for _, v := range histog.histo {
		totalFreq += float64(v.Freq)
	}

	pos := 0
	for ; pos < len(histog.histo) - 1; pos++ {
		if (histog.histo[pos].Value < offset || common.IsValueEqual(histog.histo[pos].Value, offset)) && histog.histo[pos + 1].Value > offset {
			break
		}
	}

	if pos == len(histog.histo) - 1  {
		frac := HistoFrac{0, 1.0}
		histogfrac.Frac = append(histogfrac.Frac, frac)
	} else {
		remainFreq := totalFreq - histog.SumFromNInf(offset)
		if offset > histog.histo[pos].Value {
			pos++
			freq := float64(histog.histo[pos].Freq)/2 + histog.SumBetweenRange(offset, histog.histo[pos].Value)
			frac := HistoFrac{histog.histo[pos].Value-offset, freq/remainFreq}
			histogfrac.Frac = append(histogfrac.Frac, frac)
		} else if common.IsValueEqual(histog.histo[pos].Value, offset) {
			freq := float64(histog.histo[pos].Freq)/2
			frac := HistoFrac{histog.histo[pos].Value-offset, freq/remainFreq}
			histogfrac.Frac = append(histogfrac.Frac, frac)
		}
		pos++

		for i:=pos;i<len(histog.histo);i++ {
			frac := HistoFrac{histog.histo[i].Value-offset,float64(histog.histo[i].Freq)/remainFreq}
			histogfrac.Frac = append(histogfrac.Frac, frac)
		}
	}

	return histogfrac
}

// joint probability distribution of X <= Y
func (X *HistogramFrac) JointProbaByEqualOrLess( Y *HistogramFrac) (proba float64, extra float64) {
	//proba := 0.0
	if Y == nil || Y.Frac == nil || len(Y.Frac) == 0 {
		return  0, 0
	}
	for _, PX := range X.Frac {
		for _, PY := range Y.Frac {
			if PX.Value <= PY.Value {
				proba += PX.Fraction*PY.Fraction
				extra += PX.Fraction*PY.Fraction*PX.Value
			}
		}
	}
	return proba, extra
}

//P(start<=x<end)
func (X *HistogramFrac) CalcProbabilityInRange(start float64, end float64) float64 {
	if X.Frac == nil || len(X.Frac) == 0 {
		return 0
	}
	if X.Frac[0].Value >= end {
		return 0
	}
	if X.Frac[len(X.Frac)-1].Value < start {
		return 0
	}
	sumFrac := 0.0
	for _, v := range X.Frac {
		if (v.Value > start || common.IsValueEqual(v.Value, start)) && v.Value < end {
			sumFrac += v.Fraction
		}
	}
	return  sumFrac
}