# Point 类型 SRID 4326 修改说明

## 概述

`types.Point` 类型已更新为使用 SRID 4326（WGS84 坐标系），符合地理信息系统（GIS）标准。

## SRID 4326

SRID 4326 是全球最常用的空间参考标识符（Spatial Reference System Identifier），对应 WGS84 坐标系。这是 GPS 设备使用的标准坐标系。

### 坐标系说明

- **X 坐标（纬度 Latitude）**：范围 -90 到 +90
- **Y 坐标（经度 Longitude）**：范围 -180 到 +180

### 常用坐标示例

| 地点 | 纬度 (X) | 经度 (Y) |
|------|-----------|-----------|
| 北京天安门 | 39.915 | 116.404 |
| 上海东方明珠 | 31.239 | 121.499 |
| 纽约时代广场 | 40.758 | -73.985 |
| 伦敦大本钟 | 51.507 | -0.127 |
| 东京塔 | 35.658 | 139.745 |

## EWKB 格式

SRID 4326 使用 EWKB（Extended Well-Known Binary）格式，相比标准 WKB 格式增加了 SRID 信息。

### EWKB 结构（25 字节）

```
偏移    大小    说明
-----  ------  ------------------------
0      1       字节序（1 = 小端序）
1-4    4       几何类型（POINT + SRID 标志位 = 0x20000001）
5-8    4       SRID（4326）
9-17   8       X 坐标（纬度，float64）
17-25  8       Y 坐标（经度，float64）
```

### 几何类型标志位

- 基础 POINT 类型：1
- SRID 标志位：0x20000000
- 组合值：0x20000001

## API 变更

### Point 结构体

```go
type Point struct {
    X float64 // 纬度 (Latitude)
    Y float64 // 经度 (Longitude)
}
```

### Scan 方法

从数据库读取 EWKB 格式的 Point 数据：

```go
func (p *Point) Scan(value interface{}) error
```

**验证逻辑：**
1. 检查字节序（仅支持小端序）
2. 验证几何类型（必须是 0x20000001）
3. 验证 SRID（必须是 4326）
4. 解析 X 和 Y 坐标

### Value 方法

将 Point 转换为 EWKB 格式：

```go
func (p Point) Value() (driver.Value, error)
```

**生成的格式：**
- 25 字节的 EWKB 格式
- 包含 SRID 4326
- 使用小端序

## 使用示例

### 基本使用

```go
import (
    "github.com/linbaozhong/gentity/pkg/types"
)

// 创建 Point（北京天安门）
p := types.Point{
    X: 39.915,  // 纬度
    Y: 116.404, // 经度
}

// 在数据库操作中使用
result, err := ace.
    Table("locations").
    Set(dialect.F("location").Set(p)).
    Create(db).
    Exec(context.Background())
```

### 从数据库读取

```go
var location types.Point
err := db.QueryRow("SELECT location FROM locations WHERE id = ?", 1).Scan(&location)
if err != nil {
    // 处理错误
}

fmt.Printf("纬度: %f, 经度: %f\n", location.X, location.Y)
```

### 计算距离

使用 MySQL 的空间函数计算两点之间的距离：

```go
// 计算两点之间的距离（单位：米）
// 参数顺序：纬度, 经度
distanceFunc := dialect.F("location").Distance(39.915, 116.404, "distance")

result, err := ace.
    Table("locations").
    Cols(dialect.F("*"), distanceFunc).
    Where(distanceFunc.Lt(1000)). // 距离小于1000米
    Select(db).
    Maps(context.Background())
```

## 测试

完整的测试用例在 `pkg/types/point_test.go` 中：

```bash
cd pkg/types
go test -v -run TestPoint
```

测试内容包括：
- ✅ Value() 方法（EWKB 格式生成）
- ✅ Scan() 方法（EWKB 格式解析）
- ✅ 零值处理
- ✅ nil 值处理
- ✅ 往返转换
- ✅ 地理坐标验证
- ✅ 接口实现验证

## 数据库要求

### MySQL 版本

- MySQL 5.7 或更高版本
- 建议使用 MySQL 8.0+ 以获得更好的空间函数支持

### 表定义

```sql
CREATE TABLE locations (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100),
    location POINT SRID 4326 NOT NULL,
    SPATIAL INDEX(location)
) ENGINE=InnoDB;
```

### 插入数据

```sql
-- 注意：SRID 4326 的 POINT 格式是 POINT(纬度, 经度)
INSERT INTO locations (name, location)
VALUES ('北京天安门', ST_GeomFromText('POINT(39.915 116.404)', 4326));
```

### 查询数据

```sql
-- 查询所有点
SELECT id, name, ST_AsText(location), ST_X(location) AS latitude, ST_Y(location) AS longitude
FROM locations;

-- 计算距离
-- 注意：POINT 格式是 POINT(纬度, 经度)
SELECT
    id,
    name,
    ST_Distance_Sphere(
        location,
        ST_GeomFromText('POINT(39.915 116.404)', 4326)
    ) AS distance
FROM locations
WHERE ST_Distance_Sphere(
    location,
    ST_GeomFromText('POINT(39.915 116.404)', 4326)
) < 1000; -- 距离小于1000米
```

## 迁移指南

如果你的项目使用了旧版本的 Point 类型：

### 1. 更新表结构

如果表中已有数据，需要迁移到 SRID 4326：

```sql
-- 添加新的 SRID 4326 列
ALTER TABLE locations ADD COLUMN location_srid4326 POINT SRID 4326;

-- 迁移数据
UPDATE locations
SET location_srid4326 = ST_GeomFromText(ST_AsText(location), 4326);

-- 删除旧列
ALTER TABLE locations DROP COLUMN location;
ALTER TABLE locations CHANGE COLUMN location_srid4326 location POINT SRID 4326;
```

### 2. 更新代码

无需修改代码，只需确保使用新的 Point 类型即可。

## 注意事项

1. **SRID 4326 专用**：当前实现仅支持 SRID 4326，如需其他 SRID 需要扩展。

2. **坐标顺序**：X 是纬度，Y 是经度，不要混淆。这与数学坐标系不同，符合地理信息系统（GIS）标准。

3. **字节序**：仅支持小端序（little-endian）。

4. **浮点数精度**：使用 float64 存储坐标，精度约为 15 位有效数字。

5. **零值处理**：当 X 和 Y 都为 0 时，Value() 方法返回 nil（SQL NULL）。

6. **空间索引**：建议为 Point 列创建 SPATIAL INDEX 以提高查询性能。

## 参考资料

- [WGS84 坐标系](https://en.wikipedia.org/wiki/World_Geodetic_System)
- [SRID 4326](https://epsg.io/4326)
- [MySQL 空间数据](https://dev.mysql.com/doc/refman/8.0/en/spatial-types.html)
- [EWKB 格式](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry#Well-known_binary)
