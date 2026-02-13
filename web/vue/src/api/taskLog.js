import httpClient from '../utils/httpClient'

export default {
  list(query, callback) {
    httpClient.get('/task/log', query, callback)
  },

  output(id, callback) {
    httpClient.get('/task/log/output', { id }, callback)
  },

  clear(callback) {
    httpClient.post('/task/log/clear', {}, callback)
  },

  stop(id, taskId, callback) {
    httpClient.post('/task/log/stop', { id, task_id: taskId }, callback)
  }
}
