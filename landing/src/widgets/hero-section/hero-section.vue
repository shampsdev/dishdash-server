<template>
  <div ref="animatedDiv" class="w-full h-screen flex flex-col justify-center text-white fade-in">
    <div class="-translate-y-10 flex h-4/5 flex-col justify-around items-center text-center">
      <ConfettiExplosion :particleCount="500" :force="0.9"/>
      <img src="./assets/man.png"  />
      <p>Переходите в бота</p>
      <a href="https://t.me/dishdash_bot?start=landing">
        <Button class="rounded-3xl" size="lg" variant="secondary">В БОТА</Button>
      </a>
    </div>
  </div>
</template>


<script setup lang="ts">
import { onMounted, ref } from 'vue';
import Button from '@/components/ui/button/Button.vue';
import ConfettiExplosion from "vue-confetti-explosion";

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
