(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-6f35c550"],{"06c5":function(t,e,o){"use strict";o.d(e,"a",(function(){return a}));o("fb6a"),o("d3b7"),o("b0c0"),o("a630"),o("3ca3"),o("ac1f"),o("00b4");var r=o("6b75");function a(t,e){if(t){if("string"===typeof t)return Object(r["a"])(t,e);var o=Object.prototype.toString.call(t).slice(8,-1);return"Object"===o&&t.constructor&&(o=t.constructor.name),"Map"===o||"Set"===o?Array.from(t):"Arguments"===o||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(o)?Object(r["a"])(t,e):void 0}}},"0f66":function(t,e,o){"use strict";o.d(e,"a",(function(){return s})),o.d(e,"b",(function(){return n}));var r=o("b775"),a={nodeList:"node/list",nodeRemove:"node/remove"};function s(t){return Object(r["b"])({url:a.nodeList,method:"post",data:t})}function n(t){return Object(r["b"])({url:a.nodeRemove,method:"post",data:t})}},"371f":function(t,e,o){"use strict";o("848e")},5143:function(t,e,o){"use strict";o.r(e);var r=function(){var t=this,e=this,o=e.$createElement,r=e._self._c||o;return r("page-header-wrapper",{attrs:{breadcrumb:{},title:e.title},on:{back:e.goback}},[r("a-card",{attrs:{bordered:!1}},[r("a-form-model",{ref:"processEdit",attrs:{model:e.form,rules:e.rules,"label-col":e.labelCol,"wrapper-col":e.wrapperCol}},[r("a-form-model-item",{ref:"name",attrs:{label:"名称",prop:"name","wrapper-col":{span:6}}},[r("a-input",{attrs:{placeholder:"程序名称"},model:{value:e.form.name,callback:function(t){e.$set(e.form,"name",t)},expression:"form.name"}})],1),r("a-form-model-item",{attrs:{label:"执行目录","wrapper-col":{span:6}}},[r("a-input",{attrs:{placeholder:"非必填，默认为节点启动目录"},model:{value:e.form.dir,callback:function(t){e.$set(e.form,"dir",t)},expression:"form.dir"}})],1),e._l(e.form.config,(function(t,o){return r("a-form-model-item",{key:o,attrs:{"wrapper-col":0===o?e.wrapperCol:{span:14,offset:3},label:0===o?"启动配置":""}},[r("a-input",{staticStyle:{width:"60%","margin-right":"8px"},attrs:{placeholder:"配置名称"},model:{value:t.name,callback:function(o){e.$set(t,"name",o)},expression:"cfg.name"}}),r("a-input",{staticStyle:{width:"90%","margin-right":"8px"},attrs:{placeholder:"配置内容",type:"textarea","auto-size":{minRows:6,maxRows:16}},model:{value:t.context,callback:function(o){e.$set(t,"context",o)},expression:"cfg.context"}}),e.form.config.length>0?r("a-icon",{staticClass:"dynamic-delete-button",attrs:{type:"minus-circle-o"},on:{click:function(o){return e.removeCfg(t)}}}):e._e()],1)})),r("a-form-model-item",{attrs:{label:0===e.form.config.length?"启动配置":"","wrapper-col":0===e.form.config.length?e.wrapperCol:{span:14,offset:3}}},[r("a-button",{staticStyle:{width:"60%"},attrs:{type:"dashed"},on:{click:e.addCfg}},[r("a-icon",{attrs:{type:"plus"}}),e._v(" 添加配置 ")],1)],1),r("a-form-model-item",{attrs:{prop:"command"}},[r("span",{attrs:{slot:"label"},slot:"label"},[e._v(" 启动命令  "),r("a-tooltip",{attrs:{title:e.cmd_question_title}},[r("a-icon",{attrs:{type:"question-circle-o"}})],1)],1),r("a-input",{attrs:{placeholder:"启动命令"},model:{value:e.form.command,callback:function(t){e.$set(e.form,"command",t)},expression:"form.command"}})],1),r("a-form-model-item",{attrs:{"wrapper-col":{span:6}}},[r("span",{attrs:{slot:"label"},slot:"label"},[e._v(" 优先级  "),r("a-tooltip",{attrs:{title:"子进程启动关闭优先级，优先级低的，最先启动，关闭的时候最后关闭"}},[r("a-icon",{attrs:{type:"question-circle-o"}})],1)],1),r("a-input",{model:{value:e.form.priority,callback:function(t){e.$set(e.form,"priority",t)},expression:"form.priority"}})],1),r("a-form-model-item",{attrs:{"wrapper-col":{span:6}}},[r("span",{attrs:{slot:"label"},slot:"label"},[e._v(" 启动检测时长  "),r("a-tooltip",{attrs:{title:"启动进程一段时间后没有异常退出，就表示进程正常启动了。"}},[r("a-icon",{attrs:{type:"question-circle-o"}})],1)],1),r("a-input",{model:{value:e.form.start_secs,callback:function(t){e.$set(e.form,"start_secs",t)},expression:"form.start_secs"}})],1),r("a-form-model-item",{attrs:{"wrapper-col":{span:6}}},[r("span",{attrs:{slot:"label"},slot:"label"},[e._v(" 停机等待时长  "),r("a-tooltip",{attrs:{title:"这个是当我们向子进程发送stop信号后，到系统返回信息所等待的最大时间。超过这个时间会向该子进程发送一个强制kill的信号。"}},[r("a-icon",{attrs:{type:"question-circle-o"}})],1)],1),r("a-input",{model:{value:e.form.stop_wait_secs,callback:function(t){e.$set(e.form,"stop_wait_secs",t)},expression:"form.stop_wait_secs"}},[r("span",{attrs:{slot:"addonAfter"},slot:"addonAfter"},[e._v("秒")])])],1),r("a-form-model-item",{attrs:{"wrapper-col":{span:6}}},[r("span",{attrs:{slot:"label"},slot:"label"},[e._v(" 自动重启次数  "),r("a-tooltip",{attrs:{title:"进程状态为 Exited时，自动重启"}},[r("a-icon",{attrs:{type:"question-circle-o"}})],1)],1),r("a-input",{model:{value:e.form.auto_start_times,callback:function(t){e.$set(e.form,"auto_start_times",t)},expression:"form.auto_start_times"}})],1),r("a-form-model-item",{attrs:{"wrapper-col":{span:14},prop:"groups"}},[r("span",{attrs:{slot:"label"},slot:"label"},[e._v(" 分组  "),r("a-tooltip",{attrs:{title:"程序分组，批量管理"}},[r("a-icon",{attrs:{type:"question-circle-o"}})],1)],1),e._l(e.form.groups,(function(t,o){return[t.length>10?r("a-tooltip",{key:o,attrs:{title:t}},[r("a",{staticClass:"group_elem"},[e._v(" "+e._s(t.slice(0,10)+"...")+"  "),r("a-icon",{attrs:{type:"close"},on:{click:function(){return e.deleteGroup(t)}}})],1)]):r("a",{key:o,staticClass:"group_elem"},[e._v(" "+e._s(t)+"  "),r("a-icon",{attrs:{type:"close"},on:{click:function(){return e.deleteGroup(t)}}})],1)]})),e.groupInputVisible?r("a-input",{staticClass:"group_elem_in",attrs:{placeholder:"新增分组"},on:{blur:e.addGroup,keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.addGroup.apply(null,arguments)}},model:{value:e.groupInputValue,callback:function(t){e.groupInputValue=t},expression:"groupInputValue"}}):r("a-button-group",[r("a-tree-select",{attrs:{"tree-data":e.treeData,value:"Select Group ","show-search":""},on:{select:e.selectGroup}}),r("a-button",{attrs:{icon:"plus"},on:{click:function(){t.groupInputVisible=!0}}},[e._v("Add Group")])],1)],2),r("a-form-model-item",{attrs:{label:"节点",prop:"node","wrapper-col":{span:6}}},[r("a-select",{attrs:{placeholder:"选择一个节点"},model:{value:e.form.node,callback:function(t){e.$set(e.form,"node",t)},expression:"form.node"}},[e._l(e.nodes,(function(t,o){return[r("a-select-option",{key:o,attrs:{value:t.name}},[e._v(" "+e._s(t.name)+" ")])]}))],2)],1),r("a-form-model-item",{attrs:{"wrapper-col":{span:14,offset:6}}},[r("a-button",{attrs:{type:"primary"},on:{click:e.submitForm}},[e._v(" "+e._s(e.submitText)+" ")]),r("a-button",{staticStyle:{"margin-left":"10px"},on:{click:e.goback}},[e._v(" 取消 ")])],1)],2)],1)],1)},a=[],s=o("6b75");function n(t){if(Array.isArray(t))return Object(s["a"])(t)}o("a4d3"),o("e01a"),o("d3b7"),o("d28b"),o("3ca3"),o("ddb0"),o("a630");function i(t){if("undefined"!==typeof Symbol&&null!=t[Symbol.iterator]||null!=t["@@iterator"])return Array.from(t)}var l=o("06c5");o("d9e2");function c(){throw new TypeError("Invalid attempt to spread non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method.")}function u(t){return n(t)||i(t)||Object(l["a"])(t)||c()}var p=o("5530"),m=(o("a434"),o("99af"),o("4de4"),o("644f")),f=o("0f66"),d={name:"ProcessEdit",data:function(){return{labelCol:{span:3},wrapperCol:{span:14},cmd_question_title:"命令中存在配置文件时，路径前加上{{path}}，自动填充",rules:{name:[{required:!0,message:"请输入名称",trigger:"blur"},{min:3,message:"名称至少3个字",trigger:"blur"}],command:[{required:!0,message:"请输入启动命令",trigger:"blur"}],groups:[{validator:this.checkGroups,message:"至少选择一个分组",trigger:"blur"}],node:[{required:!0,message:"选择节点",trigger:"blur"}]},form:{id:0,name:"",dir:"",config:[],command:"",priority:10,start_secs:3,stop_wait_secs:10,auto_start_times:3,node:"",groups:[]},title:"新建配置",path:"",nodes:[],treeData:[],groupInputVisible:!1,groupInputValue:"",option:"",submitText:"创建"}},mounted:function(){this.path=this.$route.params.path,this.option=this.$route.params.option,this.treeData=this.$route.params.treeData,this.loadNodeList(),"edit"!==this.option&&"copy"!==this.option||(this.form=Object(p["a"])({},this.$route.params.item)),"edit"===this.option&&(this.submitText="修改",this.title="修改配置"),console.log(this.option,this.$route.params.item,this.form)},methods:{removeCfg:function(t){console.log(t);var e=this.form.config.indexOf(t);-1!==e&&this.form.config.splice(e,1)},addCfg:function(){this.form.config.push({name:"",context:""})},loadNodeList:function(){var t=this,e={pageNo:1,pageSize:1e3};Object(f["a"])(e).then((function(e){t.nodes=e.nodeList}))},addGroup:function(){this.groupInputValue&&-1===this.form.groups.indexOf(this.groupInputValue)&&(console.log(this.groupInputValue),this.form.groups=[].concat(u(this.form.groups),[this.groupInputValue])),console.log(this.form.groups),this.groupInputVisible=!1,this.groupInputValue=""},deleteGroup:function(t){var e=this.form.groups.filter((function(e){return e!==t}));this.form.groups=e,console.log(t,this.form.groups)},selectGroup:function(t){t&&-1===this.form.groups.indexOf(t)&&(console.log(t),this.form.groups=[].concat(u(this.form.groups),[t])),console.log(this.form.groups),this.tagGroupType="button"},checkGroups:function(){return this.form.groups.length>0},submitForm:function(){var t=this;this.$refs.processEdit.validate((function(e){if(e){console.log(t.form);var o=Object(p["a"])({},t.form);o.auto_start_times=parseInt(t.form.auto_start_times),o.priority=parseInt(t.form.priority),o.start_secs=parseInt(t.form.start_secs),o.stop_wait_secs=parseInt(t.form.stop_wait_secs),"create"===t.option||"copy"===t.option?Object(m["c"])(o).then((function(e){t.goback()})):Object(m["h"])(o).then((function(e){t.goback()}))}}))},goback:function(){this.$router.push({name:"plist",params:{path:this.path}})}}},b=d,g=(o("371f"),o("2877")),h=Object(g["a"])(b,r,a,!1,null,"15004490",null);e["default"]=h.exports},"644f":function(t,e,o){"use strict";o.d(e,"a",(function(){return s})),o.d(e,"b",(function(){return n})),o.d(e,"e",(function(){return i})),o.d(e,"c",(function(){return l})),o.d(e,"h",(function(){return c})),o.d(e,"d",(function(){return u})),o.d(e,"f",(function(){return p})),o.d(e,"g",(function(){return m}));var r=o("b775"),a={processList:"process/list",processCreate:"process/create",processUpdate:"process/update",processDelete:"process/delete",processStart:"process/start",processStop:"process/stop",groupList:"process/glist",groupAdd:"process/gadd",groupRemove:"process/gremove"};function s(t){return Object(r["b"])({url:a.groupList,method:"post",data:t})}function n(t){return Object(r["b"])({url:a.groupRemove,method:"post",data:t})}function i(t){return Object(r["b"])({url:a.processList,method:"post",data:t})}function l(t){return Object(r["b"])({url:a.processCreate,method:"post",data:t})}function c(t){return Object(r["b"])({url:a.processUpdate,method:"post",data:t})}function u(t){return Object(r["b"])({url:a.processDelete,method:"post",data:t})}function p(t){return Object(r["b"])({url:a.processStart,method:"post",data:t})}function m(t){return Object(r["b"])({url:a.processStop,method:"post",data:t})}},"6b75":function(t,e,o){"use strict";function r(t,e){(null==e||e>t.length)&&(e=t.length);for(var o=0,r=new Array(e);o<e;o++)r[o]=t[o];return r}o.d(e,"a",(function(){return r}))},"848e":function(t,e,o){}}]);
//# sourceMappingURL=chunk-6f35c550.8e1b1fbb.js.map