<template>
  <div>
    <a-card :bordered="false">
      <div style="marginBottom: 24px " >
        <a-row
          justify="end"
          type="flex">
          <a-col :span="4">
            <a-button type="primary" icon="plus" @click="openAddSetModal">AddSet</a-button>&nbsp;
            <a-popconfirm title="确定要排空所有Kv吗？" @confirm="drainKv()">
              <a-button
                type="primary"
              >DrainKv</a-button>
            </a-popconfirm>
          </a-col>
        </a-row>
      </div>
      <a-table
        bordered
        :columns="columns"
        :data-source="sets"
        size="small"
        :pagination="false"
      >
        <template slot="markClear" slot-scope="set">
          <a-switch checked-children="T" un-checked-children="F" :checked="set.markClear" @change="handleSetMarkClear(set)" />
        </template>
        <template slot="action" slot-scope="set">
          <div >
            <a @click="openAddNodeModal(set.setID)">AddNode</a>
            <a-divider type="vertical" />
            <a-popconfirm title="确定要删除吗？" @confirm="handleRemSet(set.setID)" >
              <a-icon slot="icon" type="question-circle-o" style="color: red" />
              <a style="color:red;">RemSet</a>
            </a-popconfirm>
          </div>
        </template>

        <a-tabs
          slot="expandedRowRender"
          slot-scope="set"
        >
          <a-tab-pane key="1" tab="KvNode">
            <a-table
              bordered
              :columns="nodeColumns"
              :data-source="set.nodes"
              size="small"
              :pagination="false"
            >
              <template slot="nodeAction" slot-scope="node" >
                <div >
                  <a-select
                    style="width: 100px"
                    value="AddStore"
                    @focus="openAddStoreSelect(set.setID,node.nodeID,set.stores,node.stores)"
                    @select="handleAddStore"
                  >
                    <a-select-option v-for="d in addStoreSelectData" :key="d.value">
                      {{ d.title }}
                    </a-select-option>
                  </a-select>
                </div>
              </template>

              <a-table
                bordered
                slot="expandedRowRender"
                slot-scope="node"
                :columns="nodeStoreColumns"
                :data-source="node.stores"
                size="small"
                :pagination="false"
              >
                <template slot="nodeStoreIsLeader" slot-scope="isLeader">
                  <a-badge v-if="isLeader" status="processing" text="Leader" />
                  <a-badge v-else status="default" text="Follower" />
                </template>
                <template slot="nodeStoreAction" slot-scope="store" >
                  <div >
                    <a-popconfirm title="确定要删除吗？" @confirm="handleRemStore(set.setID,node.nodeID,store.storeID)">
                      <a-icon slot="icon" type="question-circle-o" style="color: red" />
                      <a style="color:red;">RemStore</a>
                    </a-popconfirm>
                  </div>
                </template>
              </a-table>

            </a-table>
          </a-tab-pane>
          <a-tab-pane key="2" tab="Store">
            <a-table
              bordered
              :columns="storeColumns"
              :data-source="set.stores"
              size="small"
              :pagination="false"
            >
              <template slot="storeAction" >
              </template>

              <a-table
                bordered
                slot="expandedRowRender"
                slot-scope="store"
                :columns="storeNodeColumns"
                :data-source="store.nodes"
                size="small"
                :pagination="false"
              >
                <template slot="storeNodeIsLeader" slot-scope="isLeader">
                  <a-badge v-if="isLeader" status="processing" text="Leader" />
                  <a-badge v-else status="default" text="Follower" />
                </template>
                <template slot="storeNodeAction" >

                </template>
              </a-table>

            </a-table>
          </a-tab-pane>
        </a-tabs>
      </a-table>
    </a-card>

    <a-modal
      :visible="addSetModal"
      title="AddSet"
      @ok="handleAddSet"
      @cancel="()=> addSetModal = false"
      :width="1000"
    >
      <a-form-model
        :label-col="{ span: 4 }"
        :wrapper-col="{ span: 18 }">
        <a-form-model-item label="SetID">
          <a-input-number v-model="addSetForm.setID" placeholder="input SetID" style="width:100%"/>
        </a-form-model-item>
        <a-form-model-item label="Nodes">
          <a-row style="background:#F0F0F0">
            <a-col :span="4">&nbsp;&nbsp;NodeID</a-col>
            <a-col :span="10">Host</a-col>
            <a-col :span="4">ServicePort</a-col>
            <a-col :span="4">RaftPort</a-col>
            <a-col :span="2"></a-col>
          </a-row>
          <template v-for="(node, index) in addSetForm.nodes" >
            <a-row :key="index">
              <a-col :span="4">
                <a-input-number v-model="node.nodeID" style="width:90%"/>
              </a-col>
              <a-col :span="10">
                <a-input v-model="node.host" style="width:90%" />
              </a-col>
              <a-col :span="4">
                <a-input-number v-model="node.servicePort" style="width:90%"/>
              </a-col>
              <a-col :span="4">
                <a-input-number v-model="node.raftPort" style="width:90%"/>
              </a-col>
              <a-col :span="2">
                <a-icon
                  type="minus-circle-o"
                  @click="remAddSetFormNodes(index)"
                />
              </a-col>
            </a-row>
          </template>
        </a-form-model-item>
        <a-form-model-item
          :wrapper-col="{ span: 18,offset:4}"
        >
          <a-button type="dashed" style="width:90%" @click="addAddSetFormNodes">
            <a-icon type="plus" /> 添加一行
          </a-button>
        </a-form-model-item>
      </a-form-model>
    </a-modal>
    <a-modal
      :visible="addNodeModal"
      title="AddNode"
      @ok="handleAddNode"
      @cancel="()=> addNodeModal = false"
      :width="800"
    >
      <a-form-model
        :label-col="{ span: 4 }"
        :wrapper-col="{ span: 18 }">
        <a-form-model-item label="SetID">
          <a-input-number v-model="addNodeForm.setID" style="width:100%" disabled/>
        </a-form-model-item>
        <a-form-model-item label="NodeID">
          <a-input-number v-model="addNodeForm.nodeID" style="width:100%"/>
        </a-form-model-item>
        <a-form-model-item label="Host">
          <a-input v-model="addNodeForm.host" style="width:100%"/>
        </a-form-model-item>
        <a-form-model-item label="ServicePort">
          <a-input-number v-model="addNodeForm.servicePort" style="width:100%"/>
        </a-form-model-item>
        <a-form-model-item label="RaftPort">
          <a-input-number v-model="addNodeForm.raftPort" style="width:100%"/>
        </a-form-model-item>
      </a-form-model>
    </a-modal>
  </div>
