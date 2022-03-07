<template>

  <page-header-wrapper
    :title="false"
    :breadcrumb="{}"
  >
    <template slot="content">
      <div style="height: 40px;width:540px;margin:0 auto;">
        <a-tree-select
          style="width: 260px;"
          :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
          :tree-data="treeData"
          @select="onChange"
          show-search
          allow-clear
          v-model="path"
          placeholder="选择分组"
        />
        <a-popconfirm
          title="确定要删除当前分组吗？"
          @confirm="deleteGroup">
          <a-button icon="delete" type="danger" style="color: white" />
        </a-popconfirm>
        &nbsp;
        <a-button icon="caret-right" type="primary" @click="startAllProcess">全部启动</a-button>
        &nbsp;
        <a-button icon="poweroff" type="primary" @click="stopAllProcess"> 全部停止</a-button>
      </div>
      <a-divider dashed orientation="left"></a-divider>
      <div style="padding: 10px 20px;">
        <a-row >
          <a-col :xs="8" :sm="8" :xl="4" class="header-col">
            <span class="header-card-title"><a-icon type="team"/> 总计</span><br/>
            <span class="header-card-value">{{ status.length }}</span><br/>
          </a-col>
          <a-col :xs="8" :sm="8" :xl="4" class="header-col">
            <span class="header-card-title"><a-icon type="alert" /> 预警中</span><br/>
            <span class="header-card-value">{{ status.alert }}</span><br/>
          </a-col>
          <a-col :xs="8" :sm="8" :xl="4" class="header-col">
            <span class="header-card-title">
              <a-tooltip title="节点离线，服务正处于Starting、Running、Stopping，不确定将要转换的状态">
                <a-icon type="question-circle"/>
              </a-tooltip> 未知
            </span><br/>
            <span class="header-card-value">{{ status.unknown }}</span><br/>
          </a-col>
          <a-col :xs="8" :sm="8" :xl="4" class="header-col">
            <span class="header-card-title"><a-icon type="loading"/> 运行中</span><br/>
            <span class="header-card-value" style="color:#1ABB9C;">{{ status.running }}</span>
            <span v-show="status.starting > 0 "><a-icon type="caret-left" />{{ status.starting }}</span>
          </a-col>
          <a-col :xs="8" :sm="8" :xl="4" class="header-col">
            <span class="header-card-title"><a-icon type="stop"/> 已停止</span><br/>
            <span class="header-card-value">{{ status.stopped }}</span>
            <span v-show="status.stopping > 0 "><a-icon type="caret-left" />{{ status.stopping }}</span>
          </a-col>
          <a-col :xs="8" :sm="8" :xl="4" class="header-col">
            <span class="header-card-title"><a-icon type="close"/> 已失败</span><br/>
            <span class="header-card-value" style="color:red">{{ status.exited }}</span><br/>
          </a-col>
        </a-row>

      </div>
    </template>

    <a-list
      :grid="{ gutter: 16, xs: 1, sm: 1, md: 2, lg: 3, xl:3, xxl :4 }"
      :data-source="status.process"
    >
      <a-list-item slot="renderItem" slot-scope="item">
        <template v-if="!item || item.id === undefined">
          <a-button class="new-btn" type="dashed" @click="openEdit(item,'create')">
            <a-icon type="plus"/>
            新增程序
          </a-button>
        </template>
        <template v-else>
          <a-card :title="item.name" size="small" :bordered="false">
            <template slot="extra" >
              <a-icon type="bell" style="color:#2894FF" @click="processBell" />
            </template>

            <a-row>
              <a-col :span="3">节点</a-col>
              <a-col :span="20">{{ item.node }}</a-col>
            </a-row>
            <a-row>
              <a-col :span="3">命令</a-col>
              <a-tooltip v-if="item.command.length > 30" :title="item.command">
                <a-col :span="20" >{{ `${item.command.slice(0, 30)}...` }}</a-col>
              </a-tooltip>
              <a-col :span="20" v-else>{{ item.command }}</a-col>
            </a-row>
            <a-row>
              <a-col :span="3">分组</a-col>
              <a-tooltip v-if="item.groups.length > 30" :title="item.groups">
                <a-col :span="20" >{{ `${item.groups.slice(0, 30)}...` }}</a-col>
              </a-tooltip>
              <a-col :span="20" v-else>{{ item.groups }}</a-col>
            </a-row>
            <a-row>
              <a-col :span="3">状态</a-col>
              <a-col :span="20">
                <template v-if="item.state.status === 'Exited'">
                  <a-popconfirm placement="topRight" style="width:200px">
                    <span slot="title" style="white-space:pre-wrap;">{{ item.state.exit_msg }}</span>
                    <a href="#" class="state_info" style="background:#FF7575;color:white">Exited</a>
                  </a-popconfirm>
                </template>
                <template v-else-if="item.state.status === 'Running'">
                  <span class="state_info" style="background:#1ABB9C;">{{ item.state.status }}</span>
                  <span class="state_desc">Pid:{{ item.state.pid }},Age: {{ item.state.timestamp | showAge }} </span>
                </template>
                <template v-else-if="item.state.status === 'Starting'">
                  <span class="state_info" style="background:#01B468;">{{ item.state.status }}</span>
                </template>
                <template v-else-if="item.state.status === 'Stopping'">
                  <span class="state_info" style="background:#D0D0D0;">{{ item.state.status }}</span>
                </template>
                <template v-else-if="item.state.status === 'Unknown'">
                  <span class="state_info" style="background:#E0E0E0;">{{ item.state.status }}</span>
                  <a-tooltip title="节点离线，服务正处于Starting、Running、Stopping，不确定将要转换的状态">
                    <a-icon type="question-circle-o" />
                  </a-tooltip>
                </template>
                <template v-else>
                  <span class="state_info" style="background:#F0F0F0;">{{ item.state.status }}</span>
                </template>
              </a-col>
            </a-row>
            <a-row>
              <a-col :span="3">CPU</a-col>
              <a-col :span="15">
                <a-progress
                  :stroke-color="progressColor(item.state.cpu)"
                  :percent="parseFloat(item.state.cpu.toFixed(2))" />
              </a-col>
            </a-row>
            <a-row>
              <a-col :span="3">内存</a-col>
              <a-col :span="15">
                <a-progress
                  :stroke-color="progressColor(item.state.mem)"
                  :percent="parseFloat(item.state.mem.toFixed(2))" />
              </a-col>
            </a-row>
            <!-- <a-row type="flex" justify="space-around">
              <a-col :span="4" >
                <a-progress
                  type="dashboard"
                  :width="80"
                  :stroke-color="progressColor(item.state.mem)"
                  :percent="parseFloat(item.state.cpu.toFixed(2))" />
              </a-col>
              <a-col :span="4">
                <a-progress
                  type="dashboard"
                  :width="80"
                  :stroke-color="progressColor(item.state.mem)"
                  :percent="parseFloat(item.state.mem.toFixed(2))" />
              </a-col>
            </a-row> -->
            <template slot="actions" >
              <a v-if="item.state.status === 'Exited' || item.state.status === 'Stopped'" @click="startProcess(item.id)">启动</a>
              <a v-else-if="item.state.status === 'Running'" @click="stopProcess(item.id)">停止</a>
              <a v-else>停止</a>
              <a @click="openEdit(item,'edit')"><span v-if="item.state.status === 'Running'">查看</span><span v-else>修改</span></a>
              <a-popconfirm
                v-if="item.state.status === 'Exited' || item.state.status === 'Stopped'"
                title="确定要删除吗？"
                @confirm="deleteProcess(item.id)">
                <a-icon slot="icon" type="question-circle-o" style="color: red" />
                <a href="#"> 删除</a>
              </a-popconfirm>
              <a-dropdown>
                <a href="javascript:;">
                  <a-icon type="ellipsis" />
                </a>
                <a-menu slot="overlay">
                  <a-menu-item>
                    <a @click="openEdit(item,'copy')">拷贝</a>
                  </a-menu-item>
                </a-menu>
              </a-dropdown>
            </template>
          </a-card>
        </template>
      </a-list-item>
    </a-list>
  </page-header-wrapper>
