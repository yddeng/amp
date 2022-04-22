<template>
  <div>
    <a-card :bordered="false">
      发一封邮件
      <a-form-model
        ref="processEdit"
        :model="form"
        :label-col="labelCol"
        :wrapper-col="wrapperCol"
      >
        <a-form-model-item label="标题" :wrapper-col="{ span: 6}">
          <a-input v-model="form.mail.title"/>
        </a-form-model-item>
        <a-form-model-item label="发送人" :wrapper-col="{ span: 6}">
          <a-input v-model="form.mail.sender"/>
        </a-form-model-item>
        <a-form-model-item label="邮件内容" :wrapper-col="{ span: 14}">
          <a-textarea v-model="form.mail.content" :auto-size="{ minRows: 5, maxRows: 8 }"/>
        </a-form-model-item>
        <a-form-model-item label="过期时间" :wrapper-col="{ span: 6}">
          <a-date-picker show-time @change="onTimeOk" />
        </a-form-model-item>
        <a-form-model-item
          v-for="(cfg, index) in awardInfos"
          :key="index"
          :wrapper-col="index === 0 ? wrapperCol:{ span: 14,offset:3}"
          :label="index === 0 ? '邮件奖励' : ''"
        >
          <a-input-number
            v-model="cfg.Type"
            style="width: 30%; margin-right: 8px"
          />
          <a-input-number
            v-model="cfg.ID"
            style="width: 30%; margin-right: 8px"
          />
          <a-input-number
            v-model="cfg.Count"
            style="width: 30%; margin-right: 8px"
          />
          <a-icon
            v-if="awardInfos.length > 0"
            class="dynamic-delete-button"
            type="minus-circle-o"
            @click="remAward(cfg)"
          />
        </a-form-model-item>
        <a-form-model-item
          :label="awardInfos.length === 0 ? '邮件奖励' : ''"
          :wrapper-col="awardInfos.length === 0 ? wrapperCol:{ span: 14,offset:3}"
        >
          <a-button type="dashed" style="width:60%" @click="addAward">
            <a-icon type="plus" /> 添加奖励
          </a-button>
        </a-form-model-item>
        <a-form-model-item label="邮件投递到" :wrapper-col="{ span: 6}">
          <a-switch checked-children="全局" un-checked-children="玩家" default-checked @change="onSwitchChange"/>
        </a-form-model-item>
        <a-form-model-item label="" :wrapper-col="{ span: 14,offset:3}" v-show="this.form.type==='user'">
          <a-input
            v-model="gameIDStr"
            placeholder="玩家GameID,多个玩家中间用 ';' 分割"
            style="width: 80%"
          />
        </a-form-model-item>
        <a-form-model-item :wrapper-col="{ span: 14,offset:6}">
          <a-button type="primary" @click="submitForm">
            发送
          </a-button>
        </a-form-model-item>
      </a-form-model>
    </a-card>

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
      labelCol: { span: 3 },
      wrapperCol: { span: 14 },
      form: {
        type: 'global',
        gameID: [],
        mail: {
          title: '',
          sender: '',
          createTime: 0,
          expireTime: 0,
          content: ''
        }
      },
      awardInfos: [],
      gameIDStr: ''
    }
  },
  mounted () {

  },
  watch: {

  },
  methods: {
    addAward () {
      this.awardInfos.push({
        Type: 0,
        ID: 0,
        Count: 0
      })
    },
    remAward (item) {
      const index = this.awardInfos.indexOf(item)
      if (index !== -1) {
        this.awardInfos.splice(index, 1)
      }
    },
    onTimeOk (value) {
      if (value) {
        this.form.mail.expireTime = value.unix()
      } else {
        this.form.mail.expireTime = 0
      }
    },
    onSwitchChange (checked) {
      if (checked) {
        this.form.type = 'global'
      } else {
        this.form.type = 'user'
        this.gameIDStr = ''
      }
    },
    submitForm () {
      if (this.host === '') {
        this.$message.info('请选择一个web节点')
        return
      }

      if (this.form.mail.title === '' || this.form.mail.sender === '' || this.form.mail.content === '') {
        this.$message.error('必填未填')
          return
      }

      for (const v of this.awardInfos) {
        if (v.Type <= 0 || v.ID <= 0 || v.Count <= 0) {
          this.$message.error('非法的邮件奖励')
          return
        }
      }
      this.form.mail.awards = {}
      if (this.awardInfos.length > 0) {
        this.form.mail.awards = { AwardInfos: this.awardInfos }
      }

      if (this.form.type === 'user') {
        if (this.gameIDStr === '') {
          this.$message.error('玩家ID不能为空')
          return
        }
        var numRe = new RegExp(/^[0-9]*$/)
        const gameIDs = this.gameIDStr.split(';')
        const gameID = []
        for (const v of gameIDs) {
          if (!numRe.test(v)) {
            this.$message.error('非法的玩家ID,只能为数字')
            return
          } else {
            gameID.push(parseInt(v))
          }
        }
        console.log(gameIDs, gameID)
        if (gameID.length === 0) {
          this.$message.error('玩家ID不能为空')
          return
        }
        this.form.gameID = gameID
      }

      console.log(this.form)

      const url = 'http://' + this.host + '/mail/add'
      axios({ url: url, method: 'post', data: this.form })
      .then(res => {
        const data = res.data
        if (data.success) {
          this.$message.success('发送成功')
        } else {
          this.$message.error(data.message)
        }
      })
    }
  }
}
</script>
