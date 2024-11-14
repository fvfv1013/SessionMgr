package conf

import (
	"encoding/json"
	"fmt"
	"sessionmgr/dbg"
	"github.com/pion/webrtc/v4"
	"os"
)

type Configuration struct {
	WebrtcConf       webrtc.Configuration `json:"WebRTC"`
	CacheSize        int                  `json:"CacheSize"`
	SessionLifeCycle int                  `json:"SessionLifeCycle"`
}

func LoadConfig(ConfPath string) (*Configuration, error) {
	data, err := os.ReadFile("conf.json")
	if err != nil {
		dbg.Println(dbg.CONFIG, err)
		return nil, err
	}

	config := &Configuration{}
	err = json.Unmarshal(data, config)
	if err != nil {
		fmt.Println(dbg.CONFIG, err)
		return nil, err
	}

	return config, nil
}
