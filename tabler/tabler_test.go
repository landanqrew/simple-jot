package tabler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

)

type exampleStruct struct{
	Name string `json:name`
	Email string `json:email`
	Birthday string `json:birthday`
	Index int `json:index`
	Confidence float64 `json:confidence`
}

func (s exampleStruct) PrepRow() []string {
	return []string{s.Name, s.Email, s.Birthday, strconv.Itoa(s.Index), strconv.FormatFloat(s.Confidence, 'f', -1, 64)}
}

const testJSON string = `[
	{
		"name": "test1",
		"email": "test@test.ca",
		"birthday": "1990-01-01",
		"index": 435,
		"confidence": 0.965
	},
	{
		"name": "test2",
		"email": "test2@test.ca",
		"birthday": "1984-05-01",
		"index": 1299560,
		"confidence": 0.05
	},
	{
		"name": "test3",
		"email": "tes3@test.ca",
		"birthday": "1995-06-21",
		"index": 256,
		"confidence": 0.452
	}
]`


func TestPrepStructRow(t *testing.T) {
	fmt.Println()
	dateTime := time.Now()
	// test strings, ints and floats
	testStruct := exampleStruct{
		Name: "Tester",
		Email: "test@test.ca",
		Birthday: dateTime.Format("2006-01-02"),
		Index: 435,
		Confidence: 0.965,
	}


	expectedStrCount := 5
	bSlice, err := json.MarshalIndent(testStruct, "", "  ")
	if err != nil {
		t.Errorf("could not marshal test struct: %v", err.Error())
		fmt.Println("test struct:")
		fmt.Println(testStruct)
	} else {
		fmt.Println("test struct:")
		fmt.Println(string(bSlice))
	}
	
	actualResult := PrepStructRow(&testStruct)
	for i, s := range actualResult {
		fmt.Printf("res %v: '%v'\n", i, s)
	}
	
	if expectedStrCount != len(actualResult) {
		t.Errorf("expected string count (%v) is different than actual (%v)", expectedStrCount , len(actualResult))
		t.Fail()
	}
	
}

func TestPrepTable(t *testing.T) {
	var structs []exampleStruct
	err := json.Unmarshal([]byte(testJSON), &structs)
	if err != nil {
		t.Fail()
		t.Fatalf("cannot unmarshal json: %v", err.Error())
		return
	}

	fmt.Println("INPUT:\n", testJSON)
	// you have to translate the structs to the interface type in order to use the prep table function
	preppedRows := make([]RowPrepper,len(structs))
	for i, s := range structs {
		preppedRows[i] = s
	}
	headers := []string{"Name", "Email", "Birthday", "Index", "Confidence"}
	dataFrame := PrepTable(preppedRows, headers)
	// add one row for headers
	expectedLen := len(preppedRows) + 1

	if len(dataFrame) != expectedLen {
		t.Errorf("resulting slice is of a different length (%v) than expected (%v)", len(dataFrame), expectedLen)
		t.Fail()
	}


	// render table
	if !t.Failed() {
		fmt.Println("\nResulting Table:")
		err = RenderTable(dataFrame[1:], headers)
		if err != nil {
			t.Fail()
			t.Fatalf("Failed to render table test. %v", err)
		}
	}
}