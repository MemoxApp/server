package resource

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
)

func ParseResources(content string) []string {
	r, _ := regexp.Compile("#\\S[^\\n]+? ")
	return r.FindAllString(content, -1)
}

// UniqueArr 去重
func UniqueArr(m []primitive.ObjectID) []primitive.ObjectID {
	d := make([]primitive.ObjectID, 0)
	tempMap := make(map[primitive.ObjectID]bool, len(m))
	for _, v := range m { // 以值作为键名
		if tempMap[v] == false {
			tempMap[v] = true
			d = append(d, v)
		}
	}
	return d
}

// RemoveFromArr 删除
func RemoveFromArr(m []primitive.ObjectID, v primitive.ObjectID) []primitive.ObjectID {
	for i, val := range m {
		if val == v {
			m = append(m[:i], m[i+1:]...)
			break
		}
	}
	return m
}

// DiffArray 列出 left 与 right 的差集
func DiffArray(left []string, right []string) ([]string, []string) {
	var leftNotContains = left
	var rightNotContains []string
	for _, val := range right {
		var result bool
		leftNotContains, result = TryRemoveFromArr(leftNotContains, val)
		if !result {
			// left中不包含该元素，加入rightNotContains
			rightNotContains = append(rightNotContains, val)
		}
	}
	return leftNotContains, rightNotContains
}

// TryRemoveFromArr 尝试删除
func TryRemoveFromArr(m []string, v string) ([]string, bool) {
	for i, val := range m {
		if val == v {
			return append(m[:i], m[i+1:]...), true
		}
	}
	return m, false
}
