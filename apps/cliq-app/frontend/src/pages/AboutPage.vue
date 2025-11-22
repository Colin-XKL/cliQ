<template>
  <div class="about-page">
    <div class="about-content">
      <h2 class="text-2xl font-bold mb-6">关于 cliQ</h2>

      <div class="section mb-6">
        <h3 class="text-xl font-semibold mb-3">应用概述</h3>
        <p class="text-gray-700 mb-3">
          cliQ 是一个将复杂的命令行界面（CLI）命令转化为直观、易用图形用户界面的应用程序。
          通过使用 cliQ，您可以轻松地创建、管理和执行模板化的命令，而无需记住复杂的命令行参数。
        </p>
      </div>

      <div class="section mb-6">
        <h3 class="text-xl font-semibold mb-3">核心概念</h3>
        <div class="space-y-4">
          <div>
            <h4 class="font-medium text-gray-800 mb-2">模板（Template）</h4>
            <p class="text-gray-700 ml-4">
              模板是 cliQ 的核心概念，它定义了一个或多个命令及其参数。模板文件通常以 .cliqfile.yaml 格式保存，
              包含命令的名称、描述、参数定义等信息。
            </p>
          </div>
          <div>
            <h4 class="font-medium text-gray-800 mb-2">命令（Command）</h4>
            <p class="text-gray-700 ml-4">
              命令是模板中定义的可执行操作，通常对应于一个 CLI 命令。命令可以包含多个参数，
              这些参数在 UI 中以表单形式呈现，方便用户输入。
            </p>
          </div>
          <div>
            <h4 class="font-medium text-gray-800 mb-2">参数（Parameter）</h4>
            <p class="text-gray-700 ml-4">
              参数是命令的输入变量，可以是字符串、数字、布尔值等类型。cliQ
              会根据参数定义自动生成相应的输入控件。
            </p>
          </div>
        </div>
      </div>

      <div class="section mb-6">
        <h3 class="text-xl font-semibold mb-3">主要功能</h3>
        <ul class="list-disc pl-6 text-gray-700 space-y-2">
          <li>导入和管理 CLI 命令模板</li>
          <li>通过直观界面执行命令</li>
          <li>模板生成器 - 轻松创建新的命令模板</li>
          <li>收藏夹功能 - 保存常用模板以便快速访问</li>
          <li>URL 导入 - 从远程链接导入模板</li>
          <li>实时命令输出 - 直观查看命令执行结果</li>
        </ul>
      </div>

      <div class="section mb-6">
        <h3 class="text-xl font-semibold mb-3">如何使用</h3>
        <ol class="list-decimal pl-6 text-gray-700 space-y-2">
          <li>在主界面导入模板文件或从 URL 导入</li>
          <li>选择要执行的命令</li>
          <li>填写命令所需的参数</li>
          <li>点击执行查看结果</li>
        </ol>
      </div>

      <div class="section">
        <h3 class="text-xl font-semibold mb-3">帮助资源</h3>
        <div class="help-links grid grid-cols-1 md:grid-cols-2 gap-4">
          <div @click="OpenUrlUseSystemBrowser('https://github.com/Colin-XKL/cliq')"
            class="help-link p-4 border rounded-lg hover:bg-gray-50 transition-colors">
            <h4 class="font-medium text-blue-600">GitHub 仓库</h4>
            <p class="text-sm text-gray-600 mt-1">访问项目源代码和贡献指南</p>
          </div>

          <div @click="OpenUrlUseSystemBrowser('https://github.com/Colin-XKL/cliq/issues')"
            class="help-link p-4 border rounded-lg hover:bg-gray-50 transition-colors">
            <h4 class="font-medium text-blue-600">问题报告</h4>
            <p class="text-sm text-gray-600 mt-1">提交 bug 和功能请求</p>
          </div>

          <div @click="OpenUrlUseSystemBrowser('https://github.com/Colin-XKL/cliq/wiki')"
            class="help-link p-4 border rounded-lg hover:bg-gray-50 transition-colors">
            <h4 class="font-medium text-blue-600">使用文档</h4>
            <p class="text-sm text-gray-600 mt-1">详细使用说明和教程</p>
          </div>

          <div @click="OpenUrlUseSystemBrowser('https://github.com/Colin-XKL/cliq/discussions')"
            class="help-link p-4 border rounded-lg hover:bg-gray-50 transition-colors">
            <h4 class="font-medium text-blue-600">社区讨论</h4>
            <p class="text-sm text-gray-600 mt-1">与其他用户交流经验</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { OpenUrlUseSystemBrowser } from '@/utils/open_url'
</script>

<style scoped>
@reference "tailwindcss";

.about-page {
  @apply p-4;
}

.about-content {
  @apply max-w-4xl mx-auto;
}

.section {
  @apply p-4 bg-white rounded-lg shadow-sm;
}

.help-link {
  @apply block no-underline;
}
</style>