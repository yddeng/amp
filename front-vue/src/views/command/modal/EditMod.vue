<template>
  <a-modal
    :title="isEdit?'修改命令':'创建命令'"
    :width="1200"
    :visible="visible"
    :confirmLoading="confirmLoading"
    @ok="submitForm"
    @cancel="cancel"
    class="form"
  >
    <a-form-model
      ref="editMod"
      :model="form"
      :rules="formRule"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 18 }"
    >

      <a-form-model-item>
        <span slot="label">命令类型</span>
        SHELL
      </a-form-model-item>

      <a-form-model-item
        prop="name"
      >
        <span slot="label">命令名称</span>
        <a-input
          v-model="form.name"
          placeholder="请输入名称"
        />
      </a-form-model-item>

      <a-form-model-item
        label="执行目录"
      >
        <a-input
          v-model="form.dir"
          placeholder="非必填，默认为节点启动目录"
        />
      </a-form-model-item>

      <a-form-model-item prop="context" >
        <span slot="label">
          命令内容&nbsp;
          <a-tooltip :title="ctx_question_title">
            <a-icon type="question-circle-o" />
          </a-tooltip>
        </span>
        <a-textarea
          v-model="form.context"
          :auto-size="{ minRows: 6, maxRows: 10 }"
          style="background: #0e0e0e; color: white"
        />
      </a-form-model-item>

      <a-form-model-item
        label="命令参数"
        v-show="hasArgs"
      >
        <template v-for="(v,kk) in form.args" >
          {{ kk }}
          <a-input
            :key="kk"
            v-model="form.args[kk]"
            placeholder="参数值选填，执行时可调整"
          />
        </template>
      </a-form-model-item>
    </a-form-model>
  </a-modal>
</template>

<script>
import { cmdCreate, cmdUpdate } from '@/api/command'
import { getCmdContextVar } from '@/utils/strReg'
export default {
  name: 'EditMod',
  props: {
    visible: {
      type: Boolean,
      required: true
    },
    isEdit: {
      type: Boolean,
      required: true
    },
    model: {
      type: Object,
      default: () => null
    }
  },
  data: function () {
  return {
        ctx_question_title: '可在命令中设置变量值，以{{key}}的形式表示',
        confirmLoading: false,
        hasArgs: false,
        form: {
          id: 0,
          name: '',
          dir: '',
          context: '',
          args: {}
        },
        formRule
      }
  },
  watch: {
    model: {
      handler (nv) {
        this.hasArgs = false
        if (nv && nv.id) { // 编辑
          this.form = { ...nv }
        } else {
          this.form = { name: '', dir: '', context: '', args: {} }
        }
      },
      deep: true,
      immediate: true
    },
    'form.context' (nval, oval) {
      if (nval) {
      const nm = getCmdContextVar(nval)
      const args = {}
      for (const k in nm) {
        args[k] = ''
        if (this.form.args.hasOwnProperty(k)) {
          args[k] = this.form.args[k]
        }
      }
      this.form.args = args
      // console.log(nm, args, this.form.args, Object.keys(this.form.args).length)
      if (Object.keys(this.form.args).length > 0) {
          this.hasArgs = true
      } else {
        this.hasArgs = false
      }
      }
    }
  },
  methods: {
    submitForm () {
      this.$refs.editMod.validate(valid => {
        if (valid) {
          this.confirmLoading = true
          const args = { ...this.form }
          if (this.isEdit) {
            cmdUpdate(args)
              .then(res => {
                this.$emit('success')// 通知外部页面 添加成功
              })
              .catch(e => { })
              .finally(e => {
                this.confirmLoading = false
              })
          } else {
            cmdCreate(args)
              .then(res => {
                this.$emit('success')// 通知外部页面 添加成功
              })
              .catch(e => { })
              .finally(e => {
                this.confirmLoading = false
              })
            }
        } else {
          this.confirmLoading = false
        }
      })
    },
    cancel () {
      this.$refs.editMod.resetFields()
      this.$emit('cancel')
    }
  }
}

const formRule = {
  name: [{ required: true, message: '命令名称', trigger: 'change' }],
  context: [{ required: true, message: '命令内容', trigger: 'change' }]
}
</script>
<style lang="less" scoped>
</style>
