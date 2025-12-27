<template>
  <v-dialog v-model="dialog" max-width="600">
    <template v-slot:activator="{ props }">
      <v-btn icon style='margin-right:15px;' v-bind="props">
        <v-icon>mdi-account-plus</v-icon>
      </v-btn>
    </template>
    <v-card>
      <v-card-title class="headline">Invite user</v-card-title>
      <v-card-text>
        <v-text-field v-model="email" label="Email" type="email" variant="outlined"></v-text-field>
        <v-select v-model="role" :items="roles" label="Role" variant="outlined"></v-select>
        <v-alert v-if="message" type="success" density="compact">{{ message }}</v-alert>
        <v-alert v-if="error" type="error" density="compact">{{ error }}</v-alert>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="green-darken-1" variant="text" @click="dialog = false">Close</v-btn>
        <v-btn color="green-darken-1" variant="text" @click="sendInvite">Send</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: 'InviteDialog',
  data: () => ({
    dialog: false,
    email: '',
    role: 'member',
    roles: ['member', 'owner'],
    message: '',
    error: ''
  }),
  computed: {
    tenantId() {
      return this.$store.state.auth.selectedTenantId
    }
  },
  methods: {
    sendInvite() {
      this.message = ''
      this.error = ''
      if (!this.email || !this.tenantId) {
        this.error = 'Email and tenant are required'
        return
      }
      this.$store.dispatch('auth/inviteUser', {
        tenantId: this.tenantId,
        email: this.email,
        role: this.role
      }).then(() => {
        this.message = 'Invite sent'
        this.email = ''
      }).catch(() => {
        this.error = 'Invite failed'
      })
    }
  }
}
</script>
