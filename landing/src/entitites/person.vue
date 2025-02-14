<template>
  <div class="flex h-80 sm:h-96 justify-center flex-col items-center space-y-4">
    <img class="rounded-full w-1/2" :src="image" />
    <div>
      <p class="text-center text-wrap">{{ name }}</p>
      <p class="text-center text-wrap">{{ surname }}</p>
    </div>
    <a :href="telegram">
      <Button size="lg" variant="secondary" class="rounded-3xl">Telegram</Button>
    </a>
  </div>
</template>

<script setup lang="ts">
import Button from '@/components/ui/button/Button.vue';
import { IPerson } from './person.interface';

defineProps<IPerson>()
</script>