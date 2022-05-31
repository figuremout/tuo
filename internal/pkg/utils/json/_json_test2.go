package main

import (
	"encoding/json"
	"fmt"
)

// 结构体定义
type robot struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type resp struct {
	Robots []robot `json:"robots"`
}

func test1() {
	r := resp{
		Robots: []robot{
			{Name: "a", Amount: 1},
			{Name: "b", Amount: 2},
		},
	}
	bs, _ := json.Marshal(r)
	fmt.Println(string(bs))
}

func test2() {
	str := "{\"robots\":[{\"name\":\"a\",\"amount\":1},{\"name\":\"b\",\"amount\":2}]}"
	r := resp{}
	json.Unmarshal([]byte(str), &r)
	bs, _ := json.Marshal(r)
	fmt.Println(string(bs))
	fmt.Println(r.Robots[0])
}

// 解析到结构体数组
func parse_array() {
	fmt.Println("解析json字符串为结构体数组")
	str := "[{\"name\":\"name1\",\"amount\":100},{\"name\":\"name2\",\"amount\":200},{\"name\":\"name3\",\"amount\":300},{\"name\":\"name4\",\"amount\":400}]"
	all := []robot{}
	err := json.Unmarshal([]byte(str), &all)
	if err != nil {
		fmt.Printf("err=%v", err)
	}
	for _, one := range all {
		fmt.Printf("name=%v, amount=%v\n", one.Name, one.Amount)
	}
}

// 解析到结构体指针的数组
func parse_pointer_array() {
	fmt.Println("解析json字符串为结构体指针的数组")
	str := "[{\"name\":\"name1\",\"amount\":100},{\"name\":\"name2\",\"amount\":200},{\"name\":\"name3\",\"amount\":300},{\"name\":\"name4\",\"amount\":400}]"
	all := []*robot{}
	err := json.Unmarshal([]byte(str), &all)
	if err != nil {
		fmt.Printf("err=%v", err)
	}
	for _, one := range all {
		fmt.Printf("name=%v, amount=%v\n", one.Name, one.Amount)
	}
}
func main() {
	test1()
	test2()
	// 解析为结构体数组
	parse_array()

	// 解析为结构体指针的数组
	parse_pointer_array()
}
