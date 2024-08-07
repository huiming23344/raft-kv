package client

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

type RecordJosn struct {
	RecordKey   string `json:"record_key"`
	RecordValue string `json:"record_value"`
}

func Test_Set(t *testing.T) {
	client, err := NewClient("127.0.0.1:2317")
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

func Test_Set_FromFile(t *testing.T) {
	var data []RecordJosn
	jsonData, err := os.ReadFile("../testFiles/data.json")
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", data)
	for i, record := range data {
		client, err := NewClient("127.0.0.1:2317")
		if err != nil {
			t.Fatal(err)
		}
		// RESP Arrays: *3\r\n$3\r\nset\r\n$4\r\nname\r\n$4\r\nmars\r\n
		fmt.Printf("SET%d: RecordKey: %s, RecordValue: %s\n", i, record.RecordKey, record.RecordValue)
		rsp, err := client.Set(record.RecordValue, record.RecordValue)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(rsp)
	}
}

func Test_Get(t *testing.T) {
	client, err := NewClient("127.0.0.1:2317")
	if err != nil {
		t.Fatal(err)
	}
	// RESP Arrays: *2\r\n$3\r\nget\r\n$4\r\nname\r\n
	rsp, err := client.Get("6")
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
