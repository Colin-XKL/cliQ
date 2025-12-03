---
title: cliqfile YAML Syntax
description: Detailed documentation for cliqfile.yaml syntax
---

## Overview

A cliqfile is a YAML configuration file with the `.cliqfile.yaml` extension that defines command-line templates for the cliQ application. These files allow users to transform complex CLI commands into user-friendly GUI forms with appropriate input components.

## File Structure

A cliqfile follows this basic structure:

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

## Metadata Fields

### `name` (required)
- **Type:** String
- **Description:** A human-readable name for the template that appears in the UI
- **Example:** `FFmpeg Video Processing Tool`

### `description` (required)
- **Type:** String
- **Description:** A brief explanation of what the template does
- **Example:** `Perform video format conversion, audio extraction, compression, and resolution adjustment using FFmpeg`

### `version` (required)
- **Type:** String
- **Description:** The version of this specific template
- **Example:** `"1.0"`

### `author` (required)
- **Type:** String
- **Description:** The creator of the template
- **Example:** `colin`

### `cliq_template_version` (required)
- **Type:** String
- **Description:** The version of the cliqfile specification used by this template. This helps cliQ parse the file correctly.
- **Example:** `"1.0"`

## Commands Section

The `cmds` field is a list of command definitions. Each template can define multiple related commands.

### Command Fields

#### `id` (optional)
- **Type:** String
- **Description:** Unique identifier for the command (auto-generated if not provided)

#### `name` (required)
- **Type:** String
- **Description:** Name of the command that appears in the UI
- **Example:** `Format Conversion`

#### `description` (required)
- **Type:** String
- **Description:** Brief explanation of what this specific command does
- **Example:** `Convert video files to other formats`

#### `command` (required)
- **Type:** String
- **Description:** The CLI command template string. Variables are referenced using `{{variable_name}}` syntax.
- **Example:** `"ffmpeg -i {{input_file}} -codec copy {{output_file}}"`

#### `variables` (required)
- **Type:** List of variable definitions
- **Description:** Defines the variables used in the command template and their corresponding UI components.

## Variable Definitions

Each variable in the `variables` list has the following fields:

### `name` (required)
- **Type:** String
- **Description:** The variable name used in the command template. Must be unique within the command.
- **Example:** `input_file`

### `type` (required)
- **Type:** String
- **Description:** Defines the type of UI component to display. See Variable Types section below.
- **Example:** `file_input`

### `arg_name` (optional)
- **Type:** String
- **Description:** Alternative name to use when the variable is used as a command-line argument. If not specified, uses the `name` field.
- **Example:** `skip-if-larger`

### `label` (required)
- **Type:** String
- **Description:** The display label shown in the UI for this variable
- **Example:** `Input File`

### `description` (required)
- **Type:** String
- **Description:** A longer description explaining the purpose of this variable
- **Example:** `Select the video file to convert`

### `required` (required)
- **Type:** Boolean
- **Description:** Whether the user must provide a value for this variable
- **Example:** `true`

### `options` (optional)
- **Type:** Map
- **Description:** Type-specific configuration options. The content varies depending on the variable type.

## Variable Types and Options

The following variable types are supported:

### `string`
- **UI Component:** Text input field
- **Purpose:** General text input
- **Options:**
  - `default`: Default text value (string)
  - `placeholder`: Placeholder text shown in the input (string)

### `file_input`
- **UI Component:** File picker dialog for input files
- **Purpose:** Selecting existing files
- **Options:**
  - `file_types`: List of allowed file extensions (array of strings, e.g., `[".png", ".jpg"]`)
  - `default`: Default file path (string)

### `file_output`
- **UI Component:** File save dialog for output files
- **Purpose:** Selecting destination for output files
- **Options:**
  - `file_types`: List of allowed file extensions (array of strings)
  - `default`: Default file path, can include variable interpolation (string)

### `number`
- **UI Component:** Number input field with validation
- **Purpose:** Numeric inputs
- **Options:**
  - `default`: Default number value (number)
  - `min`: Minimum allowed value (number)
  - `max`: Maximum allowed value (number)
  - `step`: Step increment (number)

### `boolean`
- **UI Component:** Checkbox
- **Purpose:** Boolean flags that can be turned on/off
- **Options:**
  - `default`: Default checked state (boolean)
  - When used in command: if true, the flag may be included in the command; if false, it may be omitted

### `select`
- **UI Component:** Dropdown selection
- **Purpose:** Choose from predefined options
- **Options:**
  - `options`: List of available choices (array of strings)
  - `default`: Default selected option (string)

## Advanced Features

### Variable Interpolation in Defaults

Default values can reference other variables using the `{{variable_name}}` syntax:

```yaml
- name: output_file
  type: file_output
  label: Output File
  description: Select the location and format to save after conversion
  required: true
  options:
    file_types: [".mp4", ".mkv", ".avi", ".mov", ".webm"]
    default: "{{input_file}}_converted.mp4"  # Uses the value of input_file variable
```

### Multiple Commands in One Template

A single template can define multiple related commands:

```yaml
cmds:
  - name: Format Conversion
    description: Convert video files to other formats
    command: "ffmpeg -i {{input_file}} -codec copy {{output_file}}"
    variables:
      # ... variable definitions
  - name: Audio Extraction
    description: Extract audio tracks from video files
    command: "ffmpeg -i {{input_file}} -vn {{output_file}}"
    variables:
      # ... variable definitions
```

## Validation Rules

1. **Template Level:**
   - Name, version, and cliq_template_version cannot be empty
   - Must contain at least one command

2. **Command Level:**
   - Name and command strings cannot be empty
   - All variable names must be unique within each command

3. **Variable Level:**
   - Name and label cannot be empty
   - Type must be one of the supported types
   - Variable names must follow valid identifier rules (alphanumeric characters and underscores)

## Best Practices

1. **Descriptive Labels:** Use clear, user-friendly labels for variables
2. **Helpful Descriptions:** Provide detailed descriptions to help users understand each variable's purpose
3. **File Type Restrictions:** Limit file types when possible to ensure compatibility
4. **Default Values:** Provide sensible default values to reduce user input
5. **Variable Naming:** Use consistent and descriptive variable names (e.g., `input_file`, `output_file`, `max_size`)
6. **Command Safety:** Ensure the command templates are safe to execute and validate user inputs where possible

## Example Template

Here's a complete example showing various variable types:

```yaml
# Template metadata
name: PNGQuant Compression Tool
description: Efficiently compress PNG images using pngquant
version: "1.0"
author: user123
cliq_template_version: "1.0"

cmds:
  - name: Compress
    description: Compress PNG file
    command: "pngquant {{input_file}} --output {{output_file}}"
    variables:
      - name: input_file
        type: file_input
        label: Input File
        description: Select the PNG file to compress
        required: true
        options:
          file_types: [".png"]

      - name: output_file
        type: file_output
        label: Output File
        description: Select the location to save after compression
        required: true
        options:
          file_types: [".png"]
          default: "{{input_file}}_compressed.png"
```