# HasteGUI / HasteCLI

## 概述
- 本项目旨在将重复文件扫描器升级为 GUI + CLI 的一体化工具。
- 目前已提供：
  - CLI 可执行：基础参数解析与最小扫描引擎（按快速哈希聚类）
  - GUI 骨架：使用 Fyne

## 运行要求
- Go 1.20+
- Windows 10/11（GUI 构建需要 C 工具链）

## 快速开始（CLI）
```powershell
# 构建 CLI
go build -o hastecli.exe .\cmd\hastecli

# 运行示例
./hastecli.exe --paths "D:\\,E:\\docs" --exclude "*.tmp;node_modules" --mode basic --concurrency 4
```

## GUI 构建（方案A：VS Build Tools）
1. 安装 Microsoft Visual Studio Build Tools 2022（含“使用 C++ 的桌面开发”与 Windows SDK）。
2. 打开 “x64 Native Tools Command Prompt for VS 2022”，执行：
```bat
cd /d D:\bak\GoDuplicate
set CGO_ENABLED=1
go clean -cache
go build -o hastegui.exe .\cmd\hastegui
```
3. 若临时不安装 C 工具链，可尝试软件渲染：
```bat
go build -tags nogl -o hastegui.exe .\cmd\hastegui
```

## 目录结构
```
cmd/
  hastecli/        # CLI 入口
  hastegui/        # GUI 入口（Fyne）
internal/
  core/            # 核心模型与扫描引擎接口/实现
  gui/             # GUI 代码
```

## 路线图（摘自需求）
- 扫描配置界面/监控/结果/策略/执行/设置 六大模块
- 处理策略系统与撤销机制
- 性能与可靠性优化

## 字体使用说明
- 项目中的`font`目录包含用于中文显示的字体文件（如simhei.ttf、simsunb.ttf）
- 这些字体文件仅用于开发测试目的，请确保在实际部署时遵守相关字体的版权协议
- 如需在商业环境中使用，请自行获取合法的字体授权

## 免责声明
- 当前扫描实现为最小可用版本，后续将替换为高性能实现。
- 项目中包含的字体文件仅供开发和测试使用，不代表对这些字体的版权拥有或授权。

