<template>
  <div>
    登陆界面公告，暂未实现
  </div>
</template>
<script>

import axios from 'axios'

export default {
  name: 'Announcement',
  props: { host: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      data: [],
      columns: [
        { title: '类型', dataIndex: 'type' },
        { title: '标题', dataIndex: 'title' },
        { title: '开始时间', dataIndex: 'startTime' },
        { title: '到期时间', dataIndex: 'expireTime' },
        { title: 'Action', scopedSlots: { customRender: 'action' } }
      ],
      version: 0,
      visible: false,
      modalTitle: '添加公告',
      modalContent: '',
      modalType: 'add'
    }
  },
  mounted () {
    this.getAnnouncement()
  },
  watch: {
    host (val) {
      this.getAnnouncement()
    }
  },
  methods: {
    getAnnouncement () {
      if (this.host === '') {
        return
      }
      const url = 'http://' + this.host + '/announcement/get'
      axios({ url: url, method: 'post', data: { version: this.version } })
      .then(res => {
        const data = res.data
        if (data.success) {
          const ret = data.data
          this.version = ret.version
          this.data = ret.announcement
        } else {
          this.$message.error(data.message)
        }
      })
    },
    addAnnouncement (content) {
      const data = JSON.parse(content)
      const url = 'http://' + this.host + '/announcement/add'
      axios({ url: url, method: 'post', data: data })
      .then(res => {
        const data = res.data
        if (data.success) {
          this.getAnnouncement()
        } else {
          this.$message.error(data.message)
        }
      })
      .finally(() => {
        this.visible = false
      })
    },
    delAnnouncement (id) {
      const url = 'http://' + this.host + '/announcement/delete'
      axios({ url: url, method: 'post', data: { id: id } })
      .then(res => {
        const data = res.data
        if (data.success) {
          this.getAnnouncement()
        } else {
          this.$message.error(data.message)
        }
      })
    },
    openModal (value) {
      this.modalContent = value
      this.visible = true
    },
    handleModalOk () {
      if (this.modalType === 'add') {
        this.addAnnouncement(this.modalContent)
      } else {

      }
    }
  }
}
</script>
