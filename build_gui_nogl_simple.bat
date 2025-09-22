@echo off

REM 简化版nogl构建脚本 - 只包含最基本的构建命令

REM 清理构建缓存
go clean -cache -modcache

REM 完全禁用CGO并使用nogl标签构建
echo 设置CGO_ENABLED=0并使用nogl标签构建...
set CGO_ENABLED=0
go build -tags nogl -o hastegui.exe ./cmd/hastegui

REM 检查结果
if exist "hastegui.exe" (
    echo 构建成功！hastegui.exe已生成
    echo 请运行hastegui.exe启动GUI版本
    echo 注意：此版本使用软件渲染，性能可能略低
) else (
    echo 构建失败
    echo 请尝试以下命令手动构建CLI版本（已验证可用）：
    echo go build -o hastecli.exe ./cmd/hastecli
)

pause