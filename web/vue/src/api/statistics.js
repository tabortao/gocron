import httpClient from '../utils/httpClient'

export default {
  getOverview (callback) {
    httpClient.get('/statistics/overview', {}, callback)
  }
}
