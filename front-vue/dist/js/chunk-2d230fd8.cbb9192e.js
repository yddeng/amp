(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-2d230fd8"],{ef45:function(t,e,a){"use strict";a.r(e);var o=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("page-header-wrapper",{attrs:{title:"Flyfish",breadcrumb:{},"tab-list":t.tabList,"tab-active-key":t.tabActiveKey,"tab-change":t.handleTabChange}},[a("template",{slot:"content"},[a("div",{staticStyle:{height:"40px",width:"540px",margin:"0 auto"}},[a("a-select",{staticStyle:{width:"300px"},attrs:{"show-search":"",placeholder:"选择Pd地址,加载数据"},on:{change:t.handleChange}},t._l(t.data,(function(e){return a("a-select-option",{key:e.title,attrs:{value:e.value}},[t._v(" "+t._s(e.field)+" ")])})),1),a("a-popconfirm",{attrs:{title:"确定要删除当前Pd地址吗？"},on:{confirm:t.deleteKv}},[a("a-button",{staticStyle:{color:"white"},attrs:{disabled:t.disabled,icon:"delete",type:"danger"}})],1),t._v("   "),a("a-button",{attrs:{type:"primary",icon:"plus"},on:{click:t.openModel}},[t._v("新增Pd地址")])],1)]),"meta"===t.tabActiveKey?a("a-meta",{attrs:{host:t.selectValue}}):"setStatus"===t.tabActiveKey?a("a-set-status",{attrs:{host:t.selectValue}}):t._e(),a("a-modal",{attrs:{title:"新增","ok-text":"确认","cancel-text":"取消"},on:{ok:t.handleOk},model:{value:t.visible,callback:function(e){t.visible=e},expression:"visible"}},[a("a-form-model",{attrs:{"label-col":{span:4},"wrapper-col":{span:18}}},[a("a-form-model-item",{attrs:{label:"名称"}},[a("a-input",{attrs:{placeholder:"input placeholder"},model:{value:t.inputName,callback:function(e){t.inputName=e},expression:"inputName"}})],1),a("a-form-model-item",{attrs:{label:"地址"}},[a("a-input",{attrs:{placeholder:"input placeholder"},model:{value:t.inputAddress,callback:function(e){t.inputAddress=e},expression:"inputAddress"}})],1)],1)],1)],2)},s=[],n=(a("ac1f"),a("1276"),a("a434"),a("b775")),r={kvSet:"kv/set",kvGet:"kv/get",kvDelete:"kv/delete"};function d(t){return Object(n["b"])({url:r.kvSet,method:"post",data:t})}function i(t){return Object(n["b"])({url:r.kvGet,method:"post",data:t})}function l(t){return Object(n["b"])({url:r.kvDelete,method:"post",data:t})}var c=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[a("a-card",{attrs:{bordered:!1}},[a("div",{staticStyle:{marginBottom:"24px"}},[a("a-row",{attrs:{justify:"space-between",type:"flex"}},[a("a-col",{attrs:{span:4}},[a("span",{staticStyle:{"font-size":"20px","font-weight":"bold"}},[t._v("Vsersion: "+t._s(t.version))])]),a("a-col",{attrs:{span:2}},[a("a-button",{attrs:{type:"primary",icon:"plus"},on:{click:function(e){return t.openModel("")}}},[t._v("AddTable")])],1)],1)],1),a("a-table",{attrs:{bordered:"",columns:t.columns,"data-source":t.data,size:"small"},on:{change:t.onPageChange},scopedSlots:t._u([{key:"action",fn:function(e,o){return[a("div",[a("a",{on:{click:function(e){return t.openModel(o.name)}}},[t._v("AddField")]),a("a-divider",{attrs:{type:"vertical"}}),a("a-popconfirm",{attrs:{title:"确定要删除吗？"}},[a("a-icon",{staticStyle:{color:"red"},attrs:{slot:"icon",type:"question-circle-o"},slot:"icon"}),a("a",{staticStyle:{color:"red"}},[t._v("Delete")])],1)],1)]}},{key:"expandedRowRender",fn:function(e,o){return a("a-table",{attrs:{bordered:"",size:"small",columns:t.innerColumns,"data-source":t.innerData[t.currentPageIndex*t.pageSize+o],pagination:!1}},[a("template",{slot:"innerAction"},[a("div",[a("a-popconfirm",{attrs:{title:"确定要删除吗？"}},[a("a-icon",{staticStyle:{color:"red"},attrs:{slot:"icon",type:"question-circle-o"},slot:"icon"}),a("a",{staticStyle:{color:"red"}},[t._v("Delete")])],1)],1)])],2)}}])})],1),a("a-modal",{attrs:{title:t.modelTitle,"ok-text":"确认","cancel-text":"取消",width:800},on:{ok:t.handleOk},model:{value:t.visible,callback:function(e){t.visible=e},expression:"visible"}},[a("a-form-model",{attrs:{"label-col":{span:4},"wrapper-col":{span:18}}},[a("a-form-model-item",{attrs:{label:"TableName"}},[a("a-input",{attrs:{placeholder:"input tableName",disabled:t.disabled},model:{value:t.form.tableName,callback:function(e){t.$set(t.form,"tableName",e)},expression:"form.tableName"}})],1),a("a-form-model-item",{attrs:{label:"Fields"}},[a("a-row",{staticStyle:{background:"#F0F0F0"}},[a("a-col",{attrs:{span:10}},[t._v("  FieldName")]),a("a-col",{attrs:{span:4}},[t._v("Type")]),a("a-col",{attrs:{span:8}},[t._v("DefaultValue")]),a("a-col",{attrs:{span:2}})],1),t._l(t.form.fields,(function(e,o){return[a("a-row",{key:o},[a("a-col",{attrs:{span:10}},[a("a-input",{staticStyle:{width:"90%"},model:{value:e.name,callback:function(a){t.$set(e,"name",a)},expression:"field.name"}})],1),a("a-col",{attrs:{span:4}},[a("a-input",{staticStyle:{width:"90%"},model:{value:e.type,callback:function(a){t.$set(e,"type",a)},expression:"field.type"}})],1),a("a-col",{attrs:{span:8}},[a("a-input",{staticStyle:{width:"90%"},model:{value:e.default,callback:function(a){t.$set(e,"default",a)},expression:"field.default"}})],1),a("a-col",{attrs:{span:2}},[a("a-icon",{attrs:{type:"minus-circle-o"},on:{click:function(e){return t.removeFieldCol(o)}}})],1)],1)]}))],2),a("a-form-model-item",{attrs:{"wrapper-col":{span:18,offset:4}}},[a("a-button",{staticStyle:{width:"90%"},attrs:{type:"dashed"},on:{click:t.addFieldCol}},[a("a-icon",{attrs:{type:"plus"}}),t._v(" 添加一行 ")],1)],1)],1)],1)],1)},u=[],h={getMeta:"flyfish/getMeta",addTable:"flyfish/addTable",addField:"flyfish/addField",getSetStatus:"flyfish/getSetStatus"};function f(t){return Object(n["b"])({url:h.getMeta,method:"post",data:t})}function m(t){return Object(n["b"])({url:h.addTable,method:"post",data:t})}function p(t){return Object(n["b"])({url:h.addField,method:"post",data:t})}function b(t){return Object(n["b"])({url:h.getSetStatus,method:"post",data:t})}function v(t){return Object(n["b"])({url:"/flyfish/setMarkClear",method:"post",data:t})}function S(t){return Object(n["b"])({url:"/flyfish/addSet",method:"post",data:t})}function y(t){return Object(n["b"])({url:"/flyfish/remSet",method:"post",data:t})}function k(t){return Object(n["b"])({url:"/flyfish/addNode",method:"post",data:t})}function g(t){return Object(n["b"])({url:"/flyfish/remNode",method:"post",data:t})}function I(t){return Object(n["b"])({url:"/flyfish/addLeaderStoreToNode",method:"post",data:t})}function D(t){return Object(n["b"])({url:"/flyfish/removeNodeStore",method:"post",data:t})}var N={name:"FlyfishMeta",props:{host:{type:String,default:""}},data:function(){return{currentPageIndex:0,pageSize:10,data:[],columns:[{title:"TableName",dataIndex:"name"},{title:"Version",dataIndex:"version"},{title:"FieldCount",dataIndex:"fields"},{title:"Action",scopedSlots:{customRender:"action"}}],innerData:[],innerColumns:[{title:"FieldName",dataIndex:"name"},{title:"Version",dataIndex:"version"},{title:"Type",dataIndex:"type"},{title:"DefaultValue",dataIndex:"defaultValue"},{title:"Action",scopedSlots:{customRender:"innerAction"}}],version:0,visible:!1,modelTitle:"AddTable",disabled:!1,form:{tableName:"",fields:[]}}},mounted:function(){this.getMeta()},watch:{host:function(t){this.getMeta()}},methods:{onPageChange:function(t){this.currentPageIndex=t.current-1,this.pageSize=t.pageSize},getMeta:function(){var t=this;if(""===this.host)return"";f({host:this.host}).then((function(e){var a=JSON.parse(e.meta);t.version=a.Version;for(var o=[],s=[],n=0;n<a.TableDefs.length;n++){var r=a.TableDefs[n];o.push({key:n,name:r.Name,version:r.Version,fields:r.Fields.length});for(var d=[],i=0;i<r.Fields.length;i++){var l=r.Fields[i];d.push({key:i,name:l.Name,version:l.TabVersion,type:l.Type,defaultValue:l.DefaultValue})}s.push(d)}t.data=o,t.innerData=s})).catch((function(e){t.version=0,t.data=[],t.innerData=[]}))},addFieldCol:function(){this.form.fields.push({name:"",type:"",default:""})},removeFieldCol:function(t){this.form.fields.splice(t,1)},openModel:function(t){this.visible=!0,this.form={tableName:"",fields:[]},""===t?(this.modelTitle="AddTable",this.disabled=!1):(this.form.tableName=t,this.modelTitle="AddField",this.disabled=!0)},handleOk:function(){var t=this;if(0===this.version)return"";var e={name:this.form.tableName,fields:this.form.fields,version:this.version,host:this.host};"AddTable"===this.modelTitle?m(e).then((function(e){t.visible=!1,t.getMeta()})):p(e).then((function(e){t.visible=!1,t.getMeta()}))}}},F=N,x=a("2877"),A=Object(x["a"])(F,c,u,!1,null,"01e9f15a",null),w=A.exports,_=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[a("a-card",{attrs:{bordered:!1}},[a("div",{staticStyle:{marginBottom:"24px"}},[a("a-button",{attrs:{type:"primary",icon:"plus"},on:{click:t.openAddSetModal}},[t._v("AddSet")])],1),a("a-table",{attrs:{bordered:"",columns:t.columns,"data-source":t.sets,size:"small",pagination:!1},scopedSlots:t._u([{key:"markClear",fn:function(e){return[a("a-switch",{attrs:{"checked-children":"T","un-checked-children":"F",checked:e.markClear},on:{change:function(a){return t.handleSetMarkClear(e)}}})]}},{key:"action",fn:function(e){return[a("div",[a("a",{on:{click:function(a){return t.openAddNodeModal(e.setID)}}},[t._v("AddNode")]),a("a-divider",{attrs:{type:"vertical"}}),a("a-popconfirm",{attrs:{title:"确定要删除吗？"},on:{confirm:function(a){return t.handleRemSet(e.setID)}}},[a("a-icon",{staticStyle:{color:"red"},attrs:{slot:"icon",type:"question-circle-o"},slot:"icon"}),a("a",{staticStyle:{color:"red"}},[t._v("RemSet")])],1)],1)]}},{key:"expandedRowRender",fn:function(e){return a("a-tabs",{},[a("a-tab-pane",{key:"1",attrs:{tab:"KvNode"}},[a("a-table",{attrs:{bordered:"",columns:t.nodeColumns,"data-source":e.nodes,size:"small",pagination:!1},scopedSlots:t._u([{key:"nodeAction",fn:function(o){return[a("div",[a("a-select",{staticStyle:{width:"100px"},attrs:{value:"AddStore"},on:{focus:function(a){return t.openAddStoreSelect(e.setID,o.nodeID,e.stores,o.stores)},select:t.handleAddStore}},t._l(t.addStoreSelectData,(function(e){return a("a-select-option",{key:e.value},[t._v(" "+t._s(e.title)+" ")])})),1)],1)]}},{key:"expandedRowRender",fn:function(o){return a("a-table",{attrs:{bordered:"",columns:t.nodeStoreColumns,"data-source":o.stores,size:"small",pagination:!1},scopedSlots:t._u([{key:"nodeStoreIsLeader",fn:function(t){return[a("a-badge",t?{attrs:{status:"processing",text:"Leader"}}:{attrs:{status:"default",text:"Follower"}})]}},{key:"nodeStoreAction",fn:function(s){return[a("div",[a("a-popconfirm",{attrs:{title:"确定要删除吗？"},on:{confirm:function(a){return t.handleRemStore(e.setID,o.nodeID,s.storeID)}}},[a("a-icon",{staticStyle:{color:"red"},attrs:{slot:"icon",type:"question-circle-o"},slot:"icon"}),a("a",{staticStyle:{color:"red"}},[t._v("RemStore")])],1)],1)]}}],null,!0)})}}],null,!0)})],1),a("a-tab-pane",{key:"2",attrs:{tab:"Store"}},[a("a-table",{attrs:{bordered:"",columns:t.storeColumns,"data-source":e.stores,size:"small",pagination:!1},scopedSlots:t._u([{key:"expandedRowRender",fn:function(e){return a("a-table",{attrs:{bordered:"",columns:t.storeNodeColumns,"data-source":e.nodes,size:"small",pagination:!1},scopedSlots:t._u([{key:"storeNodeIsLeader",fn:function(t){return[a("a-badge",t?{attrs:{status:"processing",text:"Leader"}}:{attrs:{status:"default",text:"Follower"}})]}}],null,!0)},[a("template",{slot:"storeNodeAction"})],2)}}],null,!0)},[a("template",{slot:"storeAction"})],2)],1)],1)}}])})],1),a("a-modal",{attrs:{visible:t.addSetModal,title:"AddSet",width:1e3},on:{ok:t.handleAddSet,cancel:function(){return t.addSetModal=!1}}},[a("a-form-model",{attrs:{"label-col":{span:4},"wrapper-col":{span:18}}},[a("a-form-model-item",{attrs:{label:"SetID"}},[a("a-input-number",{staticStyle:{width:"100%"},attrs:{placeholder:"input SetID"},model:{value:t.addSetForm.setID,callback:function(e){t.$set(t.addSetForm,"setID",e)},expression:"addSetForm.setID"}})],1),a("a-form-model-item",{attrs:{label:"Nodes"}},[a("a-row",{staticStyle:{background:"#F0F0F0"}},[a("a-col",{attrs:{span:4}},[t._v("  NodeID")]),a("a-col",{attrs:{span:10}},[t._v("Host")]),a("a-col",{attrs:{span:4}},[t._v("ServicePort")]),a("a-col",{attrs:{span:4}},[t._v("RaftPort")]),a("a-col",{attrs:{span:2}})],1),t._l(t.addSetForm.nodes,(function(e,o){return[a("a-row",{key:o},[a("a-col",{attrs:{span:4}},[a("a-input-number",{staticStyle:{width:"90%"},model:{value:e.nodeID,callback:function(a){t.$set(e,"nodeID",a)},expression:"node.nodeID"}})],1),a("a-col",{attrs:{span:10}},[a("a-input",{staticStyle:{width:"90%"},model:{value:e.host,callback:function(a){t.$set(e,"host",a)},expression:"node.host"}})],1),a("a-col",{attrs:{span:4}},[a("a-input-number",{staticStyle:{width:"90%"},model:{value:e.servicePort,callback:function(a){t.$set(e,"servicePort",a)},expression:"node.servicePort"}})],1),a("a-col",{attrs:{span:4}},[a("a-input-number",{staticStyle:{width:"90%"},model:{value:e.raftPort,callback:function(a){t.$set(e,"raftPort",a)},expression:"node.raftPort"}})],1),a("a-col",{attrs:{span:2}},[a("a-icon",{attrs:{type:"minus-circle-o"},on:{click:function(e){return t.remAddSetFormNodes(o)}}})],1)],1)]}))],2),a("a-form-model-item",{attrs:{"wrapper-col":{span:18,offset:4}}},[a("a-button",{staticStyle:{width:"90%"},attrs:{type:"dashed"},on:{click:t.addAddSetFormNodes}},[a("a-icon",{attrs:{type:"plus"}}),t._v(" 添加一行 ")],1)],1)],1)],1),a("a-modal",{attrs:{visible:t.addNodeModal,title:"AddNode",width:800},on:{ok:t.handleAddNode,cancel:function(){return t.addNodeModal=!1}}},[a("a-form-model",{attrs:{"label-col":{span:4},"wrapper-col":{span:18}}},[a("a-form-model-item",{attrs:{label:"SetID"}},[a("a-input-number",{staticStyle:{width:"100%"},attrs:{disabled:""},model:{value:t.addNodeForm.setID,callback:function(e){t.$set(t.addNodeForm,"setID",e)},expression:"addNodeForm.setID"}})],1),a("a-form-model-item",{attrs:{label:"NodeID"}},[a("a-input-number",{staticStyle:{width:"100%"},model:{value:t.addNodeForm.nodeID,callback:function(e){t.$set(t.addNodeForm,"nodeID",e)},expression:"addNodeForm.nodeID"}})],1),a("a-form-model-item",{attrs:{label:"Host"}},[a("a-input",{staticStyle:{width:"100%"},model:{value:t.addNodeForm.host,callback:function(e){t.$set(t.addNodeForm,"host",e)},expression:"addNodeForm.host"}})],1),a("a-form-model-item",{attrs:{label:"ServicePort"}},[a("a-input-number",{staticStyle:{width:"100%"},model:{value:t.addNodeForm.servicePort,callback:function(e){t.$set(t.addNodeForm,"servicePort",e)},expression:"addNodeForm.servicePort"}})],1),a("a-form-model-item",{attrs:{label:"RaftPort"}},[a("a-input-number",{staticStyle:{width:"100%"},model:{value:t.addNodeForm.raftPort,callback:function(e){t.$set(t.addNodeForm,"raftPort",e)},expression:"addNodeForm.raftPort"}})],1)],1)],1)],1)},M=[],P=(a("4e82"),{name:"FlyfishSetStatus",props:{host:{type:String,default:""}},data:function(){return{columns:[{title:"SetID",dataIndex:"setID"},{title:"MarkClear",scopedSlots:{customRender:"markClear"}},{title:"NodeCount",customRender:function(t){return t.nodes.length}},{title:"StoreCount",customRender:function(t){return t.stores.length}},{title:"Action",scopedSlots:{customRender:"action"}}],nodeColumns:[{title:"NodeID",dataIndex:"nodeID"},{title:"StoreCount",customRender:function(t){return t.stores.length}},{title:"Action",scopedSlots:{customRender:"nodeAction"}}],nodeStoreColumns:[{title:"StoreID",dataIndex:"storeID"},{title:"Type",dataIndex:"type"},{title:"Value",dataIndex:"value"},{title:"IsLeader",dataIndex:"isLeader",scopedSlots:{customRender:"nodeStoreIsLeader"}},{title:"Action",scopedSlots:{customRender:"nodeStoreAction"}}],storeColumns:[{title:"StoreID",dataIndex:"storeID"},{title:"NodeCount",customRender:function(t){return t.nodes.length}},{title:"Action",scopedSlots:{customRender:"storeAction"}}],storeNodeColumns:[{title:"NodeID",dataIndex:"nodeID"},{title:"IsLeader",dataIndex:"isLeader",scopedSlots:{customRender:"storeNodeIsLeader"}},{title:"Action",scopedSlots:{customRender:"storeNodeAction"}}],sets:[],addSetModal:!1,addSetForm:{setID:0,nodes:[]},addNodeModal:!1,addNodeForm:{setID:0,nodeID:0,host:"",servicePort:0,raftPort:0},addStoreSelectData:[],addStoreForm:{setID:0,nodeID:0,store:0}}},mounted:function(){this.getSetStatus()},watch:{host:function(t){this.getSetStatus()}},methods:{getSetStatus:function(){var t=this;if(""===this.host)return"";b({host:this.host}).then((function(e){t.sets=e.sets,t.sets.sort((function(t,e){return t.setID-e.setID}));for(var a=0;a<t.sets.length;a++){var o=t.sets[a];o.stores.sort((function(t,e){return t.storeID-e.storeID}));for(var s=0;s<o.stores.length;s++)o.stores[s].nodes=[];o.nodes.sort((function(t,e){return t.nodeID-e.nodeID}));for(var n=0;n<o.nodes.length;n++){var r=o.nodes[n];r.stores||(r.stores=[]),r.stores.sort((function(t,e){return t.storeID-e.storeID}));for(var d=0;d<r.stores.length;d++)for(var i=r.stores[d],l=0;l<o.stores.length;l++)if(o.stores[l].storeID===i.storeID){o.stores[l].nodes.push({nodeID:r.nodeID,isLeader:i.isLeader});break}}}})).catch((function(e){t.sets=[]}))},handleSetMarkClear:function(t){if(""===this.host||t.markClear)return this.$message.info("不允许操作",1),"";v({host:this.host,setMarkClear:{setID:t.setID}}).then((function(e){t.markClear=!0}))},openAddSetModal:function(){this.addSetForm={setID:0,nodes:[]},this.addSetModal=!0},addAddSetFormNodes:function(){this.addSetForm.nodes.push({nodeID:0,host:"",servicePort:0,raftPort:0})},remAddSetFormNodes:function(t){this.addSetForm.nodes.splice(t,1)},handleAddSet:function(){var t=this;if(""===this.host)return this.$message.info("请选择Pd节点",1),"";S({host:this.host,addSet:{set:this.addSetForm}}).then((function(e){t.addSetModal=!1,t.getSetStatus()})).catch((function(e){t.addSetModal=!1}))},handleRemSet:function(t){var e=this;if(""===this.host)return this.$message.info("请选择Pd节点",1),"";y({host:this.host,remSet:{setID:t}}).then((function(t){e.getSetStatus()}))},openAddNodeModal:function(t){this.addNodeForm={setID:t,nodeID:0,host:"",servicePort:0,raftPort:0},this.addNodeModal=!0},handleAddNode:function(){var t=this;if(""===this.host)return this.$message.info("请选择Pd节点",1),"";k({host:this.host,addNode:this.addNodeForm}).then((function(e){t.addNodeModal=!1,t.getSetStatus()})).catch((function(e){t.addNodeModal=!1}))},handleRemNode:function(t,e){var a=this;if(""===this.host)return this.$message.info("请选择Pd节点",1),"";g({host:this.host,remNode:{setID:t,nodeID:e}}).then((function(t){a.getSetStatus()}))},openAddStoreSelect:function(t,e,a,o){this.addStoreForm.setID=t,this.addStoreForm.nodeID=e,this.addStoreSelectData=[];for(var s=0;s<a.length;s++){for(var n=!1,r=0;r<o.length;r++)if(a[s].storeID===o[r].storeID){n=!0;break}n||this.addStoreSelectData.push({title:"Store "+a[s].storeID,value:a[s].storeID})}},handleAddStore:function(t){var e=this;if(this.addStoreForm.store=t,""===this.host)return this.$message.info("请选择Pd节点",1),"";I({host:this.host,addLearnerStoreToNode:this.addStoreForm}).then((function(t){e.getSetStatus()}))},handleRemStore:function(t,e,a){var o=this,s={setID:t,nodeID:e,store:a};if(""===this.host)return this.$message.info("请选择Pd节点",1),"";D({host:this.host,removeNodeStore:s}).then((function(t){o.getSetStatus()}))}}}),C=P,R=Object(x["a"])(C,_,M,!1,null,"0edc0561",null),T=R.exports,V={name:"Flyfish",components:{"a-meta":w,"a-set-status":T},data:function(){return{key:"fly_pd",data:[],selectValue:"",tabList:[{key:"meta",tab:"Meta"},{key:"setStatus",tab:"SetStatus"}],tabActiveKey:"meta",visible:!1,inputName:"",inputAddress:"",disabled:!0}},beforeMount:function(){this.getValue()},methods:{getValue:function(){var t=this;i({key:this.key}).then((function(e){if(t.selectValue="",e.exist){t.data=[];var a=e.value.split(";");for(var o in a){var s=a[o],n=s.split("@"),r=n[0],d=n[1];t.data.push({title:r,value:d,field:s})}}}))},valueToString:function(){if(0===this.data.length)return"";if(1===this.data.length)return this.data[0].field;for(var t=this.data[0].field,e=1;e<this.data.length;e++)t+=";"+this.data[e].field;return t},deleteKv:function(){for(var t=this,e=-1,a=0;a<this.data.length;a++){var o=this.data[a];o.value===this.selectValue&&(e=a)}if(-1!==e&&this.data.splice(e,1),0===this.data.length)l({key:this.key}).then((function(e){t.getValue()}));else{var s=this.valueToString();d({key:this.key,value:s}).then((function(e){t.getValue()}))}},handleChange:function(t){this.selectValue=t,""===this.selectValue?this.disabled=!0:this.disabled=!1},openModel:function(){this.inputName="",this.inputAddress="",this.visible=!0},handleOk:function(){var t=this;if(""!==this.inputName&&""!==this.inputAddress){for(var e=-1,a=0;a<this.data.length;a++){var o=this.data[a];o.title===this.inputName&&(o.value=this.inputAddress,o.field=this.inputName+"@"+this.inputAddress,e=a)}-1===e&&this.data.push({title:this.inputName,value:this.inputAddress,field:this.inputName+"@"+this.inputAddress});var s=this.valueToString();d({key:this.key,value:s}).then((function(e){t.visible=!1,t.getValue()}))}else this.visible=!1},handleTabChange:function(t){this.tabActiveKey=t}}},$=V,O=Object(x["a"])($,o,s,!1,null,"1b6ea979",null);e["default"]=O.exports}}]);