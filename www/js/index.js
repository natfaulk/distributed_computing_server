const UPDATE_INTERVAL = 10*1000

window.onload=()=>{
  var app = new Vue({
    el: '#app',
    data: {
      tasks: [],
      clients: [],
      humanTime: function(_time) {
        if (_time === -1) return ''
        return moment.unix(_time).fromNow()
      },
      truncateID:function(_id) {
        if (_id === null) return 'None'
        if (_id.length <= 4) return _id
        return _id.slice(0,4)+'...'
      },
      updateStatus:async function() {
        let t_status=await getStatus()
        console.log(t_status)
        this.tasks=t_status.tasks
        this.clients=t_status.clients
      }
    },

    mounted() {
      this.updateStatus()
      setInterval(async ()=>{
        this.updateStatus()
      }, UPDATE_INTERVAL)
    }
  })
}

let getStatus = async () => {
  let response = await fetch('/status')
  let data=await response.json()
  return data
}
