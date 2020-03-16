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
		fileName = strings.ReplaceAll(fileName, " ", "")
		line := 0
		ignore := map[int]bool{}
		//遍历行读取
		for _, row := range sheet.Rows {
			if len(row.Cells) == 0 {
				continue
			}
			if strings.HasPrefix(row.Cells[0].String(), "#") {
				continue
			}
			//遍历每行的列读取
			col := 0
			linecontent := ""
			for _, cell := range row.Cells {
				col++
				tmp := cell.String()
				tmp = strings.ReplaceAll(tmp, ",", "/")
				if line == 0 {
					if strings.HasPrefix(tmp, "#") || len(strings.TrimSpace(tmp)) == 0 {
						ignore[col] = false
					} else {
						ignore[col] = true
						linecontent = linecontent + tmp + splitSign
					}
				} else {
					if ignore[col] {
						linecontent = linecontent + tmp + splitSign
					}
				}
			}
			line++
			linecontent = strings.ReplaceAll(linecontent, "\n", "")
			linecontent = strings.ReplaceAll(linecontent, "\r", "")
			linecontent = strings.ReplaceAll(linecontent, "\x0a", "")
			linecontent = strings.ReplaceAll(linecontent, " ", "")
			length := len(strings.ReplaceAll(linecontent, ",", ""))
			//fmt.Println(">>cnt:", linecontent)
			//fmt.Println(">>len:", length)
			if len(linecontent) != 0 && length != 0 {
				content = content + strings.TrimRight(linecontent, splitSign) + "\n"
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
	reg := findReg.FindAllStringSubmatch(oldName, -1)
	if len(reg) > 0 {
		return reg[0][1]
	}
	return ""
}
