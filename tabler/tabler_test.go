package tabler

import (
	"fmt"
	"testing"
	"time"
)



func TestPrepStructRow(t *testing.T) {
	dateTime := time.Now()
	example := struct{
		name string
		email string
		birthday string
	}{
		name: "Tester",
		email: "test@test.ca",
		birthday: dateTime.Format("2006-01-02"),
	}
	expectedStrCount := 3
	fmt.Println("example struct:")
	fmt.Println(example)
	actualResult := PrepStructRow(&example)
	for i, s := range actualResult {
		fmt.Printf("res %v: '%v'", i, s)
	}
	
	if expectedStrCount != len(actualResult) {
		t.Errorf("expected string count (%v) is different than actual (%v)", expectedStrCount , len(actualResult))
	}
	
}