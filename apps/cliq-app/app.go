package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"cliq/config"
	"cliq/handlers"
	"repo/cliqfile"
)

// App struct
type App struct {
	ctx             context.Context
	template        *cliqfile.TemplateFile
	fileHandler     *handlers.FileHandler
	settingsService *config.SettingsService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		fileHandler: handlers.NewFileHandler(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.fileHandler.Startup(ctx)
	// init settings service
	ss, err := config.NewSettingsService()
	if err == nil {
		a.settingsService = ss
	}
}

// OpenFileDialog opens a file dialog and returns the selected file path
func (a *App) OpenFileDialog() (string, error) {
	return a.fileHandler.OpenFileDialog()
}

// OpenFileDialogWithFilters opens a file dialog with specific file type filters and returns the selected file path
func (a *App) OpenFileDialogWithFilters(filters []runtime.FileFilter) (string, error) {
	return a.fileHandler.OpenFileDialogWithFilters(filters)
}

// SaveFileDialog opens a save file dialog and returns the selected file path
func (a *App) SaveFileDialog() (string, error) {
	return a.fileHandler.SaveFileDialog()
}

// ExecuteCommand executes a shell command with the given input and output file paths
func (a *App) ExecuteCommand(commandID string, variables map[string]interface{}) (string, error) {
	if a.fileHandler == nil {
		return "", fmt.Errorf("fileHandler 未初始化")
	}
	if a.template == nil {
		return "", fmt.Errorf("模板未加载")
	}
	return a.fileHandler.ExecuteCommand(a.template, commandID, variables)
}

func (a *App) GetCommandText(commandID string, variables map[string]interface{}) (string, error) {
	if a.fileHandler == nil {
		return "", fmt.Errorf("fileHandler is nil")
	}
	if a.template == nil {
		return "", fmt.Errorf("template is nil")
	}
	return a.fileHandler.GetCommandText(a.template, commandID, variables)
}

// ParseCommandToTemplate 将命令字符串解析为模板
func (a *App) ParseCommandToTemplate(commandStr string) (*cliqfile.TemplateFile, error) {
	return cliqfile.GenerateFromCommand(commandStr)
}

// GenerateYAMLFromTemplate 将模板对象转换为YAML字符串
func (a *App) GenerateYAMLFromTemplate(template *cliqfile.TemplateFile) (string, error) {
	return cliqfile.GenerateYAML(template)
}

// ValidateYAMLTemplate 验证YAML模板格式
func (a *App) ValidateYAMLTemplate(yamlStr string) error {
	validationErrors, err := cliqfile.Validate([]byte(yamlStr))
	if err != nil {
		return err
	}
	if len(validationErrors) > 0 {
		var errMsg strings.Builder
		errMsg.WriteString("模板文件存在错误:\n")
		for _, ve := range validationErrors {
			errMsg.WriteString(fmt.Sprintf("- 行 %d, 字段 '%s': %s\n", ve.Line, ve.Field, ve.Message))
		}
		return errors.New(errMsg.String())
	}
	return nil
}

// ParseYAMLToTemplate 解析YAML字符串为模板对象
func (a *App) ParseYAMLToTemplate(yamlStr string) (*cliqfile.TemplateFile, error) {
	// ParseYAMLToTemplate in frontend expects validation as well
	err := a.ValidateYAMLTemplate(yamlStr)
	if err != nil {
		return nil, err
	}
	return cliqfile.Parse([]byte(yamlStr))
}

// ExportTemplateToFile 将模板导出为文件
func (a *App) ExportTemplateToFile(template *cliqfile.TemplateFile, filePath string) error {
	return a.fileHandler.ExportTemplateToFile(template, filePath)
}

// SaveYAMLToFile opens a save file dialog and saves the YAML content to the selected file
func (a *App) SaveYAMLToFile(yamlContent string) error {
	return a.fileHandler.SaveYAMLToFile(yamlContent)
}

// SaveFavTemplate 将模板保存到收藏目录
func (a *App) SaveFavTemplate(template *cliqfile.TemplateFile) error {
	return a.fileHandler.SaveFavTemplate(template)
}

// ListFavTemplates 列出所有收藏的模板文件
func (a *App) ListFavTemplates() ([]*cliqfile.TemplateFile, error) {
	return a.fileHandler.ListFavTemplates()
}

// DeleteFavTemplate 从收藏目录删除指定模板文件
func (a *App) DeleteFavTemplate(templateName string) error {
	return a.fileHandler.DeleteFavTemplate(templateName)
}

// GetFavTemplate 读取指定收藏模板文件内容
func (a *App) GetFavTemplate(templateName string) (*cliqfile.TemplateFile, error) {
	template, err := a.fileHandler.GetFavTemplate(templateName)
	if err != nil {
		return nil, err
	}

	// Set the loaded template to the app's template field so that
	// commands like GetCommandText can access it.
	a.template = template

	return template, nil
}

// UpdateFavTemplate 更新指定收藏模板文件内容
func (a *App) UpdateFavTemplate(oldTemplateName string, newTemplateName string, updatedTemplate *cliqfile.TemplateFile) error {
	return a.fileHandler.UpdateFavTemplate(oldTemplateName, newTemplateName, updatedTemplate)
}

func (a *App) GetAppSettings() (*config.AppSettings, error) {
	if a.settingsService == nil {
		ss, err := config.NewSettingsService()
		if err != nil {
			return nil, err
		}
		a.settingsService = ss
	}
	return a.settingsService.Load()
}

func (a *App) UpdateAppSettings(partial map[string]any) error {
	if a.settingsService == nil {
		ss, err := config.NewSettingsService()
		if err != nil {
			return err
		}
		a.settingsService = ss
	}
	if err := a.settingsService.Update(partial); err != nil {
		return err
	}
	return nil
}
