package main

import "github.com/al6nlee/obsidian/filemanager"

func main() {
	// 指定目录，建议使用绝对路径
	path := "/Users/alan/Library/Mobile Documents/com~apple~CloudDocs/Obsidian/NoteHub/OCR技术"
	// path := "./note"
	var tag []string
	tag = append(tag, "OCR")
	fileAtt := filemanager.FileAttribute{Author: "alan", Tag: tag, ISDraft: true}

	filemanager.ProcessFiles(path, &fileAtt)
}
