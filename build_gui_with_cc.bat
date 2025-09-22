@echo off

REM 设置CC环境变量为用户提供的cl.exe路径
set "CC=C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Tools\MSVC\14.44.35207\bin\Hostx64\x64\cl.exe"
set "CXX=C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Tools\MSVC\14.44.35207\bin\Hostx64\x64\cl.exe"

REM 检查cl.exe是否存在
if exist "%CC%" (
    echo 找到cl.exe编译器：%CC%
    
    REM 切换到项目目录
    cd /d "e:/Mike/GoDuplicate"
    
    REM 清理缓存
go clean -cache
    
    REM 设置CGO_ENABLED=1并构建GUI版本
    set CGO_ENABLED=1
    echo 正在使用指定的编译器构建GUI版本...
    go build -o hastegui.exe ./cmd/hastegui
    
    REM 检查构建结果
    if exist "hastegui.exe" (
        echo 构建成功！hastegui.exe已生成。
        echo 请运行hastegui.exe启动GUI版本。
    ) else (
        echo 构建失败，请检查错误信息。
        echo 您也可以尝试使用nogl标签进行软件渲染模式构建：
        echo go build -tags nogl -o hastegui.exe ./cmd/hastegui
    )
) else (
    echo 未找到cl.exe文件，请确认路径是否正确。
    echo 提供的路径：%CC%
)

pause