# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Go PDF 读取库，fork 自 `rsc.io/pdf`。提供 PDF 文本提取（纯文本/带样式/按行按列分组）、页面遍历、字体解码、大纲提取等功能。模块路径：`github.com/ledongthuc/pdf`。

## 常用命令

```bash
go test -v ./...          # 运行所有测试
go build -v ./...         # 构建
go mod download           # 下载依赖（仅标准库，无外部依赖）
```

目前没有 `_test.go` 测试文件，CI 通过 `go test -v ./...` 验证编译通过。

## 架构

核心数据模型是 `Value` 类型（read.go），它用统一的访问器暴露 PDF 的八种对象类型（Null/Bool/Integer/Real/String/Name/Dict/Array/Stream）。设计哲学是零值容错：类型不匹配时返回零值而非报错，允许快速遍历 PDF 结构但可能静默吞掉错误。

**数据流**：`Open()` → `Reader` 解析 xref 表 → `Reader.Page(n)` 遍历页面树 → `Page` 方法提取内容。

**关键文件与职责**：
- `read.go` — 核心引擎：`Value`/`Reader` 类型、xref 解析、对象解析、加密解密（RC4/AES）。这是最大的文件（~2100 行）。
- `page.go` — 页面处理：`Page`/`Font`/`Text`/`Content` 类型、内容流解析（PDF 操作符解释）、文本提取（按行/列/样式）、`Outline` 大纲提取。
- `lex.go` — 词法分析：PDF token 解析（关键字/名称/字符串/数字）、`buffer` IO、对象指针解析。
- `text.go` — 字符编码：PDFDocEncoding/WinAnsiEncoding/MacRomanEncoding/UTF-16BE 解码。
- `name.go` — PDF 名称常量映射。
- `ps.go` — PostScript 解释器（栈式执行），用于处理 PDF 中的 PostScript 计算表达式。
- `ascii85.go` — ASCII85 解码。

**核心类型**：
- `Reader` — PDF 文件读取器，持有 xref 表、trailer、解密密钥。
- `Value` — PDF 对象的统一表示，所有访问通过 `.Key()/.Index()/.Int64()/.String()` 等方法。
- `Page` — 页面包装器，提供 `GetPlainText()`/`GetTextByRow()`/`GetTextByColumn()`/`GetStyledTexts()`。
- `Font` — 字体包装器，内含 `TextEncoding` 用于字符映射。
- `TextEncoding` 接口 — 字符编码抽象（nopEncoder/byteEncoder/cmap 等实现）。

## 注意事项

- `DebugOn = true` 可开启调试日志，排查 PDF 解析问题时使用。
- 页码从 1 开始（非 0）。
- 加密文件支持通过 `NewReaderEncrypted()` 处理，传入密码回调函数。
- `page.go` 中的内容流解析器实现了 PDF 图形状态机（CTM 矩阵变换、文本矩阵），修改文本提取逻辑时需注意矩阵运算。
