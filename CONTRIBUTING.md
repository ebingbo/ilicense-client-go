# 贡献指南

感谢你参与 `ilicense-client-go`。

## 本地开发准备

1. 安装 Go `1.24+`。
2. 克隆仓库。
3. 执行基础检查：

```bash
make fmt
make test
make vet
```

## 分支与提交规范

- 从 `main` 拉取功能分支。
- 单个 PR 保持聚焦，避免一次改动过大。
- 提交信息应清晰描述变更意图。

## 提交 PR 前检查

- 行为变更已补充或更新测试。
- API 或配置变更已同步更新 README/文档。
- 本地 `make test` 与 `make vet` 已通过。
- 若涉及破坏性改动，已先在 issue 或 PR 中说明。

## 代码风格

- 遵循 Go 惯例和 `gofmt`。
- 导出 API 需要保持稳定并补充注释。
- 避免输出敏感许可证数据。

## Bug 反馈建议

请在 GitHub Issue 中提供：

- 复现步骤。
- 期望行为与实际行为。
- Go 版本与操作系统。
- 最小复现代码。

## 安全问题

安全漏洞请勿公开提交 issue，请按 [SECURITY.md](SECURITY.md) 私下报告。
