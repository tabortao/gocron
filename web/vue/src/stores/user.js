import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: '',
    uid: '',
    username: '',
    isAdmin: false
  }),
  
  getters: {
    isLogin: (state) => state.token !== ''
  },
  
  actions: {
    setUser(user) {
      this.token = user.token || ''
      this.uid = user.uid || ''
      this.username = user.username || ''
      this.isAdmin = user.isAdmin || false
    },
    
    logout() {
      this.token = ''
      this.uid = ''
      this.username = ''
      this.isAdmin = false
    }
  },
  
  persist: {
    key: 'gocron-user',
    storage: localStorage,
    paths: ['token', 'uid', 'username', 'isAdmin']
  }
})
