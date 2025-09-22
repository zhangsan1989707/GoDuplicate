# 检查vcvarsall.bat文件是否存在
$vcvarsallPath = 'C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat'

if (Test-Path $vcvarsallPath) {
    Write-Host '找到vcvarsall.bat，正在设置编译器环境...'
    
    # 运行vcvarsall.bat设置编译器环境（amd64表示64位编译器）
    & $vcvarsallPath amd64
    
    # 等待环境变量设置完成
    Start-Sleep -Seconds 2
    
    # 切换到项目目录
    Set-Location 'e:/Mike/GoDuplicate'
    
    # 设置CGO_ENABLED环境变量
    $env:CGO_ENABLED = 1
    
    # 清理缓存
    go clean -cache
    
    # 构建GUI版本
    Write-Host '开始构建GUI版本...'
    go build -o hastegui.exe ./cmd/hastegui
    
    # 检查构建结果
    if (Test-Path 'e:/Mike/GoDuplicate/hastegui.exe') {
        Write-Host '构建成功！hastegui.exe已生成。'
    } else {
        Write-Host '构建失败，请检查错误信息。'
    }
} else {
    Write-Host '未找到vcvarsall.bat文件，请确认VS Build Tools安装路径。'
    Write-Host '默认路径应为：C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat'
}