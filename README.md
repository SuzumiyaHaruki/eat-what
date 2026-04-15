# Eat What

一个使用 Go + Fyne 编写的桌面小工具，用来帮你从菜单里随机决定“今天吃什么”。

项目现在支持：
- 新建菜单并自动生成对应的 `.txt` 文件
- 从 `.txt` 导入菜单
- 手动添加、删除、清空菜单项
- 记住当前菜单和选项
- 一键随机抽取今天吃什么
- 打包为 Windows 可执行文件

## 预览

- 主界面聚焦在“今日推荐”结果区
- 菜单支持本地 `.txt` 文件管理
- 可通过弹窗进行添加、导入、管理和重命名

## 技术栈

- Go `1.25.8`
- [Fyne](https://fyne.io/) `v2.7.3`

## 运行环境

建议先安装：
- Go
- 桌面运行所需的图形环境

在项目目录执行：

```bash
go run .
```

如果只想编译当前平台：

```bash
go build .
```

## 菜单文件说明

菜单会保存成项目根目录下的 `.txt` 文件。

例如你创建一个叫 `食堂` 的菜单，就会在当前项目目录生成：

```text
食堂.txt
```

文件内容是一行一个选项，例如：

```text
麻辣烫
牛肉面
寿司
```

## 功能说明

### 1. 新建菜单

- 点击 `新建菜单`
- 输入菜单名称
- 创建后会立刻在项目目录生成对应的 `.txt` 文件

### 2. 添加选项

- 点击 `添加选项`
- 一行输入一个食物
- 已存在的重复项会自动去重

### 3. 导入 TXT

- 支持从现有 `.txt` 文件导入
- 可选择：
  - 在原菜单基础上增加
  - 覆盖原菜单

### 4. 管理菜单

- 查看当前所有选项
- 删除选中项
- 清空全部项
- 重命名菜单

### 5. 随机抽取

- 点击 `吃什么`
- 程序会从当前菜单里随机给出一个结果

## Windows 打包

项目里已经提供了一键打包脚本：

```bash
./build-windows.sh
```

默认会使用项目根目录下的图标文件：

```text
eat_what.png
```

如果想指定其他图标：

```bash
./build-windows.sh ./your-icon.png
```

### 打包前需要准备

1. 安装 `fyne` 打包工具

脚本会自动尝试安装：

```bash
go install fyne.io/tools/cmd/fyne@latest
```

2. 安装 Windows 交叉编译器

在 Ubuntu / Debian 上：

```bash
sudo apt install gcc-mingw-w64-x86-64
```

### 打包命令实际执行内容

脚本内部会调用：

```bash
fyne package --target windows --icon eat_what.png --release --app-id com.example.eatwhat
```

## 项目结构

```text
eat-what/
├── main.go
├── ui.go
├── food_manager.go
├── file_loader.go
├── persistence.go
├── build-windows.sh
├── eat_what.png
└── *.txt
```

各文件职责：

- `main.go`：程序入口
- `ui.go`：界面与交互逻辑
- `food_manager.go`：菜单数据管理
- `file_loader.go`：读取 txt 文件
- `persistence.go`：状态保存、菜单文件路径与 txt 写入
- `build-windows.sh`：Windows 打包脚本

## 开发说明

格式化代码：

```bash
gofmt -w *.go
```

检查能否编译：

```bash
go build ./...
```

## 后续可扩展方向

- 抽取动画优化
- 历史记录
- 多菜单切换
- 导出/备份功能
- 更完整的 Windows 发布目录

## License

如果你准备公开发布到 GitHub，建议补一个 LICENSE 文件，例如 MIT。
