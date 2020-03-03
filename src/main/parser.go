package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var sourcePath string
var targetPath string
var splitSign string

func ParseXlsx() {

	fmt.Println(" ---------------------------------------------")
	fmt.Println(" |           Xlsx 2 Csv for Thoruni          |")
	fmt.Println(" ---------------------------------------------")

	files, _ := GetFilesAndDirs(sourcePath)
	for _, file := range files {
		fmt.Println(file)
		Import(file)
	}
}

func Import(inFile string) {
	// 打开文件
	xlFile, err := xlsx.OpenFile(inFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 遍历sheet页读取
	for _, sheet := range xlFile.Sheets {
		content := ""
		fileName := FindRealName(sheet.Name)
		// 如果没有获取到正确的名字
		if strings.HasPrefix(sheet.Name, "#") || fileName == "" {
			continue
		}
		fileName = strings.ReplaceAll(fileName," ","")
		line := 0
		ignore := map[int]bool{}
		//遍历行读取
		for _, row := range sheet.Rows {
			if strings.HasPrefix(row.Cells[0].String(),"#") {
				continue
			}
			//遍历每行的列读取
			col := 0
			for _, cell := range row.Cells {
				col++;
				if line == 0 {
					tmp := cell.String()
					if strings.HasPrefix(tmp,"#") {
						ignore[col] = false
					} else {
						ignore[col] = true
						content = content + cell.String() + splitSign
					}
				} else {
					if ignore[col] {
						content = content + cell.String() + splitSign
					}
				}
			}
			line++
			if len(content) != 0 {
				content = content + "\n"
			}
		}
		WriteWithIoutil(target+"/"+fileName+".csv", strings.TrimRight(content, splitSign))

	}
	fmt.Println("\n\nimport success")
}

//获取指定目录下的所有文件和目录
func GetFilesAndDirs(dirPth string) (files []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if !fi.IsDir() { // 目录, 递归遍历
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".xlsx")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	return files, nil
}

func WriteWithIoutil(name, content string) {
	data := []byte(content)
	if ioutil.WriteFile(name, data, 0644) == nil {
		fmt.Println("写入文件成功:", name)
	}
}

func FindRealName(oldName string) string {
	findReg := regexp.MustCompile(`(?s)\((.*)\)`)
	reg := findReg.FindAllStringSubmatch(oldName,-1)
	if len(reg)> 0 {
		return reg[0][1]
	}
	return ""
}
