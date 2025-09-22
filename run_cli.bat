@echo off

REM GoDuplicate CLI版本运行脚本

REM 显示当前目录
echo 当前目录: %cd%

REM 检查hastecli.exe是否存在
if not exist "hastecli.exe" (
    echo 错误: 未找到hastecli.exe文件！
    echo 正在尝试重新构建CLI版本...
    go build -o hastecli.exe ./cmd/hastecli
    if not exist "hastecli.exe" (
        echo 构建失败，请检查Go环境是否正确安装
        pause
        exit /b 1
    )
)

echo.
echo ===================================================
echo GoDuplicate CLI版本已准备就绪！
echo ===================================================
echo 以下是一些常用的命令示例：
echo.
echo 1. 扫描当前目录，使用基本模式和2个并发度：
echo    .\hastecli.exe --paths . --mode basic --concurrency 2

echo 2. 扫描多个目录，排除临时文件：
echo    .\hastecli.exe --paths "D:\;E:\docs" --exclude "*.tmp;*.bak"

echo 3. 使用SHA256哈希算法扫描：
echo    .\hastecli.exe --paths . --hash sha256

echo 4. 只扫描大于1MB的文件：
echo    .\hastecli.exe --paths . --min-size 1048576

echo ===================================================
echo.

REM 询问用户是否直接运行示例命令
echo 要直接运行默认扫描吗？(y/n)
set /p choice=

if /i "%choice%" equ "y" (
    echo 正在扫描当前目录...
    .\hastecli.exe --paths . --mode basic --concurrency 2
) else (
    echo 您可以手动输入上述命令之一来运行CLI版本。
)

pause