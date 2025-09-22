@echo off

REM 构建脚本：使用nogl标签进行软件渲染模式构建（避免CGO依赖）
REM 此脚本适用于遇到GCC编译器兼容性问题的情况

setlocal enabledelayedexpansion

REM 设置项目根目录为当前目录
set PROJECT_DIR=%~dp0
cd /d "%PROJECT_DIR%"

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

REM 清理之前的构建缓存
echo 清理Go构建缓存...
call go clean -cache -modcache -i -r
if %errorlevel% neq 0 (
    echo 警告: 清理缓存过程中出现错误，但将继续尝试构建
)

echo.

REM 使用nogl标签进行构建（软件渲染模式，避免CGO依赖）
echo 正在使用nogl标签构建GUI版本（软件渲染模式）...
echo 构建命令: go build -tags nogl -o hastegui.exe ./cmd/hastegui
call go build -tags nogl -o hastegui.exe ./cmd/hastegui

REM 检查构建结果
if %errorlevel% equ 0 (
    if exist "hastegui.exe" (
        echo.
        echo [92m构建成功！[0m
echo 可执行文件已生成: %cd%\hastegui.exe
        echo.
        echo 由于使用了nogl标签，此版本使用软件渲染模式，不依赖CGO和图形库。
        echo 功能上可能略有限制，但可以避免编译器兼容性问题。
    ) else (
        echo [91m错误: 构建命令执行成功，但未找到生成的可执行文件！[0m
    )
) else (
    echo.
    echo [91m构建失败！[0m
echo 错误码: %errorlevel%
    echo.
    echo 可能的解决方案:
    echo 1. 确保您的Go环境正确安装: go version
    echo 2. 尝试更新Go到最新版本
    echo 3. 检查项目依赖: go mod tidy
    echo 4. 如果问题持续存在，可以尝试直接使用CLI版本: go build -o hastecli.exe ./cmd/hastecli
)

echo.
pause