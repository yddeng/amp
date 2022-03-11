<template>
  <div>
    <a-card :bordered="false">
      <div style="marginBottom: 24px " >
        <a-row
          justify="space-between"
          type="flex">
          <a-col :span="4"><span style="font-size:20px;font-weight:bold">Vsersion: {{ version }}</span></a-col>
          <a-col :span="4">
            <a-button
              type="primary"
              @click="openModel('')"
              icon="plus"
            >AddTable</a-button>&nbsp;
            <a-popconfirm title="确定要清除所有表吗？" @confirm="clearDBData('')">
              <a-button
                type="primary"
              >ClearDBData</a-button>
            </a-popconfirm>
          </a-col>
        </a-row>
      </div>
      <a-table
        bordered
        :columns="columns"
        :data-source="data"
        size="small"
        @change="onPageChange">
        <template slot="action" slot-scope="text, record">
          <div >
            <a @click="openModel(record.name)">AddField</a>
            <a-divider type="vertical" />
            <a-popconfirm title="确定要清除当前表吗？" @confirm="clearDBData(record.name)">
              <a>ClearDBData</a>
            </a-popconfirm>
            <a-divider type="vertical" />
            <a-popconfirm title="确定要删除吗？">
              <a style="color:red;">Delete</a>
            </a-popconfirm>
          </div>
        </template>
        <a-table
          bordered
          size="small"
          slot="expandedRowRender"
          slot-scope="record,index"
          :columns="innerColumns"
          :data-source="innerData[currentPageIndex*pageSize +index]"
          :pagination="false"
        >
          <template slot="innerAction" >
            <div >
              <a-popconfirm title="确定要删除吗？">
                <a-icon slot="icon" type="question-circle-o" style="color: red" />
                <a style="color:red;">Delete</a>
              </a-popconfirm>
            </div>
          </template>
        </a-table>
      </a-table>
    </a-card>

    <a-modal
      v-model="visible"
      :title="modelTitle"
      ok-text="确认"
      cancel-text="取消"
      @ok="handleOk"
      :width="800"
    >

      <a-form-model
        :label-col="{ span: 4 }"
        :wrapper-col="{ span: 18 }">
        <a-form-model-item label="TableName">
          <a-input v-model="form.tableName" placeholder="input tableName" :disabled="disabled"/>
        </a-form-model-item>
        <a-form-model-item label="Fields">
          <a-row style="background:#F0F0F0">
            <a-col :span="10">&nbsp;&nbsp;FieldName</a-col>
            <a-col :span="4">Type</a-col>
            <a-col :span="8">DefaultValue</a-col>
            <a-col :span="2"></a-col>
          </a-row>
          <template v-for="(field, index) in form.fields" >
            <a-row :key="index">
              <a-col :span="10">
                <a-input v-model="field.name" style="width:90%"/>
              </a-col>
              <a-col :span="4">
                <a-input v-model="field.type" style="width:90%" />
              </a-col>
              <a-col :span="8">
                <a-input v-model="field.default" style="width:90%"/>
              </a-col>
              <a-col :span="2">
                <a-icon
                  type="minus-circle-o"
                  @click="removeFieldCol(index)"
                />
              </a-col>
            </a-row>
          </template>
        </a-form-model-item>
        <a-form-model-item
          :wrapper-col="{ span: 18,offset:4}"
        >
          <a-button type="dashed" style="width:90%" @click="addFieldCol">
            <a-icon type="plus" /> 添加一行
          </a-button>
        </a-form-model-item>
      </a-form-model>
    </a-modal>
  </div>
</template>

<script>
import { getMeta, addTable, addField, clearDBData } from '@/api/flyfish'

export default {
  name: 'FlyfishMeta',
  props: { host: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      currentPageIndex: 0,
      pageSize: 10,
      data: [],
      columns: [
        { title: 'TableName', dataIndex: 'name' },
        { title: 'Version', dataIndex: 'version' },
        { title: 'FieldCount', dataIndex: 'fields' },
        { title: 'Action', scopedSlots: { customRender: 'action' } }
      ],
      innerData: [],
      innerColumns: [
        { title: 'FieldName', dataIndex: 'name' },
        { title: 'Version', dataIndex: 'version' },
        { title: 'Type', dataIndex: 'type' },
        { title: 'DefaultValue', dataIndex: 'defaultValue' },
        { title: 'Action', scopedSlots: { customRender: 'innerAction' } }
      ],
      version: 0,
      visible: false,
      modelTitle: 'AddTable',
      disabled: false,
      form: {
        tableName: '',
        fields: []
      }
    }
  },
  mounted () {
    this.getMeta()
  },
  watch: {
    host (val) {
      this.getMeta()
    }
  },
  methods: {
    onPageChange (page) {
      this.currentPageIndex = page.current - 1
      this.pageSize = page.pageSize
    },
    getMeta () {
      if (this.host === '') {
        return ''
      }
      getMeta({ host: this.host })
        .then(res => {
          console.log(res)
          const mObj = JSON.parse(res.meta)
          this.version = mObj.Version

          const data = []
          const innerData = []
          for (let i = 0; i < mObj.TableDefs.length; i++) {
              const table = mObj.TableDefs[i]
              data.push({
                key: i,
                name: table.Name,
                version: table.Version,
                fields: table.Fields.length
              })
              const fields = []
              for (let j = 0; j < table.Fields.length; j++) {
                const field = table.Fields[j]
                fields.push({
                  key: j,
                  name: field.Name,
                  version: field.TabVersion,
                  type: field.Type,
                  defaultValue: field.DefaultValue
                })
              }
              innerData.push(fields)
          }
          this.data = data
          this.innerData = innerData
        })
        .catch(e => {
          this.version = 0
          this.data = []
          this.innerData = []
          })
    },
    addFieldCol () {
      this.form.fields.push({
        name: '',
        type: '',
        default: ''
      })
    },
    removeFieldCol (index) {
      this.form.fields.splice(index, 1)
    },
    openModel (tableName) {
      console.log('openModel', tableName)
      this.visible = true
      this.form = { tableName: '', fields: [] }
      if (tableName === '') {
        this.modelTitle = 'AddTable'
        this.disabled = false
      } else {
        this.form.tableName = tableName
        this.modelTitle = 'AddField'
        this.disabled = true
      }
    },
    handleOk () {
      if (this.version === 0) {
        return ''
      }
      const args = {
        name: this.form.tableName,
        fields: this.form.fields,
        version: this.version,
        host: this.host
      }
      console.log(args)
        if (this.modelTitle === 'AddTable') {
          addTable(args)
            .then(res => {
            this.visible = false
            this.getMeta()
          })
        } else {
          addField(args)
            .then(res => {
            this.visible = false
            this.getMeta()
          })
        }
    },
    clearDBData (name) {
      var tables = []
      if (name !== '') {
        tables.push(name)
      }
      clearDBData({ host: this.host, clearDBData: { tables: tables } })
        .then(res => {
          this.$message.info('操作成功')
          this.getMeta()
        })
    }
  }

}
</script>

<style lang="less" scoped>

</style>
