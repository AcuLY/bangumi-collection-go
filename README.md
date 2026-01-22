# bangumi-collection-go

获取 [Bangumi](https://bgm.tv) 用户收藏列表的 Go 客户端

## 安装

```bash
go get github.com/AcuLY/bangumi-collection-go
```

## 快速开始

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AcuLY/bangumi-collection-go"
)

func main() {
	client := collection.NewClient("AcuL/bangumi-collection-go")

	subjects, err := client.Fetch(
		context.Background(),
		"lucay126",                          
		collection.SubjectTypeAnime,    
		collection.CollectionTypeDoing,
		collection.CollectionTypeDone, 
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("共 %d 部动画\n\n", len(subjects))
	for _, s := range subjects {
		name := s.NameCn
		if name == "" {
			name = s.Name
		}
		fmt.Printf("ID: %d | %s | 评分: %d | 标签: %v\n", s.ID, name, s.Rate, s.Tags)
	}
}
```

## API

### 创建客户端

```go
client := collection.NewClient(userAgent, options...)
```

`userAgent` 必填，参考 [Bangumi API UserAgent 建议](https://github.com/bangumi/api/blob/master/docs-raw/user%20agent.md)。

### 可选配置

```go
collection.WithAccessToken(token)              // 设置访问令牌
collection.WithConcurrencyLimit(10)            // 并发数限制（默认 10）
collection.WithRequestTimeout(30*time.Second)  // 请求超时（默认 30s）
collection.WithMaxRetries(3)                   // 最大重试次数（默认 3）
collection.WithRetryInterval(time.Second)      // 重试间隔（默认 1s，指数退避）
collection.WithHTTPClient(client)              // 自定义 HTTP 客户端
```

### 获取全部收藏

```go
subjects, err := client.Fetch(ctx, userID, subjectType, collectionTypes...)
```

支持同时获取多个收藏类型：

```go
subjects, err := client.Fetch(
    ctx,
    "sai",
    collection.SubjectTypeAnime,
    collection.CollectionTypeDoing,  // 在看
    collection.CollectionTypeDone,   // 看过
)
```

### 分页获取

```go
page, err := client.FetchPage(ctx, userID, subjectType, collectionType, limit, offset)
// page.Data   - 当前页数据
// page.Total  - 总数
// page.Limit  - 每页数量
// page.Offset - 当前偏移
```

### 条目类型

| 常量 | 值 | 说明 |
|------|---|------|
| `SubjectTypeBook` | 1 | 书籍 |
| `SubjectTypeAnime` | 2 | 动画 |
| `SubjectTypeGame` | 3 | 游戏 |
| `SubjectTypeMusic` | 4 | 音乐 |
| `SubjectTypeReal` | 6 | 三次元 |

### 收藏类型

| 常量 | 值 | 说明 |
|------|---|------|
| `CollectionTypeWish` | 1 | 想看 |
| `CollectionTypeDone` | 2 | 看过 |
| `CollectionTypeDoing` | 3 | 在看 |
| `CollectionTypeOnHold` | 4 | 搁置 |
| `CollectionTypeDropped` | 5 | 抛弃 |

### 错误类型

| 错误 | 说明 |
|------|------|
| `ErrInvalidUserID` | 用户 ID 不存在或无效 |
| `ErrUnauthorized` | 访问令牌无效或已过期 |
| `ErrForbidden` | 没有权限访问该资源 |
| `ErrRateLimited` | 请求过于频繁 |
| `ErrServerError` | 服务器内部错误 |
| `ErrEmptyUserID` | 用户 ID 为空 |
| `ErrNoCollectionTypes` | 未指定收藏类型 |
| `*NetworkError` | 网络错误 |
| `*HTTPError` | 非预期的 HTTP 状态码 |
