---
uri: test
title: "Markdown 语法测试"
createdAt: 1717521516
category: 探索日志
tags:
  - python
  - Go
desc: '
Markdown是一种轻量级标记语言，设计目的是以简洁的语法使文档易于书写和阅读，同时可以转换成丰富的HTML格式。


该帖子用于markdown转成html后的效果测试


访问该项目的github地址[static blog base](https://github.com/qmstar0/static_blog_base)以了解该网站是如何部署、生成的
'
---

# 一级标题

## 二级标题

### 三级标题

#### 四级标题

##### 五级标题

###### 六级标题

## 段落

这是一个普通的段落。

这是另一个段落，包含**加粗**、*斜体*、~~删除线~~以及`行内代码`。

## 链接

[普通链接(https://github.com/qmstar0)](https://github.com/qmstar0)

[引用式链接(https://github.com/qmstar0)][example]

[example]: https://github.com/qmstar0

## 列表

### 无序列表

- 项目 1
    - 子项目 1
    - 子项目 2
- 项目 2
    - 子项目 1
    - 子项目 2

### 有序列表

1. 第一项
2. 第二项
    1. 子项 1
    2. 子项 2
3. 第三项

### 任务列表

- [x] 已完成任务
- [ ] 未完成任务

## 引用

> 这是一个引用。
>
> 这是同一个引用的第二段。

## 代码块

### 行内代码

这是 `行内代码` 示例。

### 多行代码

``` javascript
function helloWorld() {
    console.log("Hello, World!");
}
```

## GO

``` go
func main() {
	fmt.Println("Hello, World!")
}
```

### 代码块（使用缩进）

``` python
def hello_world():
	print("Hello, World!")
```

## 表格

|         头 2 |   **头 1**   | **头 3**     |
|------------:|:-----------:|:------------|
|       单元格 5 |    元格 4     | 单元格 6       |
| 单元      格 2 | 单元      格 1 | 单元      格 3 |
|       单元格 8 |    单元格 7    | 单元格 9       |

## 图片

![替代文字](https://avatars.githubusercontent.com/u/98627866?s=400&u=ebea71733d13a7764f7ed0f6d03c3ddc2a8e5a51&v=4 "图片标题")

## 水平线
---

***

___

## HTML 元素

<p>这是一个段落，用 HTML 标签表示。</p>

<strong>这个是加粗文本，用 HTML 标签表示。</strong>

## 转义字符

\*这段文字不会被渲染成斜体\*

test[test][1]

[1]: http://baidu.com

### 脚注

这是一个带有脚注的句子。[^1]

[^1]: 这是脚注的内容。

### 锚点

这是一个带有 [内部锚点](#标题) 的链接。

## 表情符号

🐚 🐲 👍
