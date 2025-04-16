# devc - 开发者剪切板监听与翻译工具

## 简介

`devc` 是一款运行在 Windows 系统下的剪切板监听小工具，使用 Go 语言开发，通过调用 Windows API
实现功能。它能够记录每天的剪切板内容，并提供翻译功能，支持中英互译，特别适用于开发人员在编写代码或数据库表字段时，快速将中文字段翻译为英文字段。

## 功能特点

- **剪切板内容记录**：每天自动记录剪切板中的内容。
- **翻译功能**：
    - 支持中英互译。
    - 支持关键字翻译，可自定义关键字长度。
    - 翻译后的内容可以重新写入剪切板。
- **命名规则转换**：
    - 支持大驼峰命名（PascalCase）。
    - 支持小驼峰命名（camelCase）。
    - 支持下划线命名（snake_case）。
- **命令行配置**：
    - 支持通过命令行参数进行个性化配置。
    - 支持将 CMD 窗口置顶。
- **应用场景**：
    - 特别适合数据库表字段设置，可将中文字段自动翻译为英文字段。

## 命令行参数

运行 `devc -h` 可查看所有支持的命令行参数及其说明：

```bash
devc -h
Usage of devc:
  -han int
        [关键字翻译] 汉语关键字长度，默认支持4个汉字，最少2个汉字。例如：订单编号、订单状态、创建时间 (default 4)
  -hump int
        [关键字翻译] 将翻译内容写入剪切板，命名规则：1.大驼峰 2.小驼峰 3.下划线 (default 1)
  -top
        [窗口] 置顶当前 CMD 窗口
  -tran
        [翻译工具] 开启自动检测翻译，关键字翻译将关闭
  -w
        [翻译工具] 将翻译内容写入剪切板
```

[**注意**]：关键字翻译和翻译工具命令不能同时，同时使用优先翻译工具，未指定参数默认关键字翻译

### 参数说明

- `-han`：设置汉语关键字的长度，默认为 4 个汉字，最少为 2 个汉字。
- `-hump`：设置翻译内容的命名规则：
    - `1`：大驼峰命名（PascalCase）。
    - `2`：小驼峰命名（camelCase）。
    - `3`：下划线命名（snake_case）。
- `-top`：将当前 CMD 窗口置顶。
- `-tran`：开启自动检测翻译，关键字翻译将关闭。
- `-w`：将翻译后的内容写入剪切板。


## 数据保存

### 历史记录

- **默认保存路径**：`%APPDATA%\devc`。
- **自定义保存路径**：通过配置 `DEVC` 环境变量指定路径。
  - 打开系统环境变量设置。
  - 添加一个名为 `DEVC` 的环境变量，值为自定义路径（例如 `D:\devc_data`）。

### 字典文件

- **字典文件**：`dict.txt`
- **文件格式**：
  ```lua
  协程|coroutine|Coroutine|coroutine|coroutine
  格式转换|Format conversion|FormatConversion|formatConversion|format_conversion
  解析文件|Parsing file|ParsingFile|parsingFile|parsing_file
  文件路径|File path|FilePath|filePath|file_path
  ```
  - 每行表示一个关键字及其翻译。
  - 格式为：`中文|英文|大驼峰|小驼峰|下划线`。
- **自定义字典文件**：将自定义的 `dict.txt` 文件放置在历史记录目录中，或通过环境变量指定路径。

## 使用示例

### 基础使用

```bash
devc -han=3 -hump=2
```

- 指定关键字长度为 3 个汉字。
- 使用小驼峰命名规则（camelCase）。

### 置顶 CMD 窗口

```bash
devc -top
```

- 将当前 CMD 窗口置顶。

### 自动翻译

```bash
devc -tran -w
```

- 开启自动检测翻译。
- 将翻译后的内容写入剪切板。

## 安装与运行

1. **下载**：从 [GitHub 仓库](https://github.com/systemmin/devc) 下载最新版本的 `devc`。
2. **安装**：将下载的文件解压到任意目录，或丢到 `C:\Windows\System32` 目录下可直接 `cmd` 使用。
3. **运行**：打开 CMD 窗口，切换到解压目录，运行以下命令：
   ```bash
   devc -h
   ```

## 注意事项

- 本工具仅支持 Windows 系统。
- 如果遇到问题，请查看日志文件或联系开发者。
- 关键字翻译和翻译工具命令不能同时，同时使用优先翻译工具，未指定参数默认关键字翻译

## 贡献与反馈

欢迎通过 GitHub 提交 Issue 或 Pull Request，帮助改进 `devc`。

## 参考资料

- [Window API](https://learn.microsoft.com/zh-cn/windows/win32/api/winuser/nf-winuser-openclipboard)
- [Github atotto/clipboard](https://github.com/atotto/clipboard/blob/master/clipboard_windows.go)