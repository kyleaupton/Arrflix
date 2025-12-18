<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Divider from 'primevue/divider'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const email = ref('')
const password = ref('')

async function onSubmit() {
  const ok = await auth.loginWithPassword(email.value, password.value)
  if (ok) {
    const redirectParam = route.query.redirect
    const redirectValue = Array.isArray(redirectParam) ? redirectParam[0] : redirectParam
    const target =
      typeof redirectValue === 'string' && redirectValue.startsWith('/') ? redirectValue : '/'
    router.replace(target)
  }
}

function onPlex() {
  auth.startPlexSso()
}
</script>

<template>
  <div class="login-wrap">
    <Card class="login-card">
      <template #title>
        <h1 class="text-3xl font-semibold mb-2">Sign in</h1>
      </template>
      <template #subtitle>
        <p class="mb-4">
          Don't have an account? <a href="/register" class="text-primary underline">Register</a>
        </p>
      </template>
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
            <Button class="w-full" type="submit" :loading="auth.isLoading" label="Sign in" />
            <!-- <Button type="button" severity="secondary" label="Sign in with Plex" @click="onPlex" /> -->
          </div>

          <p v-if="auth.errorMessage" class="error">{{ auth.errorMessage }}</p>
        </form>

        <Divider align="center" :dt="{ horizontal: { margin: '1.25em 0' } }">
          <b class="text-muted-color">or</b>
        </Divider>

        <Button class="w-full" severity="secondary" label="Sign in with Plex" @click="onPlex" />
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
  width: min(80vw, 420px);
}

.form {
  display: grid;
  gap: 0.75rem;
}

.label {
  font-size: 0.875rem;
}

.actions {
  margin-top: 0.75rem;
}

.error {
  color: var(--p-red-500);
  margin-top: 0.5rem;
}
</style>
