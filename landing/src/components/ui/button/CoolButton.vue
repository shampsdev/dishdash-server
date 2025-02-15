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
    <Primitive :as="as" :as-child="asChild" :class="cn(buttonVariants({ variant, size }), props.class, 'button')">
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
  transform: translateY(-4px);
  transition: transform 0.2s;
}

.button-wrapper::before {
  content: '';
  position: absolute;
  top: 4px;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #208ada;
  border-radius: 10px;
  z-index: 0;
  box-shadow: 0px 4px 12px rgba(0, 0, 0, 0.3);
}

/* Убираем изменение цвета на hover */
.button:hover {
    background-color: #2EA5FF;
}

.button:active {
  transform: translateY(0);
}

</style>
