'use strict';

import 'whatwg-fetch'

var $ = require('jQuery')
var Vue = require('vue')
Vue.use(require('vue-resource'))

var app = new Vue({
  el: '#app',
  data: {
    token: null,
    items: [],
  },
  created() {
    this.$http.get('/api/fetch')
      .then((response) => response.json())
      .then((json) => this.items = json)
      .catch((json) => console.error(json.reason))
  },
})

var uploader = new Vue({
  el: '#uploader',
  data: {
    images: null
  },
  methods: {
    select: function(e) {
       e.preventDefault()
      $('.select').click()
    },
    chooseImages: function(e) {
      e.preventDefault()
      let files = e.target.files
      if (!files.length) return
      this.images = files
      $('#selects').val($('.select').val().replace("C:\\fakepath\\", ""))
    },
    submit: function(e) {
      e.preventDefault()

      let data = new FormData()
      $.each(this.images, function(i, v) {
        data.append('files', v)
      })
      this.$http.post('/api/v1/upload', data, {
        headers: {
          'Content-Type': 'multipart/form-data',
          'Authorization': 'Bearer ' + nav.token,
        }
      }).then((response) => {
        this.images = null
        location.reload()
        return response.json()
      })
      .catch((json) => console.error(json.reason))
    }
  }
})

var nav = new Vue({
    el: '#navbar',
    data: {
      token: null
    },
    created() {
      this.token = localStorage.getItem('token')
      if (this.token)
        $('#uploader').removeClass('unvisible')
    },
    methods: {
        login: function(e) {
          e.preventDefault()
          let which = 'login'
          if (modal.active !== null) {
            $('#form-' + modal.active).removeClass('active')
            $('#' + modal.active + '-form').removeClass('active')
          }

          $('#login-modal').addClass('active')
          $('#form-' + which).addClass('active')
          $('#' + which + '-form').addClass('active')
          modal.active = which
        },
        logout: function(e) {
          e.preventDefault()
          localStorage.removeItem("token")
          location.reload()
        }
    }
});

var modal_submit_register = 'Register'
var modal_submit_password = 'Reset Password'
var modal_submit_login = 'Login'

var modal = new Vue({
    el: '#login-modal',
    data: {
      active: null,

      // Submit button text
      registerSubmit: modal_submit_register,
      passwordSubmit: modal_submit_password,
      loginSubmit: modal_submit_login,

      // Modal text fields
      registerName: '',
      registerEmail: '',
      registerPassword: '',
      loginUser: '',
      loginPassword: '',
      passwordEmail: '',

      // Modal error messages
      registerError: '',
      loginError: '',
      passwordError: '',
    },
    methods: {
      close: function(e) {
        e.preventDefault()
        if (e.target === this.$el)
          $('#login-modal').removeClass('active')
      },
      flip: function(which, e) {
        e.preventDefault()
        if (which !== this.active) {
          $('#form-' + this.active).removeClass('active')
          $('#form-' + which).addClass('active')
          $('#' + which + '-form').addClass('active')
          $('#' + this.active + '-form').removeClass('active')

          this.active = which
        }
      },
      submit: function(which, e) {
        var instance = this
        e.preventDefault()
        $('#'+which+'Submit').addClass('disabled')
        var data = {
            form: which
        }

        switch (which) {
          case 'register':
            data.username = this.registerName
            data.password = this.registerPassword
            this.registerSubmit = 'Registering...'
            break
          case 'login':
            data.username = this.loginUser
            data.password = this.loginPassword
            this.loginSubmit = 'Logging In...'
            break
        }

        instance.$http.post('/'+which, data)
        .then((response) => {
          $('#login-modal').removeClass('active')
          alert('Signed successed: '+ data.username)
          return response.json()
        })
        .then((json) => {
          if (typeof json.token !== 'undefined') {
            console.log(json.token)
            nav.token = json.token
            localStorage.setItem("token", json.token)
            $('#navbar').addClass('disabled')
            $('#uploader').removeClass('unvisible')
          }
          return json
        })
        .catch((json) => console.error(json.reason))

        switch (which) {
          case 'register':
            this.registerSubmit = modal_submit_register
            break
          case 'login':
            this.loginSubmit = modal_submit_login
            break
        }
        $('#'+which+'Submit').removeClass('disabled')
      }
    }
});

/*
var app = new Vue({
  el: '#app',
  data: {
    progress: 0,
    URL: 'https://www.youtube.com/watch?v=1jCo-B0FKoc'
  },
  methods: {
    reverseMessage: function() {
      fetch('/api/v1/download', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({url: this.URL})
      }).then(function(response) {
        var jsonStream = response.body.pipeThrough(new TextDecoder()).getReader()

        jsonStream.read().then(function process(result) {
          const json = JSON.parse(result.value)
          this.progress = json.percent
        })
      }).then(function() {
        console.log("Complete")
      })
    }
  }
})
*/