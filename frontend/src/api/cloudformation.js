import http from './http'

export default {
  getStacks(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/cloudformation/stacks`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getEvents(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/cloudformation/events?name=${data.name}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getTemplate(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/cloudformation/template?name=${data.name}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getDescription(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/cloudformation/description?name=${data.name}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  createStack(data, cb, errorCb) {
    http.put(`${data.serverPath}/service/cloudformation/stack`, data.formData)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  deleteStack(data, cb, errorCb) {
    http.delete(`${data.serverPath}/service/cloudformation/stack?name=${data.name}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  updateTerminationProtection(data, cb, errorCb) {
    http.post(`${data.serverPath}/service/cloudformation/termination_protection?name=${data.name}&status=${data.status}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  updateStack(data, cb, errorCb) {
    http.post(`${data.serverPath}/service/cloudformation/stack`, data.formData)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
}
