# goreasoner

[![Version](https://img.shields.io/badge/version-v0.4.0-blue)](https://github.com/beyondcivic/goreasoner/releases/tag/v0.4.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/doc/devel/release.html)
[![Go Reference](https://pkg.go.dev/badge/github.com/beyondcivic/goreasoner.svg)](https://pkg.go.dev/github.com/beyondcivic/goreasoner)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md)

一个高性能的 Go 语言 RDF/OWL 本体前向推理引擎。本库同时提供命令行工具和 Go 库接口，用于语义推理——解析 Turtle 格式输入，并应用 RDFS/OWL 推理规则，从 TBox（术语/模式）和 ABox（断言/实例）中推导新事实。

## ⚡ 性能（v0.4.0）

> 测试环境：Apple M3 Max。测试数据：Person/Employee/Manager 类层次结构 + friendOf/colleagueOf→knows 属性层次结构 + friendOf 对称属性。

| 规模 | 输入三元组 | goreasoner (Go) | Apache Jena RDFS (Java) | Go vs Jena |
|---|---|---|---|---|
| 1千实体 | 4,451 | **1 ms** | 6 ms | **快 6 倍** |
| 1万实体 | 44,451 | **12 ms** | 36 ms | **快 3 倍** |
| 10万实体 | 444,451 | **173 ms** | 526 ms | **快 3 倍** |
| **100万实体** | **4,444,451** | **2.15 秒** | **7.80 秒** | **快 3.6 倍** |

**吞吐量**：200万–600万 三元组/秒（因规模而异）。零外部依赖（无需 JVM，无需 HTTP 服务）。

### v0.4.0 优化

内部引擎已为性能进行重写，公共 API 保持向后兼容：

| 优化技术 | 说明 | 效果 |
|---|---|---|
| **字符串驻留** | 所有 URI 映射为 uint32 ID；三元组从约 200 字节压缩到 12 字节 | 3–5 倍加速，大幅降低内存 |
| **半朴素求值** | 每轮仅处理新增（增量）三元组；避免每次迭代的全量扫描 | 大图 2–10 倍加速 |
| **开放寻址哈希集** | Robin Hood 哈希做三元组去重；每条约 16 字节（Go map 约 90 字节） | 去重内存减少 50% |
| **组合索引** | SP（主语+谓语）和 PO（谓语+宾语）索引，O(1) 双键查找 | 连接查询 1.2 倍加速 |
| **规则直接写入** | 规则直接写入存储+增量集；无中间切片分配 | 降低 GC 压力 |
| **`AddTriple()` API** | 当已有字符串 URI 时跳过 Turtle 解析 | 编程加载最高 2 倍加速 |

### 与其他方案对比

| 方案 | 100万实体耗时 | 语言 | 外部依赖 |
|---|---|---|---|
| **goreasoner v0.4.0** | **2.15 秒** | Go | **无** |
| Apache Jena RDFS（内嵌） | 7.80 秒 | Java | JVM |
| Apache Jena Fuseki（HTTP） | ~10 秒 | Java | JVM + HTTP 服务 |
| rdflib + OWL-RL | >60 秒（估计） | Python | Python 运行时 |

自行运行基准测试：

```bash
go test -run TestBenchmarkReport ./pkg/reasoner/ -v -timeout=600s
```

## 概述

语义网依靠本体和推理从显式事实中推导隐式知识。本工具通过以下方式简化 RDF/OWL 推理：

- **解析 Turtle 格式 RDF 数据**，完整支持前缀和 IRI
- **应用前向推理规则**，基于 RDFS 和 OWL 规范
- **推导新知识**，通过传递性、对称性和层级继承推理
- **提供高效三元组存储**，支持索引快速查询
- **提供命令行和库两种接口**，满足不同集成需求

本项目同时提供命令行工具和 Go 库，用于对 RDF/OWL 数据进行语义推理。

## 主要特性

- ✅ **Turtle 解析器**：自定义解析器（无外部依赖），支持前缀、IRI、空白节点和字面量
- ✅ **前向推理**：完整的 RDFS/OWL 推理规则实现
- ✅ **类层次结构**：传递性子类关系和类型继承
- ✅ **属性推理**：领域/值域推理和属性层次结构
- ✅ **OWL 支持**：等价类、同一性推理、逆属性和传递属性
- ✅ **Datalog 推理**：内置 Datalog 解析器和求值器，支持规则、事实和布尔查询
- ✅ **多种输出格式**：支持 N-Triples 和 Datalog 输出格式
- ✅ **命令行和库**：同时提供命令行工具和 Go 库接口
- ✅ **跨平台**：支持 Linux、macOS 和 Windows

## 快速开始

### 前置条件

- Go 1.24 或更高版本
- Nix 2.25.4 或更高版本（可选但推荐）
- PowerShell v7.5.1 或更高版本（用于构建）

### 安装

#### 方式一：从源码安装

1. 克隆仓库：

```bash
git clone https://github.com/beyondcivic/goreasoner.git
cd goreasoner
```

2. 构建应用：

```bash
go build -o goreasoner .
```

#### 方式二：使用 Nix（推荐）

1. 克隆仓库：

```bash
git clone https://github.com/beyondcivic/goreasoner.git
cd goreasoner
```

2. 使用 Nix flakes 准备环境：

```bash
nix develop
```

3. 构建应用：

```bash
./build.ps1
```

#### 方式三：Go Install

```bash
go install github.com/beyondcivic/goreasoner@latest
```

## 快速上手

### 命令行界面

`goreasoner` 工具提供对 RDF 数据进行语义推理的命令：

```bash
# 对 RDF 数据执行前向推理（N-Triples 输出）
goreasoner run instances.ttl schema.ttl -o results.nt

# 使用 Datalog 格式输出执行前向推理
goreasoner run instances.ttl schema.ttl --outputType=datalog -o results.dl

# 查询 Datalog 程序
goreasoner dlquery results.dl "?- type(myTesla, Vehicle)."

# 显示版本信息
goreasoner version
```

### Go 库使用

```go
package main

import (
    "fmt"
    "log"

    "github.com/beyondcivic/goreasoner/pkg/reasoner"
)

func main() {
    tbox := `
@prefix ex: <http://example.org/> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .

ex:Car rdfs:subClassOf ex:Vehicle .
ex:Vehicle rdfs:subClassOf ex:Transport .
`

    abox := `
@prefix ex: <http://example.org/> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .

ex:myCar rdf:type ex:Car .
`

    // 获取所有三元组（包括推理得出的）
    triples, err := reasoner.ForwardReason(abox, tbox)
    if err != nil {
        log.Fatalf("Error running forward reasoning: %v", err)
    }

    fmt.Printf("Inferred %d total triples\n", len(triples))
    for _, t := range triples {
        fmt.Println(t)
    }
}
```

## 详细命令参考

### `run` - 执行前向推理

使用 TBox（模式）和 ABox（实例）对 RDF 数据执行前向推理。

```bash
goreasoner run [ABOX_FILE] [TBOX_FILE] [OPTIONS]
```

**选项：**

- `-o, --output`：输出文件路径（默认：`[abox_filename]_inferred.nt`）
- `--outputType`：输出格式 - `ntriple` 或 `datalog`（默认：`ntriple`）

**示例：**

```bash
# 基本推理（N-Triples 输出）
goreasoner run instances.ttl schema.ttl

# 自定义输出路径
goreasoner run instances.ttl schema.ttl -o my-results.nt

# Datalog 输出格式
goreasoner run instances.ttl schema.ttl --outputType=datalog

# Datalog 输出到自定义文件
goreasoner run instances.ttl schema.ttl --outputType=datalog -o results.dl
```

### `dlquery` - 查询 Datalog 程序

对 Datalog 程序（事实和规则）执行布尔查询。

```bash
goreasoner dlquery [DATALOG_FILE] [QUERY]
```

**参数：**

- `DATALOG_FILE`：包含 Datalog 事实和规则的 `.dl` 文件路径
- `QUERY`：`?- predicate(args).` 格式的 Datalog 查询字符串

**示例：**

```bash
# 查询一个基础事实
goreasoner dlquery data.dl "?- type(myTesla, Car)."

# 查询一个推导事实
goreasoner dlquery data.dl "?- Ancestor(john, jane)."

# 使用变量查询（如果存在任何绑定则返回 true）
goreasoner dlquery data.dl "?- type(X, Vehicle)."
```

### `version` - 显示版本信息

显示版本、构建信息和系统详情。

```bash
goreasoner version
```

## 支持的推理规则

本推理引擎实现了完整的 RDFS/OWL 推理规则：

| 规则类别                    | 说明                                        | 示例                       |
| -------------------------------- | ------------------------------------------ | ------------------------- |
| **rdfs:subClassOf 传递性** | 若 A ⊑ B 且 B ⊑ C，则 A ⊑ C             | Car ⊑ Vehicle ⊑ Transport |
| **rdf:type 继承**         | 若 x:A 且 A ⊑ B，则 x:B                 | myCar:Car → myCar:Vehicle |
| **rdfs:domain 推理**        | 若 P domain C 且 x P y，则 x:C          | hasOwner domain Person    |
| **rdfs:range 推理**         | 若 P range C 且 x P y，则 y:C           | hasAge range Integer      |
| **rdfs:subPropertyOf**           | 属性层次结构推理               | drives ⊑ operates         |
| **owl:equivalentClass**          | 类等价（对称/传递）   | Vehicle ≡ Automobile      |
| **owl:sameAs**                   | 个体同一性（对称/传递） | person1 ≡ person2         |
| **owl:inverseOf**                | 逆属性推理                 | owns ⟷ isOwnedBy          |
| **owl:TransitiveProperty**       | 传递属性链                 | locatedIn 传递性    |
| **owl:SymmetricProperty**        | 对称属性推理               | marriedTo 对称性        |

## 输出格式

本工具支持两种推理结果输出格式：

### N-Triples 格式（默认）

标准 N-Triples 格式，使用完整 IRI：

```
<http://example.org/myCar> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://example.org/Car> .
<http://example.org/myCar> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://example.org/Vehicle> .
<http://example.org/myCar> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://example.org/Transport> .
<http://example.org/Car> <http://www.w3.org/2000/01/rdf-schema#subClassOf> <http://example.org/Vehicle> .
<http://example.org/Vehicle> <http://www.w3.org/2000/01/rdf-schema#subClassOf> <http://example.org/Transport> .
```

### Datalog 格式

Datalog 事实格式，使用简化名称，适用于 Datalog 推理系统：

```
type(myCar, Car).
type(myCar, Vehicle).
type(myCar, Transport).
subClassOf(Car, Vehicle).
subClassOf(Vehicle, Transport).
```

Datalog 格式将 RDF 三元组 `<subject> <predicate> <object>` 转换为事实 `predicate(subject, object).`，并通过提取本地名称来简化 IRI。

## Datalog 推理

除了 RDF/OWL 前向推理之外，goreasoner 还内置了 Datalog 求值器，支持事实、带变量的规则、递归规则和布尔查询。求值器使用**朴素自底向上（前向链）求值**和不动点计算，在回答查询之前推导所有可能的事实。

### Datalog 语法

#### 事实

以句号结尾的基础原子：

```
Parent(john, mary).
Type(myTesla, Car).
```

#### 规则

包含头部和一个或多个体原子的 Horn 子句，用 `:-` 分隔：

```
Ancestor(X, Y) :- Parent(X, Y).
Ancestor(X, Z) :- Parent(X, Y), Ancestor(Y, Z).
```

#### 变量

变量通过以下约定识别：

- 单个大写字母：`X`、`Y`、`Z`
- 全大写标识符：`VAR_X`、`PERSON`、`NODE1`
- `?` 前缀标识符：`?x`、`?person`

其他所有内容视为常量。

#### 查询

查询使用 `?-` 前缀，返回布尔值（如果存在匹配事实则为 true）：

```
?- Ancestor(john, jane).
?- type(X, Vehicle).
```

#### 注释

支持 Prolog 风格（`%`）和 C 风格（`//`）行注释：

```
% 这是注释
Parent(john, mary).  // 这也是注释
```

### Datalog Go 库使用

```go
package main

import (
    "fmt"
    "log"

    "github.com/beyondcivic/goreasoner/pkg/reasoner"
)

func main() {
    program := `
Parent(john, mary).
Parent(mary, jane).
Ancestor(X, Y) :- Parent(X, Y).
Ancestor(X, Z) :- Parent(X, Y), Ancestor(Y, Z).
`

    result, err := reasoner.DLQuery(program, "?- Ancestor(john, jane).")
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    fmt.Println(result) // true
}
```

如需更细粒度的控制，可以直接使用底层 API：

```go
program, err := reasoner.ParseDatalog(input)
if err != nil {
    log.Fatal(err)
}

// 通过不动点求值推导所有事实
derivedFacts := program.Reason()

// 解析并求值查询
query, err := reasoner.ParseQuery("?- Ancestor(john, jane).")
if err != nil {
    log.Fatal(err)
}

satisfied := program.EvaluateQuery(query, derivedFacts)
```

### Datalog API 参考

#### `DLQuery(datalogContent, queryStr string) (bool, error)`

Datalog 查询的主 API 函数。解析程序，运行推理至不动点，然后求值查询。

#### `ParseDatalog(input string) (*DatalogProgram, error)`

将 Datalog 程序字符串解析为包含事实和规则的 `DatalogProgram`。

#### `ParseQuery(s string) (DLAtom, error)`

将查询字符串（带或不带 `?-` 前缀）解析为 `DLAtom`。

#### `(*DatalogProgram) Reason() []DLAtom`

运行前向链求值直到不再推导出新事实。返回所有基础事实（原始的和推理的）。

#### `(*DatalogProgram) EvaluateQuery(query DLAtom, derivedFacts []DLAtom) bool`

检查查询是否匹配任何推导事实。查询中的变量作为通配符。

### Datalog 限制

Datalog 求值器设计用于简单的正 Datalog 程序。请注意以下限制：

- **不支持否定**：不支持规则体中的失败即否定（`not`、`\+`）。规则中只能出现正原子。
- **不支持内置比较或算术**：不支持 `!=`、`<`、`>` 等运算符和算术表达式。所有项都是符号常量或变量。
- **仅支持布尔查询**：`DLQuery` 返回 `true`/`false`，不返回变量绑定。例如，查询 `?- Ancestor(john, X).` 会告诉你是否存在祖先，但不会枚举它们。
- **不检查规则安全性**：解析器接受头部包含不在体中出现的变量的规则（如 `Foo(X) :- Bar(Y).`）。这类规则不会产生错误结果（未接地的头部会被静默丢弃），但不会发出警告。
- **事实无索引**：求值器在匹配规则体原子时对所有事实进行线性扫描。这对于中小规模程序足够用，但在事实达到数千条时可能变慢。
- **不支持聚合或约束**：不支持扩展 Datalog 系统中的 `count`、`min`、`max` 或完整性约束等功能。

## 示例

### 示例 1：基本类层次结构推理

```bash
# 创建模式文件 (tbox.ttl)
echo '@prefix ex: <http://example.org/> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .

ex:SportsCar rdfs:subClassOf ex:Car .
ex:Car rdfs:subClassOf ex:Vehicle .' > tbox.ttl

# 创建实例文件 (abox.ttl)
echo '@prefix ex: <http://example.org/> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .

ex:myFerrari rdf:type ex:SportsCar .' > abox.ttl

# 执行推理（N-Triples 输出）
goreasoner run abox.ttl tbox.ttl -o results.nt

# 执行推理（Datalog 输出）
goreasoner run abox.ttl tbox.ttl --outputType=datalog -o results.dl

# 查看结果
cat results.nt
cat results.dl
```

### 示例 2：复杂本体推理

给定一个包含领域/值域限制、属性层次结构和等价类的本体，推理引擎将根据 RDFS/OWL 语义推导出所有隐式知识。

## API 参考

### 核心函数

#### `ForwardReason(abox, tbox string) ([]string, error)`

执行 RDF 数据前向推理的主 API 函数。

**参数：**

- `abox`：包含实例数据（断言）的 Turtle 字符串
- `tbox`：包含模式/本体定义的 Turtle 字符串

**返回：**

- `[]string`：所有三元组（包括推理的）的 N-Triples 格式列表
- `error`：任何解析或处理错误

#### `ForwardReasonWithDetails(abox, tbox string) (*ReasoningResult, error)`

返回带有原始/推理三元组分离的详细推理结果。

**参数：**

- `abox`：包含实例数据的 Turtle 字符串
- `tbox`：包含模式定义的 Turtle 字符串

**返回：**

- `*ReasoningResult`：详细结果结构
- `error`：任何解析或处理错误

### 直接使用推理引擎

```go
r := reasoner.NewReasoner()

// 加载 TBox（模式）
if err := r.LoadTurtle(tbox); err != nil {
    log.Fatalf("Error loading TBox: %v", err)
}

// 加载 ABox（实例）
if err := r.LoadTurtle(abox); err != nil {
    log.Fatalf("Error loading ABox: %v", err)
}

// 执行推理
inferredCount := r.RunForwardReasoning()
fmt.Printf("Inferred %d new triples\n", inferredCount)

// 按模式查询（使用 "" 作为通配符）
vehicles := r.Query("", reasoner.RDFType, "http://example.org/Vehicle")
for _, t := range vehicles {
    fmt.Printf("%s is a Vehicle\n", t.Subject)
}

// 获取特定实例的所有类型
types := r.GetInferredTypes("http://example.org/myCar")
fmt.Printf("Types of myCar: %v\n", types)
```

### 数据结构

#### `ReasoningResult`

表示详细推理结果：

```go
type ReasoningResult struct {
    OriginalTriples []string // 输入中的三元组
    InferredTriples []string // 新推理的三元组
    AllTriples      []string // 所有三元组合并
    OriginalCount   int      // 原始三元组数量
    InferredCount   int      // 推理三元组数量
    TotalCount      int      // 三元组总数
}
```

#### `Triple`

表示一个 RDF 三元组：

```go
type Triple struct {
    Subject   string
    Predicate string
    Object    string
}
```

#### `Reasoner`

主推理引擎结构体及其方法：

```go
type Reasoner struct {
    store  *TripleStore
    rules  []Rule
    parser *TurtleParser
}
```

### 推理引擎方法

| 方法                                                  | 说明                                                                |
| --------------------------------------------------- | ----------------------------------------------------------------- |
| `NewReasoner() *Reasoner`                           | 创建使用默认规则的新推理引擎                          |
| `NewReasonerWithRules(rules []Rule) *Reasoner`      | 创建使用自定义规则的推理引擎                               |
| `LoadTurtle(content string) error`                  | 解析并加载 Turtle 内容                                     |
| `RunForwardReasoning() int`                         | 应用所有规则直到不动点，返回推理三元组数量 |
| `GetAllTriples() []string`                          | 获取所有三元组的 N-Triples 字符串                              |
| `GetInferredTypes(subject string) []string`         | 获取主语的所有 rdf:type 值                             |
| `Query(subject, predicate, object string) []Triple` | 模式匹配查询（使用 "" 作为通配符）                       |
| `GetStore() *TripleStore`                           | 访问底层三元组存储                                |

## 架构

本库由以下几个核心组件组成：

### 核心包（`pkg/reasoner`）

- **前向推理引擎**：基于规则的推理与不动点计算
- **Turtle 解析器**：完整的 Turtle 格式解析器，支持前缀
- **三元组存储**：带索引的内存存储，用于高效查询
- **规则系统**：模块化的 RDFS/OWL 推理规则
- **Datalog 求值器**：Datalog 程序的解析器和朴素自底向上推理引擎
- **查询接口**：模式匹配、类型推理和 Datalog 查询

### 命令行界面（`cmd/goreasoner`）

- **基于 Cobra 的 CLI**，包含推理操作的子命令
- **文件 I/O 处理**，支持 Turtle 输入和 N-Triples 输出
- **完善的帮助系统**，包含详细的使用示例
- **灵活的输出选项**和错误处理

### 贡献

1. Fork 本仓库
2. 创建功能分支：`git checkout -b feature/new-feature`
3. 修改代码并添加测试
4. 确保所有测试通过：`go test ./...`
5. 提交更改：`git commit -am 'Add new feature'`
6. 推送到分支：`git push origin feature/new-feature`
7. 提交 pull request

### 测试

运行测试套件：

```bash
go test ./...
```

运行带覆盖率的测试：

```bash
go test -cover ./...
```

## 构建环境

### 使用 Nix（推荐）

使用 Nix flakes 设置构建环境：

```bash
nix develop
```

### 手动构建

查看 `build.ps1` 中的构建参数：

```bash
# 构建带有版本信息的静态二进制文件
$env:CGO_ENABLED = "1"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
```

然后运行：

```bash
./build.ps1
```

或手动构建：

```bash
go build -o goreasoner .
```

## 文件结构

```
goreasoner/
├── go.mod                    # 模块定义
├── main.go                   # 主入口
├── build.ps1                 # 构建脚本
├── flake.nix                 # Nix flake 配置
├── cmd/
│   └── goreasoner/
│       ├── main.go           # 命令行界面
│       └── commands.go       # 命令定义
├── pkg/
│   ├── reasoner/
│   │   ├── core.go           # 主 API 和 Reasoner 类型
│   │   ├── parser.go         # Turtle 格式解析器
│   │   ├── store.go          # 内存三元组存储
│   │   ├── rules.go          # 前向推理规则
│   │   ├── datalog.go        # Datalog 解析器和推理引擎
│   │   ├── utils.go          # 工具函数
│   │   └── error.go          # 错误处理
│   └── version/
│       └── version.go        # 版本信息
└── docs/
    └── docs-gen.go           # 文档生成器
```

## 许可证

本项目基于 MIT 许可证授权——详见 [LICENSE](LICENSE) 文件。
