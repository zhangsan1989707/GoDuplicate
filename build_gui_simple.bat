@echo off

REM 切换到项目目录
cd /d "e:/Mike/GoDuplicate"

REM 清理缓存
go clean -cache

REM 使用nogl标签进行软件渲染模式构建
set CGO_ENABLED=1
echo 正在使用nogl标签构建GUI版本...
go build -tags nogl -o hastegui.exe ./cmd/hastegui

REM 检查构建结果
if exist "hastegui.exe" (
    echo 构建成功！hastegui.exe已生成。
    echo 请运行hastegui.exe启动GUI版本。
) else (
    echo 构建失败。
    echo 尝试其他方法：
    echo 方法1: 在Developer Command Prompt for VS 2022中运行：
    echo    cd e:\Mike\GoDuplicate
    echo    set CGO_ENABLED=1
    echo    go build -o hastegui.exe ./cmd/hastegui
    echo 方法2: 或在命令提示符中先运行：
    echo    "C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat" amd64
    echo    然后执行上述构建命令
)

pause