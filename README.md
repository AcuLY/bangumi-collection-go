# bangumi-collection-go

Bangumi 公共用户收藏的匿名、只读 Go 客户端。

> `v0.1.2` 修复收藏记录与嵌套条目元数据的合法类型不同导致整份收藏获取失败的问题，保留收藏记录原有类型及其他字段。

## 范围

- 只读取无需登录即可访问的公共收藏。
- 不接受或发送 access token、Authorization、Cookie，也不提供私有收藏或写接口。
- `Fetch` 自动获取所有计划内页面，去重后按 `(SubjectType, SubjectID, Type)` 排序。
- `FetchPage` 只获取一页，并保留该页的上游顺序。
- 同一个 `Client` 的所有 `Fetch`、`FetchPage` 和重试共享 QPS 与在途请求上限。

模块的 Go 语言与兼容性下限是 Go 1.26.0；本地开发和
发布验收必须使用 Go 1.26.5，默认的 `GOTOOLCHAIN=auto` 会按
`go.mod` 中的 toolchain 声明选择该补丁版本。

安装：

```bash
go get github.com/AcuLY/bangumi-collection-go@v0.1.2
```

## 快速开始

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	collection "github.com/AcuLY/bangumi-collection-go"
)

func main() {
	client := collection.NewClient("example-app/1.0 (contact: you@example.com)")
	subjects, err := client.Fetch(
		context.Background(),
		"lucay126",
		collection.SubjectTypeAnime,
		collection.CollectionTypeDoing,
		collection.CollectionTypeDone,
	)
	if err != nil {
		if errors.Is(err, collection.ErrRateLimited) {
			log.Fatal("Bangumi 请求频率受限")
		}
		log.Fatal(err)
	}

	for _, subject := range subjects {
		name := subject.NameCn
		if name == "" {
			name = subject.Name
		}
		fmt.Printf(
			"%d | %s | 收藏状态=%d | 评分=%d | 更新时间=%s\n",
			subject.SubjectID,
			name,
			subject.Type,
			subject.Rate,
			subject.UpdatedAt.Format("2006-01-02 15:04:05Z07:00"),
		)
	}
}
```

## 客户端配置

```go
client := collection.NewClient(
	userAgent,
	collection.WithConcurrencyLimit(10),
	collection.WithRateLimit(3, 1),
	collection.WithRequestTimeout(30*time.Second),
	collection.WithMaxRetries(3),
	collection.WithRetryInterval(time.Second),
	collection.WithMaxRetryDelay(30*time.Second),
	collection.WithHTTPClient(httpClient),
)
```

配置项：

| Option | 默认值 | 说明 |
|---|---:|---|
| `WithConcurrencyLimit(int)` | `10` | 每个 Client 共享的最大在途请求数 |
| `WithRateLimit(float64, int)` | `3 req/s, burst 1` | 每个 Client 共享的 token bucket |
| `WithRequestTimeout(time.Duration)` | `30s` | 每次 HTTP attempt 的独立超时 |
| `WithMaxRetries(int)` | `3` | 首次 attempt 之后的最大重试次数 |
| `WithRetryInterval(time.Duration)` | `1s` | full-jitter 指数退避基数 |
| `WithMaxRetryDelay(time.Duration)` | `30s` | 本地退避及 `Retry-After` 的共同上限 |
| `WithHTTPClient(*http.Client)` | 安全默认值 | 浅复制客户端并清除 Jar/Client.Timeout、拒绝 redirect |
| `WithEndpoint(string)` | `https://api.bgm.tv` | HTTPS 根地址；HTTP 仅接受 loopback 测试地址 |

既有非认证 option 的无效值继续保留默认值。新的 endpoint、rate、max-retry-delay 配置或 nil `Option` 无效时，Client 会被固定标记为 `ErrInvalidConfiguration`，所有操作都在 transport 前失败。

## API

### 获取完整收藏

```go
subjects, err := client.Fetch(
	ctx,
	"user-id",
	collection.SubjectTypeAnime,
	collection.CollectionTypeDoing,
	collection.CollectionTypeDone,
)
```

`Fetch` 要求至少一个收藏类型。重复类型会被合并；首个 50 条页面决定固定 page plan；任何页面失败都不会返回 partial data。

### 获取单页

```go
page, err := client.FetchPage(
	ctx,
	"user-id",
	collection.SubjectTypeAnime,
	collection.CollectionTypeDone,
	50,
	0,
)
```

