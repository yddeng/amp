<template>
  <div>
    <a-card :bordered="false">
      <div style="marginBottom: 24px " >
        <a-button
          type="primary"
          @click="openModal('')"
        >添加公告</a-button>
      </div>
      <a-table
        bordered
        :columns="columns"
        :data-source="data"
        :pagination="false"
        size="small"
        :rowKey="(record,index) => index"
      >
        <template slot="action" slot-scope="text, record">
          <div >
            <a-popconfirm title="确定要清除当前表吗？" @confirm="delAnnouncement(record.id)">
              <a>删除</a>
            </a-popconfirm>
          </div>
        </template>

      </a-table>
    </a-card>

    <a-modal
      v-model="visible"
      :title="modalTitle"
      :width="800"
      centered
      @ok="handleModalOk"
    >
      <div style="height:170px">
        <a-textarea
          v-model="modalContent"
          placeholder="josn格式的文本"
          :auto-size="{ minRows: 6, maxRows: 8 }"/>
      </div>
    </a-modal>
  </div>
</template>
<script>

import axios from 'axios'
import moment from 'moment'

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
        { title: '开始时间', dataIndex: 'startTime', customRender: (text, record) => this.momentTime(text) },
        { title: '到期时间', dataIndex: 'expireTime', customRender: (text, record) => this.momentTime(text) },
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
    momentTime (text) {
        return text === 0 ? '-' : moment.unix(text).format('YYYY-MM-DD HH:mm:ss')
    },
    getAnnouncement () {
      if (this.host === '') {
        this.$message.info('请选择一个web节点')
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
    addAnnouncement (data) {
      if (this.host === '') {
        this.$message.info('请选择一个web节点')
        return
      }
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
      if (this.modalContent === '') {
        this.$message.info('内容不能为空')
        return
      }

      try {
        const data = JSON.parse(this.modalContent)
        if (this.modalType === 'add') {
          this.addAnnouncement(data)
        } else {

        }
      } catch (error) {
        this.$message.error('不是一个json字符串')
      }
    }
  }
}
</script>
