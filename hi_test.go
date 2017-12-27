package hi

import (
	"testing"
	"time"
)

func TestRegiste(t *testing.T) {
	hi := NewHi([]string{"localhost:2379"}, "goushuyun")

	err := hi.Register("order", "127.0.0.1:8848")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("registe OK ...")
	time.Sleep(time.Second * 30)

	err = hi.Unregister("order")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Hour * 4)
}

func TestTargetParse(t *testing.T) {
	target := parseTarget("hi://kai/serverA")

	t.Logf("%+v", target)
}
