package handlers

import (
    "context"
    "crypto/md5"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"

    "github.com/wailsapp/wails/v2/pkg/runtime"

    "repo/shared-go-lib/models"

    "gopkg.in/yaml.v3"
)

// FileHandler handles file-related operations
type FileHandler struct {
	ctx context.Context
}

// NewFileHandler creates a new file handler
func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (fh *FileHandler) Startup(ctx context.Context) {
	fh.ctx = ctx
}

// OpenFileDialog opens a file dialog and returns the selected file path
func (fh *FileHandler) OpenFileDialog() (string, error) {
	// Default to all files if no specific filters are provided
	options := runtime.OpenDialogOptions{
		Title: "选择输入文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
	}

	return runtime.OpenFileDialog(fh.ctx, options)
}

// OpenFileDialogWithFilters opens a file dialog with specific file type filters and returns the selected file path
func (fh *FileHandler) OpenFileDialogWithFilters(filters []runtime.FileFilter) (string, error) {
	options := runtime.OpenDialogOptions{
		Title:   "选择输入文件",
		Filters: filters,
	}

	if len(filters) == 0 {
		// If no filters are provided, default to all files
		options.Filters = []runtime.FileFilter{
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		}
	}

	return runtime.OpenFileDialog(fh.ctx, options)
}

// SaveFileDialog opens a save file dialog and returns the selected file path
func (fh *FileHandler) SaveFileDialog() (string, error) {
	options := runtime.SaveDialogOptions{
		Title: "保存输出文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
	}

	return runtime.SaveFileDialog(fh.ctx, options)
}

// SaveYAMLToFile opens a save file dialog and saves the YAML content to the selected file
func (fh *FileHandler) SaveYAMLToFile(yamlContent string) error {
	if yamlContent == "" {
		return fmt.Errorf("YAML内容不能为空")
	}

	// Open save file dialog
	filePath, err := runtime.SaveFileDialog(fh.ctx, runtime.SaveDialogOptions{
		Title: "保存模板文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "CLIQ模板文件 (*.yaml, *.yml)",
				Pattern:     "*.yaml;*.yml",
			},
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
		DefaultFilename: "template.cliqfile.yaml",
	})

	if err != nil {
		return fmt.Errorf("打开保存对话框失败: %w", err)
	}

	// 用户取消选择
	if filePath == "" {
		return fmt.Errorf("未选择保存路径")
	}

	// Write content to file
	err = os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// ExportTemplateToFile 将模板导出为文件
