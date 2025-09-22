@echo off

REM 检查是否在Developer Command Prompt中运行
cl /? >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo 错误：请在Developer Command Prompt for VS 2022中运行此脚本！
    echo 或先运行vcvarsall.bat设置编译器环境。
    pause
    exit /b 1
)

REM 显示当前目录，帮助用户确认
echo 当前目录：%CD%

REM 切换到项目目录
echo 正在切换到项目目录...
cd /d "e:/Mike/GoDuplicate"

REM 再次显示当前目录，确认切换成功
echo 切换后的目录：%CD%

REM 检查go.mod文件是否存在
if not exist "go.mod" (
    echo 错误：在e:\Mike\GoDuplicate目录下未找到go.mod文件！
    echo 请确认项目路径是否正确。
    pause
    exit /b 1
)

REM 设置CGO_ENABLED=1
echo 设置CGO_ENABLED=1
set CGO_ENABLED=1

REM 清理缓存
echo 清理缓存...
go clean -cache

REM 构建GUI版本
echo 开始构建GUI版本...
go build -o hastegui.exe ./cmd/hastegui

REM 检查构建结果
if exist "hastegui.exe" (
    echo 构建成功！hastegui.exe已生成。
    echo 请运行hastegui.exe启动GUI版本。
) else (
    echo 构建失败，请检查错误信息。
    echo 尝试以下解决方案：
    echo 1. 确保Visual Studio Build Tools安装了"使用C++的桌面开发"组件
    echo 2. 确认Windows SDK已正确安装
    echo 3. 尝试使用nogl标签进行软件渲染模式构建
    echo    go build -tags nogl -o hastegui.exe ./cmd/hastegui
)

pause