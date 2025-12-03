<template>
  <div class="max-w-md mx-auto mt-10">
    <Card>
      <template #title>Register</template>
      <template #content>
        <div class="flex flex-col gap-4">
          <div class="flex flex-col gap-2">
            <label for="username">Username</label>
            <InputText id="username" v-model="username" />
          </div>
          <div class="flex flex-col gap-2">
            <label for="email">Email</label>
            <InputText id="email" v-model="email" />
          </div>
          <div class="flex flex-col gap-2">
            <label for="password">Password</label>
            <InputText id="password" v-model="password" type="password" />
          </div>
          <Button label="Register" @click="handleRegister" :loading="loading" />
          <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useUserStore } from '../stores/user';
import { useRouter } from 'vue-router';
import Card from 'primevue/card';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';

const username = ref('');
const email = ref('');
const password = ref('');
const loading = ref(false);
const error = ref('');
const userStore = useUserStore();
const router = useRouter();

const handleRegister = async () => {
  loading.value = true;
  error.value = '';
  const success = await userStore.register(username.value, email.value, password.value);
  loading.value = false;
  if (success) {
    router.push('/login');
  } else {
    error.value = 'Registration failed.';
  }
};
</script>
