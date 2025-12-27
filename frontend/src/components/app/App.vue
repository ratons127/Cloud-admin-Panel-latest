<template>
  <v-app id="inspire">
    <router-view />
  </v-app>
</template>

<script>
export default {
  name: 'App',
  created() {
    this.$store.dispatch('core/init')
  }
}
</script>

<style>
.v-list-item__title {
  font-size: 1.0em !important;
}
th,
td {
  width: 100px;
  overflow: hidden;
  white-space:nowrap;
  text-overflow:ellipsis;
}
.actions {
  position: absolute;
  right: 10px;
}
.action-btn {
  margin-left: 4px;
  margin-right: 4px;
  font-size:0.9em !important;
}
.actionBlock {
  margin-top: 10px;
  margin-bottom :10px;
}
.v-data-footer__select {
  display: none !important;
}
html{
  overflow-y: auto !important;
}
.title {
  margin-left:20px;
}
.titleLink {
  text-decoration: none;
  color: inherit !important;
}
.v-main__wrap .container {
  padding-bottom: 0;
}
</style>
