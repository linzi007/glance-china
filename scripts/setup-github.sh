#!/bin/bash

# Glance 中国版 GitHub 设置脚本

set -e

echo "🚀 设置 Glance 中国版 GitHub 仓库..."

# 检查是否已经是 git 仓库
if [ ! -d ".git" ]; then
    echo "📁 初始化 Git 仓库..."
    git init
    git add .
    git commit -m "Initial commit: Glance China version"
else
    echo "✅ Git 仓库已存在"
fi

# 检查是否有远程仓库
if ! git remote get-url origin > /dev/null 2>&1; then
    echo "❓ 请输入您的 GitHub 仓库 URL (例如: https://github.com/username/glance-china.git):"
    read -r REPO_URL
    
    echo "🔗 添加远程仓库..."
    git remote add origin "$REPO_URL"
else
    echo "✅ 远程仓库已配置"
fi

# 检查当前分支
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "🔄 切换到 main 分支..."
    git checkout -b main 2>/dev/null || git checkout main
fi

# 推送到 GitHub
echo "📤 推送代码到 GitHub..."
git push -u origin main

echo "✅ GitHub 仓库设置完成！"
echo ""
echo "🎉 下一步："
echo "1. 访问 https://github.com/your-username/glance-china"
echo "2. 在仓库设置中确认仓库为 Public"
echo "3. 添加仓库描述和标签"
echo "4. 设置 GitHub Pages (可选)"
echo ""
echo "📚 文档链接："
echo "- README: https://github.com/your-username/glance-china/blob/main/README.md"
echo "- 部署指南: https://github.com/your-username/glance-china/blob/main/docs/DEPLOYMENT.md"
