<template>
  <v-menu
    v-model="dateMenu"
    :close-on-content-click="false"
    transition="scale-transition"
    offset="8"
  >
    <template v-slot:activator="{ props }">
      <v-text-field
      :value="getDateFormatted"
      label="Date filter"
      prepend-inner-icon="mdi-calendar"
      readonly
      hide-details
      style='margin-left:25px; margin-top:0; padding-top:0;max-width:350px;'
      v-bind="props"
      ></v-text-field>
    </template>
    <v-date-picker v-model="startDateStr" no-title scrollable>
      <v-spacer></v-spacer>
      <v-btn  variant="text" color="primary" @click="dateMenu = false">Cancel</v-btn>
      <v-btn  variant="text" color="primary" @click="startDateStrChange">OK</v-btn>
    </v-date-picker>
  </v-menu>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'DateMenu',
  computed: {
    ...mapGetters('cloudwatchlogs', {
      getDateFormatted: 'getDateFormatted'
    }),
    dateMenu: {
      get() {
        return this.$store.state.cloudwatchlogs.dateMenu
      },
      set(value) {
        this.$store.commit('cloudwatchlogs/updateDateMenu', value)
      }
    },
    startDateStr: {
      get() {
        return this.$store.state.cloudwatchlogs.startDateStr
      },
      set(value) {
        this.$store.commit('cloudwatchlogs/updateStartDateStr', value)
      }
    }
  },
  methods: {
    startDateStrChange(){
      this.$store.commit('cloudwatchlogs/updateDateMenu', false)
      this.$store.dispatch('cloudwatchlogs/onStartDateChange')
    }
  }
}
</script>


