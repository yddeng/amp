<template>
  <a-modal
    title="执行日志"
    :width="1000"
    centered
    footer=""
    :maskClosable="true"
    :visible="visible"
    @cancel="cancel"
  >
    <a-card :bordered="false" size="small">
      <a-descriptions title="执行信息">
        <a-descriptions-item label="命令名称">{{ cmdName }}</a-descriptions-item>
        <a-descriptions-item label="执行ID">{{ record.id }}</a-descriptions-item>
        <a-descriptions-item label="执行用户">{{ record.user }}</a-descriptions-item>
      </a-descriptions>
      <a-descriptions>
        <a-descriptions-item label="状态">
          <span v-if="record.result_at >0" style="color:#3CB371">已完成</span>
          <span v-else style="color:#48D1CC">执行中</span>
        </a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ record.create_at | formatDate }}</a-descriptions-item>
        <a-descriptions-item v-if="record.result_at > 0" label="执行时长">{{ record.result_at - record.create_at }} 秒</a-descriptions-item>
      </a-descriptions>
      <a-descriptions>
        <a-descriptions-item label="执行节点">{{ record.node }}</a-descriptions-item>
        <a-descriptions-item label="执行目录">"{{ record.dir }}"</a-descriptions-item>
      </a-descriptions>
      <a-divider style="margin-bottom:12px"/>
      <a-descriptions title="脚本内容">
        <div>
          <a-textarea
            :value="record.context"
            :auto-size="{ minRows: 4, maxRows: 10 }"
            style="background: #0e0e0e; color: white" />
        </div>
      </a-descriptions>
      <a-divider style="margin-bottom:12px"/>
      <a-descriptions title="执行结果">
        <div>
          <a-textarea
            :value="record.result"
            :auto-size="{ minRows: 4, maxRows: 10 }"
            style="background: #0e0e0e; color: white" />
        </div>
      </a-descriptions>
    </a-card>
  </a-modal>
</template>

<script>
import moment from 'moment'

export default {
  name: 'LogInfo',
  components: {
    moment
  },
  props: {
    visible: {
      type: Boolean,
      required: true
    },
    cmdName: {
      type: String,
      default: 'null'
    },
    record: {
      type: Object,
      default: null
    }
  },
  data () {
    return {

    }
  },
  filters: {
    formatDate (time) {
      return moment.unix(time).format('YYYY-MM-DD HH:mm:ss')
    }
  },
  mounted () {
    console.log(this.cmdName, this.record)
  },
  methods: {
    cancel () {
      this.$emit('cancel')
    }
  }
}
</script>
