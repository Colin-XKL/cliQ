package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"repo/cliqfile"
)

// ImportTemplate 导入模板文件
func (a *App) ImportTemplate() (*cliqfile.TemplateFile, error) {
	// 打开文件选择对话框
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择模板文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "CLIQ模板文件 (*.yaml, *.yml)",
				Pattern:     "*.yaml;*.yml",
			},
		},
	})

	if err != nil {
		return nil, err
	}

	// 用户取消选择
	if filePath == "" {
		return nil, errors.New("未选择文件")
	}

	// 解析并验证模板文件
	template, err := a.parseAndValidateTemplateFromFile(filePath)
	if err != nil {
		return nil, err
	}

	// 设置应用的模板
	a.setTemplate(template)

	return template, nil
}

// ImportTemplateFromURL 从URL导入模板文件
func (a *App) ImportTemplateFromURL(url string) (*cliqfile.TemplateFile, error) {
	// 从URL下载内容
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("下载模板文件失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 限制读取大小为 1MB (1024 * 1024 bytes) to prevent out-of-memory errors
	const maxSize = 1024 * 1024 // 1MB in bytes
	limitReader := io.LimitReader(resp.Body, maxSize+1)

	// 读取响应内容 with size limit
	data, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, fmt.Errorf("读取模板内容失败: %w", err)
	}
	if len(data) > maxSize {
		return nil, fmt.Errorf("模板文件过大，超过 1MB 限制")
	}

	// 解析并验证YAML内容
	template, err := a.parseAndValidateTemplateFromData(data)
	if err != nil {
		return nil, err
	}

	// 设置应用的模板
	a.setTemplate(template)

	return template, nil
}

// parseAndValidateTemplateFromFile 解析并验证文件中的模板
func (a *App) parseAndValidateTemplateFromFile(filePath string) (*cliqfile.TemplateFile, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	return a.parseAndValidateTemplateFromData(data)
}

// parseAndValidateTemplateFromData 解析并验证数据中的模板
func (a *App) parseAndValidateTemplateFromData(data []byte) (*cliqfile.TemplateFile, error) {
	// 验证模板
	validationErrors, err := cliqfile.Validate(data)
	if err != nil {
		return nil, fmt.Errorf("模板验证过程出错: %w", err)
	}

	if len(validationErrors) > 0 {
		var errMsg strings.Builder
		errMsg.WriteString("模板文件存在错误:\n")
		for _, ve := range validationErrors {
			errMsg.WriteString(fmt.Sprintf("- 行 %d, 字段 '%s': %s\n", ve.Line, ve.Field, ve.Message))
		}
		return nil, errors.New(errMsg.String())
	}

	// 解析YAML
	template, err := cliqfile.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("解析YAML失败: %w", err)
	}

	return template, nil
}

// setTemplate 设置应用的当前模板
func (a *App) setTemplate(template *cliqfile.TemplateFile) {
	a.template = template
}
