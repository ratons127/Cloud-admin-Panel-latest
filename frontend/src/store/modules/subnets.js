import ec2 from '../../api/ec2'

const state = {
  search: '',
  headers: [{
    title: 'SubnetId',
    key:'SubnetId',
    align: 'start',
    filterable: true
  },
  {
    title: 'Name',
    key: 'Name',
    filterable: true
  },
  {
    title: 'ARN',
    key: 'SubnetArn',
    filterable: true
  },
  {
    title: 'VpcId',
    key: 'VpcId',
    filterable: false
  },
  {
    title: 'CidrBlock',
    key: 'CidrBlock',
    filterable: false
  },
  {
    title: 'Stack Name',
    key: 'StackName',
    filterable: true
  }],
  subnets: []
}

const mutations = {
  updateSearch(state, value){
    state.search = value
  },
  updateSubnets(state, value){
    state.subnets = value
  }
}

const actions = {
  initSubnets({state, commit, rootState}){
    ec2.getSubnets({
      serverPath: rootState.core.serverPath
    },
    (response) => {
      if (response.data.Subnets){
        response.data.Subnets.map(it => {
          if (it.Tags){
            var nameEntry = it.Tags.filter(it => it.Key === "Name");
            var stackNameEntry = it.Tags.filter(it => it.Key === "aws:cloudformation:stack-name");
            it.Name = (nameEntry.length > 0) ? nameEntry[0].Value : "";
            it.StackName = (stackNameEntry.length > 0) ? stackNameEntry[0].Value : "";
          }
          return it
        });
        commit('updateSubnets', response.data.Subnets)
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