func (fh *FileHandler) ExportTemplateToFile(template *models.TemplateFile, filePath string) error {
	if template == nil {
		return fmt.Errorf("模板不能为空")
	}

	if filePath == "" {
		return fmt.Errorf("文件路径不能为空")
	}

	// 序列化模板为YAML
	data, err := yaml.Marshal(template)
	if err != nil {
		return fmt.Errorf("序列化模板失败: %w", err)
	}

	// 写入文件
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

func getCommandParts(cmdTemplateStr string, variables map[string]interface{}) ([]string, error) {
	// 替换命令模板中的变量
	commandStr := cmdTemplateStr
	for name, value := range variables {
		// TODO: 根据变量类型进行格式化，例如布尔值转换为 --flag 或空
		commandStr = strings.ReplaceAll(commandStr, fmt.Sprintf("{{%s}}", name), fmt.Sprintf("%v", value))
	}

	parts := strings.Fields(commandStr)
	return parts, nil
}

// ExecuteCommand executes a shell command with the given input and output file paths
func (fh *FileHandler) ExecuteCommand(template *models.TemplateFile, commandID string, variables map[string]interface{}) (string, error) {
	if template == nil {
		return "", fmt.Errorf("模板未加载")
	}
	if variables == nil {
		variables = make(map[string]interface{})
	}
	// 根据 commandID 查找对应的命令
	var selectedCommand models.Command
	found := false
	for _, cmd := range template.Cmds {
		if cmd.ID == commandID {
			selectedCommand = cmd
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("未找到命令: %s", commandID)
	}

	// 替换命令模板中的变量
	cmdTemplateStr := selectedCommand.Command
	parts, err := getCommandParts(cmdTemplateStr, variables)
	if err != nil {
		return "", fmt.Errorf("获取命令文本失败: %w", err)
	}
	if len(parts) == 0 {
		return "", fmt.Errorf("命令为空")
	}

	name := parts[0]
	args := parts[1:]

	cmd := exec.Command(name, args...)

	// 执行命令并获取输出
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("执行命令失败: %w\n%s", err, string(out))
	}

	return strings.TrimSpace(string(out)), nil
}

func (fh *FileHandler) GetCommandText(template *models.TemplateFile, commandID string, variables map[string]interface{}) (string, error) {
	if template == nil {
		return "", fmt.Errorf("template is nil")
	}
	if variables == nil {
		variables = make(map[string]interface{})
	}
	var selectedCommand models.Command
	found := false
	for _, cmd := range template.Cmds {
		if cmd.ID == commandID {
			selectedCommand = cmd
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("未找到命令: %s", commandID)
	}

	// 替换命令模板中的变量
	cmdTemplateStr := selectedCommand.Command
	parts, err := getCommandParts(cmdTemplateStr, variables)
	if err != nil {
		return "", fmt.Errorf("获取命令文本失败: %w", err)
	}
	return strings.Join(parts, " "), nil
}

// getHashForTemplateName 生成基于模板名称的安全哈希值，防止路径遍历和特殊字符问题
func getHashForTemplateName(templateName string) string {
	hash := md5.Sum([]byte(templateName))
	return fmt.Sprintf("%x", hash)
}

// getFileNamesForTemplate 根据模板名称生成两种可能的文件名（.cliqfile.yaml 和 .cliqfile.yml）
func getFileNamesForTemplate(hashedName string) []string {
	return []string{
		fmt.Sprintf("%s.cliqfile.yaml", hashedName),
		fmt.Sprintf("%s.cliqfile.yml", hashedName),
	}
}

// getFileNamesForTemplateOriginal 根据原始模板名称生成两种可能的文件名（.cliqfile.yaml 和 .cliqfile.yml）
func getFileNamesForTemplateOriginal(originalName string) []string {
	return []string{
		fmt.Sprintf("%s.cliqfile.yaml", originalName),
		fmt.Sprintf("%s.cliqfile.yml", originalName),
	}
}

// findTemplateFile 在给定目录中查找指定模板的文件，支持两种后缀格式
func (fh *FileHandler) findTemplateFile(dirPath, hashedName, originalName string) (string, error) {
	// 首先检查使用哈希名称的文件
	hashFileNames := getFileNamesForTemplate(hashedName)
	for _, fileName := range hashFileNames {
		filePath := filepath.Join(dirPath, fileName)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil // 找到文件
		} else if !os.IsNotExist(err) {
			// 如果不是 "文件不存在" 错误，则返回具体错误 (如权限问题)
			return "", fmt.Errorf("检查模板文件状态时出错 (%s): %w", filePath, err)
		}
	}

	// 然后检查使用原始名称的文件（向后兼容）
	origFileNames := getFileNamesForTemplateOriginal(originalName)
	for _, fileName := range origFileNames {
		filePath := filepath.Join(dirPath, fileName)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil // 找到文件
		} else if !os.IsNotExist(err) {
			// 如果不是 "文件不存在" 错误，则返回具体错误 (如权限问题)
			return "", fmt.Errorf("检查模板文件状态时出错 (%s): %w", filePath, err)
		}
	}

	// If we get here, the file was not found with either naming convention
	// Return a proper error instead of a potential non-existent path with os.ErrNotExist
	return "", fmt.Errorf("未找到模板文件 (哈希名: %s, 原始名: %s, 目录: %s)", hashedName, originalName, dirPath)
}

// getFavTemplatesDirPath 获取收藏模板的存储路径
func (fh *FileHandler) getFavTemplatesDirPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("获取用户配置目录失败: %w", err)
	}
	favTemplatesDir := filepath.Join(configDir, "cliq", "fav_templates")
	return favTemplatesDir, nil
}

// ensureFavTemplatesDirExists 确保收藏模板目录存在，如果不存在则创建
func (fh *FileHandler) ensureFavTemplatesDirExists() (string, error) {
	dirPath, err := fh.getFavTemplatesDirPath()
	if err != nil {
		return "", err
	}

	// 检查目录是否存在，如果不存在则创建
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return "", fmt.Errorf("创建收藏模板目录失败: %w", err)
		}
	}

	return dirPath, nil
}

