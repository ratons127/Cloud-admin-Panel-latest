<template>
  <v-dialog v-model="deleteDialog" max-width="1000" v-if="showActions">
    <template v-slot:activator="{ props }">
      <v-btn class="action-btn"  size="small" color="error" v-bind="props">Delete</v-btn>
    </template>
    <v-card>
      <v-card-title class="headline">Delete AWS account {{selectedAccount.name}}</v-card-title>
      <v-card-text>Delete this AWS account configuration?</v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="green-darken-1"  variant="text" @click="deleteDialog = false">Cancel</v-btn>
        <v-btn color="green-darken-1"  variant="text" @click="deleteAccount">Delete</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'DeleteCredentialsDialog',
  computed: {
    ...mapState({
      selectedAccount: state => state.credentials.selectedAccount,
      showActions: state => state.credentials.showActions,
    }),
    deleteDialog: {
      get() {
        return this.$store.state.credentials.deleteDialog
      },
      set(value) {
        this.$store.commit('credentials/updateDeleteDialog', value)
      }
    }
  },
  methods: {
    deleteAccount(){
      this.$store.dispatch('credentials/deleteAccount')
    }
  }
}
</script>
