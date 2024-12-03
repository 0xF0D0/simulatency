package util

import (
	"github.com/hashicorp/go-set/v3"
	"gonum.org/v1/gonum/stat/combin"
)

type CombinationMap interface {
	Insert(val int) []int
	Remove(val int) []int
}

//This map calculates absolute value between 2 distinct values
type combinationMap struct {
	internalVals set.Set[int]
	internalAbss set.Set[int]
}

func NewCombinationMap() CombinationMap {
	return &combinationMap{
		internalVals: *set.New[int](0),
		internalAbss: *set.New[int](0),
	}
}

//Returns new added abss
func (c *combinationMap) Insert(val int) []int {
	isNewValue := c.internalVals.Insert(val)
	if !isNewValue || c.internalVals.Size() == 1 {
		return []int{}
	}

	//Get old abss and diff new abss
	newAbss := c.getCombination(c.internalVals)
	result := newAbss.Difference(&c.internalAbss).Slice()
	c.internalAbss = *newAbss
	return result
}

//Returns deleted abss
func (c *combinationMap) Remove(val int) []int {
	somethingDeleted := c.internalVals.Remove(val)
	if !somethingDeleted {
		return []int{}
	}

	//Get old abss and diff new abss
	newAbss := c.getCombination(c.internalVals)
	result := c.internalAbss.Difference(newAbss).Slice()
	c.internalAbss = *newAbss
	return result
}

func (c *combinationMap) getCombination(vals set.Set[int]) *set.Set[int] {
	
	valsList := vals.Slice()
	combList := combin.Combinations(vals.Size(), 2)
	newSet := set.New[int](len(valsList) * (len(valsList) - 1) / 2)
	for i, iList := range combList {
		for _, j := range iList {
			newSet.Insert(getAbs(valsList[i], valsList[j]))
		}
	}
	return newSet
}

func getAbs(a, b int) int {
	r := a - b
	if r < 0 {
		return -r
	}
	return r
}