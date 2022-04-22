<template>
  <page-header-wrapper
    title="GameMater"
    :breadcrumb="{}"
    :tab-list="tabList"
    :tab-active-key="tabActiveKey"
    :tab-change="handleTabChange"
  >
    <template slot="content">
      <div style="height: 40px;width:540px;margin:0 auto;">
        <a-select
          show-search
          placeholder="选择gm地址"
          style="width: 300px"
          @change="handleChange"
        >
          <a-select-option v-for="d in data" :key="d.title" :value="d.value">
            {{ d.field }}
          </a-select-option>
        </a-select>
        <a-popconfirm
          title="确定要删除当前gm地址吗？"
          @confirm="deleteKv">
          <a-button :disabled="disabled" icon="delete" type="danger" style="color: white" />
        </a-popconfirm>
        &nbsp; <a-button type="primary" icon="plus" @click="openModel">新增gm地址</a-button>
      </div>
    </template>

    <a-announcement v-if="tabActiveKey==='announcement'" :host="selectValue"></a-announcement>
    <a-notification v-else-if="tabActiveKey==='notification'" :host="selectValue"></a-notification>
    <a-mail v-else-if="tabActiveKey==='mail'" :host="selectValue"></a-mail>

    <a-modal v-model="visible" title="新增" ok-text="确认" cancel-text="取消" @ok="handleOk">
      <a-form-model
        :label-col="{ span: 4 }"
        :wrapper-col="{ span: 18 }">
        <a-form-model-item label="名称">
          <a-input v-model="inputName" placeholder="input placeholder" />
        </a-form-model-item>
        <a-form-model-item label="地址">
          <a-input v-model="inputAddress" placeholder="input placeholder" />
        </a-form-model-item>
      </a-form-model>
    </a-modal>
  </page-header-wrapper>
</template>

<script>
import { kvSet, kvGet, kvDelete } from '@/api/kv'
import Announcement from './Announcement'
import Notification from './Notification'
import Mail from './Mail'

export default {
  name: 'Flyfish',
  components: {
    'a-announcement': Announcement,
    'a-notification': Notification,
    'a-mail': Mail
  },
  data () {
    return {
      key: 'gm_addr',
      data: [],
      selectValue: '',
      tabList: [
        { key: 'announcement', tab: '公告' },
        { key: 'notification', tab: '通知' },
        { key: 'mail', tab: '邮件' }
      ],
      tabActiveKey: 'announcement',
      visible: false,
      inputName: '',
      inputAddress: '',
      disabled: true
    }
  },
  beforeMount () {
    this.getValue()
  },
  methods: {
    getValue () {
      kvGet({ key: this.key })
        .then(res => {
          console.log(res)
          this.selectValue = ''
          if (res.exist) {
            this.data = []
            const vals = res.value.split(';')
            for (const idx in vals) {
              const field = vals[idx]
              const fs = field.split('@')
              const k = fs[0]
              const v = fs[1]
              this.data.push({ title: k, value: v, field: field })
            }
          }
       })
    },
    valueToString () {
      if (this.data.length === 0) {
        return ''
      } else if (this.data.length === 1) {
        return this.data[0].field
      } else {
        let s = this.data[0].field
        for (let i = 1; i < this.data.length; i++) {
          s += ';' + this.data[i].field
        }
        return s
      }
    },
    deleteKv () {
      let index = -1
      for (let i = 0; i < this.data.length; i++) {
        const elem = this.data[i]
        if (elem.value === this.selectValue) {
          index = i
        }
      }
      if (index !== -1) {
        this.data.splice(index, 1)
      }
      // console.log(this.selectValue, index)
      if (this.data.length === 0) {
        kvDelete({ key: this.key })
        .then(res => {
          this.getValue()
        })
      } else {
        const value = this.valueToString()
        kvSet({ key: this.key, value: value })
        .then(res => {
          this.getValue()
        })
      }
    },
    handleChange (value) {
      this.selectValue = value
      if (this.selectValue === '') {
        this.disabled = true
      } else {
        this.disabled = false
      }
    },
    openModel () {
      this.inputName = ''
      this.inputAddress = ''
      this.visible = true
    },
    handleOk () {
      // console.log(this.inputName, this.inputAddress)
      if (this.inputName !== '' && this.inputAddress !== '') {
        let index = -1
        for (let i = 0; i < this.data.length; i++) {
          const elem = this.data[i]
          if (elem.title === this.inputName) {
            elem.value = this.inputAddress
            elem.field = this.inputName + '@' + this.inputAddress
            index = i
          }
        }
        if (index === -1) {
          this.data.push({ title: this.inputName, value: this.inputAddress, field: this.inputName + '@' + this.inputAddress })
        }
        const value = this.valueToString()
        kvSet({ key: this.key, value: value })
        .then(res => {
          this.visible = false
          this.getValue()
        })
      } else {
        this.visible = false
      }
    },
    handleTabChange (key) {
     this.tabActiveKey = key
    }
  }
}
</script>

<style lang="less" scoped>
</style>