// SaveFavTemplate 保存收藏模板文件
func (fh *FileHandler) SaveFavTemplate(template *models.TemplateFile) error {
	if template == nil {
		return fmt.Errorf("模板不能为空")
	}

	// 确保收藏目录存在
	dirPath, err := fh.ensureFavTemplatesDirExists()
	if err != nil {
		return err
	}

	// 使用模板名称的哈希值作为文件名，并添加.cliqfile.yaml后缀
	hashedName := getHashForTemplateName(template.Name)
	fileName := fmt.Sprintf("%s.cliqfile.yaml", hashedName)
	filePath := filepath.Join(dirPath, fileName)

	// 序列化模板为YAML
	data, err := yaml.Marshal(template)
	if err != nil {
		return fmt.Errorf("序列化模板失败: %w", err)
	}

	// 写入文件
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("写入收藏模板文件失败: %w", err)
	}

	return nil
}

// ListFavTemplates 列出所有收藏的模板文件
func (fh *FileHandler) ListFavTemplates() ([]*models.TemplateFile, error) {
	dirPath, err := fh.ensureFavTemplatesDirExists()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("读取收藏模板目录失败: %w", err)
	}

	var templates []*models.TemplateFile
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".cliqfile.yaml") || strings.HasSuffix(file.Name(), ".cliqfile.yml")) {
			// 验证文件名是否是MD5哈希格式 (32位十六进制字符) + 后缀
			filename := file.Name()
			var isHashFormat bool

			if strings.HasSuffix(filename, ".cliqfile.yaml") {
				namePart := strings.TrimSuffix(filename, ".cliqfile.yaml")
				isHashFormat = isValidMD5Hash(namePart)
			} else if strings.HasSuffix(filename, ".cliqfile.yml") {
				namePart := strings.TrimSuffix(filename, ".cliqfile.yml")
				isHashFormat = isValidMD5Hash(namePart)
			}

			// 只处理符合哈希格式的文件
			if isHashFormat {
				filePath := filepath.Join(dirPath, file.Name())
				data, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("读取文件 %s 失败: %v\n", filePath, err)
					continue
				}

				var template models.TemplateFile
				err = yaml.Unmarshal(data, &template)
				if err != nil {
					fmt.Printf("解析文件 %s 失败: %v\n", filePath, err)
					continue
				}
				templates = append(templates, &template)
			}
		}
	}

	return templates, nil
}

// isValidMD5Hash 验证字符串是否为有效的MD5哈希值格式
func isValidMD5Hash(hash string) bool {
	if len(hash) != 32 {
		return false
	}

	for _, c := range hash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}

	return true
}

// getFilePathForTemplate 根据模板名称获取文件路径，同时兼容新旧格式和两种后缀
func (fh *FileHandler) getFilePathForTemplate(templateName string) (string, error) {
	dirPath, err := fh.getFavTemplatesDirPath()
	if err != nil {
		return "", err
	}

	hashedName := getHashForTemplateName(templateName)

	// 使用新的查找函数
	filePath, err := fh.findTemplateFile(dirPath, hashedName, templateName)
	if err != nil && os.IsNotExist(err) {
		// 如果文件不存在，返回哈希名称的默认路径（用于保存操作）
		defaultFileName := fmt.Sprintf("%s.cliqfile.yaml", hashedName)
		return filepath.Join(dirPath, defaultFileName), nil
	}

	return filePath, err
}

// DeleteFavTemplate 从收藏目录删除指定模板文件
func (fh *FileHandler) DeleteFavTemplate(templateName string) error {
	if templateName == "" {
		return fmt.Errorf("模板名称不能为空")
	}

	dirPath, err := fh.getFavTemplatesDirPath()
	if err != nil {
		return err
	}

	hashedName := getHashForTemplateName(templateName)

	// 先尝试查找存在的文件
	filePath, err := fh.findTemplateFile(dirPath, hashedName, templateName)
	if err != nil {
		return fmt.Errorf("模板文件不存在: %s", templateName)
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("模板文件不存在: %s", templateName)
	}

	// 删除文件
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("删除收藏模板文件失败: %w", err)
	}

	return nil
}

