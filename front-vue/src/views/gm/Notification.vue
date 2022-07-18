<template>
  <div>
    <a-card :bordered="false">
      <div style="marginBottom: 24px;font-size:16px " >
        服务状态:
        <a-switch v-model="isOpen" checked-children="在线" un-checked-children="离线" @change="onSwitchChange" />

      </div>
      <a-tabs :default-active-key="activeKey" @change="tabChange">
        <a-tab-pane key="online" tab="在线公告">
          <a-input v-model="onlineData.title" placeholder="公告标题" size="large" style="width: 50%;"/></br></br>
          <a-textarea
            v-model="onlineData.content"
            placeholder="公告内容"
            :auto-size="{ minRows: 6, maxRows: 8 }"/>
        </a-tab-pane>
        <a-tab-pane key="offline" tab="离线公告" force-render>
          <a-input v-model="offlineData.title" placeholder="公告标题" size="large" style="width: 50%;"/></br></br>
          <a-textarea
            v-model="offlineData.content"
            placeholder="公告内容"
            :auto-size="{ minRows: 6, maxRows: 8 }"/>
        </a-tab-pane>
      </a-tabs>
      </br>
      <a-button @click="handleUpdate">更改当前公告</a-button>
    </a-card>

  </div>
</template>
<script>

import axios from 'axios'

export default {
  name: 'Notification',
  props: { host: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      isOpen: false,
      switchLoading: false,
      onlineData: {},
      offlineData: {},
      activeKey: 'online'
    }
  },
  mounted () {
    this.getNotification()
  },
  watch: {
    host (val) {
      this.getNotification()
    }
  },
  methods: {
    onSwitchChange (checked) {
      if (this.host === '') {
        this.$message.info('请选择一个web节点')
        return
      }
      this.isOpen = checked
      // console.log(this.isOpen)
      this.switchLoading = true
      const url = 'http://' + this.host + '/serverstatus/set'
      axios({ url: url, method: 'post', data: { closed: !this.isOpen } })
      .then(res => {
        const data = res.data
        if (data.success) {
          this.$message.success('状态更改成功')
          setTimeout(() => {
            this.getNotification()
            // console.log('fffff')
          }, 2000)
        } else {
          this.isOpen = !this.isOpen
          this.$message.error(data.message)
        }
      })
      .finally(() => {
        this.switchLoading = false
      })
    },
    getNotification () {
      if (this.host === '') {
        return
      }
      const url = 'http://' + this.host + '/notification/getAll'
      axios({ url: url, method: 'post' })
      .then(res => {
        const data = res.data
        if (data.success) {
          const ret = data.data
          this.isOpen = !ret.isClosed
          this.onlineData = ret.notifications[0]
          this.offlineData = ret.notifications[1]
          // console.log(ret)
        } else {
          this.$message.error(data.message)
        }
      })
    },
    tabChange (activeKey) {
      this.activeKey = activeKey
    },
    handleUpdate () {
      if (this.host === '') {
        this.$message.info('请选择一个web节点')
        return
      }
      try {
        let _type = 'offline'
        let notify = this.offlineData

        if (this.activeKey === 'online') {
          _type = 'online'
          notify = this.onlineData
        }

        const url = 'http://' + this.host + '/notification/update'
        axios({ url: url, method: 'post', data: { type: _type, notification: notify } })
        .then(res => {
          const data = res.data
          if (data.success) {
            this.$message.success('更改公告成功')
            this.getNotification()
          } else {
            this.$message.error(data.message)
          }
        })
      } catch (error) {
        this.$message.error('不是一个json字符串')
      }
    }

  }
}
</script>
