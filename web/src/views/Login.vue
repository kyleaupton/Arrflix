<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Card from 'primevue/card'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')

async function onSubmit() {
  const ok = await auth.loginWithPassword(email.value, password.value)
  if (ok) {
    router.replace({ path: '/' })
  }
}

function onPlex() {
  auth.startPlexSso()
}
</script>

<template>
  <div class="login-wrap">
    <Card class="login-card">
      <template #title>Sign in</template>
      <template #subtitle>Use your Snaggle account or Plex</template>
      <template #content>
        <form class="form" @submit.prevent="onSubmit">
          <label class="label" for="email">Email</label>
          <InputText
            id="email"
            v-model="email"
            type="email"
            autocomplete="email"
            placeholder="you@example.com"
          />

          <label class="label" for="password">Password</label>
          <Password id="password" v-model="password" toggleMask :feedback="false" />

          <div class="actions">
            <Button type="submit" :loading="auth.isLoading" label="Sign in" />
            <Button type="button" severity="secondary" label="Sign in with Plex" @click="onPlex" />
          </div>

          <p v-if="auth.errorMessage" class="error">{{ auth.errorMessage }}</p>
        </form>
      </template>
    </Card>
  </div>
</template>

<style scoped>
.login-wrap {
  display: grid;
  place-items: center;
  min-height: 80vh;
}

.login-card {
  width: 100%;
  max-width: 420px;
}

.form {
  display: grid;
  gap: 0.75rem;
}

.label {
  font-size: 0.875rem;
}

.actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.error {
  color: var(--p-red-500);
  margin-top: 0.5rem;
}
</style>
