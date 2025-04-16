---
title: "go中怎么向数据库存入切片"
date: "2025-04-16"
categories: ["go"]
---

```go
type DataIndicators struct {
    ProjectID             int
    ProjectPage           string
    PvTrackingPoint       []*string
    PvOther               []*int
    UvTrackingPoint       []*string
    UvOther               []*int
    // 其他字段...
}
```

```go
type DataIndicators struct {
    ProjectID             int
    ProjectPage           string
    PvTrackingPoint       []string // 普通切片
    PvOther               []int    // 普通切片
    UvTrackingPoint       []string // 普通切片
    UvOther               []int    // 普通切片
    // 其他字段...
}
```

用上诉两种方式存储都会报错,因为不支持存储空切片（作为指针类型，相当于出现了空指针）

以下是两种解决办法:
1. 以json格式存储
```go
type DataIndicators struct {
    ID                    uint   `json:"id"`
    ProjectID             int    `json:"project_id"`
    ProjectPage           string `json:"project_page"`
    PvTrackingPoint       []string `gorm:"type:json" json:"pv_tracking_point"`
    PvOther               []int    `gorm:"type:json" json:"pv_other"`
    // 其他字段...
}
```
2. 自定义类型，使用 GORM 的 Scanner 和 Valuer 接口
```go
type StringArray []string

// Scan 从数据库读取数据时将数据反序列化
func (s *StringArray) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("failed to scan value")
    }
    return json.Unmarshal(bytes, &s)
}

// Value 将数据序列化成数据库存储格式
func (s StringArray) Value() (driver.Value, error) {
    return json.Marshal(s)
}

type DataIndicators struct {
    ID                    uint       `json:"id"`
    ProjectID             int        `json:"project_id"`
    ProjectPage           string     `json:"project_page"`
    PvTrackingPoint       StringArray `json:"pv_tracking_point"`
    PvOther               StringArray `json:"pv_other"`
    // 其他字段...
}
```