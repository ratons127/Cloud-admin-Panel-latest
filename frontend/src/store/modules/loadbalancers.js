import elbv2 from '../../api/elbv2'

const state = {
  search: '',
  headers: [{
    title: 'LoadBalancerName',
    align: 'start',
    key: 'LoadBalancerName',
    filterable: true
  },
  {
    title: 'LoadBalancerArn',
    key: 'LoadBalancerArn',
    filterable: true
  },
  {
    title: 'DNSName',
    key: 'DNSName',
    filterable: false
  },
  {
    title: 'Scheme',
    key: 'Scheme',
    filterable: false
  },
  {
    title: 'VpcId',
    key: 'VpcId',
    filterable: false
  },
  {
    title: 'Type',
    key: 'Type',
    filterable: false
  },{
    title: 'CanonicalHostedZoneId',
    key:  'CanonicalHostedZoneId',
    filterable: false
  },
  {
    title: 'Subnets',
    key: 'Subnets',
    filterable: false
  }],
  lb: []
}

const mutations = {
  updateSearch(state, value){
    state.search = value
  },
  updateLoadBalancers(state, value){
    state.lb = value
  }
}

const actions = {
  initLoadBalancers({state, commit, rootState}){
    elbv2.getLoadBalancers({
      serverPath: rootState.core.serverPath
    },
    (response) => {
      if (response.data.LoadBalancers){
        response.data.LoadBalancers.map(it => {
          it.AZ = it.AvailabilityZones.map(it => it.ZoneName).join(",")
          it.Subnets = it.AvailabilityZones.map(it => it.SubnetId).join(",")
          return it
        })
        commit('updateLoadBalancers', response.data.LoadBalancers)
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
