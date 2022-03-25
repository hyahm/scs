package server

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/status"
)

// 比较大于等于
func ge(one, two, sep string) bool {
	if sep == "" {
		return one >= two
	}
	if one == two {
		return true
	}
	l1 := strings.Split(one, sep)
	l2 := strings.Split(two, sep)
	length := len(l1)
	if len(l1) >= len(l2) {
		length = len(l2)
	}
	for i := 0; i < length; i++ {
		if len(l1[i]) != len(l2[i]) {
			return len(l1[i]) > len(l2[i])
		}
		if l1[i] > l2[i] {
			return true
		}
		if l1[i] < l2[i] {
			return false
		}
	}
	if len(l1) != len(l2) {
		return len(one) > len(l2)
	}
	return true
}

// 比较大于
func gt(one, two, sep string) bool {
	if sep == "" {
		return one > two
	}
	if one == two {
		return false
	}
	l1 := strings.Split(one, sep)
	l2 := strings.Split(two, sep)
	length := len(l1)
	if len(l1) >= len(l2) {
		length = len(l2)
	}
	for i := 0; i < length; i++ {
		if len(l1[i]) != len(l2[i]) {
			return len(l1[i]) > len(l2[i])
		}

		if l1[i] > l2[i] {
			return true
		}
		if l1[i] < l2[i] {
			return false
		}

	}
	if len(l1) != len(l2) {
		return len(one) > len(l2)
	}
	return true
}

// 比较小于等于
func le(one, two, sep string) bool {
	if sep == "" {
		return one <= two
	}
	if one == two {
		return true
	}
	l1 := strings.Split(one, sep)
	l2 := strings.Split(two, sep)
	length := len(l1)
	if len(l1) >= len(l2) {
		length = len(l2)
	}
	for i := 0; i < length; i++ {
		if len(l1[i]) != len(l2[i]) {
			return len(l1[i]) < len(l2[i])
		}
		if l1[i] < l2[i] {
			return true
		}
		if l1[i] > l2[i] {
			return false
		}
	}
	if len(l1) != len(l2) {
		return len(one) < len(l2)
	}
	return true
}

// 比较小于
func lt(one, two, sep string) bool {
	if sep == "" {
		return one < two
	}
	if one == two {
		return false
	}
	l1 := strings.Split(one, sep)
	l2 := strings.Split(two, sep)
	length := len(l1)
	if len(l1) >= len(l2) {
		length = len(l2)
	}
	for i := 0; i < length; i++ {
		if len(l1[i]) != len(l2[i]) {
			return len(l1[i]) < len(l2[i])
		}
		if l1[i] < l2[i] {
			return true
		}
		if l1[i] > l2[i] {
			return false
		}
	}
	if len(l1) != len(l2) {
		return len(one) < len(l2)
	}
	return true
}

func (svc *Server) Install() (err error) {

	svc.Status.Status = status.INSTALL
	for _, v := range svc.PreStart {
		if strings.Trim(v.Path, " ") == "" &&
			strings.Trim(v.Command, " ") == "" &&
			strings.Trim(v.ExecCommand, " ") == "" {
			continue
		}
		if strings.Trim(v.Path, " ") != "" {
			v.Path = internal.Format(v.Path, svc.Env)
			golog.Info("check path: ", v.Path)
			_, err := os.Stat(filepath.Join(svc.Dir, v.Path))
			if err == nil || !os.IsNotExist(err) {
				continue
			}

		}
		if strings.Trim(v.Command, " ") != "" {
			v.Command = internal.Format(v.Command, svc.Env)
			golog.Info("check command: ", v.Command)
			_, err := exec.LookPath(v.Command)
			if err == nil {
				continue
			}
		}
		if strings.Trim(v.ExecCommand, " ") != "" {
			v.ExecCommand = internal.Format(v.ExecCommand, svc.Env)
			golog.Info("exec command: ", v.ExecCommand)
			result, err := exec.Command(v.ExecCommand).Output()
			if err == nil {
				// 查看是否有结果比较， 如果没有就跳过
				// EQ string `yaml:"eq,omitempty" json:"eq,omitempty"`
				// NE string `yaml:"ne,omitempty" json:"ne,omitempty"`
				// GT string `yaml:"gt,omitempty" json:"gt,omitempty"`
				// LT string `yaml:"lt,omitempty" json:"lt,omitempty"`
				// GE string `yaml:"ge,omitempty" json:"ge,omitempty"`
				// LE string `yaml:"le,omitempty" json:"le,omitempty"`
				golog.Info("confidition")
				if v.EQ != "" && string(bytes.TrimSpace(result)) != v.EQ {
					goto install
				}
				if v.NE != "" && string(bytes.TrimSpace(result)) == v.EQ {
					goto install
				}
				if v.GT != "" && gt(string(bytes.TrimSpace(result)), v.EQ, v.Separation) {
					goto install
				}
				if v.LT != "" && lt(string(bytes.TrimSpace(result)), v.EQ, v.Separation) {
					goto install
				}
				if v.GE != "" && ge(string(bytes.TrimSpace(result)), v.EQ, v.Separation) {
					goto install
				}
				if v.LE != "" && le(string(bytes.TrimSpace(result)), v.EQ, v.Separation) {
					goto install
				}
				continue
			}
		}
	install:
		err = nil
		v.Install = internal.Format(v.Install, svc.Env)
		if v.Install != "" {
			golog.Info("install: ", v.Install)
			if err := svc.shellWithOutDir(v.Install); err != nil {
				golog.Error(err)
				return err
			}
		}
		if v.Template != "" {
			golog.Info("template")
			v.Template = filepath.Clean(v.Template)
			v.Path = filepath.Clean(v.Path)
			// 将文件内容模板重新写入
			b, err := os.ReadFile(filepath.Join(svc.Dir, v.Template))
			if err != nil {
				golog.Error(err)
				return err
			}
			fc := string(b)
			fc = internal.Format(fc, svc.Env)
			err = os.WriteFile(filepath.Join(svc.Dir, v.Path), []byte(fc), 0644)
			if err != nil {
				golog.Error(err)
				return err
			}

		}
	}
	return nil
}
