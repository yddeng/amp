<template>
  <div>
    <a-card :bordered="false">
      <div style="marginBottom: 24px ">
        <a-button
          type="primary"
          @click="showEdit('add')"
          icon="plus"
        >新增</a-button>
      </div>

      <s-table
        rowKey="username"
        ref="table"
        size="default"
        data-name="userList"
        :columns="columns"
        :handleMethods="handleDataList"
        :data="loadData"
      >
        <template #password="password">
          <span v-show="password.showPassword">{{ password.password }}</span>
          <span v-show="!password.showPassword">******</span>
          <a-icon class="icon-right" type="eye" v-if="!password.showPassword" @click="password.showPassword=true"/>
          <a-icon class="icon-right" type="eye-invisible" v-else @click="password.showPassword=false"/>
        </template>

        <template #action="action">
          <div class="table-operation">
            <!-- <a-tooltip>
            <template slot="title">编辑</template>
            <a-icon
              type="edit"
              @click="showEdit('edit',action)"
            />
          </a-tooltip> -->
            <a-popconfirm
              title="确定要删除吗？"
              @confirm="dele(action)"
            >
              <a-tooltip>
                <template slot="title">删除</template>
                <a-icon type="delete" />
              </a-tooltip>
            </a-popconfirm>
          </div>
        </template>
      </s-table>
    </a-card>

    <edit-model
      ref="editMod"
      :visible="visible"
      :model="mdl"
      :is-edit="isEdit"
      @cancel="handleCancel"
      @success="handleOk"
    ></edit-model>
  </div>
</template>

<script>
import { STable } from '@/components'
import { userList, userDelete } from '@/api/user'
import EditModel from './modules/EditModel'

export default {
  name: 'UserGroup',
  components: {
    STable,
    EditModel
  },
  data () {
    return {
      columns,
      // 查询条件参数
      // name: '',
      // 加载数据方法 必须为 Promise 对象
      loadData: parameter => {
       // const requestParameters = { keyword: this.name, ...parameter } //如果要查询要这样写
        return userList(parameter).then(res => {
          return res
        })
      },
      // 编辑/新增模态框相关
      visible: false,
      confirmLoading: false,
      isEdit: true,
      mdl: null
    }
  },
  methods: {
    handleDataList (list) {
      console.log(666, list)
      const tabList = list.map((item) => {
        return {
          ...item,
          showPassword: false
        }
      })
      return tabList
    },
    handleOk () {
      this.$refs.table.refresh()
      this.visible = false
      this.mdl = {}
    },

    handleCancel () {
      this.visible = false
      this.mdl = {}
    },
    showEdit (opation, record) {
      this.visible = true
      if (opation === 'edit') {
        this.isEdit = true
        this.mdl = { ...record }
      } else {
        this.isEdit = false
        this.mdl = {}
      }
    },
    dele (record) {
      userDelete({ username: [record.username] })
        .then(res => {
          this.$refs.table.refresh()
        })
        .catch(e => { })
        .finally(e => {
          this.loading = false
        })
    }
  }
}
const columns = [
  {
    title: '用户名称',
    dataIndex: 'username'
  },
  {
    title: '密码',
    width: '400px',
    scopedSlots: { customRender: 'password' },
    needTotal: true
  },
  {
    title: '操作',
    scopedSlots: { customRender: 'action' }
  }
]
</script>

<style lang="less" scoped>
.icon-right{
  display: inline-block;
  margin-left: 20px;
  color: #1890ff;
}
</style>
