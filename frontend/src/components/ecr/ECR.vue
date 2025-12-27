<template>
  <v-container fluid>
    <v-row no-gutters>
      <v-col cols="12" class="actionBlock">
        <div style='float: left;margin-left:5px;display: flex;'>
          <h2 style='margin-right: 5px;'>ECR</h2> <v-btn v-blur v-on:click="refresh" icon size="small"><v-icon>mdi-refresh</v-icon></v-btn>
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
          <DeleteDialog v-if="showActions"/>
          <CreateDialog />
        </div>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col cols="3">
        <template>
          <v-data-table
            @click:row="onClickRow"
            @update:selected="onItemSelected"
            v-model:selected="selectedItems"
            :headers="repositoryHeaders"
            :items-per-page="14"
            :items="repositories"
            class="elevation-1"
            item-value="RepositoryArn"
            select-strategy="single"
            show-select
            :search="search"
          ></v-data-table>
        </template>
      </v-col>
      <v-col cols="9">
        <template>
          <v-tabs fixed-tabs @update:modelValue="selectTab">
            <v-tab>Information</v-tab>
            <v-tab>Images</v-tab>
          </v-tabs>
        </template>
        <div v-if="tabState === 0"><TabInformation /></div>
        <div v-if="tabState === 1"><TabImages /></div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import DeleteDialog from './DeleteDialog'
import CreateDialog from './CreateDialog'
import TabImages from './TabImages'
import TabInformation from './TabInformation'

import { mapState } from 'vuex'

export default {
  name: 'ECR',
  components: {
    DeleteDialog,
    CreateDialog,
    TabImages,
    TabInformation
  },
  computed: {
    ...mapState({
      repositoryHeaders: state => state.ecr.repositoryHeaders,
      repositories: state => state.ecr.repositories,
      tabState: state => state.ecr.tabState,
      showActions: state => state.ecr.showActions,
    }),
    search: {
      get() {
        return this.$store.state.ecr.search
      },
      set(value) {
        this.$store.commit('ecr/updateSearch', value)
      }
    },
    selectedItems: {
      get() {
        return this.$store.state.ecr.selectedItems
      },
      set(value) {
        this.$store.commit('ecr/updateSelected', value)
      }
    }
  },
  methods: {
    onClickRow(event, row){
      const item = row && row.item ? row.item : event
      this.$store.dispatch('ecr/onClickRow', item)
    },
    onItemSelected(selected){
      this.$store.dispatch('ecr/onItemSelected', selected)
    },
    selectTab(tabState) {
      this.$store.commit('ecr/updateTabState', tabState)
    },
    refresh(){
      this.$store.dispatch('ecr/refresh')
    }
  },
  created() {
    this.$store.dispatch('ecr/getRepositories')
  }
}
</script>

<style>
.tabTable {
  margin-left: 25px;
  margin-top: 25px;
}
</style>







