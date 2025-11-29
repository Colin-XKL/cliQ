---
title: cliqfile 语法
description: cliqfile.yaml 语法详解
---

## 概述

`cliqfile` 是一个带有 `.cliqfile.yaml` 扩展名的 YAML 配置文件，用于为 `cliQ` 应用程序定义命令行模板。这些文件允许用户将复杂的 CLI 命令转换为带有适当输入组件的用户友好型 GUI 表单。

## 文件结构

`cliqfile` 遵循以下基本结构：

```yaml
# Template metadata
name:            # Name of the template
description:     # Description of the template
version:         # Version of the template
author:          # Author of the template
cliq_template_version:  # Specification version for parsing (currently "1.0")

# Commands section
cmds:            # List of command definitions
  - name:        # Name of the command
    description: # Description of what the command does
    command:     # The actual command template with variables
    variables:   # List of variables for the command
      - name:    # Name of the variable (used in command template)
        type:    # Type of input component
        label:   # Display label for the input
        description: # Description of the variable
        required:    # Whether the variable is required (true/false)
        options:     # Additional configuration options specific to the type
```

## 元数据字段

### `name` (必填)
- **类型:** String
- **描述:** 模板的名称，将显示在 cliQ 界面的模板列表中。
- **示例:** `FFmpeg 视频处理工具`

### `description` (必填)
- **类型:** String
- **描述:** 模板的详细描述，帮助用户理解其功能。
- **示例:** `使用 FFmpeg 进行视频格式转换、提取音频、压缩和调整分辨率等操作`

### `version` (必填)
- **类型:** String
- **描述:** 模板的版本号，遵循语义化版本规范。
- **示例:** `"1.0"`

### `author` (必填)
- **类型:** String
- **描述:** 模板作者的姓名或联系方式。
- **示例:** `colin`

### `cliq_template_version` (必填)
- **类型:** String
- **描述:** 此模板使用的 cliqfile 规范的版本。这有助于 cliQ 正确解析文件。
- **示例:** `"1.0"`

## 命令部分

`cmds` 字段是命令定义的列表。每个模板可以定义多个相关命令。

### 命令字段

#### `id` (可选)
- **类型:** String
- **描述:** 命令的唯一标识符（如果未提供，则自动生成）

#### `name` (必填)
- **类型:** String
- **描述:** 显示在 UI 中的命令名称
- **示例:** `格式转换`

#### `description` (必填)
- **类型:** String
- **描述:** 对此特定命令功能的简要说明
- **示例:** `将视频文件转换为其他格式`

#### `command` (必填)
- **类型:** String
- **描述:** CLI 命令模板字符串。使用 `{{variable_name}}` 语法引用变量。
- **示例:** `"ffmpeg -i {{input_file}} -codec copy {{output_file}}"`

#### `variables` (必填)
- **类型:** 变量定义列表
- **描述:** 定义命令模板中使用的变量及其对应的 UI 组件。

## 变量定义

`variables` 列表中的每个变量都包含以下字段：

### `name` (必填)
- **类型:** String
- **描述:** 命令模板中使用的变量名。在命令内部必须是唯一的。
- **示例:** `input_file`

### `type` (必填)
- **类型:** String
- **描述:** 定义要显示的 UI 组件类型。请参阅下面的变量类型部分。
- **示例:** `file_input`

### `arg_name` (可选)
- **类型:** String
- **描述:** 当变量用作命令行参数时使用的备用名称。如果未指定，则使用 `name` 字段。
- **示例:** `skip-if-larger`

### `label` (必填)
- **类型:** String
- **描述:** 此变量在 UI 中显示的标签
- **示例:** `输入文件`

### `description` (必填)
- **类型:** String
- **描述:** 解释此变量用途的较长描述
- **示例:** `选择要转换的视频文件`

### `required` (必填)
- **类型:** Boolean
- **描述:** 用户是否必须为此变量提供值
- **示例:** `true`

### `options` (可选)
- **类型:** Map
- **描述:** 特定于类型的配置选项。内容因变量类型而异。

## 变量类型和选项

支持以下变量类型：

