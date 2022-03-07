<template>
  <div>
    <a-row :gutter="24">
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '24px' }" >
        <a-card class="header-card" :bordered="false">
          <a-row>
            <a-col :span="18">
              <span class="header-card-value">{{ onlineNodeCnt }}/{{ nodes.length }}</span><br/>
              <span class="header-card-title">在线/总数</span><br/>
              <span class="header-card-desc">已注册的节点数</span>
            </a-col>
            <a-col :span="6"><a-icon class="header-card-icon" type="cloud-server"/> </a-col>
          </a-row>
        </a-card>
      </a-col>
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '24px' }">
        <a-card class="header-card" :bordered="false">
          <a-row>
            <a-col :span="18">
              <span class="header-card-value">{{ cmdCount }}</span><br/>
              <span class="header-card-title">任务数</span><br/>
              <span class="header-card-desc">已注册的任务</span>
            </a-col>
            <a-col :span="6"><a-icon class="header-card-icon" type="calendar"/> </a-col>
          </a-row>
        </a-card>
      </a-col>
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '24px' }">
        <a-card class="header-card" :bordered="false">
          <a-row>
            <a-col :span="18">
              <span class="header-card-value">{{ cmdSuccess }}</span><br/>
              <span class="header-card-title">成功</span><br/>
              <span class="header-card-desc">任务成功次数</span>
            </a-col>
            <a-col :span="6"><a-icon class="header-card-icon" type="check"/> </a-col>
          </a-row>
        </a-card>
      </a-col>
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '24px' }">
        <a-card class="header-card" :bordered="false">
          <a-row>
            <a-col :span="18">
              <span class="header-card-value">{{ cmdFailed }}</span><br/>
              <span class="header-card-title">失败</span><br/>
              <span class="header-card-desc">任务失败次数</span>
            </a-col>
            <a-col :span="6"><a-icon class="header-card-icon" type="close"/> </a-col>
          </a-row>
        </a-card>
      </a-col>
    </a-row>

    <a-card :bordered="false">
      <div :style="{ marginBottom: '24px' }">
        <a-button
          type="primary"
          @click="showEdit('create')"
          icon="plus"
        >创建命令</a-button>
      </div>

      <s-table
        rowKey="name"
        ref="table"
        size="default"
        data-name="cmdList"
        :columns="columns"
        :handleMethods="handleDataList"
        :data="loadData"
      >
        <template slot="name" slot-scope="text" >
          <a-tooltip v-if="text.length > 10" :title="text">
            {{ text.slice(0, 10) + '...' }}
          </a-tooltip>
          <span v-else>{{ text }}</span>
        </template>
        <template slot="context" slot-scope="text" >
          <a-tooltip v-if="text.length > 30" :title="text" style="color:#00BFFF">
            {{ text.slice(0, 30) + '...' }}
          </a-tooltip>
          <a-tooltip v-else :title="text" style="color:#00BFFF">
            {{ text }}
          </a-tooltip>
          <a-tooltip >
            <template slot="title">{{ copyTitle }}</template>
            <a-icon
              type="copy"
              v-clipboard:copy="text"
              @click="copyClick"/>
          </a-tooltip>
        </template>

        <template slot="action" slot-scope="text, record">
          <div >
            <a @click="execPage(record)">执行</a>
            <a-divider type="vertical" />
            <a @click="showEdit('edit',record)">修改</a>
            <a-divider type="vertical" />
            <a @click="logPage(record)">日志</a>
            <a-divider type="vertical" />
            <a-popconfirm title="确定要删除吗？" @confirm="deleteCmd(record)">
              <a-icon slot="icon" type="question-circle-o" style="color: red" />
              <a style="color:red;">删除</a>
            </a-popconfirm>
          </div>
        </template>
      </s-table>
    </a-card>

    <edit-mod
      ref="editMod"
      :visible="visible"
      :model="mdl"
      :is-edit="isEdit"
      @cancel="handleCancel"
      @success="handleOk"
    ></edit-mod>
  </div>
</template>

<script>
import moment from 'moment'
import { STable } from '@/components'
import { cmdList, cmdDelete } from '@/api/command'
import { nodeList } from '@/api/node'
import EditMod from './modal/EditMod'

export default {
  name: 'CmdList',
  components: {
    STable,
    EditMod
  },
  data () {
    return {
      columns,
      // 查询条件参数
      // name: '',
      // 加载数据方法 必须为 Promise 对象
      loadData: parameter => {
       // const requestParameters = { keyword: this.name, ...parameter } //如果要查询要这样写
        return cmdList(parameter).then(res => {
          console.log(res)
          this.cmdCount = res.totalCount
          this.cmdSuccess = res.success
          this.cmdFailed = res.failed
          return res
        })
      },
      // 编辑/新增模态框相关
      visible: false,
      confirmLoading: false,
      isEdit: true,
      mdl: null,
      copyTitle: '复制',

      nodes: [],
      onlineNodeCnt: 0,
      cmdCount: 0,
      cmdSuccess: 0,
      cmdFailed: 0
    }
  },
  mounted () {
    this.getNodeList()
  },
  methods: {
    getNodeList () {
      this.nodes = []
      this.onlineNodeCnt = 0
      const args = { pageNo: 1, pageSize: 1000 }
      nodeList(args)
        .then(res => {
          this.nodes = res.nodeList
          for (const idx in res.nodeList) {
            const v = res.nodeList[idx]
            if (v.online) {
              this.onlineNodeCnt += 1
            }
          }
        })
    },
    handleDataList (list) {
      const tabList = list.map((item) => {
        return {
          ...item
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
    execPage (record) {
      this.$router.push({ name: 'cmdexec', params: { record: { ...record } } })
    },
    logPage (record) {
      this.$router.push({ name: 'cmdlog', params: { id: record.id, name: record.name } })
    },
    deleteCmd (record) {
      cmdDelete({ id: record.id })
        .then(res => {
          this.$refs.table.refresh()
        })
        .catch(e => { })
        .finally(e => {
          this.loading = false
        })
    },
    copyClick () {
      this.copyTitle = '复制成功'
      setTimeout(() => {
        this.copyTitle = '复制'
      }, 1500)
    }
  }
}
const columns = [
  {
    title: '命令名称',
    dataIndex: 'name',
    scopedSlots: { customRender: 'name' }
  },
  {
    title: '命令内容',
    dataIndex: 'context',
    scopedSlots: { customRender: 'context' }
  },
  {
    title: '执行次数',
    align: 'center',
    customRender: (text) => text + ' 次',
    dataIndex: 'call_no'
  },
  {
    title: '创建时间',
    dataIndex: 'create_at',
    customRender: (text) => moment.unix(text).format('YYYY-MM-DD HH:mm:ss')
  },
  {
    title: '最后修改人',
    dataIndex: 'user'
  },
  {
    title: '最后修改时间',
    dataIndex: 'update_at',
    customRender: (text) => moment.unix(text).format('YYYY-MM-DD HH:mm:ss')
  },
  {
    title: '操作',
    scopedSlots: { customRender: 'action' }
  }
]
</script>

<style lang="less" scoped>

.header-card-value{
  font-size:30px;
  color:#7B7B7B,
}

.header-card-title{
  font-size:20px;
  color:#9D9D9D,
}

.header-card-icon{
  font-size:60px;
  color:#9D9D9D;
}
</style>
