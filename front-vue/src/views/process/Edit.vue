<template>
  <page-header-wrapper
    :breadcrumb="{}"
    :title="title"
    @back="goback"
  >
    <a-card :bordered="false">
      <a-form-model
        ref="processEdit"
        :model="form"
        :rules="rules"
        :label-col="labelCol"
        :wrapper-col="wrapperCol"
      >
        <a-form-model-item ref="name" label="名称" prop="name" :wrapper-col="{ span: 6}">
          <a-input v-model="form.name" placeholder="程序名称"/>
        </a-form-model-item>
        <a-form-model-item label="执行目录" :wrapper-col="{ span: 6}">
          <a-input v-model="form.dir" placeholder="非必填，默认为节点启动目录" />
        </a-form-model-item>
        <a-form-model-item
          v-for="(cfg, index) in form.config"
          :key="index"
          :wrapper-col="index === 0 ? wrapperCol:{ span: 14,offset:3}"
          :label="index === 0 ? '启动配置' : ''"
        >
          <a-input
            v-model="cfg.name"
            placeholder="配置名称"
            style="width: 60%; margin-right: 8px"
          />
          <a-input
            v-model="cfg.context"
            placeholder="配置内容"
            style="width: 90%; margin-right: 8px"
            type="textarea"
            :auto-size="{ minRows: 6, maxRows: 16 }"
          />
          <a-icon
            v-if="form.config.length > 0"
            class="dynamic-delete-button"
            type="minus-circle-o"
            @click="removeCfg(cfg)"
          />
        </a-form-model-item>
        <a-form-model-item
          :label="form.config.length === 0 ? '启动配置' : ''"
          :wrapper-col="form.config.length === 0 ? wrapperCol:{ span: 14,offset:3}"
        >
          <a-button type="dashed" style="width:60%" @click="addCfg">
            <a-icon type="plus" /> 添加配置
          </a-button>
        </a-form-model-item>
        <a-form-model-item prop="command">
          <span slot="label">
            启动命令&nbsp;
            <a-tooltip :title="cmd_question_title">
              <a-icon type="question-circle-o" />
            </a-tooltip>
          </span>
          <a-input v-model="form.command" placeholder="启动命令"/>
        </a-form-model-item>
        <a-form-model-item :wrapper-col="{ span: 6}">
          <span slot="label">
            优先级&nbsp;
            <a-tooltip title="子进程启动关闭优先级，优先级低的，最先启动，关闭的时候最后关闭">
              <a-icon type="question-circle-o" />
            </a-tooltip>
          </span>
          <a-input v-model="form.priority"/>
        </a-form-model-item>
        <a-form-model-item :wrapper-col="{ span: 6}">
          <span slot="label">
            启动检测时长&nbsp;
            <a-tooltip title="启动进程一段时间后没有异常退出，就表示进程正常启动了。">
              <a-icon type="question-circle-o" />
            </a-tooltip>
          </span>
          <a-input v-model="form.start_secs"/>
        </a-form-model-item>
        <a-form-model-item :wrapper-col="{ span: 6}">
          <span slot="label">
            停机等待时长&nbsp;
            <a-tooltip title="这个是当我们向子进程发送stop信号后，到系统返回信息所等待的最大时间。超过这个时间会向该子进程发送一个强制kill的信号。">
              <a-icon type="question-circle-o" />
            </a-tooltip>
          </span>
          <a-input v-model="form.stop_wait_secs">
            <span slot="addonAfter" >秒</span>
          </a-input>
        </a-form-model-item>
        <a-form-model-item :wrapper-col="{ span: 6}">
          <span slot="label">
            自动重启次数&nbsp;
            <a-tooltip title="进程状态为 Exited时，自动重启">
              <a-icon type="question-circle-o" />
            </a-tooltip>
          </span>
          <a-input v-model="form.auto_start_times"/>
        </a-form-model-item>
        <a-form-model-item :wrapper-col="{ span: 14}" prop="groups">
          <span slot="label">
            分组&nbsp;
            <a-tooltip title="程序分组，批量管理">
              <a-icon type="question-circle-o" />
            </a-tooltip>
          </span>
          <template v-for="(group, index) in form.groups">
            <a-tooltip v-if="group.length > 10" :key="index" :title="group">
              <a class="group_elem">
                {{ `${group.slice(0, 10)}...` }}&nbsp;&nbsp;<a-icon type="close" @click="() => deleteGroup(group)"></a-icon>
              </a>
            </a-tooltip>
            <a v-else :key="index" class="group_elem">
              {{ group }}&nbsp;&nbsp;<a-icon type="close" @click="() => deleteGroup(group)"></a-icon>
            </a>
          </template>
          <a-input
            v-if="groupInputVisible"
            class="group_elem_in"
            v-model="groupInputValue"
            @blur="addGroup"
            @keyup.enter="addGroup"
            placeholder="新增分组"
          />
          <a-button-group v-else>
            <a-tree-select
              :tree-data="treeData"
              @select="selectGroup"
              value="Select Group "
              show-search
            />
            <a-button @click="()=>{this.groupInputVisible = true}" icon="plus">Add Group</a-button>
          </a-button-group>
        </a-form-model-item>
        <a-form-model-item label="节点" prop="node" :wrapper-col="{ span: 6}">
          <a-select placeholder="选择一个节点" v-model="form.node" >
            <template v-for="(v,i) in nodes" >
              <a-select-option :value="v.name" :key="i">
                {{ v.name }}
              </a-select-option>
            </template>
          </a-select>
        </a-form-model-item>
        <a-form-model-item :wrapper-col="{ span: 14,offset:6}">
          <a-button type="primary" @click="submitForm">
            {{ submitText }}
          </a-button>
          <a-button style="margin-left: 10px" @click="goback">
            取消
          </a-button>
        </a-form-model-item>
      </a-form-model>

    </a-card>
  </page-header-wrapper>
