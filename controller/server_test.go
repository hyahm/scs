package controller

import (
	"testing"
)

type User struct {
	Name string
	Age  int
}

var aaa map[string]*User

func init() {
	aaa = make(map[string]*User)

}
func TestStroe(t *testing.T) {
	user := getuser("555")
	user.Age = 70

	t.Log(aaa["555"].Age)
}

func getuser(name string) *User {
	user := &User{
		Name: "cander",
		Age:  60,
	}

	aaa["555"] = user
	return user
}
