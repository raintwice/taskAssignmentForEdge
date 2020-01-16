package common

import (
	"fmt"
	"testing"
)

type StuScore struct {
	name string
	score int
}

func TestList(t *testing.T) {
	course := NewListWithHash()
	scores := []StuScore{{"zhangsan", 100}, {"lisi", 89}, {"wangwu", 60}}
	//ranks  := []string{"A", 'B', 'C'}
	for _, v := range scores {
		//fmt.Println(v)
		course.AddUnique(v.name, v)
	}

	for e := course.GetList().Front();e != nil; e = e.Next(){
		fmt.Println(e.Value)
	}

	lefte := course.Remove("lisi")
	fmt.Printf("The remain elements in the list:\n")
	for e := course.GetList().Front();e != nil; e = e.Next(){
		fmt.Println(e.Value)
	}
	fmt.Printf("The element that was deleted:\n")
	fmt.Println(lefte)

	newRecord := StuScore{"mayuin",10000}
	course.Add(newRecord.name, newRecord)
	fmt.Printf("After adding a new element, the remain elements in the list:\n")
	for e := course.GetList().Front();e != nil; e = e.Next(){
		fmt.Println(e.Value)
	}
}

type Person struct {
	name string
	gender bool
}

type Score struct {
	student Person
	grade int
}

func TestListStructure(t *testing.T) {
	course := NewListWithHash()
	scores := []Score{{Person{"zhangsan", true}, 100},{ Person{"lisi", true},99}, {Person{"wangwu", false},98}}

	for i := 0; i < len(scores); i++ {
		course.AddUnique(scores[i].student, &scores[i])
	}

	for e := course.GetList().Front();e != nil; e = e.Next(){
		fmt.Println(e.Value)
	}

	lefte := course.Remove(scores[1].student)
	fmt.Printf("The remain elements in the list:\n")
	for e := course.GetList().Front();e != nil; e = e.Next(){
		fmt.Println(e.Value)
	}
	fmt.Printf("The element that was deleted:\n")
	fmt.Println(lefte)

	newRecord := Score{Person{"mayuin",false},10000}
	course.Add(newRecord, "D")
	fmt.Printf("After adding a new element, the remain elements in the list:\n")
	for e := course.GetList().Front();e != nil; e = e.Next(){
		fmt.Println(e.Value)
	}
}