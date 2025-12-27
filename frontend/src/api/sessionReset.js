import http from './http'

export default {
  reset(data, cb, errorCb) {
    http.post(`${data.serverPath}/session_reset`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
}
