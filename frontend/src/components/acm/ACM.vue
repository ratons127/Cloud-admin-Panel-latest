<template>
  <v-container fluid>
    <v-row no-gutters>
      <v-col cols="12" class="actionBlock">
        <div style='float: left;margin-left:5px;display: flex;'>
          <h2 style='margin-right: 5px;'>Certificate Manager</h2> <v-btn v-blur v-on:click="refresh" icon size="small"><v-icon>mdi-refresh</v-icon></v-btn>
          <v-text-field
            v-model="search"
            append-inner-icon="mdi-magnify"
            label="Search"
            density="compact"
            hide-details
            style='margin-left:25px; margin-top:0; padding-top:0;min-width:350px;'
          ></v-text-field>
        </div>
        <div class="actions">
          <DeleteDialog />
          <ImportDialog /> 
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
            :headers="headers"
            :items-per-page="14"
            :items="certificates"
            class="elevation-1"
            item-value="CertificateArn"
            select-strategy="single"
            show-select
            :search="search"
          ></v-data-table>
        </template>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapState } from 'vuex'

import DeleteDialog from './DeleteDialog'
import ImportDialog from './ImportDialog'

export default {
  name: 'ACM',
  components: {
    DeleteDialog,
    ImportDialog
  },
  computed: {
    ...mapState({
      headers: state => state.acm.headers,
      certificates: state => state.acm.certificates
    }),
    search: {
      get() {
        return this.$store.state.acm.search
      },
      set(value) {
        this.$store.commit('acm/udpateSearch', value)
      }
    },
    selectedItems: {
      get() {
        return this.$store.state.acm.selectedItems
      },
      set(value) {
        this.$store.commit('acm/updateSelected', value)
      }
    }
  },
  methods: {
    onClickRow(event, row){
      const item = row && row.item ? row.item : event
      this.$store.dispatch('acm/onClickRow', item)
    },
    onItemSelected(selected){
      this.$store.dispatch('acm/onItemSelected', selected)
    },
    refresh(){
      this.$store.dispatch('acm/refresh')
    }
  },
  created(){
    this.$store.dispatch('acm/initACM')
  }
}
</script>

<style>
</style>