// GetFavTemplate 读取指定收藏模板文件内容
func (fh *FileHandler) GetFavTemplate(templateName string) (*models.TemplateFile, error) {
	if templateName == "" {
		return nil, fmt.Errorf("模板名称不能为空")
	}

	dirPath, err := fh.getFavTemplatesDirPath()
	if err != nil {
		return nil, fmt.Errorf("获取模板目录失败: %w", err)
	}

	hashedName := getHashForTemplateName(templateName)

	// 查找存在的文件
	filePath, err := fh.findTemplateFile(dirPath, hashedName, templateName)
	if err != nil {
		return nil, fmt.Errorf("查找模板文件失败 (模板名: %s, 哈希: %s): %w", templateName, hashedName, err)
	}

	// 检查文件是否存在 (this check is redundant since findTemplateFile already checks, but keeping for safety)
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("模板文件不存在: %s (路径: %s)", templateName, filePath)
	} else if err != nil {
		return nil, fmt.Errorf("检查模板文件状态失败: %w", err)
	}

	// Check if file is empty
	if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("模板文件为空: %s", filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取收藏模板文件失败 (路径: %s): %w", filePath, err)
	}

	// Check if read data is empty
	if len(data) == 0 {
		return nil, fmt.Errorf("读取的模板内容为空: %s", filePath)
	}

	var template models.TemplateFile
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		return nil, fmt.Errorf("解析收藏模板文件失败 (路径: %s): %w", filePath, err)
	}

	return &template, nil
}

// UpdateFavTemplate 更新指定收藏模板文件内容
func (fh *FileHandler) UpdateFavTemplate(oldTemplateName string, newTemplateName string, updatedTemplate *models.TemplateFile) error {
	if oldTemplateName == "" {
		return fmt.Errorf("原模板名称不能为空")
	}
	if newTemplateName == "" {
		return fmt.Errorf("新模板名称不能为空")
	}
	if updatedTemplate == nil {
		return fmt.Errorf("更新模板不能为空")
	}

	dirPath, err := fh.getFavTemplatesDirPath()
	if err != nil {
		return err
	}

	// 生成旧模板的哈希名称
	oldHashedName := getHashForTemplateName(oldTemplateName)

	// 查找旧模板文件
	oldFilePath, err := fh.findTemplateFile(dirPath, oldHashedName, oldTemplateName)
	if err != nil {
		return fmt.Errorf("查找原模板文件失败 (原模板名: %s, 哈希: %s): %w", oldTemplateName, oldHashedName, err)
	}

	// 检查旧文件是否存在
	if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
		return fmt.Errorf("原模板文件不存在: %s", oldTemplateName)
	}

	// 如果模板名称发生变更，则删除旧文件并创建新文件
	if oldTemplateName != newTemplateName {
		// 删除旧文件
		err = os.Remove(oldFilePath)
		if err != nil {
			return fmt.Errorf("删除原模板文件失败: %w", err)
		}

		// 生成新模板的哈希名称和路径
		newHashedName := getHashForTemplateName(newTemplateName)
		newFileName := fmt.Sprintf("%s.cliqfile.yaml", newHashedName)
		newFilePath := filepath.Join(dirPath, newFileName)

		// 序列化更新后的模板为YAML
		data, err := yaml.Marshal(updatedTemplate)
		if err != nil {
			return fmt.Errorf("序列化更新模板失败: %w", err)
		}

		// 写入新文件
		err = os.WriteFile(newFilePath, data, 0644)
		if err != nil {
			return fmt.Errorf("写入新模板文件失败: %w", err)
		}
	} else {
		// 名称未变更，直接更新原文件
		// 序列化更新后的模板为YAML
		data, err := yaml.Marshal(updatedTemplate)
		if err != nil {
			return fmt.Errorf("序列化更新模板失败: %w", err)
		}

		// 写入文件
		err = os.WriteFile(oldFilePath, data, 0644)
		if err != nil {
			return fmt.Errorf("写入收藏模板文件失败: %w", err)
		}
	}

	return nil
}