</template>
<script>
import { groupList, groupRemove, processList, processDelete, processStart, processStop } from '@/api/process'
import moment from 'moment'

export default {
  name: 'ProcessList',
  data () {
    return {
      path: '',
      treeData: [],
      status: {
        alert: 0,
        unknown: 0,
        starting: 0,
        running: 0,
        exited: 0,
        stopping: 0,
        stopped: 0,
        length: 0,
        process: []
      },
      ticker: null
    }
  },
  // beforeMount () {
  //   this.loadTreeData()
  //   this.getProcess()
  // },
  mounted () {
    console.log(this.$route)
    if (this.$route.params.path) {
      this.path = this.$route.params.path
    }
    this.loadTreeData()
    this.getProcess()
    this.ticker = setInterval(() => {
        this.getProcess()
      }, 2000)
  },
  filters: {
    showAge (time) {
      // const age = moment().unix() - time
      // return moment.unix(age).format('YYYY-MM-DD hh:mm:ss')
      return moment.unix(time).fromNow(true)
    }
  },
  destroyed () {
    console.log('destroyed')
    clearInterval(this.ticker)
  },
  methods: {
    loadTreeData (res) {
       groupList()
      .then(res => {
      this.treeData = []
      for (const kk in res) {
        const kkk = kk.split('/')
        let children = this.treeData
        for (const idx in kkk) {
          const name = kkk[idx]
          let has = false
          for (const ni in children) {
            if (children[ni].id === name) {
              has = true
              children = children[ni].children
              break
            }
          }

          if (!has) {
            let key = kkk[0]
            for (let i = 1; i <= idx; i++) {
              key += '/' + kkk[i]
            }
            const newChild = { id: name, title: key, key: key, value: key, children: [] }
            children.push(newChild)
            children = newChild.children
          }
        }
      }
      })
    },
    onChange (value) {
      this.path = value
      this.getProcess()
    },
    getProcess () {
      const args = { group: this.path }
      processList(args)
      .then(res => {
        this.status = { alert: 0, unknown: 0, starting: 0, running: 0, exited: 0, stopped: 0, stopping: 0, length: 0, process: [] }
        this.status.process.push({})
        for (const k in res) {
          const v = res[k]
          this.status.length += 1
          this.status.process.push(v)
          if (v.state.status === 'Unknown') {
            this.status.unknown += 1
          } else if (v.state.status === 'Running') {
            this.status.running += 1
          } else if (v.state.status === 'Starting') {
            this.status.starting += 1
          } else if (v.state.status === 'Stopping') {
            this.status.stopping += 1
          } else if (v.state.status === 'Stopped') {
            this.status.stopped += 1
          } else {
            this.status.exited += 1
          }
        }
        //  console.log(res)
        })
      .catch(e => { })
    },
    startProcess (id) {
      processStart({ id: id }).then(res => {
        this.getProcess()
      })
    },
    stopProcess (id) {
      processStop({ id: id }).then(res => {
         this.getProcess()
      })
    },
    deleteProcess (id) {
      processDelete({ id: id }).then(res => {
         this.getProcess()
      })
    },
    openEdit (item, option) {
      console.log(option, item)
      this.$router.push({ name: 'pedit', params: { option: option, path: this.path, treeData: this.treeData, item: { ...item } } })
    },
    startAllProcess () {
      for (const idx in this.status.process) {
        const item = this.status.process[idx]
        if (item.id !== undefined && (item.state.status === 'Stopped' || item.state.status === 'Exited')) {
          this.startProcess(item.id)
        }
      }
    },
    stopAllProcess () {
      for (const idx in this.status.process) {
        const item = this.status.process[idx]
        if (item.id !== undefined && item.state.status === 'Running') {
          this.stopProcess(item.id)
        }
      }
    },
    deleteGroup () {
      if (this.path !== '') {
        groupRemove({ group: this.path }).then(res => {
         this.path = ''
         this.loadTreeData()
        })
      }
    },
    progressColor (percent) {
      if (percent >= 80) {
        return 'red'
      } else if (percent >= 50) {
        return '#EAC100'
      }
    },
    processBell () {
      this.$message.info('暂未实现该功能')
    }
  }
}
</script>

<style lang="less" scoped>
  .header-card-value{
    font-size:40px;
    font-weight:bold;
    color:#7B7B7B;
  }

  .header-col{
    border-left:2px solid #ADB2B5;
    padding-left:10px;
  }

  .new-btn {
    background-color: #fff;
    border-radius: 2px;
    width: 100%;
    height: 233px;
  }

  .state_info{
  border: 1px solid #F0F0F0;
  border-radius:5px;
  font:14px;
  align:center;
  padding:0 5px;
  margin-right:5px;
}

.state_desc{
  border: 1px solid #F0F0F0;
  border-radius:5px;
  font:14px;
  align:center;
  padding:0 5px;
}

</style>
