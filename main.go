package main

import "github.com/al6nlee/obsidian/filemanager"

func main() {
	// 指定目录，建议使用绝对路径
	path := "./note"
	var tag []string
	tag = append(tag, "示例标签")

	fileAtt := filemanager.FileAttribute{Author: "alan", Tag: tag}
	filemanager.ProcessFiles(path, &fileAtt)
}
