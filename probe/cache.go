package probe

import "scs/alert"

var cache []*alert.AlertInfo

func init() {
	cache = make([]*alert.AlertInfo, 4)
}
