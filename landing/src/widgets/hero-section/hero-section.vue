<template>
  <div ref="animatedDiv" class="w-full items-center h-screen flex flex-col justify-center text-white overflow-hidden">

    <div class="-translate-y-20 h-100 flex flex-col justify-around items-center text-center relative">
      <div class="confetti-wrapper">
        <ConfettiExplosion v-if="showConfetti" :particleCount="particleCount" :force="1" :stageHeight="1000"
          :duration="5000" />
      </div>
      <img src="./assets/man.png" />

      <a href="https://t.me/dishdash_bot?start=landing" class="translate-y-10">
        <CoolButton class="rounded-3xl" size="lg" variant="secondary">
          Заходи!
          <img src="./assets/telegram-svgrepo-com.svg" color="black" class="ms-2" width="16px" />
        </CoolButton>
      </a>
    </div>
  </div>
</template>


<script setup lang="ts">

import { onMounted, ref } from 'vue';
import ConfettiExplosion from "vue-confetti-explosion";
import CoolButton from '@/components/ui/button/CoolButton.vue';

const animatedDiv = ref(null);
const particleCount = ref(100);
const showConfetti = ref(false);

onMounted(() => {
  const options = {
    root: null, // это окно просмотра
    threshold: 0.1, // процент видимости
  };

  particleCount.value = isMobile() ? 100 : 300;
  showConfetti.value = true;

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
