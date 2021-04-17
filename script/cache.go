package script

var cache []*AlertInfo

func init() {
	cache = make([]*AlertInfo, 4)
}
