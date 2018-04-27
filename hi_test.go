package hi

import (
	"testing"
	"time"
)

func TestRegiste(t *testing.T) {
	hi := NewHi([]string{"127.0.0.1:2379"}, "")

	err := hi.Register("goods", "127.0.0.1:8848")
	if err != nil {
		t.Fatal(err)
	}

	println("Registe Ok !")
	time.Sleep(time.Minute * 1)
}

func TestTargetParse(t *testing.T) {
	target := parseTarget("wonamingv3://author/project/test")

	t.Logf("%+v", target)
}
