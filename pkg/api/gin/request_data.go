// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package api

// ReadJSON Content-Type为application/json的请求
func ReadJSON(ctx Context, ptr any) error {
	return adapt(ctx).ReadJSON(ptr)
}

// ReadForm Content-Type为application/x-www-form-urlencoded的请求
func ReadForm(ctx Context, ptr any) error {
	return adapt(ctx).ReadForm(ptr)
}

// ReadQuery 读取 URL query 参数
func ReadQuery(ctx Context, ptr any) error {
	return adapt(ctx).ReadQuery(ptr)
}
