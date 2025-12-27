import http from './http'

export default {
  getBuckets(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/s3`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getObjects(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/s3/objects?name=${data.name}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  deleteObject(data, cb, errorCb) {
    http.delete(`${data.serverPath}/service/s3/objects?name=${data.name}&key=${data.key}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  deleteBucket(data, cb, errorCb) {
    http.delete(`${data.serverPath}/service/s3?name=${data.name}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  createBucket(data, cb, errorCb) {
    http.put(`${data.serverPath}/service/s3?name=${data.name}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  uploadObject(data, cb, errorCb) {
    http.put(`${data.serverPath}/service/s3/objects`, data.formData)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  downloadObject(data) {
    window.open(`${data.serverPath}/service/s3/objects/download?name=${data.name}&key=${data.key}&_=${new Date().getTime()}`, '_blank');
  },
}
