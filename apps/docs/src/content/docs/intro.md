---
title: Introduction
description: What is cliQ?
---

**cliQ** pronunciation: /klɪk/, like "click"
cliQ is a lightweight tool that transforms complex CLI commands into an intuitive and user-friendly graphical user interface (GUI).

By defining a command template with placeholder variables, cliQ automatically generates a corresponding GUI. Users can input parameters via clicks and selections—no need to memorize command syntax—greatly improving usability and efficiency.

## Key Features
- **Command Template Definition**: Input a CLI command to generate a dynamic form interface.
- **Import/Export Templates**: Share templates (`.cliqfile.yaml`) with others or import from teammates.
- **Multiple Input Components**: Supports file pickers, number inputs, dropdowns, checkboxes, and more—adapting to various parameter types.
- **Cross-Platform**: Works on Windows, macOS, and Linux.
- **Template Marketplace**: Upload or download templates for common tools (e.g., ImageMagick, ffmpeg, pngquant) to build a shared ecosystem.

## How to Use

Download the appropriate package from the Releases page and double-click to run.

**Windows**:
Run the `.exe` file directly. Compatible with Windows 10 and later.
cliQ depends on Webview. If not installed, a prompt will appear on first launch—follow the instructions to install it.

**macOS**:
Compatible with macOS 10.13 and later.
Since the app is not signed, you may see a "damaged" warning. To resolve:
1. Open Terminal and run:
```bash
sudo xattr -dr com.apple.quarantine /Users/colin/Downloads/cliq.app
# Note: Update the path as needed. You can drag the app into Terminal to auto-fill the path.
```
2. Launch cliQ. If it opens successfully, the app is now trusted.
3. (Optional) Move cliQ to your Applications folder for easier access.

**Linux**:
cliQ supports Linux, but the current release hasn't been fully tested. Feedback via GitHub Issues is welcome.

## Tech Stack
- Built with Wails for cross-platform support
- Frontend: Vue 3 + TypeScript
- UI Components: PrimeVue V4, Vite for build, partial styling with TailwindCSS
