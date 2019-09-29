package main

import (
	"sort"
	"strconv"
)

/*
	сюда вам надо писать функции, которых не хватает, чтобы проходили тесты в gotchas_test.go

	IntSliceToString
	MergeSlices
	GetMapValuesSortedByKey
*/

func IntSliceToString(array []int) string {
	var res string
	for _, v := range array {
		res += strconv.Itoa(v)
	}
	return res
}

func MergeSlices(farr []float32, iarr []int32) []int {
	arr := make([]int, len(farr)+len(iarr))
	for i := range farr {
		arr[i] = int(farr[i])
	}
	for i := range iarr {
		arr[len(farr)+i] = int(iarr[i])
	}
	return arr
}

func GetMapValuesSortedByKey(input map[int]string) []string {
	var keys []int

	for k := range input {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	val := make([]string, len(keys))
	for i, _ := range keys {
		val[i] = input[keys[i]]
	}
	return val
}
