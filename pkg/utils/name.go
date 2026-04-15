package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hyahm/scs/pkg"
)

func MakeSubName(pname string, index int) string {
	return fmt.Sprintf("%s_%d", pname, index)
}

func ParseSubName(subname string) (pname string, index int, err error) {
	parts := strings.Split(subname, "_")
	if len(parts) != 2 {
		return "", 0, pkg.ErrInvalidSubName
	}
	index, err = strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, pkg.ErrInvalidSubName
	}
	return parts[0], index, nil
}
