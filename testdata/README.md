# MarkNum - 自动生成 markdown 标题序号

自动添加/更新 markdown 标题序号，支持多级标题。

## 安装

```bash

```

## 使用

* -file: 指定文件
* -dir: 指定目录(和 -f 二选一)
* -cover: 是否覆盖原文件, 默认为 false, 新建 $filename.marknum.md 文件
* -min: 最小标题级数, 范围[min,max), 默认为 2; 生成二级, 三级标题的序号(`## 1. 标题` 和 `### 1.1. 标题`)
* -max: 最大标题级数, 范围[min,max), 默认为 4; 生成二级, 三级标题的序号(`## 1. 标题` 和 `### 1.1. 标题`)

