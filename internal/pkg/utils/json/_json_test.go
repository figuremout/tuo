// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"fmt"
)

func Stringfy(v interface{}, indent string) string {
	b, err := json.MarshalIndent(v, "", indent)
	if err != nil {
		return ""
	}
	return string(b)
}

type Base struct {
	Error string `json:"error"`
}

type Resp struct {
	Base
	Name string `json:"name"`
}

func c2c() {
	r := Resp{
		Base: Base{
			Error: "this is error",
		},
		Name: "this is name",
	}
	str := Stringfy(&r, "")
	fmt.Print(str)
	var resp Resp
	if err := json.Unmarshal([]byte(str), &resp); err != nil {
		fmt.Print(err)
	}
	fmt.Print(resp.Error)
	fmt.Print(resp.Name)
}

func c2f() {
	b := Base{
		Error: "this is error",
	}
	str := Stringfy(&b, "")
	fmt.Print(str)
	var resp Resp
	if err := json.Unmarshal([]byte(str), &resp); err != nil {
		fmt.Print(err)
	}
	//fmt.Print(resp.Error)
	fmt.Print(resp.Name)
}

func main() {
	fmt.Println("Hello, 世界")
	c2f()
	c2c()
}
