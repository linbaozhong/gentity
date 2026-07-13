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
		targets = []string{"pc", "h5", "mp"}
	} else if target != "" {
		targets = strings.Split(target, ",")
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
	pcSpecs := []feTmplSpec{
		{"fe_pc_form.tmpl", "pc",
			func(td TempData) string { return toPascal(td.StructName) + "Form.vue" },
			func(td TempData) bool { return hasTag(td, "@request") }},
	}

	for _, td := range tds {
		td.Entity = toEntity(td.StructName)
		td.Group = group
		for _, spec := range sharedSpecs {
			if spec.onlyIf != nil && !spec.onlyIf(td) {
				continue
			}
			if err := writeFeFile(out, spec, td); err != nil {
				return err
			}
		}
		for _, tg := range targets {
			if tg != "pc" {
				continue // h5/mp 在 M2/M3/M4 接入
			}
			for _, spec := range pcSpecs {
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
	full := filepath.Join(dir, spec.outName(td)) // ← 改为传 td，不再用 td.Entity
	_f := feFuncMap()
	return writeFeTextFile(full, _f, func(w io.Writer, fm template.FuncMap) error {
		tmpl := template.New("").Funcs(fm)
		if _, e := tmpl.ParseFS(resources.TemplatesFS, "templates/fe/"+spec.tmpl); e != nil {
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
	}
}

// 保留原空壳以备兼容（可删除）
func frontend(name string) error {
	return nil
}
