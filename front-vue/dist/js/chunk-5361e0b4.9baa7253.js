(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-5361e0b4"],{"0f66":function(t,e,a){"use strict";a.d(e,"a",(function(){return s})),a.d(e,"b",(function(){return o}));var n=a("b775"),r={nodeList:"node/list",nodeRemove:"node/remove"};function s(t){return Object(n["b"])({url:r.nodeList,method:"post",data:t})}function o(t){return Object(n["b"])({url:r.nodeRemove,method:"post",data:t})}},1657:function(t,e,a){"use strict";a("46a4")},"46a4":function(t,e,a){},"868c":function(t,e,a){"use strict";a.r(e);var n=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[a("a-row",{attrs:{gutter:24}},[a("a-col",{style:{marginBottom:"24px"},attrs:{sm:24,md:12,xl:6}},[a("a-card",{staticClass:"header-card",attrs:{bordered:!1}},[a("a-row",[a("a-col",{attrs:{span:18}},[a("span",{staticClass:"header-card-value"},[t._v(t._s(t.onlineNodeCnt)+"/"+t._s(t.nodes.length))]),a("br"),a("span",{staticClass:"header-card-title"},[t._v("在线/总数")]),a("br"),a("span",{staticClass:"header-card-desc"},[t._v("已注册的节点数")])]),a("a-col",{attrs:{span:6}},[a("a-icon",{staticClass:"header-card-icon",attrs:{type:"cloud-server"}})],1)],1)],1)],1),a("a-col",{style:{marginBottom:"24px"},attrs:{sm:24,md:12,xl:6}},[a("a-card",{staticClass:"header-card",attrs:{bordered:!1}},[a("a-row",[a("a-col",{attrs:{span:18}},[a("span",{staticClass:"header-card-value"},[t._v(t._s(t.cmdCount))]),a("br"),a("span",{staticClass:"header-card-title"},[t._v("任务数")]),a("br"),a("span",{staticClass:"header-card-desc"},[t._v("已注册的任务")])]),a("a-col",{attrs:{span:6}},[a("a-icon",{staticClass:"header-card-icon",attrs:{type:"calendar"}})],1)],1)],1)],1),a("a-col",{style:{marginBottom:"24px"},attrs:{sm:24,md:12,xl:6}},[a("a-card",{staticClass:"header-card",attrs:{bordered:!1}},[a("a-row",[a("a-col",{attrs:{span:18}},[a("span",{staticClass:"header-card-value"},[t._v(t._s(t.cmdSuccess))]),a("br"),a("span",{staticClass:"header-card-title"},[t._v("成功")]),a("br"),a("span",{staticClass:"header-card-desc"},[t._v("任务成功次数")])]),a("a-col",{attrs:{span:6}},[a("a-icon",{staticClass:"header-card-icon",attrs:{type:"check"}})],1)],1)],1)],1),a("a-col",{style:{marginBottom:"24px"},attrs:{sm:24,md:12,xl:6}},[a("a-card",{staticClass:"header-card",attrs:{bordered:!1}},[a("a-row",[a("a-col",{attrs:{span:18}},[a("span",{staticClass:"header-card-value"},[t._v(t._s(t.cmdFailed))]),a("br"),a("span",{staticClass:"header-card-title"},[t._v("失败")]),a("br"),a("span",{staticClass:"header-card-desc"},[t._v("任务失败次数")])]),a("a-col",{attrs:{span:6}},[a("a-icon",{staticClass:"header-card-icon",attrs:{type:"close"}})],1)],1)],1)],1)],1),a("a-card",{attrs:{bordered:!1}},[a("div",{style:{marginBottom:"24px"}},[a("a-button",{attrs:{type:"primary",icon:"plus"},on:{click:function(e){return t.showEdit("create")}}},[t._v("创建命令")])],1),a("s-table",{ref:"table",attrs:{rowKey:"name",size:"default","data-name":"cmdList",columns:t.columns,handleMethods:t.handleDataList,data:t.loadData},scopedSlots:t._u([{key:"name",fn:function(e){return[e.length>10?a("a-tooltip",{attrs:{title:e}},[t._v(" "+t._s(e.slice(0,10)+"...")+" ")]):a("span",[t._v(t._s(e))])]}},{key:"context",fn:function(e){return[e.length>30?a("a-tooltip",{staticStyle:{color:"#00BFFF"},attrs:{title:e}},[t._v(" "+t._s(e.slice(0,30)+"...")+" ")]):a("a-tooltip",{staticStyle:{color:"#00BFFF"},attrs:{title:e}},[t._v(" "+t._s(e)+" ")]),a("a-tooltip",[a("template",{slot:"title"},[t._v(t._s(t.copyTitle))]),a("a-icon",{directives:[{name:"clipboard",rawName:"v-clipboard:copy",value:e,expression:"text",arg:"copy"}],attrs:{type:"copy"},on:{click:t.copyClick}})],2)]}},{key:"action",fn:function(e,n){return[a("div",[a("a",{on:{click:function(e){return t.execPage(n)}}},[t._v("执行")]),a("a-divider",{attrs:{type:"vertical"}}),a("a",{on:{click:function(e){return t.showEdit("edit",n)}}},[t._v("修改")]),a("a-divider",{attrs:{type:"vertical"}}),a("a",{on:{click:function(e){return t.logPage(n)}}},[t._v("日志")]),a("a-divider",{attrs:{type:"vertical"}}),a("a-popconfirm",{attrs:{title:"确定要删除吗？"},on:{confirm:function(e){return t.deleteCmd(n)}}},[a("a-icon",{staticStyle:{color:"red"},attrs:{slot:"icon",type:"question-circle-o"},slot:"icon"}),a("a",{staticStyle:{color:"red"}},[t._v("删除")])],1)],1)]}}])})],1),a("edit-mod",{ref:"editMod",attrs:{visible:t.visible,model:t.mdl,"is-edit":t.isEdit},on:{cancel:t.handleCancel,success:t.handleOk}})],1)},r=[],s=a("5530"),o=(a("d81d"),a("b0c0"),a("d3b7"),a("c1df")),i=a.n(o),c=a("2af9"),d=a("bb62"),l=a("0f66"),u=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("a-modal",{staticClass:"form",attrs:{title:t.isEdit?"修改命令":"创建命令",width:1200,visible:t.visible,confirmLoading:t.confirmLoading},on:{ok:t.submitForm,cancel:t.cancel}},[a("a-form-model",{ref:"editMod",attrs:{model:t.form,rules:t.formRule,"label-col":{span:4},"wrapper-col":{span:18}}},[a("a-form-model-item",[a("span",{attrs:{slot:"label"},slot:"label"},[t._v("命令类型")]),t._v(" SHELL ")]),a("a-form-model-item",{attrs:{prop:"name"}},[a("span",{attrs:{slot:"label"},slot:"label"},[t._v("命令名称")]),a("a-input",{attrs:{placeholder:"请输入名称"},model:{value:t.form.name,callback:function(e){t.$set(t.form,"name",e)},expression:"form.name"}})],1),a("a-form-model-item",{attrs:{label:"执行目录"}},[a("a-input",{attrs:{placeholder:"非必填，默认为节点启动目录"},model:{value:t.form.dir,callback:function(e){t.$set(t.form,"dir",e)},expression:"form.dir"}})],1),a("a-form-model-item",{attrs:{prop:"context"}},[a("span",{attrs:{slot:"label"},slot:"label"},[t._v(" 命令内容  "),a("a-tooltip",{attrs:{title:t.ctx_question_title}},[a("a-icon",{attrs:{type:"question-circle-o"}})],1)],1),a("a-textarea",{staticStyle:{background:"#0e0e0e",color:"white"},attrs:{"auto-size":{minRows:6,maxRows:10}},model:{value:t.form.context,callback:function(e){t.$set(t.form,"context",e)},expression:"form.context"}})],1),a("a-form-model-item",{directives:[{name:"show",rawName:"v-show",value:t.hasArgs,expression:"hasArgs"}],attrs:{label:"命令参数"}},[t._l(t.form.args,(function(e,n){return[t._v(" "+t._s(n)+" "),a("a-input",{key:n,attrs:{placeholder:"参数值选填，执行时可调整"},model:{value:t.form.args[n],callback:function(e){t.$set(t.form.args,n,e)},expression:"form.args[kk]"}})]}))],2)],1)],1)},m=[];a("b64b"),a("4d63"),a("ac1f"),a("25f0"),a("466d"),a("1276");function f(t){var e=new RegExp(/\{\{(_*[a-zA-Z]+[_a-zA-Z0-9]*)\}\}/g),a={},n=t.match(e);if(n)for(var r=0;r<n.length;r++){var s=n[r].split("{{")[1].split("}}")[0];a.hasOwnProperty(s)||(a[s]="")}return a}var p={name:"EditMod",props:{visible:{type:Boolean,required:!0},isEdit:{type:Boolean,required:!0},model:{type:Object,default:function(){return null}}},data:function(){return{ctx_question_title:"可在命令中设置变量值，以{{key}}的形式表示",confirmLoading:!1,hasArgs:!1,form:{id:0,name:"",dir:"",context:"",args:{}},formRule:h}},watch:{model:{handler:function(t){this.hasArgs=!1,t&&t.id?this.form=Object(s["a"])({},t):this.form={name:"",dir:"",context:"",args:{}}},deep:!0,immediate:!0},"form.context":function(t,e){if(t){var a=f(t),n={};for(var r in a)n[r]="",this.form.args.hasOwnProperty(r)&&(n[r]=this.form.args[r]);this.form.args=n,Object.keys(this.form.args).length>0?this.hasArgs=!0:this.hasArgs=!1}}},methods:{submitForm:function(){var t=this;this.$refs.editMod.validate((function(e){if(e){t.confirmLoading=!0;var a=Object(s["a"])({},t.form);t.isEdit?Object(d["f"])(a).then((function(e){t.$emit("success")})).catch((function(t){})).finally((function(e){t.confirmLoading=!1})):Object(d["a"])(a).then((function(e){t.$emit("success")})).catch((function(t){})).finally((function(e){t.confirmLoading=!1}))}else t.confirmLoading=!1}))},cancel:function(){this.$refs.editMod.resetFields(),this.$emit("cancel")}}},h={name:[{required:!0,message:"命令名称",trigger:"change"}],context:[{required:!0,message:"命令内容",trigger:"change"}]},b=p,v=a("2877"),g=Object(v["a"])(b,u,m,!1,null,"395738dc",null),_=g.exports,x={name:"CmdList",components:{STable:c["a"],EditMod:_},data:function(){var t=this;return{columns:y,loadData:function(e){return Object(d["d"])(e).then((function(e){return t.cmdCount=e.totalCount,t.cmdSuccess=e.success,t.cmdFailed=e.failed,e}))},visible:!1,confirmLoading:!1,isEdit:!0,mdl:null,copyTitle:"复制",nodes:[],onlineNodeCnt:0,cmdCount:0,cmdSuccess:0,cmdFailed:0}},mounted:function(){this.getNodeList()},methods:{getNodeList:function(){var t=this;this.nodes=[],this.onlineNodeCnt=0;var e={pageNo:1,pageSize:1e3};Object(l["a"])(e).then((function(e){for(var a in t.nodes=e.nodeList,e.nodeList){var n=e.nodeList[a];n.online&&(t.onlineNodeCnt+=1)}}))},handleDataList:function(t){var e=t.map((function(t){return Object(s["a"])({},t)}));return e},handleOk:function(){this.$refs.table.refresh(),this.visible=!1,this.mdl={}},handleCancel:function(){this.visible=!1,this.mdl={}},showEdit:function(t,e){this.visible=!0,"edit"===t?(this.isEdit=!0,this.mdl=Object(s["a"])({},e)):(this.isEdit=!1,this.mdl={})},execPage:function(t){this.$router.push({name:"cmdexec",params:{record:Object(s["a"])({},t)}})},logPage:function(t){this.$router.push({name:"cmdlog",params:{id:t.id,name:t.name}})},deleteCmd:function(t){var e=this;Object(d["b"])({id:t.id}).then((function(t){e.$refs.table.refresh()})).catch((function(t){})).finally((function(t){e.loading=!1}))},copyClick:function(){var t=this;this.copyTitle="复制成功",setTimeout((function(){t.copyTitle="复制"}),1500)}}},y=[{title:"命令名称",dataIndex:"name",scopedSlots:{customRender:"name"}},{title:"命令内容",dataIndex:"context",scopedSlots:{customRender:"context"}},{title:"执行次数",align:"center",customRender:function(t){return t+" 次"},dataIndex:"call_no"},{title:"创建时间",dataIndex:"create_at",customRender:function(t){return i.a.unix(t).format("YYYY-MM-DD HH:mm:ss")}},{title:"最后修改人",dataIndex:"user"},{title:"最后修改时间",dataIndex:"update_at",customRender:function(t){return i.a.unix(t).format("YYYY-MM-DD HH:mm:ss")}},{title:"操作",scopedSlots:{customRender:"action"}}],C=x,w=(a("1657"),Object(v["a"])(C,n,r,!1,null,"e52d85c4",null));e["default"]=w.exports},bb62:function(t,e,a){"use strict";a.d(e,"d",(function(){return s})),a.d(e,"a",(function(){return o})),a.d(e,"b",(function(){return i})),a.d(e,"f",(function(){return c})),a.d(e,"c",(function(){return d})),a.d(e,"e",(function(){return l}));var n=a("b775"),r={cmdList:"cmd/list",cmdCreate:"cmd/create",cmdDelete:"cmd/delete",cmdUpdate:"cmd/update",cmdExec:"cmd/exec",cmdLog:"cmd/log"};function s(t){return Object(n["b"])({url:r.cmdList,method:"post",data:t})}function o(t){return Object(n["b"])({url:r.cmdCreate,method:"post",data:t})}function i(t){return Object(n["b"])({url:r.cmdDelete,method:"post",data:t})}function c(t){return Object(n["b"])({url:r.cmdUpdate,method:"post",data:t})}function d(t,e){return Object(n["b"])({url:r.cmdExec,method:"post",timeout:t,data:e})}function l(t){return Object(n["b"])({url:r.cmdLog,method:"post",data:t})}}}]);