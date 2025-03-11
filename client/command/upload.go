package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func upload(filename string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//要上传的文件

	//创建第一个需要上传的文件,filepath.Base获取文件的名称
	fileWriter1, _ := bodyWriter.CreateFormFile("file", filepath.Base(filename))
	//打开文件
	fd1, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd1.Close()
	//把第一个文件流写入到缓冲区里去
	_, err = io.Copy(fileWriter1, fd1)
	if err != nil {
		return err
	}
	f1, err := bodyWriter.CreateFormField("overwrite")
	if err != nil {
		return err
	}

	io.Copy(f1, bytes.NewReader([]byte(fmt.Sprintf("%t", overwrite))))

	osversion, err := bodyWriter.CreateFormField("osversion")
	if err != nil {
		return err
	}
	send, _ := json.Marshal(osDir)
	io.Copy(osversion, bytes.NewReader(send))
	//获取请求Content-Type类型,后面有用
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	//创建一个http客户端请求对象
	client := &http.Client{}
	//请求url
	//创建一个post请求
	req, _ := http.NewRequest("POST", repoUrl, nil)
	//设置请求头
	//这里的Content-Type值就是上面contentType的值
	req.Header.Set("Content-Type", contentType)
	//转换类型
	req.Body = io.NopCloser(bodyBuf)
	//发送数据
	data, err := client.Do(req)
	if err != nil {
		return err
	}
	//读取请求返回的数据
	bytes, _ := io.ReadAll(data.Body)
	defer data.Body.Close()
	//打印数据
	fmt.Println(string(bytes))
	return nil
}

var UploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload package(future)",
	Long:  `upload package(future)`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(osDir) == 0 {
			fmt.Println("must be set os,eg: -o centos")
		}
		// 选择url和os , 是否覆盖
		if repoUrl == "" {
			repoUrl = "http://localhost:8080/upload"
		}

		upload(args[0])

	},
}

func init() {
	UploadCmd.Flags().BoolVarP(&overwrite, "overwrite", "w", false, "overwrite package, default: false")
	UploadCmd.Flags().StringSliceVarP(&osDir, "os", "o", make([]string, 0), "os name, like centos, ubuntu, mac")
	UploadCmd.Flags().StringVarP(&repoUrl, "url", "l", "http://localhost:8080/upload", "upload url, default: http://localhost:8080/upload")
	UploadCmd.Flags().StringVarP(&username, "username", "u", "", "upload username")
	UploadCmd.Flags().StringVarP(&password, "password", "p", "", "upload password")
	rootCmd.AddCommand(UploadCmd)
}
