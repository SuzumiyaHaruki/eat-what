#!/usr/bin/env bash

set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ICON_PATH="${1:-$PROJECT_DIR/eat_what.png}"
APP_ID="com.example.eatwhat"
WINDOWS_CC="${WINDOWS_CC:-x86_64-w64-mingw32-gcc}"

export PATH="$PATH:$HOME/go/bin"

if ! command -v go >/dev/null 2>&1; then
  echo "错误: 未找到 go，请先安装 Go。" >&2
  exit 1
fi

if ! command -v fyne >/dev/null 2>&1; then
  echo "未找到 fyne，正在安装 fyne 打包工具..."
  go install fyne.io/tools/cmd/fyne@latest
fi

if ! command -v "$WINDOWS_CC" >/dev/null 2>&1; then
  echo "错误: 未找到 Windows 交叉编译器: $WINDOWS_CC" >&2
  echo "请先安装 MinGW-w64，例如：" >&2
  echo "  sudo apt install gcc-mingw-w64-x86-64" >&2
  echo "安装完成后再执行 ./build-windows.sh" >&2
  exit 1
fi

if [ ! -f "$ICON_PATH" ]; then
  echo "错误: 未找到图标文件: $ICON_PATH" >&2
  echo "用法: ./build-windows.sh [图标路径]" >&2
  exit 1
fi

cd "$PROJECT_DIR"

echo "开始打包 Windows 可执行文件..."
export CC="$WINDOWS_CC"
export CGO_ENABLED=1

fyne package --target windows --icon "$ICON_PATH" --release --app-id "$APP_ID"

echo "打包完成。"
