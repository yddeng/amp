<template>
  <div>
    <a-card :bordered="false">
      <div style="marginBottom: 24px " >
        <a-button
          type="primary"
          @click="openModal('add')"
        >添加公告</a-button>
      </div>
      <a-tabs default-active-key="AnnouncementType_System" >
        <a-tab-pane key="AnnouncementType_System" tab="系统公告">
          <a-table
            bordered
            :columns="columns"
            :data-source="systemData"
            :pagination="false"
            size="small"
            :rowKey="(record,index) => index"
          >
            <template slot="action" slot-scope="text, record">
              <div >
                <a @click="openModal('update',record)">更改</a>
                <a-divider type="vertical" />
                <a-popconfirm title="确定要清除当前表吗？" @confirm="delAnnouncement(record.id)">
                  <a>删除</a>
                </a-popconfirm>
              </div>
            </template>

          </a-table>
        </a-tab-pane>
        <a-tab-pane key="AnnouncementType_Activity" tab="活动公告">
          <a-table
            bordered
            :columns="columns"
            :data-source="activityData"
            :pagination="false"
            size="small"
            :rowKey="(record,index) => index"
          >
            <template slot="action" slot-scope="text, record">
              <div >
                <a @click="openModal('update',record)">更改</a>
                <a-divider type="vertical" />
                <a-popconfirm title="确定要清除当前表吗？" @confirm="delAnnouncement(record.id)">
                  <a>删除</a>
                </a-popconfirm>
              </div>
            </template>

          </a-table>
        </a-tab-pane>
      </a-tabs>

    </a-card>

    <a-modal
      v-model="visible"
      :title="modalTitle"
      :width="800"
      centered
      @ok="handleModalOk"
    >
      <a-form-model
        :model="form"
        :label-col="labelCol"
        :wrapper-col="wrapperCol"
      >
        <a-form-model-item label="类型" :wrapper-col="{ span: 6}">
          <a-select v-model="form.type" style="width: 120px">
            <a-select-option value="AnnouncementType_System">系统公告</a-select-option>
            <a-select-option value="AnnouncementType_Activity">活动公告</a-select-option>
          </a-select>
        </a-form-model-item>
        <a-form-model-item label="标题" :wrapper-col="{ span: 10}">
          <a-input v-model="form.title"/>
        </a-form-model-item>
        <a-form-model-item label="副标题" :wrapper-col="{ span: 10}">
          <a-input v-model="form.smallTitle" />
        </a-form-model-item>
        <a-form-model-item label="时间" :wrapper-col="{ span: 16}">
          <a-date-picker
            show-time
            :value="unixToMoment(form.startTime)"
            @change="onStartTime"
            placeholder="开始时间" />
          <a-date-picker
            show-time
            :value="unixToMoment(form.expireTime)"
            @change="onExpireTime"
            placeholder="结束时间" />
        </a-form-model-item>
        <a-form-model-item label="强制提醒" :wrapper-col="{ span: 14}">
          <a-switch checked-children="开" un-checked-children="关" v-model="form.remind"/>
        </a-form-model-item>
        <a-form-model-item
          v-show="form.content.length > 0"
          label="公告内容"
          :wrapper-col="{ span: 18}">
          <a-row
            v-for="(cfg, index) in form.content"
            :key="index">
            <a-col v-if="cfg.type==='1'">
              文本行
              <a-input v-model="cfg.text" style="width: 80%; margin-right: 5px" placeholder="文本"/>
              <a-icon
                v-if="form.content.length > 1"
                class="dynamic-delete-button"
                type="minus-circle-o"
                @click="remModalContent(cfg)"
              />
            </a-col>
            <a-col v-else>
              图片行
              <a-input-number v-model="cfg.imageSkip" style="width: 20%; margin-right: 5px"/>
              <a-input v-model="cfg.image" style="width: 50%; margin-right: 5px" placeholder="图片名或网络路径"/>
              <a-icon
                v-if="form.content.length > 1"
                class="dynamic-delete-button"
                type="minus-circle-o"
                @click="remModalContent(cfg)"
              />
            </a-col>
          </a-row>
        </a-form-model-item>
        <a-form-model-item
          :label="form.content.length === 0 ? '公告内容':''"
          :wrapper-col="form.content.length === 0 ? wrapperCol:{ span: 14,offset:4}">
          <a-button type="dashed" style="width:40%;margin-right: 5px" @click="addModalContent('0')">
            <a-icon type="plus" /> 添加图片行
          </a-button>
          <a-button type="dashed" style="width:40%" @click="addModalContent('1')">
            <a-icon type="plus" /> 添加文本行
          </a-button>
        </a-form-model-item>

      </a-form-model>
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
      systemData: [],
      activityData: [],
      columns: [
        { title: '标题', dataIndex: 'title' },
        { title: '开始时间', dataIndex: 'startTime', customRender: (text, record) => this.momentTime(text) },
        { title: '到期时间', dataIndex: 'expireTime', customRender: (text, record) => this.momentTime(text) },
        { title: 'Action', scopedSlots: { customRender: 'action' } }
      ],
      visible: false,
      version: 0,
      modalTitle: '添加公告',

      modalType: 'add',
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
      form: {
        id: 0,
        type: 'AnnouncementType_Activity', // AnnouncementType_System
        title: '',
        smallTitle: '',
        startTime: 0,
        expireTime: 0,
        remind: false,
        content: [] // {		"type": "0",		"imageSkip": 0,		"text": "",		"image": "Announce_ad_1"	}
      }
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
    unixToMoment (s) {
        return s === 0 ? null : moment.unix(s)
    },
    onStartTime (value) {
      this.form.startTime = value ? value.unix() : 0
    },
    onExpireTime (value) {
      this.form.expireTime = value ? value.unix() : 0
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
           console.log(ret)
          if (this.version !== ret.version) {
            this.version = ret.version
            this.systemData = []
            this.activityData = []
            for (const v of ret.announcement) {
                if (v.type === 'AnnouncementType_System') {
                  this.systemData.push({ ...v })
                } else {
                  this.activityData.push({ ...v })
                }
            }
          }
        } else {
          this.$message.error(data.message)
        }
      })
    },
    addAnnouncement (data) {
      const url = 'http://' + this.host + '/announcement/add'
      axios({ url: url, method: 'post', data: data })
      .then(res => {
        const data = res.data
        if (data.success) {
          this.getAnnouncement()
          this.$message.success('添加成功')
        } else {
          this.$message.error(data.message)
        }
      })
      .finally(() => {
        this.visible = false
      })
    },
    updateAnnouncement (data) {
      const url = 'http://' + this.host + '/announcement/update'
      axios({ url: url, method: 'post', data: data })
      .then(res => {
        const data = res.data
        if (data.success) {
          this.getAnnouncement()
          this.$message.success('修改成功')
        } else {
          this.$message.error(data.message)
        }
      })
      .finally(() => {
        this.visible = false
      })
    },
    delAnnouncement (id) {
      if (this.host === '') {
        this.$message.info('请选择一个web节点')
        return
      }
      const url = 'http://' + this.host + '/announcement/delete'
      axios({ url: url, method: 'post', data: { id: id } })
      .then(res => {
        const data = res.data
        if (data.success) {
          this.getAnnouncement()
          this.$message.success('删除成功')
        } else {
          this.$message.error(data.message)
        }
      })
    },
    addModalContent (type) {
      this.form.content.push({
        type: type,
        imageSkip: 0,
        image: '',
        text: ''
      })
    },
    remModalContent (item) {
      const index = this.form.content.indexOf(item)
      if (index !== -1) {
        this.form.content.splice(index, 1)
      }
    },
    openModal (status, record) {
      this.modalType = status
      if (status === 'update') {
        this.modalTitle = '更改公告'
        this.form = { ...record }
      } else {
        this.modalTitle = '添加公告'
        this.form = {
          id: 0,
          type: 'AnnouncementType_Activity', // AnnouncementType_System
          title: '',
          smallTitle: '',
          startTime: 0,
          expireTime: 0,
          remind: false,
          content: [] // {		"type": "0",		"imageSkip": 0,		"text": "",		"image": "Announce_ad_1"	}
        }
      }
      this.visible = true
      // console.log(status, record)
    },
    handleModalOk () {
      console.log(this.modalType, this.form)
      if (this.host === '') {
        this.$message.info('请选择一个web节点')
        return
      }
      if (this.modalType === 'add') {
        this.addAnnouncement(this.form)
      } else {
        this.updateAnnouncement(this.form)
      }
    }
  }
}
</script>
