# 变更日志

本文件记录该项目的重要变更。

格式参考 Keep a Changelog，版本遵循 Semantic Versioning。

## [Unreleased]

### 新增

- 面向 GitHub 开源发布的第一版基线能力。
- 治理文档：贡献指南、安全策略、行为准则。
- GitHub Actions CI 与 Issue/PR 模板。
- 可运行示例：`examples/basic`、`examples/offline_activate`。
- `ilicense` 包级和导出 API 的 GoDoc 注释。
- 发布策略文档：`docs/RELEASING.md`。

### 变更

- README 重写并补充使用与安全说明。
- 模块路径和 import 调整为 GitHub 发布路径。
- 模块授权匹配从子串匹配改为精确匹配。
- 内部校验逻辑迁移到 `internal/core`，不再作为公共 API 暴露。
- SDK 日志改为可注入（`Config.Logger`），默认静默。
