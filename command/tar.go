package command

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func compress(dirname string) error {
	os.RemoveAll(outfile)
	f, err := os.OpenFile(outfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	wf := tar.NewWriter(f)
	defer f.Close()
	defer wf.Close()
	return write(dirname, wf)

}

func write(dirname string, w *tar.Writer) error {
	fs, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}
	for _, fi := range fs {
		h, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}
		// 写信息头
		if err := w.WriteHeader(h); err != nil {
			return err
		}
		if !fi.IsDir() {
			f, err := os.Open(filepath.Join(dirname, fi.Name()))
			if err != nil {
				return err
			}

			_, err = io.Copy(w, f)
			if err != nil {
				return err
			}

			continue
		}
		write(filepath.Join(dirname, fi.Name()), w)

	}
	return nil
}

var outfile string
var TarCmd = &cobra.Command{
	Use:   "tar",
	Short: "archive tar package",
	Long:  `archive tar  package`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// 选择url和os , 是否覆盖
		if repoUrl == "" {
			repoUrl = "http://localhost:8080/upload"
		}
		if err := compress(filepath.Clean(args[0])); err != nil {
			fmt.Println(err)
			return
		}

	},
}

func init() {
	TarCmd.Flags().StringVarP(&outfile, "output", "o", "", "output file ***.tar")
	rootCmd.AddCommand(TarCmd)
}
