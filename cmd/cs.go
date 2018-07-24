package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/text"
	"github.com/tmacychen/ContentSearch/filePool"
	sech "github.com/tmacychen/ContentSearch/search"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var version = "0.5"
var isRecursive bool
var isDebug bool
var isSingleFile bool
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
		if err := checkArgs(args); err != nil {
			log.Errorf("Error : %v\n\n", err)
			cmd.Usage()
			os.Exit(1)
		}
		log.Debugf("path :%v\n", path)
		log.Debugf("content :%v\n", content)

		mainWork()
	},
}

func init() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.InfoLevel)
	csCmd.PersistentFlags().BoolVarP(&isRecursive, "recursive", "r", false, "recursive for directory")
	csCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "debug mode")
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
	isSingleFile = fi.Mode().IsRegular()
	if isRecursive {
		if isSingleFile {
			return errors.New("需要一个目录")
		}
	} else {
		if !isSingleFile {
			return errors.New("需要一个文件")
		}
	}
	if isDebug {
		log.Info("set debug")
		log.SetHandler(text.Default)
		log.SetLevel(log.DebugLevel)
	}
	return nil
}
func readPath(p string, fs *filePool.FileSet) {
	d, err := ioutil.ReadDir(p)
	if err != nil {
		log.Errorf("get file err :%v\n", err)
	}
	for _, info := range d {
		n := info.Name()
		subp := p + "/" + n
		if info.IsDir() {
			readPath(subp, fs)
			continue
		} else {
			if isReadableFile(subp) {
				fs.Add(subp)
				log.Debug("add file :" + subp)
			}
		}
	}
}

func isReadableFile(name string) bool {
	if strings.HasSuffix(name, ".doc") {
		return true
	}
	if strings.HasSuffix(name, ".docx") {
		return true
	}
	return false

}

func mainWork() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fs := filePool.FileSetNew()
	task := sech.TaskInit(sech.Key(content), runtime.NumCPU())
	//todo 同时启动fs获取路径下的所有文件，task启动，初始化worker，等待
	// 任务开始工作
	//file pool subroutine
	log.Debugf("sigle file :%v\n", isSingleFile)
	if isSingleFile {
		fs.Add(path)
		log.Debug("add file :" + path)
		task.SetEnd(true)
		task.Exec(fs)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		readPath(path, fs)
		task.SetEnd(true)
		log.Debugf("set end success")
		defer wg.Done()
	}()
	wg.Add(1)
	go func() {
		task.Exec(fs)
		defer wg.Done()
	}()

	wg.Wait()
	log.Infof("length of fs:%v\n", fs.Length())
}
