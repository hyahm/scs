package script

import (
	"encoding/json"
)

func GetEnv(name string) ([]byte, bool) {
	for pname := range ss.Infos {
		if _, ok := ss.Infos[pname][name]; ok {

			send, _ := json.Marshal(ss.Infos[pname][name].Env)
			return send, true
		}
	}
	return nil, false
}
