package histogram

import (
	"fmt"
	"taskAssignmentForEdge/common"
	"testing"
)

func TestHistogram(t *testing.T) {
	/*histo := [] common.HistoPoint {
		{6,2},
		{2,2},
		{1, 3},
		{5, 3},
	}

	fmt.Println(histo)
	common.SortHistoPoint(histo, func(p, q * common.HistoPoint) bool {
		return p.Value < q.Value
	})
	fmt.Println(histo)*/

	fmt.Println(common.IsValueEqual(2,2))
	fmt.Println(common.IsValueEqual(2,2.0003))
	fmt.Println(common.IsValueEqual(2,2.00000001))

	q := []int{23,19,10,16,36,2,9,32,30,45}
	his := NewHistogram(5)
	fmt.Println(his)

	for i := 0; i < 5; i++ {
		his.Update(float64(q[i]), 1)
	}
	fmt.Println(his)

	for i := 5; i < 7; i++ {
		his.Update(float64(q[i]), 1)
		fmt.Println(his)
	}

	q2 := []int{32,30,45}
	his2 := NewHistogramWithQueue(q2,5)
	fmt.Println(his2)
	his.Merge(his2, 5)
	fmt.Println(his)
	fmt.Println(his.SumFromNInf(1))
	fmt.Println(his.SumFromNInf(2))
	fmt.Println(his.SumFromNInf(15))
	fmt.Println(his.SumFromNInf(19.333333333333332))
	fmt.Println(his.SumFromNInf(19.333333333333332) - his.SumFromNInf(15))
	fmt.Println(his.SumBetweenRange(15, 19.333333333333332))

	fmt.Printf("sum to 45: %f", his.SumFromNInf(45))

	fmt.Println(his.Uniform(3))

	fmt.Printf("Fraction after time-shifting : \n")
	fmt.Println(his.CalFractionAfterShifting(0.0))
	fmt.Println(his.CalFractionAfterShifting(10.0))

	fmt.Printf("Fraction from offset : \n")
	//fmt.Println(his.SumFromNInf(15))
	res := his.CalFractionFromOffset(15)
	sum := 0.0
	for _, v := range res.Frac {
		sum += v.Fraction
	}
	if !common.IsValueEqual(sum, 1.00) {
		t.Errorf("expect:[%f], actually:[%f]", 1.0, sum)
	}
	res = his.CalFractionFromOffset(1)
	sum = 0.0
	for _, v := range res.Frac {
		sum += v.Fraction
	}
	if !common.IsValueEqual(sum, 1.00) {
		t.Errorf("expect:[%f], actually:[%f]", 1.0, sum)
	}
	fmt.Println(his.CalFractionFromOffset(2))
	res = his.CalFractionFromOffset(2)
	sum = 0.0
	for _, v := range res.Frac {
		sum += v.Fraction
	}
	if !common.IsValueEqual(sum, 1.00) {
		t.Errorf("expect:[%f], actually:[%f]", 1.0, sum)
	}
}

func TestNewHistogram(t *testing.T) {
	q := []int{23,19,10,16,36,2,9,32,30,45, 12,34,54,55,4, 6,6, 3,2,4,56,7,8,8,33 ,44}
	his := NewHistogram(5)
	fmt.Println(his)

	for i := 0; i < 5; i++ {
		his.Update(float64(q[i]), 1)
	}
	fmt.Println(his)
}

func TestJointProbaByEqualOrLess(t *testing.T) {
	X := &HistogramFrac{[]HistoFrac{{20, 0.2},{40, 0.3},{60, 0.5}}}
	Y := &HistogramFrac{[]HistoFrac{{30,0.3},{30,0.3},{50,0.3}}}
	fmt.Println(X.JointProbaByEqualOrLess(Y))
}

func TestHistogram_CalFractionFromOffset(t *testing.T) {
	q := []int{23,19,10,16,36,2,9,32,30,45}
	his := NewHistogram(5)
	fmt.Println(his)

	for i := 0; i < len(q); i++ {
		his.Update(float64(q[i]), 1)
	}
	fmt.Println(his)

	fmt.Printf("Fraction from offset : \n")
	fmt.Printf("offset = 0\n")
	fmt.Println(his.SumFromNInf(0))
	res := his.CalFractionFromOffset(0)
	fmt.Println(res)

	fmt.Printf("offset = 1\n")
	fmt.Println(his.SumFromNInf(1))
	res = his.CalFractionFromOffset(1)
	fmt.Println(res)

	fmt.Printf("offset = 2\n")
	fmt.Println(his.SumFromNInf(2))
	res = his.CalFractionFromOffset(2)
	fmt.Println(res)

	fmt.Printf("offset = 15\n")
	fmt.Println(his.SumFromNInf(15))
	res = his.CalFractionFromOffset(15)
	fmt.Println(res)

	fmt.Printf("offset = 45\n")
	fmt.Println(his.SumFromNInf(45))
	res = his.CalFractionFromOffset(45)
	fmt.Println(res)

	fmt.Printf("offset = 46\n")
	fmt.Println(his.SumFromNInf(46))
	res = his.CalFractionFromOffset(46)
	fmt.Println(res)
}