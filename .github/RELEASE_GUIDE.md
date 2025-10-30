# Release Guide

本文档说明如何发布gocron的新版本。

## 自动发布流程

### 1. 标签发布（推荐）

当推送新的版本标签时，GitHub Actions会自动构建和发布：

```bash
# 创建并推送标签
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

### 2. 自动化流程

GitHub Actions会执行以下步骤：

1. **环境准备**
   - 设置Go 1.23环境
   - 设置Node.js 18环境
   - 安装依赖

2. **构建前端**
   - 安装yarn依赖
   - 构建Vue应用
   - 生成静态资源

3. **多平台构建**
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64, arm64)

4. **生成发布说明**
   - 自动分析git提交记录
   - 按类型分类变更内容
   - 生成下载链接

5. **创建GitHub Release**
   - 上传所有平台的二进制包
   - 发布到GitHub Releases

## 手动发布流程

### 通过GitHub界面

1. 访问项目的Actions页面
2. 选择"Manual Release"工作流
3. 点击"Run workflow"
4. 填写以下信息：
   - **Version**: 版本号（如v1.2.3）
   - **Pre-release**: 是否为预发布版本
   - **Release notes**: 自定义发布说明（可选）
5. 点击"Run workflow"执行

### 通过本地脚本

```bash
# 使用本地发布脚本
./scripts/release.sh v1.2.3

# 然后推送标签
git push origin v1.2.3
```

## 版本命名规范

遵循语义化版本控制（SemVer）：

- **主版本号**：不兼容的API修改
- **次版本号**：向下兼容的功能性新增
- **修订号**：向下兼容的问题修正

示例：
- `v1.0.0` - 首个稳定版本
- `v1.1.0` - 新增功能
- `v1.1.1` - 修复bug
- `v2.0.0` - 重大更新，可能不兼容

### 预发布版本

- `v1.2.0-alpha.1` - Alpha版本
- `v1.2.0-beta.1` - Beta版本
- `v1.2.0-rc.1` - Release Candidate

## 发布检查清单

发布前请确认：

- [ ] 所有测试通过
- [ ] 文档已更新
- [ ] CHANGELOG已更新（如果有）
- [ ] 版本号符合规范
- [ ] 没有未提交的更改
- [ ] 在正确的分支上（通常是main）

## 发布后验证

1. **检查GitHub Release**
   - 确认所有平台的包都已上传
   - 验证发布说明内容正确

2. **测试下载包**
   - 下载并测试主要平台的包
   - 验证二进制文件可以正常运行

3. **更新文档**
   - 更新README中的版本信息
   - 更新安装说明

## 故障排除

### 构建失败

1. 检查GitHub Actions日志
2. 确认所有依赖都可用
3. 验证package.sh脚本权限

### 发布失败

1. 检查GITHUB_TOKEN权限
2. 确认标签格式正确
3. 验证没有重复的标签

### 包缺失

1. 检查package.sh脚本输出
2. 验证构建目录结构
3. 确认文件路径配置正确

## 回滚发布

如果需要回滚发布：

1. **删除GitHub Release**
   ```bash
   # 通过GitHub界面删除Release
   ```

2. **删除Git标签**
   ```bash
   # 删除本地标签
   git tag -d v1.2.3
   
   # 删除远程标签
   git push origin :refs/tags/v1.2.3
   ```

3. **重新发布**
   ```bash
   # 修复问题后重新创建标签
   git tag -a v1.2.3 -m "Release v1.2.3 (fixed)"
   git push origin v1.2.3
   ```