<p align="center">
<img src="images/logo.png" width="200px"/>
<br>
<p align="center">
 <img src="https://img.shields.io/github/stars/jarvanstack/marknum" />
 <img src="https://img.shields.io/github/issues/jarvanstack/marknum" />
 <img src="https://img.shields.io/github/forks/jarvanstack/marknum" />
</p>
</p>

#  MarkNum - 自动生成 markdown 标题序号

自动添加/更新 markdown 标题序号，支持多级标题。

## 1. 示例代码

输入 

```bash
$ marknum -file test.md
[成功] 输出文件: test.md.marknum.md 
```

原始文件 test.md

```markdown
# 一级标题

## 二级标题

### 三级标题

## 二级标题

### 三级标题
```

输出文件 test.md.marknum.md

```bash
#  一级标题

## 1. 二级标题

### 1.1. 三级标题

## 2. 二级标题
```

## 2. 安装

### 2.1. Go语言安装

```bash
go install github.com/jarvanstack/marknum@latest
```

### 2.2. 可执行文件

下载可执行文件: https://github.com/jarvanstack/marknum/releases

## 3. 使用

```bash
$ marknum -h
Usage of marknum:
  -cover
        是否覆盖原文件, 默认为 false, 新建 $filename.marknum.md 文件
  -dir string
        指定目录(和 -f 二选一)
  -file string
        指定文件
  -max int
        最大标题级数, 范围[min,max), 默认为 4; 生成二级, 三级标题的序号(## 1. 标题 和 ### 1.1. 标题) (default 4)
  -min int
        最小标题级数, 范围[min,max), 默认为 2; 生成二级, 三级标题的序号(## 1. 标题 和 ### 1.1. 标题) (default 2)
```

## 4. 常用命令

```bash
# 将当前目录下所有 markdown 文件添加/更新序号, 覆盖源文件
marknum -dir ./ -cover ture
```


