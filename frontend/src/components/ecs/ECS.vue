<template>
  <v-container fluid>
    <v-row no-gutters>
      <v-col cols="12" class="actionBlock">
        <div style='float: left;margin-left:5px;display: flex;'>
          <h2 style='margin-right: 5px;'>ECS</h2> <v-btn v-blur v-on:click="refresh" icon size="small"><v-icon>mdi-refresh</v-icon></v-btn>
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
          <ForceUpdateDialog />
        </div>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col cols="3">
        <template>
          <v-data-table
            @click:row="onClusterClickRow"
            @update:selected="onClusterItemSelected"
            v-model:selected="selectedClusterArr"
            :headers="clusterHeaders"
            :items-per-page="14"
            :items="clusters"
            class="elevation-1"
            item-value="ClusterName"
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
            <v-tab>Services</v-tab>
            <v-tab>Events</v-tab>
            <v-tab>Tasks</v-tab>
            <v-tab>Containers</v-tab>
          </v-tabs>
        </template>
        <!-- Tab Information-->
        <template v-if="tabState === 0">
          <TabInformation />
        </template>
        <!-- Tab Services-->
        <template v-if="(tabState === 1) && (services.length > 0)">
          <v-data-table
            @click:row="onServiceClickRow"
            @update:selected="onServiceItemSelected"
            v-model:selected="selectedServiceArr"
            :headers="serviceHeaders"
            :items-per-page="14"
            :items="services"
            class="tabTable"
            item-value="ServiceArn"
            select-strategy="single"
            show-select
          ></v-data-table>
        </template>
        <!-- Events-->
        <template v-if="(tabState === 2) && selectedService.Events && (selectedService.Events.length > 0)">
          <v-data-table
            :headers="eventHeaders"
            :items-per-page="14"
            :items="selectedService.Events"
            class="tabTable"
            item-value="Id"
          ></v-data-table>
        </template>
        <template v-if="(tabState === 2) && !selectedService.Events">
          <div class="tabTable">
          You must select a service first to see the service events
          </div>
        </template>
        <!-- Tasks-->
        <template v-if="(tabState === 3) && (tasks.length > 0)">
          <v-data-table
            @click:row="onTaskClickRow"
            @update:selected="onTaskItemSelected"
            v-model:selected="selectedTaskArr"
            :headers="taskHeaders"
            :items-per-page="14"
            :items="tasks"
            class="tabTable"
            item-value="TaskArn"
            select-strategy="single"
            show-select
          ></v-data-table>
        </template>
        <!-- Containers-->
        <template v-if="(tabState === 4) && selectedTask.Containers && (selectedTask.Containers.length > 0)">
          <v-data-table
            :headers="containerHeaders"
            :items-per-page="14"
            :items="selectedTask.Containers"
            class="tabTable"
            item-value="TaskArn"
          ></v-data-table>
        </template>
        <template v-if="(tabState === 4) && !selectedTask.Containers">
          <div class="tabTable">
          You must select a task first to see the containers
          </div>
        </template>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapState } from 'vuex'

import ForceUpdateDialog from './ForceUpdateDialog'
import TabInformation from './TabInformation'

export default {
  name: 'ECS',
  components: {
    ForceUpdateDialog,
    TabInformation
  },
  computed: {
    ...mapState({
      clusterHeaders: state => state.ecs.clusterHeaders,
      clusters: state => state.ecs.clusters,
      serviceHeaders: state => state.ecs.serviceHeaders,
      services: state => state.ecs.services,
      selectedService: state => state.ecs.selectedService,
      taskHeaders: state => state.ecs.taskHeaders,
      tasks: state => state.ecs.tasks,
      eventHeaders: state => state.ecs.eventHeaders,
      containerHeaders: state => state.ecs.containerHeaders,
      selectedTask: state => state.ecs.selectedTask,
      tabState: state => state.ecs.tabState,
    }),
    selectedClusterArr: {
      get() {
        return this.$store.state.ecs.selectedClusterArr
      },
      set(value) {
        this.$store.commit('ecs/updateSelectedClusterArr', value)
      }
    },
    selectedServiceArr: {
      get() {
        return this.$store.state.ecs.selectedServiceArr
      },
      set(value) {
        this.$store.commit('ecs/updateSelectedServiceArr', value)
      }
    },
    selectedTaskArr: {
      get() {
        return this.$store.state.ecs.selectedTaskArr
      },
      set(value) {
        this.$store.commit('ecs/updateSelectedTaskArr', value)
      }
    },
    search: {
      get() {
        return this.$store.state.ecs.search
      },
      set(value) {
        this.$store.commit('ecs/updateSearch', value)
      }
    },
  },
  methods: {
    selectTab(tabState){
      this.$store.commit('ecs/updateTabState', tabState)
    },
    refresh(){
      this.$store.dispatch('ecs/refresh')
    },
    onClusterClickRow(event, row){
      const item = row && row.item ? row.item : event
      this.$store.dispatch('ecs/onClusterClickRow', item)
    },
    onClusterItemSelected(event){
      this.$store.dispatch('ecs/onClusterItemSelected', event)
    },
    onServiceClickRow(event, row){
      const item = row && row.item ? row.item : event
      this.$store.dispatch('ecs/onServiceClickRow', item)
    },
    onServiceItemSelected(event){
      this.$store.dispatch('ecs/onServiceItemSelected', event)
    },
    onTaskClickRow(event, row){
      const item = row && row.item ? row.item : event
      this.$store.dispatch('ecs/onTaskClickRow', item)
    },
    onTaskItemSelected(event){
      this.$store.dispatch('ecs/onTaskItemSelected', event)
    },
  },
  created () {
    this.$store.dispatch('ecs/getClusters')
  }
}
</script>

<style>
.tabTable {
  margin-left: 25px;
  margin-top: 25px;
}
</style>







