package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var file = flag.String("file", "", "指定文件")
var dir = flag.String("dir", "", "深度遍历目录下所有md文件(和 -f 二选一)")
var cover = flag.Bool("cover", false, "是否覆盖原文件, 默认为 false, 新建 $filename.marknum.md 文件")
var min = flag.Int("min", 2, "最小标题级数, 范围[min,max), 默认为 2; 生成二级, 三级标题的序号(## 1. 标题 和 ### 1.1. 标题)")
var max = flag.Int("max", 4, "最大标题级数, 范围[min,max), 默认为 4; 生成二级, 三级标题的序号(## 1. 标题 和 ### 1.1. 标题)")

func main() {
	flag.Parse()

	if *file == "" && *dir == "" {
		fmt.Printf("Help:\n %s -h  \n", os.Args[0])
		fmt.Printf("Example: \n marknum -dir ./ -cover ture \n")
		os.Exit(1)
	}

	if *file != "" {
		oneFile(*file)
	}

	if *dir != "" {
		files := mdPaths(*dir)
		for _, filename := range files {
			oneFile(filename)
		}
	}

}

// 通过目录获取所有的 md 文件
func mdPaths(dir string) []string {
	var files []string
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if strings.HasSuffix(d.Name(), ".md") {
				files = append(files, path)
			}
		}
		return nil
	})
	return files
}

func oneFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("os.Open(filename) err:%v \n", err)
		os.Exit(1)
	}

	s, err := sectionNumber(f)
	if err != nil {
		fmt.Printf("sectionNumber(f) err:%v \n", err)
		os.Exit(1)
	}

	var output string
	if !*cover {
		output = filename + ".marknum.md"
	} else {
		output = filename
	}

	// 写入文件或者覆盖文件
	err = os.WriteFile(output, []byte(s), 0644)
	if err != nil {
		fmt.Printf("os.WriteFile(output, []byte(s), 0644) err:%v \n", err)
		os.Exit(1)
	}
	fmt.Printf("[成功] 输出文件: %s \n", output)
}

// add/update section numbers
// 一行行读取; 识别代码块; 识别标题 -> 删除标题序号 -> 添加标题序号 -> 写入
func sectionNumber(in io.Reader) (string, error) {
	r := bufio.NewReader(in)
	buf := bytes.Buffer{}

	// 标题序号
	sectionNumbers := make([]int, 6)
	// 是否在代码块中
	inCodeBlock := false

	// 一行行读取
	var finish bool
	for !finish {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				finish = true
			} else {
				return "", err
			}
		}

		// 识别代码块
		if isCodeBlock(line) {
			inCodeBlock = !inCodeBlock
		}

		// 不在代码块中
		if !inCodeBlock {
			// 是标题
			if level := headerLevel(line); level != 0 {
				// 更新 sectionNumbers
				updateSectionNumbers(sectionNumbers, level)

				// 删除标题序号
				line = delSectionNumber(line)

				// 添加标题序号
				line = addSectionNumber(line, sectionNumbers, level)
			}
		}

		buf.WriteString(line)

	}

	return buf.String(), nil
}

// 更新 sectionNumbers
// 如果 level = 1, 一级标题, 可以由一开始 0 到 1; => 添加 sectionNumbers[level-1]++ 即可
// 如果 level = 1, 一级标题, 可以由一开始 2 到 1; => 添加 sectionNumbers[level-1]++ 即可, 并且清理后面的, 将后面的设置为 0
// 所以每次更新新只需要清理后面的就行了
func updateSectionNumbers(sectionNumbers []int, level int) {
	sectionNumbers[level-1]++

	for level < len(sectionNumbers) {
		sectionNumbers[level] = 0
		level++
	}
}

// 添加标题序号
// 比如输入 "## 标题" 返回 "## 1. 标题"
func addSectionNumber(line string, sectionNumbers []int, level int) string {
	s := sectionNumberStr(sectionNumbers[:level])
	if s != "" {
		return fmt.Sprintf("%s %s\n", strings.TrimSpace(s), strings.TrimSpace(line))
	}
	return ""
}

// 删除标题的header和序号
// 比如输入 "## 1. 标题" 返回 "标题"
// 比如输入 "## 1 标题" 返回 "标题"
// 比如输入 "## 1.1 标题" 返回 "标题"
// 比如输入 "## 1.1. 标题" 返回 "标题"
// 比如输入 "## 1.1.1 标题" 返回 "标题"
// 比如输入 "## 1.1.1. 标题" 返回 "标题"
var delSectionNumberRe = regexp.MustCompile(`(\s*#+\s+)([\d\.]*)(\s*)`)

func delSectionNumber(line string) string {
	return delSectionNumberRe.ReplaceAllString(line, "")
}

// 获取标题级别
func headerLevel(line string) int {
	level := 0
	for _, ch := range line {
		if ch == '#' {
			level++
		} else {
			break
		}
	}
	return level
}

func isCodeBlock(line string) bool {
	return strings.HasPrefix(line, "```")
}

func sectionNumberStr(sectionNumbers []int) string {
	var buf bytes.Buffer
	// 添加 #
	for i := 0; i < len(sectionNumbers); i++ {
		buf.WriteString("#")
	}

	// 空格
	buf.WriteString(" ")

	// 序号
	for i, n := range sectionNumbers {
		// 例如一级标题不需要序号
		level := i + 1
		if level < *min {
			continue
		}

		// 例如只需要二级三级标题
		if level >= *max {
			break
		}

		buf.WriteString(fmt.Sprintf("%d.", n))
	}
	return buf.String()
}
