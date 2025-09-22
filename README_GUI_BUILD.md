# GoDuplicate GUI版本构建指南

经过多次尝试，我们发现最可靠的构建GUI版本的方法是使用Visual Studio的Developer Command Prompt。以下是详细步骤：

## 方法1：使用Developer Command Prompt for VS 2022

1. **打开Developer Command Prompt for VS 2022**
   - 点击Windows开始菜单
   - 搜索并打开"Developer Command Prompt for VS 2022"

2. **在命令提示符中执行以下命令**：
   ```cmd
   cd e:\Mike\GoDuplicate
   set CGO_ENABLED=1
   go clean -cache
   go build -o hastegui.exe ./cmd/hastegui
   ```

3. **如果构建成功**，当前目录将生成hastegui.exe文件，直接运行即可启动GUI版本。

## 方法2：使用普通命令提示符设置编译器环境

如果找不到Developer Command Prompt，可以尝试以下步骤：

1. **按Win+R，输入cmd，打开命令提示符**

2. **运行Visual Studio的vcvarsall.bat脚本设置环境变量**：
   ```cmd
   "C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat" amd64
   ```

3. **然后执行构建命令**：
   ```cmd
   cd e:\Mike\GoDuplicate
   set CGO_ENABLED=1
   go clean -cache
   go build -o hastegui.exe ./cmd/hastegui
   ```

## 方法3：尝试软件渲染模式（nogl标签）

如果以上方法都不奏效，可以尝试使用nogl标签进行软件渲染模式构建，这可能会减少对C编译器的依赖：

```cmd
cd e:\Mike\GoDuplicate
set CGO_ENABLED=1
go clean -cache
go build -tags nogl -o hastegui.exe ./cmd/hastegui
```

## 故障排除

- 如果遇到`cgo: C compiler "gcc" not found`错误，说明Go的CGO没有正确识别Visual Studio的编译器
- 确保在Visual Studio Installer中安装了"使用C++的桌面开发"和Windows SDK组件
- 确认cl.exe路径正确：`C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Tools\MSVC\14.44.35207\bin\Hostx64\x64\cl.exe`

## CLI版本（已验证可用）

如果您只需要测试核心功能，CLI版本已经成功构建并验证可以使用：

```cmd
./hastecli.exe --paths . --mode basic --concurrency 2
```

该命令将扫描当前目录，查找重复文件并显示结果。