<script setup lang="ts">
import { ref, type HTMLAttributes } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Popcorn } from 'lucide-vue-next'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
  FieldSeparator,
} from '@/components/ui/field'
import { Input } from '@/components/ui/input'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')

async function handleSubmit(e: Event) {
  e.preventDefault()
  const success = await authStore.loginWithPassword(email.value, password.value)
  if (success) {
    const redirectParam = route.query.redirect
    const redirectValue = Array.isArray(redirectParam) ? redirectParam[0] : redirectParam
    const target =
      typeof redirectValue === 'string' && redirectValue.startsWith('/') ? redirectValue : '/'
    router.replace(target)
  }
}
</script>

<template>
  <div class="flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
    <div class="flex w-full max-w-sm flex-col gap-6">
      <a href="#" class="flex items-center gap-2 self-center font-medium">
        <div
          class="bg-primary text-primary-foreground flex size-10 items-center justify-center rounded-md"
        >
          <Popcorn class="size-8" />
        </div>
        <div class="text-2xl font-semibold">Arrflix</div>
      </a>

      <div :class="cn('flex flex-col gap-6', props.class)">
        <Card>
          <CardHeader class="text-center">
            <CardTitle class="text-xl"> Welcome back </CardTitle>
            <CardDescription> Login with your Plex account? </CardDescription>
          </CardHeader>
          <CardContent>
            <form @submit="handleSubmit">
              <FieldGroup>
                <Field>
                  <Button variant="outline" type="button" @click="authStore.startPlexSso()">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                      <path
                        d="M12.152 6.896c-.948 0-2.415-1.078-3.96-1.04-2.04.027-3.91 1.183-4.961 3.014-2.117 3.675-.546 9.103 1.519 12.09 1.013 1.454 2.208 3.09 3.792 3.039 1.52-.065 2.09-.987 3.935-.987 1.831 0 2.35.987 3.96.948 1.637-.026 2.676-1.48 3.676-2.948 1.156-1.688 1.636-3.325 1.662-3.415-.039-.013-3.182-1.221-3.22-4.857-.026-3.04 2.48-4.494 2.597-4.559-1.429-2.09-3.623-2.324-4.39-2.376-2-.156-3.675 1.09-4.61 1.09zM15.53 3.83c.843-1.012 1.4-2.427 1.245-3.83-1.207.052-2.662.805-3.532 1.818-.78.896-1.454 2.338-1.273 3.714 1.338.104 2.715-.688 3.559-1.701"
                        fill="currentColor"
                      />
                    </svg>
                    Login with Plex
                  </Button>
                </Field>
                <FieldSeparator class="*:data-[slot=field-separator-content]:bg-card">
                  Or continue with
                </FieldSeparator>
                <Field v-if="authStore.errorMessage">
                  <FieldDescription class="text-destructive">
                    {{ authStore.errorMessage }}
                  </FieldDescription>
                </Field>
                <Field>
                  <FieldLabel for="email"> Email </FieldLabel>
                  <Input
                    id="email"
                    v-model="email"
                    type="email"
                    placeholder="email@example.com"
                    required
                    :disabled="authStore.isLoading"
                  />
                </Field>
                <Field>
                  <div class="flex items-center">
                    <FieldLabel for="password"> Password </FieldLabel>
                    <a href="#" class="ml-auto text-sm underline-offset-4 hover:underline">
                      Forgot your password?
                    </a>
                  </div>
                  <Input
                    id="password"
                    v-model="password"
                    type="password"
                    required
                    :disabled="authStore.isLoading"
                  />
                </Field>
                <Field>
                  <Button type="submit" :disabled="authStore.isLoading">
                    {{ authStore.isLoading ? 'Logging in...' : 'Login' }}
                  </Button>
                  <FieldDescription class="text-center">
                    Don't have an account?
                    <a href="#"> Sign up </a>
                  </FieldDescription>
                </Field>
              </FieldGroup>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>
