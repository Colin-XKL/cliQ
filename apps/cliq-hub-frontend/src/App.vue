<template>
  <div class="min-h-screen bg-gray-50 text-gray-900">
    <nav class="bg-white shadow-sm p-4 mb-4">
      <div class="container mx-auto flex justify-between items-center">
        <router-link to="/" class="text-xl font-bold text-blue-600">Cliq Hub</router-link>
        <div class="flex items-center gap-4">
          <router-link to="/" class="hover:text-blue-600">Market</router-link>
          <template v-if="userStore.isAuthenticated">
            <router-link to="/upload" class="hover:text-blue-600">Upload</router-link>
            <span class="text-gray-500">Hi, {{ userStore.user?.username }}</span>
            <button @click="logout" class="text-red-500 hover:text-red-700">Logout</button>
          </template>
          <template v-else>
            <router-link to="/login" class="hover:text-blue-600">Login</router-link>
            <router-link to="/register" class="hover:text-blue-600">Register</router-link>
          </template>
        </div>
      </div>
    </nav>
    <main class="container mx-auto p-4">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from './stores/user';
import { useRouter } from 'vue-router';

const userStore = useUserStore();
const router = useRouter();

const logout = () => {
  userStore.logout();
  router.push('/login');
};
</script>
