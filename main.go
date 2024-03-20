package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
)

func fmtStr(file fileAttribute) string {
	// 创建一个字符串切片来存储tag
	tags := make([]string, len(file.tag))
	for i, tag := range file.tag {
		tags[i] = fmt.Sprintf("  - %s", tag)
	}

	// 使用strings.Join将tag连接成一个字符串
	tagStr := strings.Join(tags, "\n")

	// 构建格式化字符串
	formatStr := "Title: %s\n" +
		"tags: \n%s\n" +
		"CreateDate: %s\n" +
		"ModDate: %s\n" +
		"Draft: false\n---\n\n"

	// 使用格式化字符串和参数生成最终结果
	return fmt.Sprintf(formatStr, file.fileName, tagStr, file.createTime, file.modTime)
}

func addAttribute(path string, file fileAttribute) error {
	// 读取原文件内容
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// 将文件内容转换为字符串
	contentStr := string(content)

	// 查找第一个 "---\n"
	startIndex := strings.Index(contentStr, "---\n")
	// fmtStr := fmt.Sprintf("title: %s\n---\n", file.fileName)
	if startIndex != -1 {
		// 查找第二个 "---\n"
		endIndex := strings.Index(contentStr[startIndex+4:], "---\n")
		if endIndex != -1 {
			// 覆盖原有的标题行，并在其后添加新内容
			contentStr = contentStr[:startIndex+4] + fmtStr(file) + contentStr[startIndex+endIndex+8:]
		}
	} else {
		// 如果文件内容不是以 "---\n" 开头，则在开头添加新内容
		contentStr = "---\n" + fmtStr(file) + contentStr
	}

	// 将新内容写入原文件
	err = os.WriteFile(path, []byte(contentStr), 0644)
	if err != nil {
		return err
	}

	return nil
}

type fileAttribute struct {
	tag        [2]string
	dir        string // 这个用不到
	fileName   string
	createTime time.Time
	modTime    time.Time
	size       int64
	mode       string
}

func main() {
	// 指定目录
	dir := "./note"
	re := regexp.MustCompile(`【(.*?)】`)
	// 遍历目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 忽略目录
		if info.IsDir() || filepath.Ext(info.Name()) != ".md" {
			return nil
		}

		var tagArr [2]string
		tagArr[0] = filepath.Base(filepath.Dir(path))
		matches := re.FindStringSubmatch(info.Name())
		if len(matches) > 1 {
			tagArr[1] = matches[1]
		}
		file := fileAttribute{tag: tagArr, dir: filepath.Dir(path), fileName: info.Name(),
			createTime: time.Unix(int64(info.Sys().(*syscall.Stat_t).Birthtimespec.Sec),
				int64(info.Sys().(*syscall.Stat_t).Birthtimespec.Nsec)),
			modTime: info.ModTime(), size: info.Size(), mode: info.Mode().String()}
		addAttribute(path, file)
		return nil
	})

	if err != nil {
		fmt.Println("遍历目录时发生错误:", err)
	}
}
