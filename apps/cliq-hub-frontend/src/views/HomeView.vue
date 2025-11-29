<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Template Market</h1>
    <div v-if="loading" class="text-center py-8">Loading...</div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card v-for="template in templateStore.templates" :key="template.ID" class="hover:shadow-lg transition-shadow">
        <template #title>{{ template.title }}</template>
        <template #subtitle>by {{ template.author?.username }}</template>
        <template #content>
          <p class="text-gray-600 line-clamp-3">{{ template.description }}</p>
        </template>
        <template #footer>
          <div class="flex justify-between items-center mt-4">
            <span class="text-sm text-gray-500">{{ template.downloads }} downloads</span>
            <Button label="View Details" size="small" @click="viewDetail(template.ID)" />
          </div>
        </template>
      </Card>
    </div>
    <div v-if="!loading && templateStore.templates.length === 0" class="text-center py-12 text-gray-500">
      No templates found. Be the first to upload one!
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useTemplateStore } from '../stores/template';
import { useRouter } from 'vue-router';
import Card from 'primevue/card';
import Button from 'primevue/button';

const templateStore = useTemplateStore();
const router = useRouter();
const loading = ref(true);

onMounted(async () => {
  await templateStore.fetchTemplates();
  loading.value = false;
});

const viewDetail = (id: number) => {
  router.push(`/templates/${id}`);
};
</script>