</template>

<script>
import { processCreate, processUpdate } from '@/api/process'
import { nodeList } from '@/api/node'

export default {
  name: 'ProcessEdit',
  data () {
    return {
      labelCol: { span: 3 },
      wrapperCol: { span: 14 },
      cmd_question_title: '命令中存在配置文件时，路径前加上{{path}}，自动填充',
      rules: {
        name: [
          { required: true, message: '请输入名称', trigger: 'blur' },
          { min: 3, message: '名称至少3个字', trigger: 'blur' }
        ],
        command: [
          { required: true, message: '请输入启动命令', trigger: 'blur' }
        ],
        groups: [
          { validator: this.checkGroups, message: '至少选择一个分组', trigger: 'blur' }
        ],
        node: [
          { required: true, message: '选择节点', trigger: 'blur' }
        ]
      },
      form: {
        id: 0,
        name: '',
        dir: '',
        config: [],
        command: '',
        priority: 10,
        start_secs: 3,
        stop_wait_secs: 10,
        auto_start_times: 3,
        node: '',
        groups: []
      },
      title: '新建配置',
      path: '',
      nodes: [],
      treeData: [],
      groupInputVisible: false,
      groupInputValue: '',
      option: '',
      submitText: '创建'
    }
  },
  mounted () {
    this.path = this.$route.params.path
    this.option = this.$route.params.option
    this.treeData = this.$route.params.treeData
    this.loadNodeList()
    if (this.option === 'edit' || this.option === 'copy') {
      this.form = { ...this.$route.params.item }
    }
    if (this.option === 'edit') {
      this.submitText = '修改'
      this.title = '修改配置'
    }
    console.log(this.option, this.$route.params.item, this.form)
  },
  methods: {
    removeCfg (item) {
      console.log(item)
      const index = this.form.config.indexOf(item)
      if (index !== -1) {
        this.form.config.splice(index, 1)
      }
    },
    addCfg () {
      this.form.config.push({
        name: '',
        context: ''
      })
    },
    loadNodeList () {
      const args = { pageNo: 1, pageSize: 1000 }
      nodeList(args)
        .then(res => {
          this.nodes = res.nodeList
        })
    },
    addGroup () {
      if (this.groupInputValue && this.form.groups.indexOf(this.groupInputValue) === -1) {
        console.log(this.groupInputValue)
        this.form.groups = [...this.form.groups, this.groupInputValue]
      }
      console.log(this.form.groups)
      this.groupInputVisible = false
      this.groupInputValue = ''
    },
    deleteGroup (group) {
      const groups = this.form.groups.filter(g => g !== group)
      this.form.groups = groups
      console.log(group, this.form.groups)
    },
    selectGroup (value) {
      if (value && this.form.groups.indexOf(value) === -1) {
        console.log(value)
        this.form.groups = [...this.form.groups, value]
      }
      console.log(this.form.groups)
      this.tagGroupType = 'button'
    },
    checkGroups () {
      if (this.form.groups.length > 0) {
        return true
      }
      return false
    },
    submitForm () {
       this.$refs.processEdit.validate(valid => {
        if (valid) {
          console.log(this.form)
          const args = { ...this.form }
          args.auto_start_times = parseInt(this.form.auto_start_times)
          args.priority = parseInt(this.form.priority)
          args.start_secs = parseInt(this.form.start_secs)
          args.stop_wait_secs = parseInt(this.form.stop_wait_secs)
          if (this.option === 'create' || this.option === 'copy') {
            processCreate(args).then(res => {
              this.goback()
            })
          } else {
            processUpdate(args).then(res => {
              this.goback()
            })
          }
        }
      })
    },
    goback () {
      // this.$router.back()
      this.$router.push({ name: 'plist', params: { path: this.path } })
    }
  }
}
</script>

<style scoped>

.group_elem_in{
  height:32px;
  width: 245px;
}

.group_elem{
  border: 1px solid #BEBEBE;
  border-radius:5px;
  height:32px;
  padding: 5px 10px;
  margin-right:5px;
  font:14px;
  align:center;
}
</style>
