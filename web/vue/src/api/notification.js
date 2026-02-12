import httpClient from '../utils/httpClient'

export default {
  slack(callback) {
    httpClient.get('/system/slack', {}, callback)
  },
  updateSlack(data, callback) {
    httpClient.post('/system/slack/update', data, callback)
  },
  createSlackChannel(channel, callback) {
    httpClient.post('/system/slack/channel', { channel }, callback)
  },
  removeSlackChannel(channelId, callback) {
    httpClient.post(`/system/slack/channel/remove/${channelId}`, {}, callback)
  },
  mail(callback) {
    httpClient.get('/system/mail', {}, callback)
  },
  updateMail(data, callback) {
    httpClient.post('/system/mail/update', data, callback)
  },
  createMailUser(data, callback) {
    httpClient.post('/system/mail/user', data, callback)
  },
  removeMailUser(userId, callback) {
    httpClient.post(`/system/mail/user/remove/${userId}`, {}, callback)
  },
  webhook(callback) {
    httpClient.get('/system/webhook', {}, callback)
  },
  updateWebHook(data, callback) {
    httpClient.post('/system/webhook/update', data, callback)
  },
  createWebhookUrl(data, callback) {
    httpClient.post('/system/webhook/url', data, callback)
  },
  removeWebhookUrl(urlId, callback) {
    httpClient.post(`/system/webhook/url/remove/${urlId}`, {}, callback)
  },
  serverchan3(callback) {
    httpClient.get('/system/serverchan3', {}, callback)
  },
  updateServerchan3(data, callback) {
    httpClient.post('/system/serverchan3/update', data, callback)
  },
  createServerchan3Url(data, callback) {
    httpClient.post('/system/serverchan3/url', data, callback)
  },
  removeServerchan3Url(urlId, callback) {
    httpClient.post(`/system/serverchan3/url/remove/${urlId}`, {}, callback)
  }
}
