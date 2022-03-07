<template>
  <page-header-wrapper
    :breadcrumb="{}"
    title="执行命令"
    :content="'( ' + cmdName + ' )'"
    @back="() => $router.go(-1)"
  >

    <a-card :bordered="false">
      <s-table
        rowKey="id"
        ref="table"
        size="middle"
        :bordered="border"
        data-name="logList"
        :columns="columns"
        :handleMethods="handleDataList"
        :data="loadData"
      >

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

        <template slot="statue" slot-scope="text, record">
          <span v-if="record.result_at > 0" style="color:#3CB371">已完成</span>
          <span v-else style="color:#48D1CC">执行中</span>
        </template>

        <template slot="action" slot-scope="text, record">
          <div >
            <a @click="showLogInfo(record)">查看详情</a>
          </div>
        </template>
      </s-table>
    </a-card>
  </page-header-wrapper>
</template>

<script>
import moment from 'moment'
import { STable } from '@/components'
import { cmdLog } from '@/api/command'
import LogInfo from './modal/LogInfo'

const columns = [
  {
    title: '执行ID',
    dataIndex: 'id'
  },
  {
    title: '执行内容',
    dataIndex: 'context',
    scopedSlots: { customRender: 'context' }
  },
  {
    title: '执行时间',
    dataIndex: 'create_at',
    customRender: (text) => moment.unix(text).format('YYYY-MM-DD HH:mm:ss')
  },
  {
    title: '执行用户',
    dataIndex: 'user',
    ellipsis: true
  },
  {
    title: '状态',
    scopedSlots: { customRender: 'statue' }
  },
  {
    title: '操作',
    scopedSlots: { customRender: 'action' }
  }
]
export default {
  name: 'CmdLog',
  components: {
    STable
  },
  data () {
    return {
      columns,
      cmdName: '',
      border: false,
      loadData: parameter => {
        const args = { id: this.$route.params.id, ...parameter }
        return cmdLog(args).then(res => {
          console.log(res)
          return res
        })
      },
      copyTitle: '复制'
    }
  },
  mounted () {
    console.log(this.$route)
    if (!this.$route.params.name) {
      this.$router.back()
    } else {
      this.cmdName = this.$route.params.name
    }
  },
  methods: {
    handleDataList (list) {
      const tabList = list.map((item) => {
        return {
          ...item
        }
      })
      return tabList
    },
    showLogInfo (record) {
    this.$dialog(LogInfo,
      // component props
      {
        cmdName: this.cmdName,
        record
      },
      // modal props
      {
        title: '执行日志',
        width: 1000,
        centered: true,
        maskClosable: true
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

</script>

<style lang="less" scoped>
</style>
