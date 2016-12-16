package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	project_name = flag.String("name", "", "项目名称")
	project_path = flag.String("path", "./", "项目路径")

	trees = FileTreeArray{
		&FileTree{
			Dir: "conf",
			Files: []FileTemplate{
				{Name: "config_default.json",
					Template: "{}",
				},
			},
		},
		&FileTree{
			Dir: "controller",
			Files: []FileTemplate{
				{Name: "controller_default.go",
					Template: controller_default_go,
				},
			},
		},
		&FileTree{
			Dir:   "static",
			Files: []FileTemplate{},
		},
		&FileTree{
			Dir:   "static/javascripts",
			Files: []FileTemplate{},
		},
		&FileTree{
			Dir: "test",
			Files: []FileTemplate{
				{Name: "default_test.go",
					Template: default_test_go,
				},
			},
		},
		&FileTree{
			Dir:   "views",
			Files: []FileTemplate{},
		},
		&FileTree{
			Dir: "",
			Files: []FileTemplate{
				{Name: "main.go"},
				{Name: "build_release.py"},
			},
		},
		&FileTree{
			Dir: "release",
			Files: []FileTemplate{
				{Name: "log_split.sh"},
				{Name: "start.sh"},
				{Name: "upgrade.sh"},
			},
		},
	}
)

func main() {

	flag.Parse()
	if flag.Parsed() == false {
		flag.PrintDefaults()
		return
	}
	if len(*project_name) <= 0 {
		flag.PrintDefaults()
		printlnf("[ *** ] 输入一个项目名称 \r\n")
		return
	}
	if len(*project_path) <= 0 {
		*project_path = "./"
	}
	if *project_path == "./" {
		fmt.Println("将在当前目录创建项目,目录结构如下:")
	}
	trees.Print()
	err := PlantTrees(*project_name, *project_path, trees)
	if err != nil {
		printlnf("[ *** ] error: %s", err)
		return
	}
}

func PlantTrees(project_name, project_path string, trees FileTreeArray) error {
	if trees == nil || len(trees) <= 0 {
		return nil
	}

	for _, tree := range trees {
		err := PlantTree(project_name, project_path, tree)
		if err != nil {
			return err
		}
	}
	return nil
}

func PlantTree(project_name, project_path string, tree *FileTree) error {
	base_path := filepath.Join(project_path, project_name)
	// printlnf("project base dir: %s", base_path)

	if len(tree.Dir) <= 0 {
		// printlnf("no dir created")
	}

	path := filepath.Join(base_path, tree.Dir)
	// printlnf("create project: %s", path)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	if tree.Files == nil || len(tree.Files) <= 0 {
		return nil
	}

	for _, file := range tree.Files {
		new_file_path := filepath.Join(path, file.Name)
		// printlnf("create file: %s", new_file_path)
		f, err := os.Create(new_file_path)
		if err != nil {
			return err
		}

		defer f.Close()

		var tmp string

		bs, err := LoadFile(file.Name) //优先使用外部模板,方便更新
		if err == nil {
			tmp = string(bs)
		} else if err == err_file_not_found { //外部模板不存在,使用内部定义模板
			tmp = file.Template
		} else {
			return err
		}

		t, err := template.New(file.Name).Parse(string(tmp))
		err = t.Execute(f, project_name)
		if err != nil {
			return err
		}
	}

	return nil
}

type FileTree struct {
	Dir   string
	Files []FileTemplate
}

type FileTreeArray []*FileTree

func (fta FileTreeArray) Print() {
	for _, tree := range fta {
		// printlnf("")
		var pre_space = ""
		if len(tree.Dir) > 0 {
			printlnf("|%s %s", strings.Repeat("-", 4), tree.Dir)
			pre_space = fmt.Sprintf("|%s|", strings.Repeat(" ", 5))
		} else {
			pre_space = fmt.Sprintf("|%s", strings.Repeat("-", 0))
		}
		if tree.Files == nil || len(tree.Files) <= 0 {
			continue
		}
		// printlnf("%s|", strings.Repeat(" ", 4))
		for _, tmp := range tree.Files {
			printlnf("%s%s %s", pre_space, strings.Repeat("-", 4), tmp.Name)
		}
	}
}

type FileTemplate struct {
	Name     string
	Template string
}

func printlnf(format string, paras ...interface{}) (int, error) {
	return fmt.Println(fmt.Sprintf(format, paras...))
}

var err_file_not_found = errors.New("file_not_found")

func LoadFile(name string) ([]byte, error) {
	path := filepath.Join("./files", name)
	if IsFileExist(path) == false {
		return nil, err_file_not_found
	}
	return ioutil.ReadFile(path)
}

// exists returns whether the given file or directory exists or not
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

var (
	controller_default_go = `
	package controller

	import (
		"github.com/gin-gonic/gin"
		"net/http"
	)
	
	func Hello(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	}

	`
	default_test_go = `
	package tests

	import (
		"testing"
	
		. "github.com/smartystreets/goconvey/convey"
	)
	
	func TestDefault(t *testing.T) {
		Convey("TestDefault ", t, func() {
			So(0, ShouldNotEqual, 0)
		})
	}
		
	`
)
