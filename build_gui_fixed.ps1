# 检查vcvarsall.bat文件是否存在
$vcvarsallPath = 'C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat'

if (Test-Path $vcvarsallPath) {
    Write-Host '找到vcvarsall.bat，正在设置编译器环境...'
    
    # 运行vcvarsall.bat设置编译器环境（amd64表示64位编译器）
    & $vcvarsallPath amd64
    
    # 等待环境变量设置完成
    Start-Sleep -Seconds 2
    
    # 设置CGO_ENABLED环境变量
    $env:CGO_ENABLED = 1
    
    # 清理缓存
    go clean -cache
    
    # 构建GUI版本，添加-ldflags="-H windowsgui"隐藏终端窗口
    Write-Host '开始构建GUI版本...'
    go build -ldflags="-H windowsgui" -o hastegui.exe ./cmd/hastegui
    
    # 检查构建结果
    if (Test-Path 'hastegui.exe') {
        Write-Host '构建成功！hastegui.exe已生成。'
        
        # 检查字体目录是否存在
        if (Test-Path 'font') {
            Write-Host '正在复制字体文件...'
            
            # 确保目标字体目录存在
            if (!(Test-Path 'dist')) {
                New-Item -ItemType Directory -Path 'dist' | Out-Null
            }
            
            # 复制可执行文件到dist目录
            Copy-Item -Path 'hastegui.exe' -Destination 'dist/' -Force
            
            # 复制字体目录到dist目录
            Copy-Item -Path 'font' -Destination 'dist/' -Recurse -Force
            
            Write-Host '字体文件已成功复制到dist目录。'
            Write-Host '请运行dist/hastegui.exe启动GUI版本。'
        } else {
            Write-Host '未找到字体目录，请确保font文件夹存在于项目根目录。'
            Write-Host '请手动复制font目录与hastegui.exe放在同一目录下。'
        }
    } else {
        Write-Host '构建失败，请检查错误信息。'
    }
} else {
    Write-Host '未找到vcvarsall.bat文件，请确认VS Build Tools安装路径。'
    Write-Host '默认路径应为：C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat'
    Write-Host '或者，您可以尝试使用其他版本的Visual Studio，调整上面的路径。'