package client

import (
	"fmt"
	"testing"
)

func Test_Set(t *testing.T) {
	client, err := NewClient("127.0.0.1:9999")
	if err != nil {
		t.Fatal(err)
	}
	// RESP Arrays: *3\r\n$3\r\nset\r\n$4\r\nname\r\n$4\r\nmars\r\n
	rsp, err := client.Set("name", "mars")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rsp)
}

func Test_Get(t *testing.T) {
	client, err := NewClient("127.0.0.1:9999")
	if err != nil {
		t.Fatal(err)
	}
	// RESP Arrays: *2\r\n$3\r\nget\r\n$4\r\nname\r\n
	rsp, err := client.Get("name")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rsp)
}

func Test_Del(t *testing.T) {
	client, err := NewClient("127.0.0.1:9999")
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := client.Del("name")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rsp)
}
