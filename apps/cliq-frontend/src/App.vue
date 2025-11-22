<script setup lang="ts">
import { ref } from 'vue'
const msg = ref('Hello World from cliQ Frontend')
</script>

<template>
  <div style="display:flex;align-items:center;justify-content:center;height:100vh;font-family:sans-serif">
    <h1>{{ msg }}</h1>
  </div>
</template>