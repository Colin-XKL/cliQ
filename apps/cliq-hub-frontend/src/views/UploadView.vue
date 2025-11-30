<template>
  <div class="max-w-2xl mx-auto">
    <Card>
      <template #title>Upload Template</template>
      <template #content>
        <div class="flex flex-col gap-6">
          <div class="flex flex-col gap-2">
            <label for="title" class="font-medium">Title</label>
            <InputText id="title" v-model="form.title" placeholder="e.g. Simple Go API Service" />
          </div>

          <div class="flex flex-col gap-2">
            <label for="description" class="font-medium">Description</label>
            <Textarea id="description" v-model="form.description" rows="3" placeholder="Describe what this template does..." />
          </div>

          <div class="flex flex-col gap-2">
            <label for="content" class="font-medium">Template Content (YAML)</label>
            <Textarea id="content" v-model="form.content" rows="15" class="font-mono text-sm" placeholder="Paste your .cliqfile.yaml content here..." />
          </div>

          <div class="flex justify-end gap-2">
            <Button label="Cancel" severity="secondary" @click="$router.back()" />
            <Button label="Upload" @click="handleUpload" :loading="loading" />
          </div>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useTemplateStore } from '../stores/template';
import { useRouter } from 'vue-router';
import Card from 'primevue/card';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import Button from 'primevue/button';

const router = useRouter();
const templateStore = useTemplateStore();
const loading = ref(false);

const form = reactive({
  title: '',
  description: '',
  content: ''
});

const handleUpload = async () => {
  if (!form.title || !form.content) {
    alert('Title and Content are required');
    return;
  }

  loading.value = true;
  const success = await templateStore.createTemplate(form);
  loading.value = false;

  if (success) {
    router.push('/');
  } else {
    alert('Failed to upload template');
  }
};
</script>
