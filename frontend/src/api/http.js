import axios from 'axios'

const http = axios.create()

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('accessToken')
  const accountId = localStorage.getItem('awsAccountId')
  const region = localStorage.getItem('region')

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  if (accountId) {
    config.headers['X-AWS-ACCOUNT-ID'] = accountId
  }
  if (region) {
    config.headers['X-AWS-REGION'] = region
  }
  return config
})

export default http

