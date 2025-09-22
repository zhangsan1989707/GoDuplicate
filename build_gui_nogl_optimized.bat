@echo off

REM 优化版nogl构建脚本：完全禁用CGO的软件渲染模式构建
REM 此脚本专为解决Windows上GCC/CL编译器兼容性问题设计

setlocal enabledelayedexpansion

REM 显示当前目录和Go版本信息
echo 当前目录: %cd%
echo Go版本: 
call go version
echo.

REM 检查go.mod文件是否存在
if not exist "go.mod" (
    echo 错误: 在当前目录未找到go.mod文件，请确保在项目根目录下运行此脚本！
    pause
    exit /b 1
)

REM 清理之前的构建缓存和临时文件
echo 清理Go构建缓存和临时文件...
call go clean -cache -modcache -i -r
if %errorlevel% neq 0 (
    echo 警告: 清理缓存过程中出现错误，但将继续尝试构建
)

echo.

REM 使用nogl标签并完全禁用CGO进行构建（强制软件渲染模式）
echo 正在使用nogl标签并禁用CGO构建GUI版本（纯软件渲染模式）...
echo 构建命令: go build -tags nogl -ldflags="-s -w" -o hastegui.exe ./cmd/hastegui

echo 设置CGO_ENABLED=0（完全禁用CGO）
set CGO_ENABLED=0

REM 执行构建
go build -tags nogl -ldflags="-s -w" -o hastegui.exe ./cmd/hastegui

REM 检查构建结果
if %errorlevel% equ 0 (
    if exist "hastegui.exe" (
        echo.
        echo [92m构建成功！[0m
echo 可执行文件已生成: %cd%\hastegui.exe
        echo.
        echo 重要说明:
        echo 1. 此版本使用nogl标签和CGO_ENABLED=0，完全依赖纯Go实现
        echo 2. 界面可能使用软件渲染，性能可能略低于硬件加速版本
        echo 3. 但避免了所有C编译器兼容性问题
        echo.
        echo 请运行hastegui.exe启动GUI版本
    ) else (
        echo [91m错误: 构建命令执行成功，但未找到生成的可执行文件！[0m
    )
) else (
    echo.
    echo [91m构建失败！[0m
echo 错误码: %errorlevel%
    echo.
    echo 可能的解决方案:
    echo 1. 确保您的Go环境版本不低于项目要求（当前项目使用go 1.20）
    echo 2. 尝试更新Go到最新稳定版本
    echo 3. 检查项目依赖: go mod tidy
    echo 4. 可以先使用已验证可用的CLI版本: ./hastecli.exe --paths . --mode basic --concurrency 2
    echo 5. 如果需要GUI功能，建议安装Visual Studio Build Tools并使用Developer Command Prompt构建
)

echo.
pause