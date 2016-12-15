package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

var (
	project_name = flag.String("name", "", "项目名称")
	project_path = flag.String("path", "./", "项目路径")
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
		fmt.Println("将在当前目录创建项目")
	}

	// tree := &FileTree{
	// 	Dir: "conf",
	// 	Files: []FileTemplate{
	// 		{Name: "config_default.json",
	// 			Template: "{}",
	// 		},
	// 	},
	// }
	// err := PlantTree(*project_name, *project_path, tree)
	// if err != nil {
	// 	printlnf("[ *** ] error: %s", err)
	// 	return
	// }

	trees := []*FileTree{
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
	}
	err := PlantTrees(*project_name, *project_path, trees)
	if err != nil {
		printlnf("[ *** ] error: %s", err)
		return
	}
}

func PlantTrees(project_name, project_path string, trees []*FileTree) error {
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
	printlnf("project base dir: %s", base_path)

	if len(tree.Dir) <= 0 {
		printlnf("no dir created")
		return nil
	}

	path := filepath.Join(base_path, tree.Dir)
	printlnf("create project: %s", path)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	if tree.Files == nil || len(tree.Files) <= 0 {
		return nil
	}

	for _, file := range tree.Files {
		new_file_path := filepath.Join(path, file.Name)
		printlnf("create file: %s", new_file_path)
		f, err := os.Create(new_file_path)
		if err != nil {
			return err
		}

		defer f.Close()

		t, err := template.New(file.Name).Parse(file.Template)
		err = t.Execute(f, project_name)
		if err != nil {
			return err
		}
		// content := fmt.Sprintf(file.Template, project_name)
		// printlnf("file [%s] content: %s", new_file_path, content)
		// _, err = io.WriteString(f, content)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

type FileTree struct {
	Dir   string
	Files []FileTemplate
}

type FileTemplate struct {
	Name     string
	Template string
}

func printlnf(format string, paras ...interface{}) (int, error) {
	return fmt.Println(fmt.Sprintf(format, paras...))
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
)
