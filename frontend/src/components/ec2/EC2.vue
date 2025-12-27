<template>
  <v-container fluid>
    <v-row no-gutters>
      <v-col cols="12" class="actionBlock">
        <div style='float: left;margin-left:5px;display: flex;'>
          <h2 style='margin-right: 5px;'>EC2 Instances</h2> <v-btn v-blur v-on:click="refresh" icon size="small"><v-icon>mdi-refresh</v-icon></v-btn>
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
          <StartInstanceDialog />
          <StopInstanceDialog />
          <RebootInstanceDialog />
          <TerminateInstanceDialog />
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
            :items="instances"
            class="elevation-1"
            item-value="InstanceId"
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
import StartInstanceDialog from './StartInstanceDialog'
import StopInstanceDialog from './StopInstanceDialog'
import RebootInstanceDialog from './RebootInstanceDialog'
import TerminateInstanceDialog from './TerminateInstanceDialog'

export default {
  name: 'EC2',
  components: {
    StartInstanceDialog,
    StopInstanceDialog,
    RebootInstanceDialog,
    TerminateInstanceDialog
  },
  computed: {
    ...mapState({
      headers: state => state.ec2.headers,
      instances: state => state.ec2.instances,
    }),
    search: {
      get() {
        return this.$store.state.ec2.search
      },
      set(value) {
        this.$store.commit('ec2/updateSearch', value)
      }
    },
    selectedItems: {
      get() {
        return this.$store.state.ec2.selectedItems
      },
      set(value) {
        this.$store.commit('ec2/updateSelected', value)
      }
    }
  },
  methods: {
    onClickRow(event, row){
      const item = row && row.item ? row.item : event
      this.$store.dispatch('ec2/onClickRow', item)
    },
    onItemSelected(selected){
      this.$store.dispatch('ec2/onItemSelected', selected)
    },
    refresh(){
      this.$store.dispatch('ec2/initInstances', true)
    }
  },
  created () {
    this.$store.dispatch('ec2/initInstances')
  },
}
</script>

<style>
</style>






