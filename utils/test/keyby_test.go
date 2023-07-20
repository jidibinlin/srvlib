package main

import (
	"fmt"
	"github.com/gzjjyz/srvlib/utils"
	"testing"
)

type A struct {
	Name string `json:"name"`
	Val  string `json:"val"`
}

func TestKeyBy(t *testing.T) {
	var as = []*A{
		{
			Name: "a",
			Val:  "a1",
		},
		{
			Name: "b",
			Val:  "b1",
		},
	}
	am := utils.KeyBy(as, "Name").(map[string]*A)
	for k, v := range am {
		fmt.Printf("k is %v, val is %v\n", k, v)
	}

	ns := utils.PluckStrings(as, "Name")
	for _, n := range ns {
		fmt.Printf("n is %v\n", n)
	}
}
