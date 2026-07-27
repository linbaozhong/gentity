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

// EntityGroup 按实体分组的 TempData 集合
// 同一实体的所有 @request/@response 结构体会合并到一个 shared/{entity}.js 文件中
type EntityGroup struct {
	Entity  string     // 实体前缀，如 "user"
	Group   string     // 接口分组前缀
	Tds     []TempData // 该实体下所有 TempData
	Targets []string   // 目标平台列表
}

// feTmplSpec 描述一个前端模板及其输出位置/文件名
type feTmplSpec struct {
	tmpl    string
	subDir  string
	outName func(td TempData) string
	onlyIf  func(td TempData) bool // 可选：满足条件才生成
}

// generateFrontend 根据 DTO 解析结果生成前端组件
func generateFrontend(tds []TempData, target, out, group string) error {
	out, _ = filepath.Abs(out)

	targets := []string{"pc"}
	if target == "all" {
		targets = []string{"pc", "mp"}
	} else if target != "" {
		targets = nil
		for _, t := range strings.Split(target, ",") {
			if t == "h5" {
				continue
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

	// 按 entity（文件名）分组
	groups := map[string]*EntityGroup{}
	for i := range tds {
		td := &tds[i]
		entity := getBaseFilename(td.FileName)
		td.Entity = entity
		td.Group = group
		if hasTag(*td, "@response") && isListView(td.StructName) {
			if reqName, ok := reqByBase[strings.TrimSuffix(td.StructName, "Resp")]; ok {
				td.SearchFormComp = toPascal(reqName) + "SearchForm.vue"
			}
		}

		eg, ok := groups[entity]
		if !ok {
			eg = &EntityGroup{
				Entity:  entity,
				Group:   group,
				Targets: targets,
			}
			groups[entity] = eg
		}
		eg.Tds = append(eg.Tds, *td)
	}

	// 1. 按 entity 生成 shared/{entity}.js（合并该实体所有 struct）
	for _, eg := range groups {
		if err := writeSharedFile(out, eg); err != nil {
			return err
		}
	}

	// 2. 为每个 TempData 生成平台组件（Form、SearchForm、Detail 等）
	for i := range tds {
		td := &tds[i]
		td.Entity = getBaseFilename(td.FileName)
		td.Group = group
		if hasTag(*td, "@response") && isListView(td.StructName) {
			if reqName, ok := reqByBase[strings.TrimSuffix(td.StructName, "Resp")]; ok {
				td.SearchFormComp = toPascal(reqName) + "SearchForm.vue"
			}
		}
		for _, tg := range targets {
			td.Platform = tg
			for _, spec := range platformSpecs(tg) {
				if spec.onlyIf != nil && !spec.onlyIf(*td) {
					continue
				}
				if err := writeFeFile(out, spec, *td); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// writeSharedFile 渲染 fe_shared.tmpl 生成 shared/{entity}.js
func writeSharedFile(out string, eg *EntityGroup) error {
	dir := filepath.Join(out, "shared")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	tmplPath := "templates/fe/fe_shared.tmpl"
	if _, e := resources.TemplatesFS.Open(tmplPath); e != nil {
		return nil
	}

	full := filepath.Join(dir, eg.Entity+".js")
	funcMap := feFuncMap()
	return writeFeTextFile(full, funcMap, func(w io.Writer, fm template.FuncMap) error {
		tmpl := template.New("").Funcs(fm)
		if _, e := tmpl.ParseFS(resources.TemplatesFS, tmplPath); e != nil {
			return e
		}
		return tmpl.ExecuteTemplate(w, "fe_shared.tmpl", eg)
	})
}

// writeFeFile 用指定模板渲染并写出（不跑 go/format）
func writeFeFile(out string, spec feTmplSpec, td TempData) error {
	dir := filepath.Join(out, spec.subDir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
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

// platformSpecs 返回某平台的前端组件模板规格
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
