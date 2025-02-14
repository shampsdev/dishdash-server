<template>
  <div class="min-h-screen w-full rounded-3xl bg-white py-10 md:py-20 -translate-y-4">
    <div class="flex gap-y-8 flex-col max-w-[600px] w-[90%] mx-auto items-center">
      <div ref="aboutProject" class="fade-in"><AboutProject/></div>
      <AboutPossibilities/>
      <div ref="aboutCommand" class="fade-in"><AboutCommand/></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import AboutCommand from '@/features/about-command/about-command.vue';
import AboutPossibilities from '@/features/about-posibilities/about-possibilities.vue';
import AboutProject from '@/features/about-project/about-project.vue';

const aboutProject = ref(null);
const aboutCommand = ref(null);

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
  if (aboutProject.value) observer.observe(aboutProject.value);
  if (aboutCommand.value) observer.observe(aboutCommand.value);
});
</script>

<style scoped>
.fade-in {
  width: 100%;
  opacity: 0;
  transform: translateY(100px);
  transition: opacity 0.3s ease-in-out, transform 0.5s ease-in-out;
}

.fade-in.visible {
  opacity: 1;
  transform: translateY(0);
}

</style>
