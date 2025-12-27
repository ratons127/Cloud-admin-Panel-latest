import elbv2 from '../../api/elbv2'

const state = {
  search: '',
  headers: [{
    title: 'TargetGroupName',
    align: 'start',
    key: 'TargetGroupName',
    filterable: true
  },
  {
    title: 'TargetGroupArn',
    key: 'TargetGroupArn',
    filterable: true
  },
  {
    title: 'TargetType',
    key: 'TargetType',
    filterable: false
  },
  {
    title: 'Protocol',
    key: 'Protocol',
    filterable: false
  },
  {
    title: 'Port',
    key: 'Port',
    filterable: false
  },
  {
    title: 'Matcher',
    key: 'Matcher',
    filterable: false
  },
  {
    title: 'HealthCheckProtocol',
    key: 'HealthCheckProtocol',
    filterable: false
  },
  {
    title: 'HealthCheckPath',
    key:'HealthCheckPath',
    filterable: false
  },
  {
    title: 'VpcId',
    key: 'VpcId',
    filterable: false
  }],
  tg: []
}

const mutations = {
  updateSearch(state, value){
    state.search = value
  },
  updateTargetGroups(state, value){
    state.tg = value
  }
}

const actions = {
  initTargetGroups({state, commit, rootState}){
    elbv2.getTargetGroups({
      serverPath: rootState.core.serverPath
    },
    (response) => {
      if (response.data.TargetGroups){
        response.data.TargetGroups.map(it => {
          if (it.Matcher){
            it.Matcher = it.Matcher.HttpCode
          }
          return it
        })
        commit('updateTargetGroups', response.data.TargetGroups)
      }
    },
    (err) => {
      console.log(err);
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
