<template>
  <div class="w-full space-y-6 flex items-center flex-col">
    <div class="space-y-2 text-center ">
      <h3 class="text-xl text-[#CBCBCB]">О проекте</h3>
    </div>

    <img src="./assets/about-trans-1.png" alt="about1" class="w-full fade-in rounded-lg" />
    <img src="./assets/about-trans-2.png" alt="about2" class="w-full fade-in rounded-lg" />


    <a href="https://t.me/dishdash_bot?start=landing">
      <CoolButton class="rounded-3xl" size="lg" variant="secondary">
        Попробуй!
        <img src="./assets/telegram-svgrepo-com.svg" color="black" class="ms-2" width="16px" />
      </CoolButton>
    </a>
  </div>
</template>

<script setup lang="ts">
import CoolButton from '@/components/ui/button/CoolButton.vue';

</script>

<style scoped></style>
