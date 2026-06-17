# goreasoner

[![Version](https://img.shields.io/badge/version-v0.4.0-blue)](https://github.com/beyondcivic/goreasoner/releases/tag/v0.4.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/doc/devel/release.html)
[![Go Reference](https://pkg.go.dev/badge/github.com/beyondcivic/goreasoner.svg)](https://pkg.go.dev/github.com/beyondcivic/goreasoner)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md)

一個高效能的 Go 語言 RDF/OWL 本體前向推理器實作。本函式庫提供命令列介面和 Go 函式庫兩種方式進行語義推理，解析 Turtle 格式輸入並套用 RDFS/OWL 推理規則，從 TBox（術語/綱要）和 ABox（斷言/實例）中推導新事實。

## ⚡ 效能 (v0.4.0)

> 在 Apple M3 Max 上進行基準測試。測試資料：Person/Employee/Manager 類別階層結構 + friendOf/colleagueOf→knows 屬性階層結構 + symmetric friendOf。

| 規模 | 輸入三元組 | goreasoner (Go) | Apache Jena RDFS (Java) | Go vs Jena |
|---|---|---|---|---|
| 1K 實體 | 4,451 | **1 ms** | 6 ms | **快 6 倍** |
| 10K 實體 | 44,451 | **12 ms** | 36 ms | **快 3 倍** |
| 100K 實體 | 444,451 | **173 ms** | 526 ms | **快 3 倍** |
| **1M 實體** | **4,444,451** | **2.15 s** | **7.80 s** | **快 3.6 倍** |

**吞吐量**：依規模不同，每秒處理 2M–6M 個三元組。零外部相依（無 JVM，無 HTTP 服務）。

### v0.4.0 最佳化

內部引擎為效能進行了重寫，公開 API 保持向後相容：

| 技術 | 描述 | 影響 |
|---|---|---|
| **字串駐留** | 所有 URI 對映為 uint32 ID；三元組僅 12 位元組而非約 200 位元組 | 3–5 倍加速，大幅減少記憶體 |
| **半樸素求值** | 每輪僅處理新增（增量）三元組；避免每次迭代全量掃描 | 大規模圖 2–10 倍加速 |
| **開放定址雜湊集** | 使用 Robin Hood 雜湊進行三元組去重；每條目約 16 位元組 vs Go map 的約 90 位元組 | 去重記憶體減少 50% |
| **複合索引** | SP（主詞+謂詞）和 PO（謂詞+賓詞）索引實現 O(1) 雙鍵查找 | 連接操作加速 1.2 倍 |
| **直接發射規則** | 規則直接寫入儲存+增量；無中間切片配置 | 降低 GC 壓力 |
| **`AddTriple()` API** | 當已有字串 URI 時略過 Turtle 剖析 | 程式化載入最高 2 倍加速 |

### 與其他方案的比較

| 方案 | 1M 實體 | 語言 | 外部相依 |
|---|---|---|---|
| **goreasoner v0.4.0** | **2.15 s** | Go | **無** |
| Apache Jena RDFS（嵌入式） | 7.80 s | Java | JVM |
| Apache Jena Fuseki（HTTP） | ~10 s | Java | JVM + HTTP 伺服器 |
| rdflib + OWL-RL | >60 s（估計） | Python | Python 執行環境 |

自行執行基準測試：

```bash
go test -run TestBenchmarkReport ./pkg/reasoner/ -v -timeout=600s
```

## 概述

語義網仰賴本體和推理從顯式事實中推導隱含知識。本工具透過以下方式簡化 RDF/OWL 推理：

- **剖析 Turtle 格式 RDF 資料**，完整支援前綴和 IRI
- **套用前向推理規則**，基於 RDFS 和 OWL 規範
- **推導新知識**，透過傳遞性、對稱性和階層推理
- **提供高效的三元組儲存**，支援索引查找以實現快速查詢
- **同時提供命令列和函式庫介面**，滿足不同整合需求

本專案同時提供命令列介面和 Go 函式庫，用於對 RDF/OWL 資料進行語義推理。

## 主要特性

- ✅ **Turtle 剖析器**：自訂剖析器（無外部相依），支援前綴、IRI、空白節點和字面值
- ✅ **前向推理**：完整的 RDFS/OWL 推理規則實作
- ✅ **類別階層結構**：傳遞性子類別關係和類型繼承
- ✅ **屬性推理**：領域/值域推理和屬性階層結構
- ✅ **OWL 支援**：等價類別、相同個體推理、逆屬性和傳遞屬性
- ✅ **Datalog 推理**：內建 Datalog 剖析器和求值器，支援規則、事實和布林查詢
- ✅ **多種輸出格式**：支援 N-Triples 和 Datalog 輸出格式
- ✅ **命令列和函式庫**：同時提供命令列工具和 Go 函式庫介面
- ✅ **跨平臺**：支援 Linux、macOS 和 Windows

## 快速開始

### 前置條件

- Go 1.24 或更高版本
- Nix 2.25.4 或更高版本（可選但建議使用）
- PowerShell v7.5.1 或更高版本（用於建置）

### 安裝

#### 方式 1：從原始碼安裝

1. 複製儲存庫：

```bash
git clone https://github.com/beyondcivic/goreasoner.git
cd goreasoner
```

2. 建置應用程式：

```bash
go build -o goreasoner .
```

#### 方式 2：使用 Nix（建議）

1. 複製儲存庫：

```bash
git clone https://github.com/beyondcivic/goreasoner.git
cd goreasoner
```

2. 使用 Nix flakes 準備環境：

```bash
nix develop
```

3. 建置應用程式：

```bash
./build.ps1
```

#### 方式 3：Go Install

```bash
go install github.com/beyondcivic/goreasoner@latest
```

## 快速上手

### 命令列介面

`goreasoner` 工具提供用於 RDF 資料語義推理的命令：

```bash
# 對 RDF 資料執行前向推理（N-Triples 輸出）
goreasoner run instances.ttl schema.ttl -o results.nt

# 以 Datalog 格式輸出進行前向推理
goreasoner run instances.ttl schema.ttl --outputType=datalog -o results.dl

# 查詢 Datalog 程式
goreasoner dlquery results.dl "?- type(myTesla, Vehicle)."

# 顯示版本資訊
goreasoner version
```

### Go 函式庫使用

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

    // 取得所有三元組（包括推理得出的）
    triples, err := reasoner.ForwardReason(abox, tbox)
    if err != nil {
        log.Fatalf("執行前向推理時出錯: %v", err)
    }

    fmt.Printf("共推理出 %d 個三元組\n", len(triples))
    for _, t := range triples {
        fmt.Println(t)
    }
}
```

## 詳細命令參考

### `run` - 執行前向推理

對 RDF 資料使用 TBox（綱要）和 ABox（實例）執行前向推理。

```bash
goreasoner run [ABOX_FILE] [TBOX_FILE] [OPTIONS]
```

**選項：**

- `-o, --output`：輸出檔案路徑（預設：`[abox_filename]_inferred.nt`）
- `--outputType`：輸出格式 - `ntriple` 或 `datalog`（預設：`ntriple`）

**範例：**

```bash
# 基本推理（N-Triples 輸出）
goreasoner run instances.ttl schema.ttl

# 自訂輸出路徑
goreasoner run instances.ttl schema.ttl -o my-results.nt

# Datalog 輸出格式
goreasoner run instances.ttl schema.ttl --outputType=datalog

# Datalog 輸出到自訂檔案
goreasoner run instances.ttl schema.ttl --outputType=datalog -o results.dl
```

### `dlquery` - 查詢 Datalog 程式

對 Datalog 程式（事實和規則）求值布林查詢。

```bash
goreasoner dlquery [DATALOG_FILE] [QUERY]
```

**參數：**

- `DATALOG_FILE`：包含 Datalog 事實和規則的 `.dl` 檔案路徑
- `QUERY`：`?- predicate(args).` 格式的 Datalog 查詢字串

**範例：**

```bash
# 查詢基礎事實
goreasoner dlquery data.dl "?- type(myTesla, Car)."

# 查詢推導事實
goreasoner dlquery data.dl "?- Ancestor(john, jane)."

# 帶變數查詢（如果存在任何繫結則傳回 true）
goreasoner dlquery data.dl "?- type(X, Vehicle)."
```

### `version` - 顯示版本資訊

顯示版本、建置資訊和系統詳情。

```bash
goreasoner version
```

## 支援的推理規則

推理器實作了完整的 RDFS/OWL 推理規則：

| 規則類別                          | 描述                                       | 範例                      |
| -------------------------------- | ------------------------------------------ | ------------------------- |
| **rdfs:subClassOf 傳遞性**        | 如果 A ⊑ B 且 B ⊑ C，則 A ⊑ C              | Car ⊑ Vehicle ⊑ Transport |
| **rdf:type 繼承**                | 如果 x:A 且 A ⊑ B，則 x:B                   | myCar:Car → myCar:Vehicle |
| **rdfs:domain 推理**             | 如果 P domain C 且 x P y，則 x:C            | hasOwner domain Person    |
| **rdfs:range 推理**              | 如果 P range C 且 x P y，則 y:C             | hasAge range Integer      |
| **rdfs:subPropertyOf**           | 屬性階層推理                                | drives ⊑ operates         |
| **owl:equivalentClass**          | 類別等價（對稱/傳遞）                        | Vehicle ≡ Automobile      |
| **owl:sameAs**                   | 個體同一性（對稱/傳遞）                      | person1 ≡ person2         |
| **owl:inverseOf**                | 逆屬性推理                                  | owns ⟷ isOwnedBy          |
| **owl:TransitiveProperty**       | 傳遞屬性鏈                                  | locatedIn 傳遞性           |
| **owl:SymmetricProperty**        | 對稱屬性推理                                | marriedTo 對稱性           |

## 輸出格式

本工具支援兩種推理結果輸出格式：

### N-Triples 格式（預設）

標準 N-Triples 格式，使用完整 IRI：

```
<http://example.org/myCar> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://example.org/Car> .
<http://example.org/myCar> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://example.org/Vehicle> .
<http://example.org/myCar> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://example.org/Transport> .
<http://example.org/Car> <http://www.w3.org/2000/01/rdf-schema#subClassOf> <http://example.org/Vehicle> .
<http://example.org/Vehicle> <http://www.w3.org/2000/01/rdf-schema#subClassOf> <http://example.org/Transport> .
```

### Datalog 格式

Datalog 事實格式，使用簡化名稱，適用於 Datalog 推理系統：

```
type(myCar, Car).
type(myCar, Vehicle).
type(myCar, Transport).
subClassOf(Car, Vehicle).
subClassOf(Vehicle, Transport).
```

Datalog 格式將 RDF 三元組 `<subject> <predicate> <object>` 轉換為事實 `predicate(subject, object).`，並透過擷取本地名稱簡化 IRI。

## Datalog 推理

除了 RDF/OWL 前向推理，goreasoner 還包含一個內建的 Datalog 求值器，支援事實、帶變數的規則、遞迴規則和布林查詢。求值器使用**樸素自底向上（前向鏈結）求值**和不動點計算，在回答查詢之前推導所有可能的事實。

### Datalog 語法

#### 事實

以句號結尾的基礎原子：

```
Parent(john, mary).
Type(myTesla, Car).
```

#### 規則

帶有頭部和一個或多個體原子的 Horn 子句，用 `:-` 分隔：

```
Ancestor(X, Y) :- Parent(X, Y).
Ancestor(X, Z) :- Parent(X, Y), Ancestor(Y, Z).
```

#### 變數

變數透過以下慣例辨識：

- 單個大寫字母：`X`、`Y`、`Z`
- 全大寫識別符號：`VAR_X`、`PERSON`、`NODE1`
- `?` 前綴識別符號：`?x`、`?person`

其他所有識別符號被視為常數。

#### 查詢

查詢使用 `?-` 前綴，傳回布林值（如果存在匹配事實則為 true）：

```
?- Ancestor(john, jane).
?- type(X, Vehicle).
```

#### 註解

同時支援 Prolog 風格（`%`）和 C 風格（`//`）行註解：

```
% 這是一條註解
Parent(john, mary).  // 這也是一條註解
```

### Datalog Go 函式庫使用

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

如需更精細的控制，可以直接使用底層 API：

```go
program, err := reasoner.ParseDatalog(input)
if err != nil {
    log.Fatal(err)
}

// 透過不動點求值推導所有事實
derivedFacts := program.Reason()

// 剖析並求值查詢
query, err := reasoner.ParseQuery("?- Ancestor(john, jane).")
if err != nil {
    log.Fatal(err)
}

satisfied := program.EvaluateQuery(query, derivedFacts)
```

### Datalog API 參考

#### `DLQuery(datalogContent, queryStr string) (bool, error)`

Datalog 查詢的主 API 函式。剖析程式，執行推理至不動點，並求值查詢。

#### `ParseDatalog(input string) (*DatalogProgram, error)`

將 Datalog 程式字串剖析為包含事實和規則的 `DatalogProgram`。

#### `ParseQuery(s string) (DLAtom, error)`

將查詢字串（帶或不帶 `?-` 前綴）剖析為 `DLAtom`。

#### `(*DatalogProgram) Reason() []DLAtom`

執行前向鏈結求值，直到不再推導出新事實。傳回所有基礎事實（原始的和推理得出的）。

#### `(*DatalogProgram) EvaluateQuery(query DLAtom, derivedFacts []DLAtom) bool`

檢查查詢是否匹配任何推導事實。查詢中的變數作為萬用字元。

### Datalog 限制

Datalog 求值器設計用於簡單的正 Datalog 程式。請注意以下限制：

- **不支援否定**：規則體中不支援失敗否定（`not`、`\+`）。規則中只能出現正原子。
- **不支援內建比較和算術**：不支援 `!=`、`<`、`>` 等運算子和算術運算式。所有項均為符號常數或變數。
- **僅支援布林查詢**：`DLQuery` 傳回 `true`/`false`。不傳回變數繫結。例如，查詢 `?- Ancestor(john, X).` 會告訴您是否存在祖先，但不會列舉它們。
- **規則無安全性檢查**：剖析器接受頭部包含體中未出現變數的規則（如 `Foo(X) :- Bar(Y).`）。此類規則不會產生錯誤結果（未接地的頭部會被靜默丟棄），但不會發出警告。
- **事實無索引**：求值器在匹配規則體原子時對所有事實進行線性掃描。這對於中小規模程式足夠，但對於數千條事實可能會變慢。
- **不支援聚合和約束**：不支援延伸 Datalog 系統中的 `count`、`min`、`max` 或完整性約束等功能。

## 範例

### 範例 1：基本類別階層推理

```bash
# 建立綱要檔案 (tbox.ttl)
echo '@prefix ex: <http://example.org/> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .

ex:SportsCar rdfs:subClassOf ex:Car .
ex:Car rdfs:subClassOf ex:Vehicle .' > tbox.ttl

# 建立實例檔案 (abox.ttl)
echo '@prefix ex: <http://example.org/> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .

ex:myFerrari rdf:type ex:SportsCar .' > abox.ttl

# 執行推理（N-Triples 輸出）
goreasoner run abox.ttl tbox.ttl -o results.nt

# 執行推理（Datalog 輸出）
goreasoner run abox.ttl tbox.ttl --outputType=datalog -o results.dl

# 檢視結果
cat results.nt
cat results.dl
```

### 範例 2：複雜本體推理

給定一個包含領域/值域限制、屬性階層結構和等價類別的本體，推理器將根據 RDFS/OWL 語義推導所有隱含知識。

## API 參考

### 核心函式

#### `ForwardReason(abox, tbox string) ([]string, error)`

執行 RDF 資料前向推理的主 API 函式。

**參數：**

- `abox`：包含實例資料（斷言）的 Turtle 字串
- `tbox`：包含綱要/本體定義的 Turtle 字串

**傳回值：**

- `[]string`：所有三元組的清單（包括推理得出的），N-Triples 格式
- `error`：任何剖析或處理錯誤

#### `ForwardReasonWithDetails(abox, tbox string) (*ReasoningResult, error)`

傳回詳細的推理結果，分離原始/推理三元組。

**參數：**

- `abox`：包含實例資料的 Turtle 字串
- `tbox`：包含綱要定義的 Turtle 字串

**傳回值：**

- `*ReasoningResult`：詳細結果結構
- `error`：任何剖析或處理錯誤

### 直接使用推理器

```go
r := reasoner.NewReasoner()

// 載入 TBox（綱要）
if err := r.LoadTurtle(tbox); err != nil {
    log.Fatalf("載入 TBox 出錯: %v", err)
}

// 載入 ABox（實例）
if err := r.LoadTurtle(abox); err != nil {
    log.Fatalf("載入 ABox 出錯: %v", err)
}

// 執行推理
inferredCount := r.RunForwardReasoning()
fmt.Printf("推理出 %d 個新三元組\n", inferredCount)

// 查詢特定模式（使用 "" 作為萬用字元）
vehicles := r.Query("", reasoner.RDFType, "http://example.org/Vehicle")
for _, t := range vehicles {
    fmt.Printf("%s 是一輛 Vehicle\n", t.Subject)
}

// 取得特定實例的所有類型
types := r.GetInferredTypes("http://example.org/myCar")
fmt.Printf("myCar 的類型: %v\n", types)
```

### 資料結構

#### `ReasoningResult`

表示詳細的推理結果：

```go
type ReasoningResult struct {
    OriginalTriples []string // 輸入中的三元組
    InferredTriples []string // 新推理出的三元組
    AllTriples      []string // 所有三元組的合集
    OriginalCount   int      // 原始三元組數量
    InferredCount   int      // 推理三元組數量
    TotalCount      int      // 三元組總數
}
```

#### `Triple`

表示一個 RDF 三元組：

```go
type Triple struct {
    Subject   string
    Predicate string
    Object    string
}
```

#### `Reasoner`

主推理器結構及其方法：

```go
type Reasoner struct {
    store  *TripleStore
    rules  []Rule
    parser *TurtleParser
}
```

### 推理器方法

| 方法                                                 | 描述                                                       |
| --------------------------------------------------- | ----------------------------------------------------------------- |
| `NewReasoner() *Reasoner`                           | 使用預設規則建立新推理器                                            |
| `NewReasonerWithRules(rules []Rule) *Reasoner`      | 使用自訂規則建立推理器                                              |
| `LoadTurtle(content string) error`                  | 剖析並載入 Turtle 內容                                             |
| `RunForwardReasoning() int`                         | 套用所有規則直到不動點，傳回推理出的三元組數量                        |
| `GetAllTriples() []string`                          | 取得所有三元組的 N-Triples 字串                                     |
| `GetInferredTypes(subject string) []string`         | 取得主詞的所有 rdf:type 值                                         |
| `Query(subject, predicate, object string) []Triple` | 模式匹配查詢（使用 "" 作為萬用字元）                                 |
| `GetStore() *TripleStore`                           | 存取底層三元組儲存                                                  |

## 架構

本函式庫組織為幾個關鍵元件：

### 核心套件 (`pkg/reasoner`)

- **前向推理引擎**：基於規則的推理，使用不動點計算
- **Turtle 剖析器**：完整的 Turtle 格式剖析器，支援前綴
- **三元組儲存**：索引化的記憶體儲存，用於高效查詢
- **規則系統**：模組化的 RDFS/OWL 推理規則
- **Datalog 求值器**：Datalog 程式的剖析器和樸素自底向上推理器
- **查詢介面**：模式匹配、類型推理和 Datalog 查詢

### 命令列介面 (`cmd/goreasoner`)

- **基於 Cobra 的命令列**，提供推理操作的子命令
- **檔案 I/O 處理**，用於 Turtle 輸入和 N-Triples 輸出
- **完善的說明系統**，提供詳細的使用範例
- **靈活的輸出選項**和錯誤處理

### 貢獻

1. Fork 本儲存庫
2. 建立功能分支：`git checkout -b feature/new-feature`
3. 進行修改並新增測試
4. 確保所有測試通過：`go test ./...`
5. 提交修改：`git commit -am 'Add new feature'`
6. 推送到分支：`git push origin feature/new-feature`
7. 提交 Pull Request

### 測試

執行測試套件：

```bash
go test ./...
```

執行帶覆蓋率的測試：

```bash
go test -cover ./...
```

## 建置環境

### 使用 Nix（建議）

使用 Nix flakes 設定建置環境：

```bash
nix develop
```

### 手動建置

檢查 `build.ps1` 中的建置參數：

```bash
# 建置帶有版本資訊的靜態二進位檔
$env:CGO_ENABLED = "1"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
```

然後執行：

```bash
./build.ps1
```

或手動建置：

```bash
go build -o goreasoner .
```

## 檔案結構

```
goreasoner/
├── go.mod                    # 模組定義
├── main.go                   # 主進入點
├── build.ps1                 # 建置腳本
├── flake.nix                 # Nix flake 組態
├── cmd/
│   └── goreasoner/
│       ├── main.go           # 命令列介面
│       └── commands.go       # 命令定義
├── pkg/
│   ├── reasoner/
│   │   ├── core.go           # 主 API 和 Reasoner 類型
│   │   ├── parser.go         # Turtle 格式剖析器
│   │   ├── store.go          # 記憶體三元組儲存
│   │   ├── rules.go          # 前向推理規則
│   │   ├── datalog.go        # Datalog 剖析器和推理器
│   │   ├── utils.go          # 工具函式
│   │   └── error.go          # 錯誤處理
│   └── version/
│       └── version.go        # 版本資訊
└── docs/
    └── docs-gen.go           # 文件產生器
```

## 授權條款

本專案基於 MIT 授權條款發布 - 詳見 [LICENSE](LICENSE) 檔案。
