# GitHub Actions CI流水线配置指南

本文档详细介绍了项目中的GitHub Actions CI流水线配置，帮助您了解如何自动构建Windows可执行文件。

## CI流水线概述

我们已经为您创建了一个完整的GitHub Actions工作流文件：`.github/workflows/build.yml`，该工作流可以：

- 自动在Windows环境下构建项目
- 分别构建无CGO依赖的CLI版本和有CGO依赖的GUI版本
- 尝试使用nogl标签构建软件渲染模式的GUI版本
- 将构建产物上传为可下载的artifact
- 可选地测试CLI版本的基本功能

## 工作流触发条件

工作流会在以下情况下自动触发：

1. 代码推送到`master`分支时
2. 有PR（Pull Request）提交到`master`分支时
3. 手动通过GitHub界面触发（`workflow_dispatch`）

## CI流水线详细步骤

### 1. 环境准备

- 使用最新的Windows服务器环境（`windows-latest`）
- 配置Go 1.20版本环境
- 安装Visual Studio Build Tools以支持CGO编译
- 设置依赖缓存以加速构建过程

### 2. 构建过程

工作流会依次构建以下版本：

#### 2.1 CLI版本构建

```bash
# CLI版本构建命令
set CGO_ENABLED=0  # 可选：完全禁用CGO
go build -o hastecli.exe ./cmd/hastecli
```

CLI版本没有CGO依赖，构建成功率最高，是项目的稳定功能版本。

#### 2.2 标准GUI版本构建

```bash
# 标准GUI版本构建命令
set CGO_ENABLED=1
call "C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat" amd64
go build -o hastegui.exe ./cmd/hastegui
```

标准GUI版本依赖CGO和Visual Studio编译器，在GitHub Actions环境中通常可以成功构建。

#### 2.3 nogl标签GUI版本构建

```bash
# nogl标签GUI版本构建命令
set CGO_ENABLED=1
call "C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat" amd64
go build -tags nogl -o hastegui_nogl.exe ./cmd/hastegui
```

使用nogl标签尝试构建软件渲染模式的GUI版本，可能具有更好的兼容性。

### 3. 构建产物处理

- 列出所有生成的`.exe`文件及其大小信息
- 将所有可执行文件上传为artifact，默认保留7天

### 4. 可选的测试步骤

工作流包含一个可选的测试阶段，用于验证CLI版本的基本功能是否正常。

## 如何使用CI构建产物

1. 登录GitHub，进入您的项目仓库
2. 点击顶部的"Actions"标签页
3. 在左侧选择"Build Windows Executables"工作流
4. 在右侧的工作流运行列表中，点击最新的一次运行
5. 滚动到页面底部，在"Artifacts"部分下载`goduplicate-windows-binaries.zip`
6. 解压后即可使用其中的可执行文件

## 自定义CI配置

如果您需要调整CI配置，可以修改`.github/workflows/build.yml`文件中的以下内容：

### 更改Go版本

```yaml
# 在setup-go步骤中修改go-version
uses: actions/setup-go@v5
with:
  go-version: '1.20'  # 更改为您需要的Go版本
```

### 调整构建参数

您可以修改构建命令中的参数，例如添加优化标志：

```yaml
# 在构建步骤中添加ldflags参数
run: |
  $env:CGO_ENABLED = "1"
  go build -ldflags="-s -w" -o hastegui.exe ./cmd/hastegui
```

### 延长artifact保留时间

```yaml
# 在upload-artifact步骤中修改retention-days
uses: actions/upload-artifact@v4
with:
  name: goduplicate-windows-binaries
  path: |
    *.exe
  retention-days: 14  # 更改为您需要的保留天数
```

## 故障排除

### GUI版本构建失败

如果GUI版本在CI环境中构建失败，可能是因为：
- Visual Studio Build Tools环境配置问题
- CGO相关依赖问题

解决方法：
1. 确保`.github/workflows/build.yml`文件中包含`ilammy/msvc-dev-cmd@v1`步骤
2. 尝试只构建CLI版本和nogl标签版本

### 依赖缓存问题

如果遇到依赖缓存问题，可以：
1. 手动删除GitHub Actions中的缓存
2. 修改缓存键名以强制重新缓存

## 总结

这个GitHub Actions CI流水线配置提供了一个完整的自动化构建解决方案，特别针对Windows平台的可执行文件构建进行了优化。通过这个配置，您可以确保每次代码变更后都能自动构建和测试项目，并且可以下载最新的构建产物进行使用。