### `string`
- **UI 组件:** 文本输入字段
- **用途:** 通用文本输入
- **选项:**
  - `default`: 默认文本值 (string)
  - `placeholder`: 输入中显示的占位符文本 (string)

### `file_input`
- **UI 组件:** 用于输入文件的文件选择器对话框
- **用途:** 选择现有文件
- **选项:**
  - `file_types`: 允许的文件扩展名列表 (字符串数组，例如 `[".png", ".jpg"]`)
  - `default`: 默认文件路径 (string)

### `file_output`
- **UI 组件:** 用于输出文件的文件保存对话框
- **用途:** 为输出文件选择目标位置
- **选项:**
  - `file_types`: 允许的文件扩展名列表 (字符串数组)
  - `default`: 默认文件路径，可以包含变量插值 (string)

### `number`
- **UI 组件:** 带验证的数字输入字段
- **用途:** 数字输入
- **选项:**
  - `default`: 默认数值 (number)
  - `min`: 允许的最小值 (number)
  - `max`: 允许的最大值 (number)
  - `step`: 步长增量 (number)

### `boolean`
- **UI 组件:** 复选框
- **用途:** 可以打开/关闭的布尔标志
- **选项:**
  - `default`: 默认选中状态 (boolean)
  - 在命令中使用时：如果为 true，则该标志可能包含在命令中；如果为 false，则可能被省略

### `select`
- **UI 组件:** 下拉选择
- **用途:** 从预定义选项中选择
- **选项:**
  - `options`: 可用选项列表 (字符串数组)
  - `default`: 默认选定选项 (string)

## 高级功能

### 默认值中的变量插值

默认值可以使用 `{{variable_name}}` 语法引用其他变量：

```yaml
- name: output_file
  type: file_output
  label: 输出文件
  description: 选择转换后保存的位置和格式
  required: true
  options:
    file_types: [".mp4", ".mkv", ".avi", ".mov", ".webm"]
    default: "{{input_file}}_converted.mp4"  # 使用 input_file 变量的值
```

### 单个模板中的多个命令

一个模板可以定义多个相关命令：

```yaml
cmds:
  - name: 格式转换
    description: 将视频文件转换为其他格式
    command: "ffmpeg -i {{input_file}} -codec copy {{output_file}}"
    variables:
      # ... 变量定义
  - name: 提取音频
    description: 从视频文件中提取音频轨道
    command: "ffmpeg -i {{input_file}} -vn {{output_file}}"
    variables:
      # ... 变量定义
```

## 验证规则

1. **模板级别:**
   - 名称、版本和 cliq_template_version 不能为空
   - 必须至少包含一个命令

2. **命令级别:**
   - 名称和命令字符串不能为空
   - 每个命令中的所有变量名必须是唯一的

3. **变量级别:**
   - 名称和标签不能为空
   - 类型必须是支持的类型之一
   - 变量名必须遵循有效的标识符规则（字母数字字符和下划线）

## 最佳实践

1. **描述性标签:** 为变量使用清晰、用户友好的标签
2. **有用的描述:** 提供详细的描述，以帮助用户理解每个变量的用途
3. **文件类型限制:** 尽可能限制文件类型以确保兼容性
4. **默认值:** 提供合理的默认值以减少用户输入
5. **变量命名:** 使用一致且描述性的变量名（例如 `input_file`、`output_file`、`max_size`）
6. **命令安全:** 确保命令模板可以安全执行，并在可能的情况下验证用户输入

## 示例模板

这是一个完整的示例，展示了各种变量类型：

```yaml
# 模板元数据
name: PNGQuant 压缩工具
description: 使用 pngquant 高效压缩 PNG 图片
version: "1.0"
author: user123
cliq_template_version: "1.0"

cmds:
  - name: 压缩
    description: 压缩 PNG 文件
    command: "pngquant {{input_file}} --output {{output_file}}"
    variables:
      - name: input_file
        type: file_input
        label: 输入文件
        description: 选择要压缩的 PNG 文件
        required: true
        options:
          file_types: [".png"]

      - name: output_file
        type: file_output
        label: 输出文件
        description: 选择压缩后保存的位置
        required: true
        options:
          file_types: [".png"]
          default: "{{input_file}}_compressed.png"
```
