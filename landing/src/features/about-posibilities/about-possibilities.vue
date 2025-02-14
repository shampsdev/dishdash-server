<template>
  <div class="w-full h-max space-y-6">
    <div ref="headerRef" class="space-y-2 text-center fade-in">
      <h3 class="text-xl text-[#CBCBCB]">Возможности</h3>
      <p>Функционал нашего приложения</p>
    </div>
    <img v-for="(image, index) in imagesList" :key="index" :src="image.src" :alt="image.alt" class="w-full fade-in"
      :ref="el => setImageRef(el, index)" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import createLink from './assets/create-link.jpg';
import inviteFriends from './assets/invite-friends.jpg';
import setupFilters from './assets/setup-filters.jpg';
import swipes from './assets/swipes.jpg';
import chooseVariant from './assets/choose-variant.png';

const imagesList = [
  { src: createLink, alt: 'create-link' },
  { src: inviteFriends, alt: 'invite-friends' },
  { src: setupFilters, alt: 'setup-filters' },
  { src: swipes, alt: 'swipes' },
  { src: chooseVariant, alt: 'choose-variant' }
];

const imageRefs = ref<any[]>([]);
const headerRef = ref<HTMLElement | null>(null);

onMounted(() => {
  const options = {
    root: null,
    threshold: 0.1,
  };

  const callback = (entries: IntersectionObserverEntry[]) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.classList.add('visible');
      }
    });
  };

  const observer = new IntersectionObserver(callback, options);
  imageRefs.value.forEach(image => {
    observer.observe(image);
  });

  if (headerRef.value) {
    observer.observe(headerRef.value);
  }
});

const setImageRef = (el: any, index: number) => {
  if (el) {
    imageRefs.value[index] = el;
  }
};
</script>

<style scoped>
.fade-in {
  opacity: 0;
  transform: translateY(50px);
  transition: opacity 0.5s ease-in-out, transform 0.5s ease-in-out;
}

.fade-in.visible {
  opacity: 1;
  transform: translateY(0);
}
</style>
