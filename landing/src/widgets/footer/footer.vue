<template>
  <footer class="max-w-[600px] w-[90%] mx-auto bg-black py-2 md:py-10 rounded-lg shadow-lg">
    <div class="flex text-center text-white flex-col items-center">
      <img class="w-48 md:w-72 mb-10" src="./assets/dishdash-footer.png" alt="Логотип DishDash" />

      <h2 class="text-xl font-semibold mb-2">Контакты:</h2>

      <div class="space-y-2">
        <p class="text-gray-300">
          Напишите в бота: 
          <a href="https://t.me/dishdash_bot?start=landing" class="text-blue-400">@dishdash_bot</a>
        </p>

        <p class="text-gray-300">
          Присоединяйтесь к нашему каналу: 
          <a href="https://t.me/shampsdev" class="text-blue-400">@shampsdev</a>
        </p>

        <p class="text-gray-300">
          Узнайте о нас больше: 
          <a href="https://shamps.dev" class="text-blue-400">Наш сайт</a>
        </p>

        <p class="text-gray-300">
          Напишите нам на почту: 
          <a href="mailto:team@shamps.dev" class="text-blue-400">team@shamps.dev</a>
        </p>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
defineProps()
</script>