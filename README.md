# ginm & gox

> 基于 Go 泛型的工具包

## 包结构

本项目包含两个独立的包：

| 包 | 导入路径 | 定位 |
|----|----------|------|
| **ginm** | `github.com/lwmacct/251219-go-pkg-ginm/pkg/ginm` | Gin 框架类型安全辅助 |
| **gox** | `github.com/lwmacct/251219-go-pkg-ginm/pkg/gox` | 通用泛型工具（无框架依赖） |

## ginm - Gin 辅助包

**ginm 是什么：**

- Gin 框架的类型安全辅助函数
- 利用 Go 1.18+ 泛型消除 `interface{}` 和类型断言
- 减少样板代码，提升开发体验

**ginm 不是什么：**

- 不是框架，不强制特定架构
- 不包含认证、数据库、缓存等业务组件

### ginm 核心特性

| 模块                 | 功能                                               |
| -------------------- | -------------------------------------------------- |
| **泛型响应**         | `Response[T]`, `PageResponse[T]` 统一 API 响应格式 |
| **泛型绑定**         | `Bind[T]`, `BindJSON[T]` 类型安全的请求绑定        |
| **Handler 包装**     | `Wrap`, `WrapJSON` 自动绑定 + 错误处理             |
| **类型安全上下文**   | `ContextKey[T]`, `Get[T]`, `Set[T]` 消除类型断言   |
| **RESTful Resource** | `Resource[T,ID,CI,UI,LQ]` 五类型参数完整 CRUD      |
| **中间件链**         | `HandlerChain` 链式中间件组合                      |

## gox - 通用工具包

**gox 是什么：**

- 纯泛型工具函数，无任何框架依赖
- 函数式编程、单子类型、数值/集合/转换工具

### gox 核心特性

| 模块             | 功能                                             |
| ---------------- | ------------------------------------------------ |
| **函数式工具**   | `Map`, `Filter`, `Reduce`, `GroupBy`, `Chunk`    |
| **Result**       | Rust 风格 `Result[T]` 显式错误处理               |
| **Optional**     | `Optional[T]` 空值安全                           |
| **数值工具**     | `Sum`, `Max`, `Min`, `Average`, `Clamp`, `Range` |
| **集合操作**     | `Intersect`, `Union`, `Difference`, `Partition`  |
| **类型转换**     | `ParseInt`, `ParseFloat`, `ParseBool` + Result   |
| **错误聚合**     | `MultiError` 批量操作错误收集                    |
| **指针工具**     | `Ptr`, `Val`, `Coalesce`                         |

## Quick Start

### Init Development Environment

```shell
pre-commit install
```

### List All Available Tasks

```shell
task -a
```

## Related Links

- Use [Taskfile](https://taskfile.dev) to manage the project's CLI
- Use [Pre-commit](https://pre-commit.com/) to manage and maintain multi-language pre-commit hooks
