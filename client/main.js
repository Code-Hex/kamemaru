import 'whatwg-fetch'

var Vue = require('vue')


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