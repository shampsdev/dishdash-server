<template>
  <div class="bg-[#F4F4F4] rounded-xl py-4 flex items-center justify-between px-8" :class="{ 'flex-row-reverse': reverse }">
    <label class="uppercase font-medium text-5xl">{{ number }}</label>
    <p>
      {{ title }}
      <br />
      {{ description }}
    </p>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  number: string;
  title: string;
  description: string;
  reverse: boolean;
}>();
</script>

<style scoped></style>
