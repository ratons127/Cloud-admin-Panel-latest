<template>
  <v-container fluid>
    <v-row no-gutters>
      <v-col cols="12" class="actionBlock">
        <div style='float: left;margin-left:5px;display: flex;'>
          <h2 style='margin-right: 5px;'>AWS Accounts</h2>
        </div>
        <div class="actions">
          <DeleteCredentialsDialog />
          <EditCredentialsDialog />
          <CreateCredentialsDialog />
        </div>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col cols="12">
        <template>
          <v-data-table
            @click:row="onClickRow"
            @update:selected="onItemSelected"
            v-model:selected="selectedItems"
            :headers="accountsHeaders"
            :items-per-page="14"
            :items="accounts"
            class="elevation-1"
            item-value="id"
            select-strategy="single"
            show-select
          >
          </v-data-table>
        </template>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapState } from 'vuex'

import CreateCredentialsDialog from './CreateCredentialsDialog'
import DeleteCredentialsDialog from './DeleteCredentialsDialog'
import EditCredentialsDialog from './EditCredentialsDialog'

export default {
  name: 'Credentials',
  components: {
    CreateCredentialsDialog,
    DeleteCredentialsDialog,
    EditCredentialsDialog
  },
  computed: {
    ...mapState({
      accounts: state => state.credentials.accounts,
      accountsHeaders: state => state.credentials.accountsHeaders,
    }),
    selectedItems: {
      get() {
        return this.$store.state.credentials.selectedItems
      },
      set(value) {
        this.$store.commit('credentials/updateSelected', value)
      }
    }
  },
  methods: {
    onClickRow(event, row){
      const item = row && row.item ? row.item : event
      this.$store.dispatch('credentials/onClickRow', item)
    },
    onItemSelected(selected){
      this.$store.dispatch('credentials/onItemSelected', selected)
    }
  },
  created() {
    this.$store.dispatch('credentials/fetchAccounts')
  }
}
</script>

<style>
</style>


