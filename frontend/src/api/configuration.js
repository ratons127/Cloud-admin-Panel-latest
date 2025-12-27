import http from './http'

export default {
  setConfiguration(data, cb, errorCb) {
    http.post(`${data.serverPath}/configuration`, data.payload)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
}
