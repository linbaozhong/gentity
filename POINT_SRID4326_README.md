# Point 类型 SRID 4326 修改完成

## 修改概述

已成功将 `types.Point` 类型更新为使用 **SRID 4326**（WGS84 坐标系），符合地理信息系统（GIS）标准。

## SRID 4326 说明

SRID 4326 是全球最常用的空间参考标识符，对应 WGS84 坐标系，这是 GPS 设备使用的标准坐标系。

### 坐标定义

- **X 坐标（纬度 Latitude）**：范围 -90 到 +90
- **Y 坐标（经度 Longitude）**：范围 -180 到 +180

## 修改的文件

### 核心修改
1. **pkg/types/point.go** - 更新为 SRID 4326 (EWKB 格式)
   - `Scan()` 方法：解析 25 字节的 EWKB 格式
   - `Value()` 方法：生成 25 字节的 EWKB 格式

### 新增文件
1. **pkg/types/point_test.go** - 完整的测试用例
2. **pkg/types/POINT_SRID4326.md** - 详细的使用文档
3. **example/point_example.go** - 实际使用示例

## EWKB 格式

SRID 4326 使用 EWKB（Extended Well-Known Binary）格式，总长度 **25 字节**：

```
偏移    大小    说明
-----  ------  ------------------------
0      1       字节序（1 = 小端序）
1-4    4       几何类型（POINT + SRID 标志位 = 0x20000001）
5-8    4       SRID（4326）
9-17   8       X 坐标（纬度，float64）
17-25  8       Y 坐标（经度，float64）
```

## 使用方法

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

// 创建 Point（北京天安门）
p := types.Point{
    X: 116.404, // 经度
    Y: 39.915,  // 纬度
}

// 插入到数据库
_, err := ace.
    Table("locations").
    Set(dialect.F("name").Set("北京天安门")).
    Set(dialect.F("location").Set(p)).
    Create(db).
    Exec(context.Background())
```

### 计算距离

```go
// 目标点
target := types.Point{X: 39.915, Y: 116.404}

// 计算距离并查询 1000 米内的位置
distanceFunc := dialect.F("location").Distance(target.X, target.Y, "distance")

result, err := ace.
    Table("locations").
    Cols(dialect.F("*"), distanceFunc).
    Where(distanceFunc.Lt(1000)).
    Select(db).
    Maps(context.Background())
```

## 数据库要求

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
### 更新数据
```sql
UPDATE locations SET location = ST_GeomFromText('POINT(39.915 116.404)', 4326) WHERE id = 1;
```
### 查询数据
```sql
SET @lng = 116.461;     //经度
SET @lat = 39.908;      //纬度
SET @radius = 5000;

SET @lat_offset = @radius / 111320;
SET @lng_offset = @radius / (111320 * COS(RADIANS(@lat)));

SELECT * FROM (
    SELECT 
        id, 
        title, 
        ST_Distance_Sphere(
            location_point, 
            ST_PointFromText(CONCAT('POINT(', @lat, ' ', @lng, ')'), 4326)
        ) AS distance_meters
    FROM events
    WHERE MBRContains(
        ST_GeomFromText(
            CONCAT(
                'POLYGON((',
                @lat - @lat_offset, ' ', @lng - @lng_offset, ', ',  -- 纬度, 经度
                @lat + @lat_offset, ' ', @lng - @lng_offset, ', ',  -- 纬度, 经度
                @lat + @lat_offset, ' ', @lng + @lng_offset, ', ',  -- 纬度, 经度
                @lat - @lat_offset, ' ', @lng + @lng_offset, ', ',  -- 纬度, 经度
                @lat - @lat_offset, ' ', @lng - @lng_offset, '))'   -- 纬度, 经度
            ), 4326
        ),
        location_point
    )
) AS evs  
WHERE distance_meters <= @radius 
ORDER BY distance_meters ASC 
LIMIT 0, 10000;
```
## 主要特性

✅ 符合 SRID 4326 (WGS84) 标准
✅ 使用 EWKB 格式（25 字节）
✅ 自动验证 SRID 和几何类型
✅ 支持 sql.Scanner 和 driver.Valuer 接口
✅ 完整的测试用例
✅ 详细的文档和示例

## 测试

运行测试：
```bash
cd pkg/types
go test -v -run TestPoint
```

测试覆盖：
- ✅ Value() 方法（EWKB 格式生成）
- ✅ Scan() 方法（EWKB 格式解析）
- ✅ 零值和 nil 值处理
- ✅ 往返转换
- ✅ 地理坐标验证
- ✅ 接口实现验证

## 常用坐标示例

| 地点 | 纬度 (X) | 经度 (Y) |
|------|-----------|-----------|
| 北京天安门 | 39.915 | 116.404 |
| 上海东方明珠 | 31.239 | 121.499 |
| 纽约时代广场 | 40.758 | -73.985 |
| 伦敦大本钟 | 51.507 | -0.127 |
| 东京塔 | 35.658 | 139.745 |

## 注意事项

1. **SRID 4326 专用**：当前实现仅支持 SRID 4326
2. **坐标顺序**：X 是纬度，Y 是经度（这与数学坐标系不同，符合地理信息系统（GIS）标准）
3. **字节序**：仅支持小端序
4. **零值处理**：当 X 和 Y 都为 0 时，返回 SQL NULL
5. **空间索引**：建议为 Point 列创建 SPATIAL INDEX

## 文档

- **详细文档**: `pkg/types/POINT_SRID4326.md`
- **测试用例**: `pkg/types/point_test.go`
- **使用示例**: `example/point_example.go`

## 向后兼容性

如果项目中已有使用旧 Point 类型的代码，需要更新表结构并迁移数据。详细迁移指南请参考 `pkg/types/POINT_SRID4326.md`。

## 参考资料

- [WGS84 坐标系](https://en.wikipedia.org/wiki/World_Geodetic_System)
- [SRID 4326](https://epsg.io/4326)
- [MySQL 空间数据](https://dev.mysql.com/doc/refman/8.0/en/spatial-types.html)
