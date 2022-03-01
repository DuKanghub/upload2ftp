/*
Copyright © 2022 DuKang <dukang@dukanghub.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/DuKanghub/upload2ftp/pkg"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "上传文件",
	Long:  `上传文件到FTP服务器`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName := time.Now().Local().Format("20060102-150405")
		files := args
		var err error
		if len(files) == 0 && autoFind != "" {
			files, err = WalkDir(autoFind, ".bak")
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if len(files) == 0 {
			fmt.Println("没有找到需要上传的文件")
			return
		}
		fileName = fileName + ".zip"
		err = pkg.ZipFiles(files, fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("打包文件成功", fileName)
		ftpConfig := pkg.FtpConfig{
			Host:     ftpHost,
			Port:     ftpPort,
			User:     ftpUser,
			Password: ftpPass,
		}
		ftpCli := pkg.NewFtpClient(ftpConfig)
		err = ftpCli.UploadFile(fileName, ftpDir)
		if err != nil {
			panic(err)
		}
		fmt.Println("上传文件成功", fileName)
		_ = os.Remove(fileName)

	},
}

func WalkDir(dirPth, suffix string) (files []string, err error) {
	suffix = strings.ToUpper(suffix)                                                     //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		//if err != nil { //忽略错误
		// return err
		//}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) && time.Now().Sub(fi.ModTime()) <= 24*time.Hour {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
