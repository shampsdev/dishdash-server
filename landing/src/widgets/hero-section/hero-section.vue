<template>
  <div ref="animatedDiv" class="w-full h-screen flex flex-col justify-center text-white fade-in">
    <div class="-translate-y-10 flex h-4/5 flex-col justify-around items-center text-center">
      <div class="space-y-4">
        <h1 class="text-5xl text-white">DishDash</h1>
        <p>Сервис быстрого поиска места<br /> для встречи с друзьями</p>
      </div>
      <img src="./assets/radar.png" class="h-64 md:h-96" />
      <p>Заполните форму,<br /> чтобы протестировать бета-версию</p>
      <a href="https://forms.gle/S2JPfT3kVrbNDGqNA">
        <Button class="rounded-3xl" size="lg" variant="secondary">Оставить заявку</Button>
      </a>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import Button from '@/components/ui/button/Button.vue';

const animatedDiv = ref(null);

onMounted(() => {
  const options = {
    root: null, // это окно просмотра
    threshold: 0.1, // процент видимости
  };

  const callback = (entries: IntersectionObserverEntry[]) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.classList.add('visible');
      }
    });
  };

  const observer = new IntersectionObserver(callback, options);
  if (animatedDiv.value) observer.observe(animatedDiv.value);
});
</script>

<style scoped>
.fade-in {
  opacity: 0;
  transform: translateY(10px);
  transition: opacity 0.5s ease-in-out, transform 0.5s ease-in-out;
}

.fade-in.visible {
  opacity: 1;
  transform: translateY(0);
}
</style>
