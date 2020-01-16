package predictor

import (
	"fmt"
	"testing"
)

func TestPredictUnit(t *testing.T) {
	dataSource := []int64{107,92,99,103,94,100,104,91,108,111,103,100,88,105,88,116,105,103,96,78,88,104,124,107,103,95,101,110,97,119,88,96,101,80,106,101,101,108,110,108,95,100,104,87,90,96,101,118,108,96}

	preUnit := NewPredictUnit(FeaturePair{FeatureType_logic_Job_Name, "MapReduce"})
	fmt.Println(preUnit.Predict())

	for _, v := range dataSource {
		preUnit.Update(v)
		fmt.Println(preUnit.Predict())
	}
	preUnit.ShowHistogram()
	fmt.Println(preUnit.Predict())

}