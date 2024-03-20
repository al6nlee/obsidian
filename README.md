## 项目说明

### 产生背景

> 市面上obsidian插件功能有限，比如一键更新指定目录下文件属性

### 解决之道

> 本项目的目的就是为了更好的使用obsidian

## 依赖环境

> go

## 初始化项目(非本项目下去更新指定文件下文件属性)

```
go mod init git.example.com
```

## 下载依赖(非本项目下去更新指定文件下文件属性)

```
go get github.com/al6nlee/obsidian
```

## 下载指定版本(非本项目下去更新指定文件下文件属性)

```
go get github.com/al6nlee/obsidian@v0.0.2
```

## 引用代码

```
package main

import "github.com/al6nlee/obsidian/filemanager"

func main() {
	// 指定目录，建议使用绝对路径
	path := "/Users/alan/Library/Mobile Documents/com~apple~CloudDocs/Obsidian"
	var tag []string
	// tag = append(tag, "示例标签")

	fileAtt := filemanager.FileAttribute{Author: "alan", Tag: tag}
	filemanager.ProcessFiles(path, &fileAtt)
}
```