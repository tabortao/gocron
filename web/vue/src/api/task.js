import request from '../utils/request'

export default {
  list: (query) => Promise.all([
    request.get('/task', { params: query }),
    request.get('/host/all')
  ]),
  
  detail: (id) => Promise.all([
    request.get(`/task/${id}`),
    request.get('/host/all')
  ]),
  
  update: (data) => request.post('/task/store', data),
  remove: (id) => request.post(`/task/remove/${id}`),
  enable: (id) => request.post(`/task/enable/${id}`),
  disable: (id) => request.post(`/task/disable/${id}`),
  run: (id) => request.get(`/task/run/${id}`)
}
