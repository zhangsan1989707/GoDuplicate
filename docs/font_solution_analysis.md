# 中文字体显示问题解决方案分析

## 问题概述
GUI界面切换至中文显示后仍出现乱码问题，经排查怀疑是GitHub构建环境中缺少中文字体所致。目前已在项目目录的font文件夹下手动存放了中文字体文件。

## 方案1：通过修改代码初始化中文字体

### 实现原理
在应用程序启动时，直接从项目目录加载中文字体文件，并配置Fyne框架使用这些字体进行中文文本渲染。

### 具体实现
需要修改`internal/gui/app.go`文件，添加字体加载逻辑：

```go
// 重写字体方法，使用项目中的中文字体
func (c customTheme) Font(s fyne.TextStyle) fyne.Resource {
    // 尝试从项目目录加载中文字体
    fontPath := "font/simhei.ttf" // 黑体，适用于大多数中文显示场景
    
    // 检查字体文件是否存在
    if _, err := os.Stat(fontPath); err == nil {
        // 加载字体文件
        fontData, err := ioutil.ReadFile(fontPath)
        if err == nil {
            // 创建字体资源
            fontRes := fyne.NewStaticResource("simhei.ttf", fontData)
            return fontRes
        }
    }
    
    // 如果加载失败，返回nil让系统自动查找
    return nil
}

// 修改ensureChineseFontSupport函数，预加载字体
func ensureChineseFontSupport() {
    // 提前检查并加载字体文件，确保在UI渲染前可用
    fontPath := "font/simhei.ttf"
    if _, err := os.Stat(fontPath); err == nil {
        // 可以在这里添加字体预加载逻辑
    }
}
```

### 可行性分析
- **优点**：
  1. 实现简单，只需修改少量代码
  2. 不依赖构建环境配置，应用程序可以在任何环境中运行
  3. 字体文件与应用程序一起分发，确保一致性

- **缺点**：
  1. 需要确保字体文件随应用程序一起打包和分发
  2. 可能会增加应用程序的大小
  3. 需注意字体的版权问题

### 适用场景
适合希望应用程序在任何环境中都能独立运行，不依赖系统字体配置的情况。

## 方案2：将字体文件上传至GitHub，并通过流水线配置中文字体

### 实现原理
在GitHub Actions构建过程中，将字体文件安装到构建环境中，确保应用程序在构建时能够访问到中文字体。

### 具体实现
需要修改`.github/workflows/go.yml`文件，添加字体安装步骤：

```yaml
# 在安装Visual Studio Build Tools之后，构建应用程序之前添加以下步骤
- name: 安装中文字体
  run: |
    # 创建字体目录
    $fontDir = "${env:WINDIR}\Fonts"
    Write-Host "字体目录: $fontDir"
    
    # 复制字体文件到系统字体目录
    Copy-Item -Path .\font\simhei.ttf -Destination $fontDir -Force
    Copy-Item -Path .\font\simsunb.ttf -Destination $fontDir -Force
    
    # 注册字体（可选，某些应用程序可能需要）
    reg add "HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts" /v "SimHei (TrueType)" /t REG_SZ /d simhei.ttf /f
    reg add "HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts" /v "SimSun-Bold (TrueType)" /t REG_SZ /d simsunb.ttf /f
  shell: pwsh
```

### 可行性分析
- **优点**：
  1. 应用程序代码不需要修改
  2. 构建环境配置集中管理
  3. 不需要在应用程序中包含字体文件

- **缺点**：
  1. 需要修改CI/CD配置，增加了构建步骤
  2. 依赖GitHub Actions构建环境的权限（安装字体可能需要管理员权限）
  3. 仅解决了构建环境中的字体问题，最终用户仍需系统中安装中文字体

### 适用场景
适合希望保持应用程序代码简洁，不希望包含额外资源文件的情况。

## 推荐方案

综合分析，**推荐方案1（通过修改代码初始化中文字体）**，原因如下：

1. **可靠性更高**：应用程序自带字体，不受系统环境限制
2. **实现简单**：只需修改少量代码，不需要复杂的环境配置
3. **用户体验一致**：所有用户都能看到相同的字体效果，无论其系统中是否安装了中文字体
4. **易于测试和验证**：可以在开发环境中直接验证效果

## 实施建议

如果选择方案1，建议按照以下步骤实施：

1. 修改`internal/gui/app.go`文件，实现字体加载逻辑
2. 更新`.gitignore`文件，确保字体文件被包含在版本控制中
3. 修改构建脚本(`build_gui.bat`和`build_gui.ps1`)，确保字体文件与可执行文件一起分发
4. 在README中添加关于字体版权的说明（如适用）

如果选择方案2，建议与方案1结合使用，以确保在不同环境中的最佳兼容性。