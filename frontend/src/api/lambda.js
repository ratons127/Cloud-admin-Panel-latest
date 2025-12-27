import http from './http'

export default {
  getLambdaList(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/lambda`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getFunction(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/lambda/function?name=${data.name}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
}
