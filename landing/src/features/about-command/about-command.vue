<script setup lang="ts">
import { Carousel, CarouselContent, CarouselItem, CarouselNext, CarouselPrevious } from '@/components/ui/carousel'
import Person from '@/entitites/person.vue';
import { Card } from '@/components/ui/card'
import { IPerson } from '@/entitites/person.interface';
import mike from './assets/mike.jpg'
import mityaiii from './assets/mityaiii.jpg'
import vika from './assets/vika.jpg'
import sasha from './assets/sasha.jpg'
import ivan from './assets/vanya.jpg'

const people: IPerson[] = [{
  name: "Мишель",
  surname: "де Джефрой",
  image: mike,
  telegram: "https://t.me/mikedegeofroy",
  role: "Frontend"
},
{
  name: "Александр",
  surname: "Дьяконов",
  image: sasha,
  telegram: "https://t.me/keine_salz",
  role: "Backend"
},
{
  name: "Вика",
  surname: "Кулешова",
  image: vika,
  telegram: "https://t.me/vktorriaa",
    role: "Design"
},
{
  name: "Дмитрий",
  surname: "Мамченко",
  image: mityaiii,
  telegram: "https://t.me/aliasmityai",
  role: "Frontend"
},
{
  name: "Иван",
  surname: "Тарасов",
  image: ivan,
  telegram: "https://t.me/vaniog",
  role: "Backend"
}
]
</script>

<template>
  <div class="w-full space-y-6">
    <div class="space-y-2 text-center ">
      <h3 class="text-xl text-[#CBCBCB]">Состав</h3>
      <p>
        Наша команда
      </p>
    </div>
    <Carousel class="mx-auto w-[80%]" :opts="{
      align: 'start',
    }">
      <CarouselContent class="-ml-1">
        <CarouselItem v-for="(person, index) in people" :key="index" class="pl-1 pt-3 basis-[75%] flex flex-cols justify-center">
          <label class="absolute bg-white rounded-full mx-auto -translate-y-2 px-2 py-1">{{ person.role }}</label>
          <div class="p-1">
            <Card>
              <Person
                :name="person.name"
                :surname="person.surname"
                :image="person.image"
                :telegram="person.telegram"
                :role="person.role"
              />
            </Card>
          </div>
        </CarouselItem>
      </CarouselContent>
      <CarouselPrevious />
      <CarouselNext />
    </Carousel>
  </div>
</template>