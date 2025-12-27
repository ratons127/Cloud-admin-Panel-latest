import ecs from '../../api/ecs'

const state = {
  selectedTask: {},
  selectedService: {},
  tasks: [],
  services: [],
  clusters: [],
  search: '',
  tabState: 0,
  showActions: false,
  forceUpdateDialog: false,
  selectedClusterArr: [],
  selectedServiceArr: [],
  selectedTaskArr: [],
  selectedCluster: {},
  clusterHeaders: [{
    title: 'Cluster Name',
    align: 'start',
    key: 'ClusterName',
    filterable: true
  }],
  serviceHeaders: [{
    title: 'ServiceName',
    key: 'ServiceName',
  },{
    title: 'ServiceArn',
    key:'ServiceArn',
    align: 'start'
  },{
    title: 'Status',
    key:'Status'
  },
  {
    title: 'LaunchType',
    key: 'LaunchType'
  },
  {
    title: 'DesiredCount',
    key: 'DesiredCount'
  },
  {
    title: 'TaskDefinition',
    key: 'TaskDefinition'
  },
  {
    title: 'CreatedAt',
    key: 'CreatedAt'
  }],
  eventHeaders: [{
    title: 'CreatedAt',
    key: 'CreatedAt'
  },{
    title: 'Message',
    key: 'Message'
  },{
    title: 'Id',
    key: 'Id'
  }],
  taskHeaders: [{
    title:'CreatedAt',
    key: 'CreatedAt'
  },{
    title: 'TaskArn',
    key: 'TaskArn'
  },{
    title: 'Connectivity',
    key: 'Connectivity'
  }, {
    title: 'Container Number',
    key: 'Containers.length'
  }],
  containerHeaders: [{
    title: 'Name',
    key: 'Name'
  },{
    title:'Image',
    key: 'Image'
  },{
    title:'LastStatus',
    key:'LastStatus'
  },{
    title:'ContainerArn',
    key: 'ContainerArn'
  },{
    title:'ContainerPort',
    key:'NetworkBindings[0].ContainerPort'
  },{
    title:'HostPort',
    key:'NetworkBindings[0].HostPort'
  }],
}

const mutations = {
  updateSelectedClusterArr(state, value){
    state.selectedClusterArr = value
  },
  updateSelectedServiceArr(state, value){
    state.selectedServiceArr = value
  },
  updateSelectedTaskArr(state, value){
    state.selectedTaskArr = value
  },
  updateSearch(state, value){
    state.search = value
  },
  updateTabState(state, value){
    state.tabState = value
  },
  setClusterSelection(state, value){
    state.selectedCluster = value;
    state.selectedClusterArr = [value];
  },
  updateClusters(state, value){
    state.clusters = value
  },
  updateServices(state, value){
    state.services = value
  },
  updateTasks(state, value){
    state.tasks = value
  },
  setServiceSelection(state, value){
    state.showActions = true
    state.selectedService = value
    state.selectedServiceArr = [value]
  },
  setTaskSelection(state, value){
    state.selectedTask = value
    state.selectedTaskArr = [value]
  },
  updateForceUpdateDialog(state, value){
    state.forceUpdateDialog = value
  }
}

const actions = {
  onClusterClickRow({commit, dispatch}, item){
    commit('setClusterSelection', item)
    dispatch('initCluster')
  },
  onClusterItemSelected({commit, dispatch}, selected){
    if (Array.isArray(selected) && selected.length > 0) {
      commit('setClusterSelection', selected[0])
      dispatch('initCluster')
    } else {
      commit('updateSelectedClusterArr', [])
    }
  },
  onServiceClickRow({commit}, item){
    commit('setServiceSelection', item)
  },
  onServiceItemSelected({commit}, selected){
    if (Array.isArray(selected) && selected.length > 0) {
      commit('setServiceSelection', selected[0])
    } else {
      commit('updateSelectedServiceArr', [])
    }
  },
  onTaskClickRow({commit}, item){
    commit('setTaskSelection', item)
  },
  onTaskItemSelected({commit}, selected){
    if (Array.isArray(selected) && selected.length > 0) {
      commit('setTaskSelection', selected[0])
    } else {
      commit('updateSelectedTaskArr', [])
    }
  },
  getClusters({commit, state, dispatch, rootState}, noFirstSelect){
    ecs.getClusters({
      serverPath: rootState.core.serverPath,
    },
    (response) => {
      if (response.data.Clusters){
        commit('updateClusters', response.data.Clusters)
        if (!noFirstSelect && state.clusters.length > 0){
          commit('setClusterSelection', state.clusters[0])
          dispatch('initCluster')
        }
      }
    },
    (err) => {
      console.log(err);
    })
  },
  initCluster({state, rootState, commit}){
    ecs.getServices({
      serverPath: rootState.core.serverPath,
      name: state.selectedCluster.ClusterName
    },
    (response) => {
      if (response.data.Services){
        commit('updateServices', response.data.Services)
      }
    },
    (err) => {
      console.log(err);
    })
    ecs.getTasks({
      serverPath: rootState.core.serverPath,
      name: state.selectedCluster.ClusterName
    },
    (response) => {
      if (response.data.Tasks){
        commit('updateTasks', response.data.Tasks)
      }
    },
    (err) => {
      console.log(err);
    })
  },
  refresh({dispatch}){
    dispatch('getClusters')
  },
  forceUpdate({state, rootState, commit, dispatch}){
    ecs.forceUpdate({
      serverPath: rootState.core.serverPath,
      cluster: state.selectedCluster.ClusterName,
      service: state.selectedService.ServiceName
    },
    (response) => {
      dispatch('core/dispatchMessage', {color: 'success', message: "Force update sent"}, {root:true})
      commit('updateForceUpdateDialog', false)
    },
    (err) => {
      console.log(err);
      if (err.response) {
        return dispatch('core/dispatchMessage', {color: 'error', message: `${err.response.status} : ${JSON.stringify(err.response.data)}`}, {root:true})
      }
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
