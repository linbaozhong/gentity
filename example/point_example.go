// Copyright © 2023 Linbaozhong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package example

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/types"
)

// PointExample 展示如何使用 SRID 4326 的 Point 类型
//
// Point 类型用于表示地理坐标（经纬度），符合 WGS84 标准
// X = 纬度 (Latitude), Y = 经度 (Longitude)
//
// SRID 4326 是全球最常用的空间参考标识符，对应 WGS84 坐标系
// 这是 GPS 设备使用的标准坐标系
func PointExample() {
	// 示例1：创建 Point（北京天安门坐标）
	beijing := types.Point{
		X: 39.915,  // 纬度
		Y: 116.404, // 经度
	}

	// 插入位置数据到数据库
	// 注意：表中需要有 location 列，类型为 POINT SRID 4326
	_, err := ace.
		Table("locations").
		Set(dialect.F("name").Set("北京天安门")).
		Set(dialect.F("location").Set(beijing)).
		Create(dbx).
		Exec(context.Background())

	if err != nil {
		// 处理错误
		panic(err)
	}
}

// PointInsertMultiple 批量插入多个位置点
func PointInsertMultiple() {
	locations := []struct {
		Name string
		Point types.Point
	}{
		{"北京天安门", types.Point{X: 39.915, Y: 116.404}},
		{"上海东方明珠", types.Point{X: 31.239, Y: 121.499}},
		{"广州塔", types.Point{X: 23.129, Y: 113.324}},
		{"深圳平安金融中心", types.Point{X: 22.523, Y: 114.055}},
	}

	// 批量插入
	for _, loc := range locations {
		_, err := ace.
			Table("locations").
			Set(dialect.F("name").Set(loc.Name)).
			Set(dialect.F("location").Set(loc.Point)).
			Create(dbx).
			Exec(context.Background())

		if err != nil {
			// 处理错误
			panic(err)
		}
	}
}

// PointDistanceQuery 查询指定距离范围内的位置
//
// 使用 MySQL 的 ST_Distance_Sphere 函数计算两点之间的球面距离
// 距离单位：米
func PointDistanceQuery() {
	// 目标点（北京天安门）
	targetPoint := types.Point{
		X: 39.915,  // 纬度
		Y: 116.404, // 经度
	}

	// 最大距离（米）
	maxDistance := 1000.0

	// 计算距离并查询
	distanceFunc := dialect.F("location").Distance(targetPoint.X, targetPoint.Y, "distance")

	result, err := ace.
		Table("locations").
		Cols(
			dialect.F("id"),
			dialect.F("name"),
			dialect.F("ST_AsText(location)").As("location_text"),
			dialect.F("ST_X(location)").As("longitude"),
			dialect.F("ST_Y(location)").As("latitude"),
			distanceFunc,
		).
		Where(distanceFunc.Lt(maxDistance)).
		OrderBy(dialect.F("distance").Asc()).
		Select(dbx).
		Maps(context.Background())

	if err != nil {
		// 处理错误
		panic(err)
	}

	// 处理结果
	for _, row := range result {
		name := row["name"]
		distance := row["distance"]
		latitude := row["latitude"]
		longitude := row["longitude"]

		// 使用结果
		_ = name
		_ = distance
		_ = latitude
		_ = longitude
	}
}

// PointNearest 查询最近的 N 个位置点
func PointNearest(n int) {
	// 目标点
	targetPoint := types.Point{
		X: 39.915,  // 纬度
		Y: 116.404, // 经度
	}

	// 计算距离并查询最近的 N 个点
	distanceFunc := dialect.F("location").Distance(targetPoint.X, targetPoint.Y, "distance")

	result, err := ace.
		Table("locations").
		Cols(
			dialect.F("id"),
			dialect.F("name"),
			dialect.F("ST_X(location)").As("latitude"),
			dialect.F("ST_Y(location)").As("longitude"),
			distanceFunc,
		).
		OrderBy(dialect.F("distance").Asc()).
		Limit(uint(n)).
		Select(dbx).
		Maps(context.Background())

	if err != nil {
		// 处理错误
		panic(err)
	}

	// 处理结果
	for _, row := range result {
		name := row["name"]
		distance := row["distance"]

		// 使用结果
		_ = name
		_ = distance
	}
}

// PointBoundingBox 查询边界框内的位置
//
// 使用 ST_Within 和 ST_MakeEnvelope 函数查询矩形区域内的点
func PointBoundingBox() {
	// 定义边界框（西南角和东北角）
	// 参数顺序：最小经度, 最小纬度, 最大经度, 最大纬度
	minLon := 116.0 // 最小经度
	minLat := 39.0  // 最小纬度
	maxLon := 117.0 // 最大经度
	maxLat := 40.0  // 最大纬度

	// 使用 ST_Within 查询边界框内的点
	result, err := ace.
		Table("locations").
		Cols(
			dialect.F("id"),
			dialect.F("name"),
			dialect.F("ST_X(location)").As("latitude"),
			dialect.F("ST_Y(location)").As("longitude"),
		).
		Where(
			dialect.Expr(
				"ST_Within(location, ST_MakeEnvelope(?, ?, ?, ?, 4326))",
				minLon, minLat, maxLon, maxLat,
			),
		).
		Select(dbx).
		Maps(context.Background())

	if err != nil {
		// 处理错误
		panic(err)
	}

	// 处理结果
	for _, row := range result {
		name := row["name"]
		latitude := row["latitude"]
		longitude := row["longitude"]

		// 使用结果
		_ = name
		_ = latitude
		_ = longitude
	}
}

// PointUpdateLocation 更新位置信息
func PointUpdateLocation(id int, newPoint types.Point) {
	_, err := ace.
		Table("locations").
		Set(dialect.F("location").Set(newPoint)).
		Where(dialect.F("id").Eq(id)).
		Update(dbx).
		Exec(context.Background())

	if err != nil {
		// 处理错误
		panic(err)
	}
}

// PointCountWithinDistance 统计指定距离内的位置数量
func PointCountWithinDistance(targetPoint types.Point, maxDistance float64) (int, error) {
	// 计算距离并统计数量
	distanceFunc := dialect.F("location").Distance(targetPoint.X, targetPoint.Y, "distance")

	count, err := ace.
		Table("locations").
		Where(distanceFunc.Lt(maxDistance)).
		Select(dbx).
		Count(context.Background())

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// PointGroupByRegion 按区域分组统计
func PointGroupByRegion() {
	// 按区域分组（例如：每 1 度 x 1 度的网格）
	result, err := ace.
		Table("locations").
		Cols(
			dialect.Expr("FLOOR(ST_X(location))").As("region_lat"),
			dialect.Expr("FLOOR(ST_Y(location))").As("region_lon"),
			dialect.Func(dialect.F("*").Count()).As("count"),
		).
		GroupBy(
			dialect.Expr("FLOOR(ST_X(location))"),
			dialect.Expr("FLOOR(ST_Y(location))"),
		).
		OrderBy(dialect.F("count").Desc()).
		Select(dbx).
		Maps(context.Background())

	if err != nil {
		// 处理错误
		panic(err)
	}

	// 处理结果
	for _, row := range result {
		regionLat := row["region_lat"]
		regionLon := row["region_lon"]
		count := row["count"]

		// 使用结果
		_ = regionLat
		_ = regionLon
		_ = count
	}
}

		regionY := row["region_y"]
		count := row["count"]

		// 使用结果
		_ = regionX
		_ = regionY
		_ = count
	}
}
