<script setup lang="ts">
import { ref } from 'vue'
import { useQuery, useMutation } from '@tanstack/vue-query'
import { Plus } from 'lucide-vue-next'
import {
  getV1UsersOptions,
  deleteV1UsersByIdMutation,
} from '@/client/@tanstack/vue-query.gen'
import type { User } from '@/components/tables/configs/userTableConfig'
import DataTable from '@/components/tables/DataTable.vue'
import {
  userColumns,
  createUserActions,
} from '@/components/tables/configs/userTableConfig'
import { useModal } from '@/composables/useModal'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import UserDialog from '@/components/modals/UserDialog.vue'

// Data queries
const { data: users, isLoading, refetch } = useQuery(getV1UsersOptions())
const modal = useModal()

// Mutations
const deleteUserMutation = useMutation(deleteV1UsersByIdMutation())

// State
const userError = ref<string | null>(null)

// Handlers
const handleAddUser = () => {
  modal.open(UserDialog, {
    props: {
      user: null,
    },
    onClose: () => {
      refetch()
    },
  })
}

const handleEditUser = (user: User) => {
  modal.open(UserDialog, {
    props: {
      user,
    },
    onClose: () => {
      refetch()
    },
  })
}

const handleDeleteUser = async (user: User) => {
  if (!user.id) return
  const confirmed = await modal.confirm({
    title: 'Delete User',
    message: `Are you sure you want to delete "${user.email}"? This action cannot be undone.`,
    severity: 'danger',
  })
  if (!confirmed) return
  try {
    await deleteUserMutation.mutateAsync({ path: { id: user.id } })
    refetch()
  } catch (err) {
    userError.value = err instanceof Error ? err.message : 'Failed to delete user'
  }
}

const userActions = createUserActions(handleEditUser, handleDeleteUser)
</script>

<template>
  <div class="flex flex-col gap-6">
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle class="text-xl font-semibold mb-2">User Management</CardTitle>
            <p class="text-sm text-muted-foreground">
              Manage application users and their roles.
            </p>
          </div>
          <Button @click="handleAddUser">
            <Plus class="mr-2 size-4" />
            Add User
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <div
          v-if="userError"
          class="p-4 bg-destructive/10 border border-destructive/30 rounded-lg text-destructive mb-4"
        >
          {{ userError }}
        </div>
        <div v-if="isLoading" class="space-y-3">
          <Skeleton class="h-12 w-full" />
          <Skeleton class="h-12 w-full" />
          <Skeleton class="h-12 w-full" />
        </div>
        <DataTable
          v-else
          :data="users || []"
          :columns="userColumns"
          :actions="userActions"
          :loading="isLoading"
          empty-message="No users found"
          searchable
          search-placeholder="Search users..."
          paginator
          :rows="10"
        />
      </CardContent>
    </Card>
  </div>
</template>