`limit` 会被限制到 `1..50`，负 `offset` 会变为 `0`。

### 完整 DTO

`Subject` 表示一条收藏记录：

| 字段 | 含义 |
|---|---|
| `ID` | 兼容别名，始终等于 `SubjectID` |
| `SubjectID` | 条目 ID |
| `SubjectType` | 收藏记录顶层的 `subject_type`，保持上游原值 |
| `Type` | 收藏状态 |
| `Name`, `NameCn` | 原名与中文名；上游省略 `subject` 时为空 |
| `Rate` | 用户评分，`0..10` |
| `Comment` | 用户评论；上游省略或明确为 `null` 时为空字符串 |
| `Tags` | 必填的用户收藏标签；空数组有效，返回值始终为非 nil slice |
| `UpdatedAt` | RFC3339 更新时间 |
| `VolStatus`, `EpStatus` | 卷数与话数进度 |
| `Private` | 上游 private 标记 |

官方条目类型映射为：书籍 `1`、动画 `2`、音乐 `3`、游戏 `4`、三次元 `6`。未打 tag 的原型曾把 `SubjectTypeGame` 与 `SubjectTypeMusic` 的名称写反；首个 `v0.1.0` 契约在发布前纠正为 `SubjectTypeMusic=3`、`SubjectTypeGame=4`，有效原始数值集合不变。收藏类型值保持不变：想看 `1`、看过 `2`、在看 `3`、搁置 `4`、抛弃 `5`。

官方响应中的 `comment` 与嵌套 `subject` 是可选字段。省略或明确为 `null` 的 `comment` 会映射为空字符串，其他已出现的值必须是字符串。省略 `subject` 时，`ID` 仍等于顶层 `SubjectID`，`Name`、`NameCn` 为空；若 `subject` 出现，则必须是完整、非 `null` 的合法值。`tags` 是必填字段，省略、`null` 或类型错误都会作为协议错误返回。

收藏记录的 `subject_type` 与嵌套条目元数据的 `subject.type` 可能不同。
只要两者均为合法条目类型、且两处条目 ID 一致，客户端会保留这条收藏，
返回的 `SubjectType` 仍取顶层 `subject_type`，名称取嵌套条目。
例如顶层为动画 `2`、嵌套条目为三次元 `6` 时，不会因此中止整份收藏获取。
不同条目 ID、非法或缺失的类型、以及顶层类型不符合查询条件仍会返回协议错误。

## 错误处理

使用 `errors.Is` 判断稳定分类，使用 `errors.As` 读取类型化元数据：

```go
var httpErr *collection.HTTPError
if errors.As(err, &httpErr) {
	fmt.Println("HTTP status:", httpErr.StatusCode)
	if errors.Is(err, collection.ErrRateLimited) {
		fmt.Println("Retry-After:", httpErr.RetryAfter)
	}
}

switch {
case errors.Is(err, collection.ErrRateLimited):
	// 429
case errors.Is(err, collection.ErrTimeout):
	// parent deadline 或单次 attempt timeout
}
```

稳定分类包括输入/配置、401、403、404、429、5xx、一般 HTTP 状态、transport、timeout、cancellation、decode、protocol、响应过大和 retry exhaustion。返回错误不会包含 UID、URL/query、headers、response body 或原始 transport 文本。

兼容字段 `HTTPError.Body` 仍然存在，但返回值始终为空字符串。404 同时匹配 `ErrNotFound` 和已弃用的 `ErrInvalidUserID`。

## `v0.1.0` 兼容边界

与未打 tag 的旧代码相比：

- 保留 `NewClient`、`Fetch`、`FetchPage`、收藏类型数值和所有非认证 options。
- **纠正** `SubjectTypeMusic=3`、`SubjectTypeGame=4`；未打 tag 的原型把这两个公开名称写反，使用这两个命名常量的调用会获得纠正后的官方语义。
- **移除** `WithAccessToken` 以及所有 Authorization/Cookie 行为。
- `Subject` 增加完整收藏字段；`ID` 保留并等于 `SubjectID`。
- `Fetch` 不再按 goroutine 完成顺序返回，改为稳定 canonical order。
- `Fetch` 的空收藏类型列表从空成功改为 `ErrNoCollectionTypes`。
- `HTTPError.Body` 不再保存或输出上游 body。
- 新增 endpoint、共享 rate limit 与最大 retry delay options。

这些变更定义首个公开版本的安全边界，后续补丁保持公开 API 兼容。
