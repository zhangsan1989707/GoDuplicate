@echo off

REM 检查Visual Studio Build Tools的vcvarsall.bat文件
set "VCVARSALL_PATH=C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat"

if exist "%VCVARSALL_PATH%" (
    echo 找到vcvarsall.bat，正在设置编译器环境...
    call "%VCVARSALL_PATH%" amd64
    
    echo
    echo 正在切换到项目目录...
    cd /d "e:/Mike/GoDuplicate"
    
    echo
    echo 正在设置CGO环境变量...
    set CGO_ENABLED=1
    
    echo
    echo 正在清理缓存...
    go clean -cache
    
    echo
    echo 开始构建GUI版本...
    go build -o hastegui.exe ./cmd/hastegui
    
    echo
    if exist "hastegui.exe" (
        echo 构建成功！hastegui.exe已生成。
        echo 请运行hastegui.exe启动GUI版本。
    ) else (
        echo 构建失败，请检查错误信息。
        echo 您也可以尝试使用nogl标签进行软件渲染模式构建：
        echo go build -tags nogl -o hastegui.exe ./cmd/hastegui
    )
) else (
    echo 未找到vcvarsall.bat文件，请确认VS Build Tools安装路径。
    echo 默认路径应为：C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat
)

pause