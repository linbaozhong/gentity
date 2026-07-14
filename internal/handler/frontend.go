package handler

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/linbaozhong/gentity/internal/resources"
)

// feTmplSpec 描述一个前端模板及其输出位置/文件名
type feTmplSpec struct {
	tmpl    string
	subDir  string
	outName func(td TempData) string
	onlyIf  func(td TempData) bool // 可选：满足条件才生成
}

// generateFrontend 根据 DTO 解析结果生成前端组件
func generateFrontend(tds []TempData, filename, target, out, group string) error {
	targets := []string{"pc"}
	if target == "all" {
		// 方案 A：砍掉 h5 端，mp 组件经 uni-app 跨端覆盖 H5
		targets = []string{"pc", "mp"}
	} else if target != "" {
		targets = nil // 先清空，再按用户传入过滤
		for _, t := range strings.Split(target, ",") {
			if t == "h5" {
				continue // 方案 A：忽略 h5
			}
			targets = append(targets, t)
		}
	}

	// M5: 建立「响应基名 → GET 请求结构体名」配对，供 ListView 自动嵌套 SearchForm
	reqByBase := map[string]string{}
	for _, td := range tds {
		if hasTag(td, "@request") && toMethod(td.StructName) == "get" {
			reqByBase[strings.TrimSuffix(td.StructName, "Req")] = td.StructName
		}
	}

	sharedSpecs := []feTmplSpec{
		{"fe_types.tmpl", "shared",
			func(td TempData) string { return strings.ToLower(td.StructName) + ".types.js" }, nil},
		{"fe_validation.tmpl", "shared",
			func(td TempData) string { return strings.ToLower(td.StructName) + ".validation.js" },
			func(td TempData) bool { return hasTag(td, "@request") }},
		{"fe_options.tmpl", "shared",
			func(td TempData) string { return strings.ToLower(td.StructName) + ".options.js" },
			func(td TempData) bool { return hasTag(td, "@request") }},
		{"fe_api.tmpl", "shared",
			func(td TempData) string { return strings.ToLower(td.StructName) + ".api.js" },
			func(td TempData) bool { return hasTag(td, "@request") }},
		{"fe_mock.tmpl", "shared",
			func(td TempData) string { return strings.ToLower(td.StructName) + ".mock.js" },
			func(td TempData) bool { return hasTag(td, "@response") }},
	}

	for _, td := range tds {
		td.Entity = toEntity(td.StructName)
		td.Group = group
		// M5: 列表型响应尝试配对对应的 GET 请求 SearchForm
		if hasTag(td, "@response") && isListView(td.StructName) {
			if reqName, ok := reqByBase[strings.TrimSuffix(td.StructName, "Resp")]; ok {
				td.SearchFormComp = toPascal(reqName) + "SearchForm.vue"
			}
		}
		for _, spec := range sharedSpecs {
			if spec.onlyIf != nil && !spec.onlyIf(td) {
				continue
			}
			if err := writeFeFile(out, spec, td); err != nil {
				return err
			}
		}
		for _, tg := range targets {
			td.Platform = tg
			for _, spec := range platformSpecs(tg) {
				if spec.onlyIf != nil && !spec.onlyIf(td) {
					continue
				}
				if err := writeFeFile(out, spec, td); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// writeFeFile 用指定模板渲染并写出（不跑 go/format）
func writeFeFile(out string, spec feTmplSpec, td TempData) error {
	dir := filepath.Join(out, spec.subDir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	// 模板文件不存在（如 M4 的 mp 模板尚未落地）时跳过，而非中断整轮生成
	tmplPath := "templates/fe/" + spec.tmpl
	if _, e := resources.TemplatesFS.Open(tmplPath); e != nil {
		return nil
	}
	full := filepath.Join(dir, spec.outName(td))
	_f := feFuncMap()
	return writeFeTextFile(full, _f, func(w io.Writer, fm template.FuncMap) error {
		tmpl := template.New("").Funcs(fm)
		if _, e := tmpl.ParseFS(resources.TemplatesFS, tmplPath); e != nil {
			return e
		}
		return tmpl.ExecuteTemplate(w, spec.tmpl, td)
	})
}

// writeFeTextFile 写文本文件（前端文件不需要 gofmt）
func writeFeTextFile(fullFilename string, funcMap template.FuncMap, fn func(io.Writer, template.FuncMap) error) error {
	var buf bytes.Buffer
	// 先渲染到内存，失败直接返回错误，不碰磁盘文件
	if e := fn(&buf, funcMap); e != nil {
		return e
	}
	if fi, e := os.Stat(fullFilename); e == nil && !fi.IsDir() {
		os.Remove(fullFilename)
	}
	_f, e := os.OpenFile(fullFilename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		return e
	}
	defer _f.Close()
	_, e = _f.Write(buf.Bytes())
	return e
}

// platformSpecs 返回某平台的前端组件模板规格。
// pc 用 fe_{kind}.tmpl（仅 class 用 {{.Platform}} 区分）；mp 用 fe_mp_{kind}.tmpl（uni-app 标签，M4 落地）。
func platformSpecs(platform string) []feTmplSpec {
	prefix := "fe_"
	if platform == "mp" {
		prefix = "fe_mp_"
	}
	kind := func(name string) string { return prefix + name + ".tmpl" }
	specs := []feTmplSpec{
		{kind("form"), platform,
			func(td TempData) string { return toPascal(td.StructName) + "Form.vue" },
			func(td TempData) bool { return hasTag(td, "@request") && toMethod(td.StructName) != "get" }},
		{kind("searchform"), platform,
			func(td TempData) string { return toPascal(td.StructName) + "SearchForm.vue" },
			func(td TempData) bool { return hasTag(td, "@request") && toMethod(td.StructName) == "get" }},
		{kind("detail"), platform,
			func(td TempData) string { return toPascal(td.StructName) + "Detail.vue" },
			func(td TempData) bool { return hasTag(td, "@response") }},
		{kind("table"), platform,
			func(td TempData) string { return toPascal(td.StructName) + "Table.vue" },
			func(td TempData) bool { return hasTag(td, "@response") && isListView(td.StructName) }},
		{kind("card"), platform,
			func(td TempData) string { return toPascal(td.StructName) + "Card.vue" },
			func(td TempData) bool { return hasTag(td, "@response") && isListView(td.StructName) }},
		{kind("columns"), platform,
			func(td TempData) string { return toPascal(td.StructName) + "Columns.js" },
			func(td TempData) bool { return hasTag(td, "@response") && isListView(td.StructName) }},
		{kind("actions"), platform,
			func(td TempData) string { return toPascal(td.StructName) + "Actions.vue" },
			func(td TempData) bool { return hasTag(td, "@response") && isListView(td.StructName) }},
		{kind("listview"), platform,
			func(td TempData) string { return toPascal(td.StructName) + "ListView.vue" },
			func(td TempData) bool { return hasTag(td, "@response") && isListView(td.StructName) }},
	}
	// 方案 A：移动端（mp）不用表格，跳过 table 规格
	if platform == "mp" {
		var filtered []feTmplSpec
		for _, s := range specs {
			if s.tmpl == kind("table") {
				continue
			}
			filtered = append(filtered, s)
		}
		specs = filtered
	}
	return specs
}

// feFuncMap 前端模板专用函数集
func feFuncMap() template.FuncMap {
	return template.FuncMap{
		"lower":       strings.ToLower,
		"toJsType":    toJsType,
		"toEntity":    toEntity,
		"toPascal":    toPascal,
		"toValidRule": toValidRule,
		"toEndpoint":  toEndpoint,
		"hasTag":      hasTag,
		"toOptions":   toOptions,
		"mockValue":   mockValue,
		"replace":     strings.ReplaceAll,
		"toMethod":    toMethod,
		"isListView":  isListView,
		"isSlice":     isSlice,
		"hasValid":    hasValid,
		"toInputType": toInputType,
		"toPkName":    toPkName,
		"toOptionsJS": toOptionsJS,
	}
}

// 保留原空壳以备兼容（可删除）
func frontend(name string) error {
	return nil
}
