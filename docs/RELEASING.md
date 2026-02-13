# 发布流程

本项目使用语义化版本（SemVer）。

## 版本规则

- `MAJOR`：不兼容的公共 API 变更。
- `MINOR`：向后兼容的新功能。
- `PATCH`：向后兼容的问题修复。

公共 API 仅指 `ilicense` 包。
`internal/` 下的实现不承诺兼容。

## 发布步骤

1. 确认 `main` 分支 CI 全绿。
2. 确认 `CHANGELOG.md` 的 `[Unreleased]` 已补齐。
3. 提交发布前整理 commit（如 changelog 定稿）。
4. 在 `main` 打 tag：`vX.Y.Z`。
5. 在 GitHub 发布 Release，附变更摘要与升级说明。

## 兼容性约束

- 在 `v1` 主版本内，避免破坏 `ilicense` 导出标识符。
- 必须破坏兼容时，需提供迁移说明并提升 MAJOR 版本。

## 安全补丁发布

若出现高危漏洞，优先发布修复版本，并在 changelog 中单独说明。
