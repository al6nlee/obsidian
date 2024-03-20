package main

import "github.com/al6nlee/obsidian/filemanager"

func main() {
	// 指定目录
	dir := "./note"
	filemanager.ProcessFiles(dir)
}
