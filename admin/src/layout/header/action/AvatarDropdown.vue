<script setup lang="ts">
import { ButtonAnimation } from '@/components'
import UserAvatar from '@/components/UserAvatar.vue'
import UserDropdown from '@/components/UserDropdown.vue'
</script>
<template>
  <div class="ml-1">
    <UserDropdown placement="bottom-end">
      <ButtonAnimation animation="beat">
        <UserAvatar src="src path" />
      </ButtonAnimation>
    </UserDropdown>
  </div>
</template>
