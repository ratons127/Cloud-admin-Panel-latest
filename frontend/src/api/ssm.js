import http from './http'

export default {
  getInstances(data, cb, errorCb){
    http.get(`${data.serverPath}/service/ssm/instances`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  startSession(data, cb, errorCb){
    http.post(`${data.serverPath}/service/ssm/session?instance=${data.instanceId}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
}