</template>

<script>
import { getSetStatus,
  setMarkClear,
  addSet,
  remSet,
  addNode,
  remNode,
  addLeaderStoreToNode,
  removeNodeStore, drainKv
} from '@/api/flyfish'
export default {
  name: 'FlyfishSetStatus',
   props: { host: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      columns: [
        { title: 'SetID', dataIndex: 'setID' },
        { title: 'KvCount', dataIndex: 'kvcount' },
        { title: 'MarkClear', scopedSlots: { customRender: 'markClear' } },
        { title: 'NodeCount', customRender: (record) => record.nodes.length },
        { title: 'StoreCount', customRender: (record) => record.stores.length },
        { title: 'Action', scopedSlots: { customRender: 'action' } }
      ],
      nodeColumns: [
        { title: 'NodeID', dataIndex: 'nodeID' },
        { title: 'StoreCount', customRender: (record) => record.stores.length },
        { title: 'Action', scopedSlots: { customRender: 'nodeAction' } }
      ],
      nodeStoreColumns: [
        { title: 'StoreID', dataIndex: 'storeID' },
        { title: 'Type', dataIndex: 'type' },
        { title: 'Value', dataIndex: 'value' },
        { title: 'Progress', dataIndex: 'progress' },
        { title: 'IsLeader', dataIndex: 'isLeader', scopedSlots: { customRender: 'nodeStoreIsLeader' } },
        { title: 'Action', scopedSlots: { customRender: 'nodeStoreAction' } }
      ],
      storeColumns: [
        { title: 'StoreID', dataIndex: 'storeID' },
        { title: 'KvCount', dataIndex: 'kvcount' },
        { title: 'NodeCount', customRender: (record) => record.nodes.length },
        { title: 'Action', scopedSlots: { customRender: 'storeAction' } }
      ],
      storeNodeColumns: [
        { title: 'NodeID', dataIndex: 'nodeID' },
        { title: 'IsLeader', dataIndex: 'isLeader', scopedSlots: { customRender: 'storeNodeIsLeader' } },
        { title: 'Action', scopedSlots: { customRender: 'storeNodeAction' } }
      ],
      sets: [],
      addSetModal: false,
      addSetForm: {
        setID: 0,
        nodes: []
      },
      addNodeModal: false,
      addNodeForm: {
        setID: 0,
        nodeID: 0,
        host: '',
        servicePort: 0,
        raftPort: 0
      },
      addStoreSelectData: [],
      addStoreForm: {
        setID: 0,
        nodeID: 0,
        store: 0
      }
    }
  },
  mounted () {
    this.getSetStatus()
  },
  watch: {
    host (val) {
      this.getSetStatus()
    }
  },
  methods: {
    getSetStatus () {
    if (this.host === '') {
      return ''
    }
    getSetStatus({ host: this.host })
      .then(res => {
        this.sets = res.sets
        console.log(this.sets)
        this.sets.sort((a, b) => (a.setID - b.setID))
        for (let i = 0; i < this.sets.length; i++) {
          const set = this.sets[i]
          set.stores.sort((a, b) => (a.storeID - b.storeID))
          for (let si = 0; si < set.stores.length; si++) {
            set.stores[si].nodes = []
          }
          set.nodes.sort((a, b) => (a.nodeID - b.nodeID))
          for (let j = 0; j < set.nodes.length; j++) {
            const node = set.nodes[j]
            if (!node.stores) {
              node.stores = []
            }
            node.stores.sort((a, b) => (a.storeID - b.storeID))
            for (let k = 0; k < node.stores.length; k++) {
              const store = node.stores[k]
              for (let si = 0; si < set.stores.length; si++) {
                if (set.stores[si].storeID === store.storeID) {
                  set.stores[si].nodes.push({
                    nodeID: node.nodeID,
                    isLeader: store.isLeader
                  })
                  break
                }
              }
            }
          }
        }
        console.log(this.sets)
      })
      .catch(e => {
        this.sets = []
      })
    },
    handleSetMarkClear (set) {
      if (this.host === '' || set.markClear) {
        this.$message.info('不允许操作', 1)
        return ''
      }
      console.log(set)
      setMarkClear({ host: this.host, setMarkClear: { setID: set.setID } })
        .then(res => {
            set.markClear = true
        })
    },
    openAddSetModal () {
      this.addSetForm = { setID: 0, nodes: [] }
      this.addSetModal = true
    },
    addAddSetFormNodes () {
      this.addSetForm.nodes.push({
        nodeID: 0,
        host: '',
        servicePort: 0,
        raftPort: 0
      })
    },
    remAddSetFormNodes (index) {
      this.addSetForm.nodes.splice(index, 1)
    },
    handleAddSet () {
      console.log('addSet', this.addSetForm)
      if (this.host === '') {
        this.$message.info('请选择Pd节点', 1)
        return ''
      }

      addSet({ host: this.host, addSet: { set: this.addSetForm } })
        .then(res => {
          this.addSetModal = false
          this.getSetStatus()
        })
        .catch(e => {
          this.addSetModal = false
        })
    },
    handleRemSet (setID) {
      console.log('remSet', setID)
      if (this.host === '') {
        this.$message.info('请选择Pd节点', 1)
        return ''
      }
      remSet({ host: this.host, remSet: { setID: setID } })
        .then(res => {
          this.getSetStatus()
        })
    },
    openAddNodeModal (setID) {
      this.addNodeForm = { setID: setID, nodeID: 0, host: '', servicePort: 0, raftPort: 0 }
      this.addNodeModal = true
    },
    handleAddNode () {
      console.log('addNode', this.addNodeForm)
      if (this.host === '') {
        this.$message.info('请选择Pd节点', 1)
        return ''
      }
      addNode({ host: this.host, addNode: this.addNodeForm })
        .then(res => {
          this.addNodeModal = false
          this.getSetStatus()
        })
        .catch(e => {
          this.addNodeModal = false
        })
    },
    handleRemNode (setID, nodeID) {
      console.log('remnode', setID, nodeID)
      if (this.host === '') {
        this.$message.info('请选择Pd节点', 1)
        return ''
      }
      remNode({ host: this.host, remNode: { setID: setID, nodeID: nodeID } })
        .then(res => {
          this.getSetStatus()
        })
    },
    openAddStoreSelect (setID, nodeID, allStores, hasStores) {
      this.addStoreForm.setID = setID
      this.addStoreForm.nodeID = nodeID
      this.addStoreSelectData = []
      for (let i = 0; i < allStores.length; i++) {
        let has = false
        for (let j = 0; j < hasStores.length; j++) {
          if (allStores[i].storeID === hasStores[j].storeID) {
            has = true
            break
          }
        }
        if (!has) {
          this.addStoreSelectData.push({
            title: 'Store ' + allStores[i].storeID,
            value: allStores[i].storeID
          })
        }
      }
    },
    handleAddStore (value) {
      this.addStoreForm.store = value
      console.log(this.addStoreForm)
      if (this.host === '') {
        this.$message.info('请选择Pd节点', 1)
        return ''
      }
      addLeaderStoreToNode({ host: this.host, addLearnerStoreToNode: this.addStoreForm })
        .then(res => {
          this.getSetStatus()
        })
    },
    handleRemStore (setID, nodeID, storeID) {
      const arg = { setID: setID, nodeID: nodeID, store: storeID }
      console.log(arg)
      if (this.host === '') {
        this.$message.info('请选择Pd节点', 1)
        return ''
      }

      removeNodeStore({ host: this.host, removeNodeStore: arg })
        .then(res => {
          this.getSetStatus()
        })
    },
    drainKv () {
      drainKv({ host: this.host })
          .then(res => {
            this.$message.info('操作成功')
            setTimeout(() => {
              this.getSetStatus()
            }, 1500)
          })
      }
  }
}
</script>

<style lang="less" scoped>

</style>
