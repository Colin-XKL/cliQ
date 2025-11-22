<template>
  <div v-if="template" class="p-4 bg-blue-50 rounded-md">
    <div class="mb-2">
      <h3 v-if="template.name" class="text-lg font-bold text-gray-800">{{ template.name }}</h3>
      <p v-if="template.description" class="text-sm text-gray-600 mt-1">{{ template.description }}</p>
    </div>
    <div class="flex flex-wrap items-center gap-4 text-xs text-blue-500 mt-2">
      <span v-if="template.author" class="flex items-center">
        <i class="pi pi-user mr-1"></i>
        作者: {{ template.author }}
      </span>
      <span v-if="template.version" class="flex items-center">
        <i class="pi pi-tag mr-1"></i>
        版本: {{ template.version }}
      </span>
      <!-- Add more metadata fields as needed -->
    </div>
  </div>
</template>

<script setup lang="ts">
import { models } from '@/wailsjs/go/models';

interface Props {
  template: models.TemplateFile | null;
}

defineProps<Props>();
</script>

<style scoped>
</style>