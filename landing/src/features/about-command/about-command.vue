<script setup lang="ts">
import { Carousel, CarouselContent, CarouselItem, CarouselNext, CarouselPrevious } from '@/components/ui/carousel'
import Person from '@/entitites/person.vue';
import { Card } from '@/components/ui/card'
import { IPerson } from '@/entitites/person.interface';
import danya from './assets/danya.png'
import mike from './assets/mike.png'
import nastya from './assets/nastya.png'
import sasha from './assets/sasha.png'
import timur from './assets/timur.png'
import ivan from './assets/vanya.png'
import vika from './assets/vika.png'

const people: IPerson[] = [{
  name: "Мишель",
  surname: "де Джефрой",
  image: mike,
  telegram: "https://t.me/mikedegeofroy",
  role: "CEO/Frontend"
},
{
  name: "Александр",
  surname: "Дьяконов",
  image: sasha,
  telegram: "https://t.me/keine_salz",
  role: "CEO/DevOps"
},
{
  name: "Анастасия",
  surname: "Богданова",
  image: nastya,
  telegram: "https://t.me/learn_how_to_use_a_cart",
  role: "PM"
},
{
  name: "Иван",
  surname: "Тарасов",
  image: ivan,
  telegram: "https://t.me/vaniog",
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
  name: "Данил",
  surname: "Зиндуль",
  image: danya,
  telegram: "https://t.me/dzindul",
  role: "PR"
},
{
  name: "Тимур",
  surname: "Валеев",
  image: timur,
  telegram: "https://t.me/Priest_Faria",
  role: "Frontend"
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