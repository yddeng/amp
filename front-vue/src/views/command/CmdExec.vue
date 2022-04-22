<template>
  <page-header-wrapper
    :breadcrumb="{}"
    title="执行命令"
    :content="'( ' + form.name + ' )'"
    @back="() => $router.go(-1)"
  >

    <a-card :bordered="false">
      <a-spin :spinning="spinning">
        <a-form-model
          ref="cmdExecForm"
          :model="form"
          :rules="formRule"
          :label-col="{ span: 2 }"
          :wrapper-col="{ span: 20 }"
        >
          <a-form-model-item label="命令类型">
            <span>SHELL</span>
          </a-form-model-item>
          <a-form-model-item label="命令内容" >
            <a-textarea
              v-model="form.context"
              :auto-size="{ minRows: 4, maxRows: 10 }"
              :readOnly="ctx_readonly"
              style="background: #0e0e0e; color: white"
            />
          </a-form-model-item>

          <a-form-model-item label="执行目录" :wrapper-col="{ span: 6}">
            <a-input v-model="form.dir" placeholder="非必填，默认为节点启动目录" />
          </a-form-model-item>

          <a-form-model-item :wrapper-col="{ span: 4}" prop="timeout">
            <span slot="label">
              超时时间&nbsp;
              <a-tooltip title="可设置范围10～86400秒，默认60秒。超时后，会强制终止进程">
                <a-icon type="question-circle-o" />
              </a-tooltip>
            </span>
            <a-input v-model="form.timeout" >
              <span slot="addonAfter" >秒</span>
            </a-input>
          </a-form-model-item>
          <a-form-model-item
            label="命令参数"
            v-show="hasArgs"
          >
            <a-row style="background:#F0F0F0">
              <a-col :span="4">&nbsp;&nbsp;参数名</a-col>
              <a-col :span="12">参数值</a-col>
            </a-row>
            <template v-for="(v,kk) in form.args" >
              <a-row :key="kk">
                <a-col :span="4">
                  <span >&nbsp;&nbsp;{{ kk }}</span>
                </a-col>
                <a-col :span="12">
                  <a-input v-model="form.args[kk]" />
                </a-col>
              </a-row>
            </template>
          </a-form-model-item>
          <a-form-model-item label="选择执行节点" prop="node" :wrapper-col="{ span: 4}">
            <a-select
              placeholder="选择一个可执行的节点"
              v-model="form.node"
            >
              <template v-for="(v,i) in nodes" >
                <a-select-option :value="v.name" :key="i">
                  {{ v.name }}
                </a-select-option>
              </template>
            </a-select>
          </a-form-model-item>
          <a-form-model-item :wrapper-col="{ span: 12,offset:2 }">
            <a-button type="primary" @click="onSubmit">
              执行命令
            </a-button>
          </a-form-model-item>
        </a-form-model>
      </a-spin>
    </a-card>

    <a-log-info
      :visible="logInfoVisible"
      :cmdName="logInfoCmd"
      :record="logInfoRecord"
      @cancel="logInfoCancel"
    >
    </a-log-info>
  </page-header-wrapper>
</template>

<script>
import { nodeList } from '@/api/node'
import { cmdExec } from '@/api/command'
import LogInfo from './modal/LogInfo'

export default {
  name: 'CmdExec',
  components: {
    'a-log-info': LogInfo
  },
  data () {
    return {
      formRule: {
        timeout: [{ validator: this.checkTimout, message: '请设置一个合适的时间', trigger: 'change' }],
        node: [{ required: true, message: '选择一个节点', trigger: 'blur' }]
      },
      ctx_readonly: true,
      hasArgs: false,
      nodes: [],
      form: {
        id: 0,
        name: '',
        context: '',
        dir: '',
        args: {},
        node: '',
        timeout: 60
      },
      spinning: false,

       logInfoVisible: false,
      logInfoCmd: '',
      logInfoRecord: {}
    }
  },
  mounted () {
    if (this.$route.params.record) {
      this.form = { ...this.$route.params.record, timeout: 60 }
      this.hasArgs = Object.keys(this.form.args).length > 0
    } else {
      this.$router.back()
    }
    this.getNodeList()
  },
  methods: {
    getNodeList () {
      const args = { pageNo: 1, pageSize: 1000 }
      nodeList(args)
        .then(res => {
          for (const idx in res.nodeList) {
            const v = res.nodeList[idx]
            if (v.online) {
              this.nodes.push(v)
            }
          }
        })
    },
    onSubmit () {
     this.$refs.cmdExecForm.validate(valid => {
        if (valid) {
          const args = {
            id: this.form.id,
            dir: this.form.dir,
            args: this.form.args,
            node: this.form.node,
            timeout: parseInt(this.form.timeout)
          }
          const tout = (this.form.timeout + 6) * 1000
          console.log(tout, this.form, args)
          this.spinning = true
          cmdExec(tout, args)
            .then(res => {
              console.log(res)
              this.spinning = false
              this.showLogInfo(res)
            })
            .catch(e => { })
            .finally(e => {
              this.spinning = false
            })
        } else {
          console.log('error submit!!')
          return false
        }
     })
    },
    checkTimout (rule, value, callback) {
      if (value < 10 || value > 86400) {
        callback(new Error('请输入一个合适的超时时长'))
      } else {
        callback()
      }
    },
    showLogInfo (record) {
      this.logInfoVisible = true
      this.logInfoCmd = this.cmdName
      this.logInfoRecord = record
    },
    logInfoCancel () {
      this.logInfoVisible = false
    }

  }
}
</script>

<style lang="less" scoped>

</style>
