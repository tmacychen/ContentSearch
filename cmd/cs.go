package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.5"
var recursive bool
var content, path string

var csCmd = &cobra.Command{
	Use:   "cs [选项] [搜索的内容] [需要搜索的文件夹路径]",
	Short: `cs is a content search tool for doc files`,
	Long: `Content Search是一个内容搜索工具，能在给定的目录下解析doc格式文件内容，
	并检索想要的内容，并显示出来。
	Content Search 工具在文本文件的搜索时可以被grep替代，与可以被shell脚本完全替代。
	实现此工具的初衷旨在联系golang的命令行编程与图形编程，本工具属于实验产品，
	没有经过大量测试，在生产环境中请慎重使用`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		//main funciton is here
		for _, i := range args {
			fmt.Printf("args :%v\n", i)
		}
		if err := checkArgs(args); err != nil {
			fmt.Printf("Error : %v\n\n", err)
			cmd.Usage()
			os.Exit(1)
		}

	},
}

func init() {
	csCmd.PersistentFlags().BoolVarP(&recursive, "recursive for directory", "r", false, "-r [dir]")
}

//Execute : execute the csCmd
func Execute() {
	if err := csCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkArgs(args []string) (err error) {
	if len(args) < 2 {
		return errors.New("需要两个参数，请检查一下命令是否正确")
	}
	content = args[0]
	path = args[1]
	if path == "" {
		return errors.New("需要提供一个路径")
	}
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			return fmt.Errorf("当前文件:%v 未找到", path)
		} else {
			// other error
			return err
		}
	}
	d := fi.Mode().IsDir()
	if recursive {
		if !d {
			return errors.New("需要一个目录")
		}
	} else {
		if d {
			return errors.New("需要一个文件")
		}
	}
	return nil
}
