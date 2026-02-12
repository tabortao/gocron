import httpClient from '../utils/httpClient'

export default {
  list(query, callback) {
    httpClient.get('/user', {}, callback)
  },

  detail(id, callback) {
    httpClient.get(`/user/${id}`, {}, callback)
  },

  update(data, callback) {
    httpClient.post('/user/store', data, callback)
  },

  login(username, password, twoFactorCode, rememberMe, callback, errorCallback) {
    const data = { username, password }
    if (twoFactorCode) {
      data.two_factor_code = twoFactorCode
    }
    if (rememberMe !== undefined) {
      data.remember_me = rememberMe
    }
    httpClient.post('/user/login', data, callback, errorCallback)
  },

  enable(id, callback) {
    httpClient.post(`/user/enable/${id}`, {}, callback)
  },

  disable(id, callback) {
    httpClient.post(`/user/disable/${id}`, {}, callback)
  },

  remove(id, callback) {
    httpClient.post(`/user/remove/${id}`, {}, callback)
  },

  editPassword(data, callback) {
    httpClient.post(
      `/user/editPassword/${data.id}`,
      {
        new_password: data.new_password,
        confirm_new_password: data.confirm_new_password
      },
      callback
    )
  },

  editMyPassword(data, callback) {
    httpClient.post(`/user/editMyPassword`, data, callback)
  },

  get2FAStatus(callback) {
    httpClient.get('/user/2fa/status', {}, callback)
  },

  setup2FA(callback) {
    httpClient.get('/user/2fa/setup', {}, callback)
  },

  enable2FA(secret, code, callback) {
    httpClient.post('/user/2fa/enable', { secret, code }, callback)
  },

  disable2FA(code, callback, errorCallback) {
    httpClient.post('/user/2fa/disable', { code }, callback, errorCallback)
  }
}
