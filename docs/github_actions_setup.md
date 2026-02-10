# GitHub Actions Secrets 配置指南

本文档介绍如何为 Gocron 项目配置 GitHub Actions 所需的 Docker Hub Secrets。

## 1. 注册与获取 Docker Hub Token

1.  登录 [Docker Hub](https://hub.docker.com/)。如果没有账号，请先注册。
2.  点击右上角头像，选择 **Account Settings**。
3.  在左侧菜单选择 **Security**。
4.  在 **Access Tokens** 区域，点击 **New Access Token** 按钮。
5.  **Description**: 输入 Token 描述，例如 `GitHub Actions Gocron`。
6.  **Access permissions**: 选择 `Read, Write, Delete`。
7.  点击 **Generate**。
8.  **复制生成的 Token**。注意：Token 只会显示一次，请妥善保存。

## 2. 配置 GitHub Repository Secrets

1.  打开项目的 GitHub 仓库页面。
2.  点击上方的 **Settings** 选项卡。
3.  在左侧菜单展开 **Secrets and variables**，选择 **Actions**。
4.  在 **Repository secrets** 区域，点击 **New repository secret** 按钮，依次添加以下三个 Secret：

### DOCKER_HUB_USERNAME

- **Name**: `DOCKER_HUB_USERNAME`
- **Secret**: 输入您的 Docker Hub 用户名（例如 `tabortao`）。
- **说明**: 用于登录 Docker Hub。

### DOCKER_HUB_TOKEN

- **Name**: `DOCKER_HUB_TOKEN`
- **Secret**: 粘贴您在第 1 步生成的 Access Token。
- **说明**: 用于 Docker Hub 认证的密码/令牌。

### DOCKER_HUB_NAMESPACE

- **Name**: `DOCKER_HUB_NAMESPACE`
- **Secret**: 输入您希望推送镜像的命名空间。
  - 如果是个人账号，通常与用户名相同（例如 `tabortao`）。
  - 如果是组织账号，输入组织名称（例如 `gocron-team`）。
- **说明**: 镜像将推送到 `${DOCKER_HUB_NAMESPACE}/gocron`。

## 3. 验证

配置完成后，当您推送 `v*.*.*` 标签或手动触发 workflow 时，GitHub Actions 将能够成功登录 Docker Hub 并推送镜像。
