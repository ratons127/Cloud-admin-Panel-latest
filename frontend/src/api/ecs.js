import http from './http'

export default {
  getClusters(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/ecs/clusters`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getServices(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/ecs/services?name=${data.name}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getTasks(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/ecs/tasks?name=${data.name}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  forceUpdate(data, cb, errorCb) {
    http.post(`${data.serverPath}/service/ecs/update?service=${data.service}&cluster=${data.cluster}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
}



