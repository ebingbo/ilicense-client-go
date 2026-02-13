# ilicense-client-go

用于客户端离线激活与许可证校验的 Go SDK。

管理端通过 `license-lite` 下发许可证与激活码，业务系统接入 `ilicense-client-go` 后可在本地完成激活、有效期校验和模块授权校验。

## 功能特性

- 离线激活码校验（RSA + SHA-256 签名验证）。
- 许可证状态校验（如 `已过期`、`未激活`）。
- 模块级权限校验。
- 激活码本地持久化，支持启动校验与定时校验。
- 内存中的许可证状态线程安全。

## 运行要求

- Go `1.24+`

## 安装

```bash
go get github.com/ebingbo/ilicense-client-go
```

## 快速开始

```go
package main

import (
	"log"

	"github.com/ebingbo/ilicense-client-go/ilicense"
)

func main() {
	cfg := ilicense.DefaultConfig()
	cfg.PublicKey = "YOUR_RSA_PUBLIC_KEY"
	cfg.ValidateOnStartup = true
	cfg.AllowStartWhenExpired = false

	client := ilicense.NewClient(&cfg)

	if err := client.Init(); err != nil {
		log.Fatalf("license init failed: %v", err)
	}

	_, err := client.Activate("YOUR_ACTIVATION_CODE")
	if err != nil {
		log.Fatalf("license activation failed: %v", err)
	}

	if err := client.CheckModule("m-a"); err != nil {
		log.Fatalf("module access denied: %v", err)
	}
}
```

## 可运行示例

- 基础激活：`go run ./examples/basic`
- 先离线校验再激活：`go run ./examples/offline_activate`

需要的环境变量：

- `ILICENSE_PUBLIC_KEY`
- `ILICENSE_ACTIVATION_CODE`（仅在需要激活时必须提供）

## 配置项

`ilicense.Config`：

- `Enabled`：是否启用许可证校验。
- `PublicKey`：用于校验激活码签名的 RSA 公钥。
- `StoragePath`：激活码本地存储路径。
- `ValidateOnStartup`：是否在启动时加载并校验许可证。
- `AllowStartWhenExpired`：许可证缺失或过期时是否允许启动。
- `Logger`：可选日志注入（`Printf`/`Println`）；默认静默。

## 对外 API

- `NewClient(config *Config) *Client`
- `(*Client).Init() error`
- `(*Client).Activate(code string) (*License, error)`
- `(*Client).CheckLicense() error`
- `(*Client).CheckModule(module string) error`
- `(*Client).GetCurrentLicense() *License`
- `(*Client).IsValid() bool`
- `(*Client).HasModule(module string) bool`

## 错误语义

- `ErrLicenseNotFound`：系统未激活。
- `ErrLicenseExpired`：许可证已过期。
- `ErrSignatureInvalid`：激活码签名校验失败。
- `LicenseError`：底层 IO 或运行时错误包装。

## 安全说明

- 私钥仅保存在管理端（`license-lite`），客户端仅下发公钥。
- 请限制 `StoragePath` 文件写入权限。
- 建议定期轮换签发密钥并支持吊销。
- 本 SDK 不覆盖受攻击客户端上的内存篡改场景。

## 开发

```bash
make test
make vet
make fmt
```

## 版本与发布

- 使用语义化版本（SemVer）。
- 兼容性承诺仅针对 `ilicense` 包。
- 发布策略见 [docs/RELEASING.md](docs/RELEASING.md)。

## 兼容性

- 管理平台：`license-lite`
- 协议：JSON 许可证数据 + RSA 签名，URL-safe Base64 二进制封装

| 组件 | 支持范围 |
| --- | --- |
| Go 运行时 | `1.24+` |
| 公共 API 稳定性 | 当前主版本内 `ilicense` 包 |
| 内部包 | 不承诺兼容性 |

## 贡献

参见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 安全策略

参见 [SECURITY.md](SECURITY.md)。

## 行为准则

参见 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)。

## 许可证

Apache-2.0，详见 [LICENSE](LICENSE)。
