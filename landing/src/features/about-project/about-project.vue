<template>
  <div class="w-full space-y-6">
    <div class="space-y-2 text-center ">
      <h3 class="text-xl text-[#CBCBCB]">О проекте</h3>
      <p>
        Какие проблемы при выборе<br />
        места решает сервис
      </p>
    </div>

    <div class="space-y-4">
      <FeatureItem
        v-for="(item, index) in items"
        :key="index"
        :number="item.number"
        :title="item.title"
        :description="item.description"
        :reverse="index % 2 !== 0"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import FeatureItem from "@/entitites/feature-item.vue";

const items = [
  { number: '01', title: 'Быстрый выбор', description: 'не тратьте время' },
  { number: '02', title: 'Персонализация', description: 'Найдите что-то особенное' },
  { number: '03', title: 'Необычные места', description: 'Пробуйте новое' },
  { number: '04', title: 'Удобные настройки', description: 'Выбирайте с друзьями' },
];

</script>

<style scoped></style>
