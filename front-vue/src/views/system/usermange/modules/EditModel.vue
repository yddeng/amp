<template>
  <a-modal
    :title="isEdit?'编辑用户':'添加用户'"
    :width="800"
    :visible="visible"
    :confirmLoading="confirmLoading"
    @ok="submitForm"
    @cancel="cancel"
    class="form"
  >
    <a-form-model
      ref="editModal"
      :model="form"
      :rules="formRule"
      :label-col="labelCol"
      :wrapper-col="wrapperCol"
    >
      <a-form-model-item
        label="用户名"
        prop="username"
      >
        <a-input
          v-model="form.username"
          placeholder="用户名"
        />
      </a-form-model-item>
      <a-form-model-item
        label="密码"
        prop="password"
      >
        <a-input
          v-model="form.password"
          placeholder="请输入密码"
        />
      </a-form-model-item>
    </a-form-model>
  </a-modal>
</template>

<script>
 import { userAdd } from '@/api/user'
export default {
  name: 'EditModel',
  props: {
    visible: {
      type: Boolean,
      required: true
    },
    isEdit: {
      type: Boolean,
      required: true
    },
    loading: {
      type: Boolean,
      default: () => false
    },
    model: {
      type: Object,
      default: () => null
    }
  },
  data () {
    return {
      labelCol: { span: 6 },
      wrapperCol: { span: 12 },
      confirmLoading: false,
      form: {
        username: null,
        code: null
      },
      formRule
    }
  },
  watch: {
    model: {
      handler (nv) {
        if (nv && nv.username) { // 编辑
          this.form = { ...nv }
        } else { // 新增
          this.form = {}
        }
      },
      deep: true,
      immediate: true
    }
  },
  created () { },
  methods: {
    submitForm () {
      this.$refs.editModal.validate(valid => {
        if (valid) {
          this.confirmLoading = true
          if (this.isEdit) {
            // 修改
            /* const args = { ...this.form }
            updateGroup(args)
              .then(res => {
                this.$refs.editModal.resetFields()
                this.$emit('success')
              })
              .catch(e => { })
              .finally(e => {
                this.confirmLoading = false
              }) */
          } else {
            // 新增
            const args = { ...this.form }
            userAdd(args)
              .then(res => {
                this.$refs.editModal.resetFields()// 清除表单数据
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
      this.$refs.editModal.resetFields()
      this.$emit('cancel')
    }
  }
}

const formRule = {
  username: [{ required: true, message: '请输入用户名', trigger: 'change' }],
  password: [{ required: true, message: '请输入密码', trigger: 'change' }]
}
</script>
<style lang="less" scoped>
</style>
