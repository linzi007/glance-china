# 贡献指南

感谢您对 Glance 中国版的关注！我们欢迎各种形式的贡献。

## 开发环境设置

1. **安装 Go 1.21+**
   \`\`\`bash
   # 检查 Go 版本
   go version
   \`\`\`

2. **克隆仓库**
   \`\`\`bash
   git clone https://github.com/your-username/glance-china.git
   cd glance-china
   \`\`\`

3. **安装依赖**
   \`\`\`bash
   go mod download
   \`\`\`

4. **运行测试**
   \`\`\`bash
   go test ./...
   \`\`\`

## 贡献类型

### 🐛 Bug 报告
- 使用 GitHub Issues 报告 bug
- 提供详细的复现步骤
- 包含系统信息和日志

### ✨ 功能请求
- 描述新功能的用途和价值
- 提供具体的使用场景
- 考虑向后兼容性

### 🔧 代码贡献
- Fork 仓库并创建功能分支
- 遵循代码规范和测试要求
- 提交 Pull Request

## 代码规范

### Go 代码风格
\`\`\`bash
# 格式化代码
go fmt ./...

# 静态检查
go vet ./...

# 运行 linter
golangci-lint run
\`\`\`

### 提交信息格式
\`\`\`
type(scope): description

[optional body]

[optional footer]
\`\`\`

类型：
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具相关

### 测试要求
- 新功能必须包含测试
- 测试覆盖率不低于 80%
- 包含单元测试和集成测试

## Pull Request 流程

1. **创建分支**
   \`\`\`bash
   git checkout -b feature/your-feature-name
   \`\`\`

2. **开发和测试**
   \`\`\`bash
   # 开发代码
   # 运行测试
   go test ./...
   # 运行基准测试
   go test -bench=. ./...
   \`\`\`

3. **提交代码**
   \`\`\`bash
   git add .
   git commit -m "feat(widget): add bilibili video widget"
   \`\`\`

4. **推送分支**
   \`\`\`bash
   git push origin feature/your-feature-name
   \`\`\`

5. **创建 Pull Request**
   - 填写详细的 PR 描述
   - 关联相关的 Issues
   - 等待代码审查

## 发布流程

1. **版本标签**
   \`\`\`bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   \`\`\`

2. **自动构建**
   - GitHub Actions 自动构建多平台二进制文件
   - 自动创建 GitHub Release

## 社区

- 💬 讨论：GitHub Discussions
- 🐛 问题：GitHub Issues
- 📧 邮件：maintainer@example.com

## 许可证

通过贡献代码，您同意您的贡献将在 MIT 许可证下授权。
