<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { Primitive, type PrimitiveProps } from 'radix-vue'
import { type ButtonVariants, buttonVariants } from '.'
import { cn } from '@/lib/utils'

interface Props extends PrimitiveProps {
  variant?: ButtonVariants['variant']
  size?: ButtonVariants['size']
  class?: HTMLAttributes['class']
}

const props = withDefaults(defineProps<Props>(), {
  as: 'button',
})
</script>

<template>
  <div class="button-wrapper">
    <Primitive :as="props.as" :class="cn(buttonVariants({ variant, size }), props.class, 'button')">
      <slot />
    </Primitive>
  </div>
</template>

<style>
.button-wrapper {
  position: relative;
  display: inline-block;
}

.button {
  position: relative;
  background-color: #2EA5FF;
  color: #fff;
  padding: 12px 24px;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  font-weight: bold;
  z-index: 1;
  animation: press-animation 1.5s infinite ease-in-out;
  transition: background-color 0.2s;
}

.button-wrapper::before {
  content: '';
  position: absolute;
  top: 10px;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #208ada;
  border-radius: 10px;
  z-index: 0;
  box-shadow: 0px 4px 12px rgba(0, 0, 0, 0.3);
}

.button:hover {
  background-color: #2EA5FF;
}

@keyframes press-animation {
  0% {
    transform: translateY(-5px);
  }
  50% {
    transform: translateY(-2px);
  }
  100% {
    transform: translateY(-5px);
  }
}
</style>
