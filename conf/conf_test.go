package conf

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf, err := LoadConfig("../conf.json")
	if err != nil {
		fmt.Println(err)
		t.Error("Failed to read ICE servers config")
	}
	fmt.Println(conf)
}
