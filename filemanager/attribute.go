package filemanager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
)

type FileAttribute struct {
	Tag        []string
	Dir        string // 这个用不到
	FileName   string
	CreateTime time.Time
	ModTime    time.Time
	Size       int64
	Mode       string
	Author     string
}

func fmtStr(file FileAttribute) string {
	// 创建一个字符串切片来存储非空的tag
	var filteredTags []string
	existingTags := make(map[string]bool) // 用于存储已存在的标签
	for _, tag := range file.Tag {
		// 将标签转换为小写
		lowercaseTag := strings.ToLower(tag)
		if tag != "" && !existingTags[lowercaseTag] {
			existingTags[lowercaseTag] = true
			filteredTags = append(filteredTags, fmt.Sprintf("  - %s", tag))
		}
	}
	// 使用strings.Join将非空tag连接成一个字符串
	tagStr := strings.Join(filteredTags, "\n")

	// 构建格式化字符串
	formatStr := "Title: %s\n" +
		"tags: \n%s\n" +
		"CreateDate: %s\n" +
		"ModDate: %s\n" +
		"Draft: false\n" +
		"Author: %s\n---\n"

	// 使用格式化字符串和参数生成最终结果
	return fmt.Sprintf(formatStr, file.FileName, tagStr, file.CreateTime, file.ModTime, file.Author)
}

func compareModTimeInFile(content string, file FileAttribute) bool {

	re := regexp.MustCompile(`ModDate: (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+.*?)\n`)
	match := re.FindStringSubmatch(content)
	if len(match) < 2 {
		// 如果没有匹配到时间数据，说明文件是首次添加属性，也返回true
		return true
	}
	// 解析时间数据为time.Time类型
	modDateStr := match[1]
	modDateTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", modDateStr)
	if err != nil {
		// 如果时间数据解析失败，返回false
		return false
	}
	// 将file.ModTime转换为UTC时间
	fileModTimeUTC := file.ModTime
	// 比较时间差
	timeDiff := fileModTimeUTC.Sub(modDateTime)
	// 如果时间差小于等于3秒，返回true，否则返回false
	return timeDiff >= 10*time.Second
}

func AddAttribute(path string, file FileAttribute) (err error) {
	// 读取原文件内容
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// 将文件内容转换为字符串
	contentStr := string(content)

	if !compareModTimeInFile(contentStr, file) {
		return errors.New("compareModTimeInFile fail")
	}

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

func ProcessFiles(path string, fileAtt *FileAttribute) error {
	re := regexp.MustCompile(`【(.*?)】`)
	// 遍历目录
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// 每一个文件都实例化出一个具体的fileAtt
		_fileAtt := *fileAtt
		if err != nil {
			return err
		}

		// 忽略目录
		if info.IsDir() || filepath.Ext(info.Name()) != ".md" {
			return nil
		}

		matches := re.FindStringSubmatch(info.Name())
		if len(matches) > 1 {
			_fileAtt.Tag = append(_fileAtt.Tag, matches[1])
		}
		_fileAtt.Tag = append(_fileAtt.Tag, filepath.Base(filepath.Dir(path)))
		_fileAtt.Dir = filepath.Dir(path)
		_fileAtt.FileName = info.Name()
		_fileAtt.CreateTime = time.Unix(int64(info.Sys().(*syscall.Stat_t).Birthtimespec.Sec),
			int64(info.Sys().(*syscall.Stat_t).Birthtimespec.Nsec))
		_fileAtt.ModTime = info.ModTime()
		_fileAtt.Size = info.Size()
		_fileAtt.Mode = info.Mode().String()
		err = AddAttribute(path, _fileAtt)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("遍历目录时发生错误: %v", err)
	}
	return nil
}
