<template>
  <div ref="animatedDiv" class="w-full h-screen flex flex-col justify-center text-white fade-in overflow-hidden">
    <div class="-translate-y-10 flex h-4/5 flex-col justify-around items-center text-center relative">
      <div class="confetti-wrapper">
        <ConfettiExplosion :particleCount="particleCount" :force="1" :stageHeight="1000" :duration="5000" />
      </div>
      <img src="./assets/man.png" />
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
const particleCount = ref(100);

onMounted(() => {
  const options = {
    root: null, // это окно просмотра
    threshold: 0.1, // процент видимости
  };

  particleCount.value = isMobile() ? 100 : 300;

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

function isMobile() {
  if (/Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)) {
    return true
  } else {
    return false
  }
}
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

.confetti-wrapper {
  width: 0;
  height: 0;
  position: absolute;
}
</style>
