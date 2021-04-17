package scs

var cache []*AlertInfo

func init() {
	cache = make([]*AlertInfo, 4)
}
