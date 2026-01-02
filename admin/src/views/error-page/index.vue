<script setup lang="ts">
import { NNumberAnimation, NButton } from 'naive-ui'
import { ref } from 'vue'

import router from '@/router'

import type { StyleValue } from 'vue'

interface ErrorPageProps {
  code?: number
  content?: string
  containerClass?: string
  containerStyle?: StyleValue
}

defineOptions({
  name: 'ErrorPage',
})

const {
  code = 404,
  content = 'Not Found',
  containerClass,
  containerStyle,
} = defineProps<ErrorPageProps>()

const prevCode = ref(0)
</script>

<template>
  <div class="grid h-dvh w-full flex-col place-items-center">
    <div
      class="flex w-full flex-col"
      :class="containerClass"
      style="max-width: 600px"
      :style="containerStyle"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
        xml:space="preserve"
        class="text-primary"
        viewBox="0 0 500 500"
        width="100%"
      >
        <path
          fill="currentColor"
          d="M51 268.5s2.4 72.2 62.6 120.3c60.3 48 143.7 51.5 211.3 56 67.6 4.7 119.5-33.6 127.3-91.4 7.8-57.8-43.8-74.8-71.8-140.6-28-65.9-25.5-87.5-84.8-128S143.6 69.2 95 138.6c-48.6 69.3-44 130-44 130z"
        />
        <path
          fill="#fff"
          d="M51 268.5s2.4 72.2 62.6 120.3c60.3 48 143.7 51.5 211.3 56 67.6 4.7 119.5-33.6 127.3-91.4 7.8-57.8-43.8-74.8-71.8-140.6-28-65.9-25.5-87.5-84.8-128S143.6 69.2 95 138.6c-48.6 69.3-44 130-44 130z"
          class="opacity-50 dark:opacity-10"
        />
        <foreignObject
          x="64"
          y="180"
          width="180"
          height="160"
        >
          <div
            xmlns="http://www.w3.org/1999/xhtml"
            style="font-size: 80px"
            class="font-bold text-neutral-700 dark:text-neutral-800"
          >
            <NNumberAnimation
              :from="prevCode"
              :to="code"
              :duration="800"
              @finish="prevCode = code"
            />
          </div>
        </foreignObject>
        <foreignObject
          x="70"
          y="288"
          width="140"
          height="66"
        >
          <div
            xmlns="http://www.w3.org/1999/xhtml"
            class="text-center leading-4 font-bold text-neutral-700 dark:text-neutral-800"
          >
            {{ content }}
          </div>
        </foreignObject>
        <path
          fill="none"
          stroke="#353535"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M438.6 283.2H42.9"
        />
        <path
          fill="currentColor"
          d="M415.7 82.1s0-4.9-7-11.8c-6.9-7-11.3-10.9-11.3-10.9s-3.4 11.4 3 18.3S415.7 91 415.7 91v-8.9zM405.6 113s0-5-6.9-11.9-11.3-10.8-11.3-10.8-3.5 11.3 3 18.2c6.4 7 15.2 13.3 15.2 13.3V113z"
        />
        <path
          fill="currentColor"
          d="M415.2 90.5s.5-4.9 8.1-11 12.5-9.6 12.5-9.6 2.2 11.6-4.9 17.8c-7.1 6.2-16.6 11.7-16.6 11.7l1-8.9zM403.1 127s1-4.9 9.2-10.2c8.2-5.3 13.4-8.2 13.4-8.2s1 11.8-6.7 17.2c-7.7 5.4-17.7 9.8-17.7 9.8l1.8-8.7zm2.3 39.6s0-5 6.9-11.8c6.9-7 11.3-10.9 11.3-10.9s3.5 11.4-3 18.3c-6.3 6.9-15.2 13.3-15.2 13.3v-8.9z"
        />
        <path
          fill="currentColor"
          d="M405.3 176.3s.7-4.9-5.2-12.7c-5.8-7.8-9.7-12.4-9.7-12.4s-5 10.8.4 18.5a115 115 0 0 0 13.2 15.4l1.3-8.8zm-6.8 28.7s0-5 6.9-11.9 11.3-10.8 11.3-10.8 3.5 11.3-3 18.2c-6.4 7-15.2 13.4-15.2 13.4v-9zm4.5 26s-1.2-4.8 3.8-13.2c4.9-8.4 8.2-13.4 8.2-13.4s6.2 10.1 1.7 18.4a119 119 0 0 1-11.4 16.8l-2.2-8.6z"
        />
        <path
          fill="none"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M416.3 80.9a118 118 0 0 1-2 15.3c-2.6 14.1-9.5 23.2-12.2 36-2.7 12.7 5.1 29 3 45.5-2 16.3-8.5 25.4-6.6 36.3 1.8 11 11.3 40.5 11.3 40.5m5.4-160.7 7.4-5.7m-7.1-.9-3.1-5.5m-9.9 49.9 12.7-9.4"
        />
        <path
          fill="none"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M405.6 121.8s-5.2-10.2-11-17m11 68.8s4.9-9.5 13.4-17.8m-14.6 24.7s-3-7.4-6.7-12.8l-4.2-6m5 52.3s2.5-6.3 5.6-11.3c3-5 7.2-9.5 7.2-9.5m-6 46.4s-.7-6 1.9-11c2.6-5 5.5-11.5 6.7-14.1"
        />
        <path
          fill="currentColor"
          d="M370.4 148s-.5-3 3-8.1l6-8s3.4 6.7.1 11.7c-3.2 5-8.1 10-8.1 10l-1-5.6zm9.6 18.3s-.5-3.1 3.1-8.2c3.6-5 6-8 6-8s3.3 6.7 0 11.7c-3.2 5-8 10-8 10l-1-5.5z"
        />
        <path
          fill="currentColor"
          d="M371.6 153.3s-.8-3-6.3-6c-5.4-3.1-8.8-4.7-8.8-4.7s-.1 7.5 5 10.6a77.3 77.3 0 0 0 11.7 5.5l-1.6-5.4zm11.6 21.5s-1.2-3-7-5.4l-9.2-3.7s.6 7.5 6 10a77.3 77.3 0 0 0 12.3 4.3l-2.1-5.2zm2.8 25.1s-.5-3.1-5.6-6.7c-5-3.6-8.3-5.6-8.3-5.6s-1 7.5 3.9 11.2c4.7 3.6 11 6.7 11 6.7l-1-5.6z"
        />
        <path
          fill="currentColor"
          d="M387.2 206s-1-3 1.8-8.6c2.8-5.5 4.7-8.8 4.7-8.8s4.4 6.2 1.8 11.7c-2.5 5.4-6.6 11-6.6 11l-1.7-5.3zm8.6 25s.4-3.2 5.3-7c4.8-4 8-6.1 8-6.1s1.3 7.4-3.3 11.3c-4.5 4-10.6 7.3-10.6 7.3l.6-5.6zm-1.3-7.8s-.5-3-5.6-6.7c-5-3.5-8.3-5.5-8.3-5.5s-.9 7.4 3.9 11c4.7 3.7 11 6.8 11 6.8l-1-5.6zm0 16.8s.2-3.1-3.8-7.9c-4-4.7-6.6-7.4-6.6-7.4s-2.8 7 1 11.7c3.6 4.7 8.9 9.3 8.9 9.3l.5-5.7z"
        />
        <path
          fill="none"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M369.9 147.3c.8 3.3 1.8 6.6 2.8 9.4 3.3 8.6 8.6 13.5 11.7 21.2 3 7.7 0 18.9 3 28.9 3.1 10 8 15 8.1 22 0 7.1-2.7 26.7-2.7 26.7M372 155.3l-5.3-2.8m4.4-1.3 1.3-3.8m11.6 30.2-9-4.5"
        />
        <path
          fill="none"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M381 171.8s2.2-6.9 5.1-11.8m.5 44.3s-4-5.4-10.2-9.7m11.8 13.9s1.1-5 2.8-8.7l2-4.3m2.3 38.7s3.5-4.3 5.7-5.9c2.2-1.6 3.7-3.2 3.7-3.2m-9.2 3.8s-2.2-3.7-4.7-6.5-5.6-5.2-5.6-5.2m8.8 28.5s-.2-3.8-2.4-6.7l-5.7-8.2"
        />
        <path
          fill="#353535"
          d="M416.5 289.2H386a2.9 2.9 0 0 1-2.9-2.9l-4.5-34.2h45.2l-4.6 34.2c0 1.6-1.2 2.9-2.8 2.9z"
        />
        <path
          fill="#353535"
          d="M423.5 258h-44.8c-2 0-3.6-1.6-3.6-3.5 0-2 1.6-3.6 3.6-3.6h44.8c2 0 3.5 1.6 3.5 3.6s-1.6 3.5-3.5 3.5zm-43.2 71.3c-.7 0-1.4-.2-2-.5a4.7 4.7 0 0 1-2.2-6.3 8 8 0 0 0-3.2-10.5c-7-4.5-23.5-5.2-29.2-5.1h-.1a4.7 4.7 0 0 1 0-9.5c2.3 0 23.7-.2 34.4 6.7 10.6 6.7 9 17.7 6.6 22.5a4.7 4.7 0 0 1-4.3 2.7z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m402.1 354.2-1.1 3.5-15.7-25.9-.4-5.6zm-1.1 3.5 5-2.4.8-3.4-4.7 2.3z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m382.8 323.3 18.2-.6 5.8 29.2-4.7 2.3zm-2.5 1.3 18.3 30.6-6.1 2.9-26.9-21.9zm12.2 33.5-.8 3.1-27.1-21.5 1-3.5z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m392.5 358.1 6.1-2.9-.9 3.1-6 2.9z"
        />
        <path
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M374.7 329a9.2 9.2 0 1 0 18.4 0 9.2 9.2 0 0 0-18.4 0z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m189.2 354.2 1 3.5 15.8-25.9.4-5.6zm1 3.5-4.9-2.4-.9-3.4 4.8 2.3z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m208.4 323.3-18.2-.6-5.8 29.2 4.8 2.3zm2.5 1.3-18.2 30.6 6 2.9 26.9-21.9z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m198.7 358.1.8 3.1 27.2-21.5-1.1-3.5z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m198.7 358.1-6-2.9.8 3.1 6 2.9z"
        />
        <path
          fill="#353535"
          d="M208 332.4a4.7 4.7 0 0 1-3-8.4c37.5-30.5 29.8-76 29.7-76.5a4.7 4.7 0 0 1 9.3-1.6 94 94 0 0 1-33 85.5c-1 .7-2 1-3 1z"
        />
        <circle
          cx="244.5"
          cy="234.5"
          r="18.5"
          fill="#adadad"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
        />
        <path
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M216.6 329a9.2 9.2 0 1 1-18.5 0 9.2 9.2 0 0 1 18.5 0z"
        />
        <ellipse
          cx="294"
          cy="283.3"
          fill="#adadad"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          rx="29.3"
          ry="27.4"
        />
        <path
          fill="#353535"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m278.7 271.6 2.8 10.7c1.6 5.4 7 8.9 12.6 8.2h0c5.7-.7 10-5.4 10.3-11l.5-10.9-26.2 3z"
        />
        <path
          fill="#353535"
          d="M336.5 217.8c-4.7 0-7.9-1-8.1-1.1l.7-2.3c.1 0 13.3 4.3 26.6-5.5l1.4 2a34.9 34.9 0 0 1-20.6 6.9zm14.8 17.7C340.5 219.7 328 221 328 221l-.2-2.3c.5 0 13.8-1.4 25.4 15.6l-2 1.3zm-2.8 11.2c-14-4.3-22.4-23-22.7-23.7l2.1-1c.1.2 8.4 18.5 21.3 22.4l-.7 2.3z"
        />
        <path
          fill="none"
          stroke="#353535"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m355.8 210.5 3.8-2.5m-8.4 24.8 2.5 5m-6.9 6.9s3.1.9 6 .9"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m240.6 214.3 3.5 24.4 86.9-8.8-.3-25.4z"
        />
        <path
          fill="#e5e5e5"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m244.1 238.7 15.6 39.3 66.8-6.6 4.5-41.5z"
        />
        <defs>
          <path
            id="a"
            d="m265 249 6.2 18a6 6 0 0 0 6.3 4l26.8-2.8a6 6 0 0 0 5.2-5.4l1.7-19c.2-2.6-2-4.8-4.7-4.5l-37.9 4.2c-2.7.3-4.4 3-3.5 5.5z"
          />
        </defs>
        <use
          xlink:href="#a"
          fill="#353535"
          overflow="visible"
        />
        <clipPath id="b">
          <use
            xlink:href="#a"
            overflow="visible"
          />
        </clipPath>
        <path
          fill="none"
          stroke="#fff"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m265.1 258.9 11-1 3.5-9.5 8.3 14.8 4.2-11 2.1 3.5 20-2.4"
          clip-path="url(#b)"
        />
        <use
          xlink:href="#a"
          fill="none"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          overflow="visible"
        />
        <circle
          cx="320.9"
          cy="241.3"
          r="4.9"
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
        />
        <path
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M324.8 256.5a4.6 4.6 0 1 1-9 0 4.6 4.6 0 0 1 9 0zm-63.1-9.9a4.9 4.9 0 1 1-9.3 2.8 4.9 4.9 0 0 1 9.3-2.8z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M265.2 261.4a4.6 4.6 0 1 1-8.7 2.6 4.6 4.6 0 0 1 8.7-2.6z"
        />
        <ellipse
          cx="276.4"
          cy="302.7"
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          rx="12.9"
          ry="12.5"
        />
        <path
          fill="#353535"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m257.4 330.7 11.8-26.6c1-2.3 3.7-3.5 6.2-2.8h0a5 5 0 0 1 3.4 6l-7.8 35-13.6-11.6z"
        />
        <path
          fill="#adadad"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m253.4 290 25.9 60.7-31 4.3s-18.6-18.8-31-54.7l36.1-10.3z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m253.4 290 1.2-5.8 25.9 62-1.2 4.5zm1.2-5.8-35.2 10.6-2.2 5.5 36.2-10.3z"
        />
        <path
          fill="none"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m219.7 306.1 35.5-11m-32.8 17.3 35.9-9.7m-33.1 16.4 36.1-9.7m-32.2 17.9 35.3-10.3m-31.6 16.9 35.2-9.3m-31.3 15.2 34.3-8.2m-29.9 14 32.1-7.6"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m244.3 350.2 31.8-5.2"
        />
        <ellipse
          cx="310.8"
          cy="302.7"
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          rx="12.9"
          ry="12.5"
        />
        <path
          fill="#353535"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M329.8 330.7 317.9 304a5.1 5.1 0 0 0-6.1-2.8h0c-2.6.8-4 3.4-3.5 6l7.8 35 13.7-11.6z"
        />
        <path
          fill="#adadad"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M333.7 290 308 350.7l31 4.3s18.5-18.8 31-54.7L333.7 290z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m333.7 290-1.2-5.8-25.8 62 1.2 4.5zm-1.2-5.8 35.3 10.6 2.1 5.5-36.2-10.3z"
        />
        <path
          fill="none"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m367.5 306.1-35.6-11m32.8 17.3-35.8-9.7m33.1 16.4-36.2-9.7m32.2 17.9L322.8 317m31.5 16.9-35.2-9.3m31.3 15.2-34.3-8.2m29.9 14-32-7.6m28.8 12.2L311 345"
        />
        <path
          fill="none"
          stroke="#353535"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m357.3 238 3.3 1.5m-5.4-.2 3.3 5.2m-5.2-5-.6 3.3m8.8-25.1h8.2m-14.5.4 11.7 3.3m-7.3-7.7 7.5-1.1m-34.6 27.1 2.1 6.1m.6-6.5 5.9 9.6m-3.8-15.3 5.7 1.9"
        />
        <path
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width=".5"
          d="m368.6 230.5-1.9 2 1.4 2.5-2.6-1.2-1.9 2 .4-2.7-2.6-1.2 2.8-.5.3-2.8 1.4 2.4zm3-27.7-1.9 2 1.3 2.4-2.5-1.1-1.9 2 .3-2.7-2.5-1.2 2.7-.6.3-2.7 1.4 2.4zm-30.7 48.4-1.9 2.1 1.4 2.4-2.6-1.1-1.9 2 .4-2.8-2.6-1.1 2.8-.6.3-2.7 1.4 2.4z"
        />
        <ellipse
          cx="250.3"
          cy="191.5"
          fill="#adadad"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          rx="14.1"
          ry="13.9"
        />
        <ellipse
          cx="311.8"
          cy="187.7"
          fill="#adadad"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          rx="14.1"
          ry="13.9"
        />
        <ellipse
          cx="283.4"
          cy="217"
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          rx="17.8"
          ry="17"
        />
        <path
          fill="#353535"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m270.4 208.3 2.5 10.8c1.5 5.5 6.7 9 12.4 8.5h0c5.7-.5 10.2-5 10.6-10.8l.8-10.7-26.3 2.2z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m254.3 163-6.7 26.6 66.3-5.4-9.9-25.3z"
        />
        <path
          fill="#e5e5e5"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m247.6 189.6 11.7 29.2 52-5.1 2.6-29.5z"
        />
        <path
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M277.4 201.5c-.2 4.1-3.8 7.4-8.1 7.7-4.4.3-8.6-2.4-9.5-6.5-1-4.7 2.5-9.3 8-9.6 5.5-.4 9.8 3.7 9.6 8.4z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M274 201.7c.1 2.7-2.1 4.8-5 5-2.8.2-5.4-1.6-5.9-4.2-.4-2.8 1.8-5.5 5-5.7 3.3-.2 6 2 6 5z"
        />
        <path
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M306.9 198.6c-.2 4.2-3.8 7.4-8.1 7.8-4.4.3-8.6-2.4-9.5-6.5-1.1-4.7 2.5-9.3 8-9.7 5.5-.3 9.8 3.7 9.6 8.4z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M303.6 198.9c0 2.6-2.3 4.7-5 5-2.9.2-5.5-1.6-6-4.3-.5-2.8 1.8-5.5 5-5.7 3.3-.2 6 2.1 6 5z"
        />
        <path
          fill="#353535"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M277.7 213.5a2.4 2.4 0 0 1-.2-4.7l12.7-1.3c1.3-.1 2.4.9 2.5 2.1.2 1.3-.8 2.5-2 2.6l-12.7 1.3h-.3z"
        />
        <path
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M260.3 171.3a4.8 4.8 0 1 1-9.6 0 4.8 4.8 0 0 1 9.6 0z"
        />
        <path
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M256.6 168.5a4 4 0 1 1-8.2 0 4 4 0 0 1 8.2 0z"
        />
        <path
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M251.7 164.1a3.2 3.2 0 1 1-6.4 0 3.2 3.2 0 0 1 6.4 0z"
        />
        <path
          fill="none"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m247.7 163-7.3-6.8"
        />
        <path
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M297.7 168a4.8 4.8 0 1 0 9.6 0 4.8 4.8 0 0 0-9.6 0z"
        />
        <path
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M301.4 165.2a4 4 0 1 0 8.2 0 4 4 0 0 0-8.2 0z"
        />
        <path
          fill="currentColor"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M306.3 160.8a3.2 3.2 0 1 0 6.4 0 3.2 3.2 0 0 0-6.4 0z"
        />
        <path
          fill="none"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m310.3 159.7 7.3-6.8"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M299.5 335.6c0 2.3-1.5 4.2-3.4 4.2s-3.4-2-3.4-4.2c0-2.3 1.5-4.2 3.4-4.2s3.4 1.9 3.4 4.2z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M306.3 339h-.3l-9.4-1.3a1.8 1.8 0 0 1 .5-3.5l9.5 1.4a1.8 1.8 0 0 1-.3 3.5z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M173.7 355.1c-1.3-2-1-4.3.6-5.3s3.9-.3 5.1 1.7c1.3 2 1 4.3-.6 5.3s-3.9.3-5.1-1.7z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m166 355.8.3-.1 8.8-4a1.8 1.8 0 0 1 1.4 3.3l-8.7 3.9a1.8 1.8 0 0 1-1.7-3.1z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M400.3 310.3c-1.2-2-1-4.3.7-5.3 1.5-1 3.8-.3 5 1.7 1.3 2 1 4.3-.6 5.3s-3.9.3-5-1.7z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m392.7 311 .3-.1 8.7-4a1.8 1.8 0 0 1 1.5 3.3l-8.8 3.9a1.8 1.8 0 0 1-1.7-3.1zm-161.5 46c-4 0-7.2 2-7.2 4.3s3.2 4.3 7.2 4.3 7.2-2 7.2-4.3c0-2.4-3.2-4.3-7.2-4.3zm0 6.5c-2 0-3.8-1-3.8-2.2 0-1.2 1.7-2.2 3.8-2.2 2 0 3.7 1 3.7 2.2 0 1.2-1.7 2.2-3.7 2.2z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m296.3 357.7-5-2-5 2-1.1 4.4 3.4 3.5h5.5l3.4-3.5-1.2-4.4zm-3.6 5.6H290l-1.7-1.8.6-2.2 2.5-1 2.5 1 .7 2.2-1.8 1.8zm97.1-69.5-5-2-4.9 2-1.2 4.4 3.4 3.5h5.5l3.4-3.5-1.2-4.4zm-3.5 5.6h-2.8l-1.7-1.8.6-2.2 2.5-1 2.5 1 .6 2.2-1.7 1.8zm-57 68.3h-.4c-.2-.1-1-.5-1.7-1.4-.5-.9-1-2.4 0-4.7.6-1.8 2-3 3.8-3.5 2.1-.6 4.7 0 6.5 1.4a1.2 1.2 0 1 1-1.5 1.9 5 5 0 0 0-4.4-1c-.8.1-1.8.7-2.3 2-.9 2.4.2 3 .4 3 .6.3.9 1 .7 1.6-.2.4-.7.7-1.1.7z"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="m353.7 372-1.1-.1-19.2-6.2a3.5 3.5 0 0 1 2.2-6.8l19.1 6.3a3.5 3.5 0 0 1-1 6.9z"
        />
        <path
          fill="none"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M332.8 365.5s.7-4.9 6.8-5m-2.7 6.1s.8-5 6.8-5m-3.1 6.4s.7-5 6.8-5.1m-3 6.1s.7-4.9 6.8-5m-3 6.4s.7-5 6.8-5"
        />
        <path
          fill="#fff"
          stroke="#353535"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-miterlimit="10"
          stroke-width="1.1"
          d="M353.6 372h-.4c-.1-.1-1-.5-1.6-1.5-.5-.8-1-2.3-.1-4.7.6-1.7 2-3 3.8-3.5 2.2-.6 4.8 0 6.6 1.5a1.2 1.2 0 1 1-1.5 1.8c-1.3-1-3-1.4-4.5-1-.7.2-1.7.7-2.2 2-.9 2.4.2 3 .4 3.1.6.3.9 1 .6 1.6-.2.4-.6.7-1 .7z"
        />
      </svg>
      <div class="flex justify-center gap-x-4">
        <NButton
          size="large"
          type="primary"
          @click="router.replace('/')"
          >回到首页</NButton
        >
        <slot />
      </div>
    </div>
  </div>
</template>
