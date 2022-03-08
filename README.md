# amp

amp 是一个自动化管理平台： 计划任务、进程控制

## 开发语言与框架

前端： vue4.5.15 + antd1.7.8

后端： go1.17.7

## 页面展示 

![image](https://github.com/yddeng/amp/blob/master/assets/cmd_list.jpg)
![image](https://github.com/yddeng/amp/blob/master/assets/cmd_exec.jpg)
![image](https://github.com/yddeng/amp/blob/master/assets/cmd_exec_result.jpg)
![image](https://github.com/yddeng/amp/blob/master/assets/process_list.jpg)
![image](https://github.com/yddeng/amp/blob/master/assets/machine_list.jpg)

## 启动

### 前端项目

切换到 `front-vue` 目录

安装依赖 `yarn install `

更改 `vue.config.js` 第7行 `const target = 'http://10.128.2.123:40156'` 地址

运行 `yarn run serve `

打包 `yarn build` , 运行后会在`front-vue`目录生成 `dist` 文件夹，里面就是构建打包好的文件

### 后端项目

1. 中心节点。部署web项目、管理子节点

切换到 `back-go/cmd` 目录

配置 `amps_config.json`
```
  "data_path": "./data",           // 数据存放目录
  "center_config": {
    "address": "0.0.0.0:40155",    // 子节点访问地址 
    "token": "token"               // 子节点登陆时验证的令牌
  },
  "web_config": {
    "address": "0.0.0.0:40156",    // 前端访问地址
    "app": "../../front-vue/dist", // 前端打包好的静态文件
    "admin": {                     // 前端控制登陆账号及密码
      "username": "admin",
      "password": "123456"
    }
  }
```

启动 `go run amps.go` 或者 `go build amps.go`+`./amps`

2. 子节点。 上报物理机、进程状态、执行脚本

切换到 `back-go/cmd` 目录

配置 `ampe_config.json`
```
{
  "name":      "executor",        // 子节点名字，唯一建
  "net":       "",                // 公网IP
  "inet":      "127.0.0.1",       // 内网IP
  "center":    "127.0.0.1:40155", // 中心节点地址
  "token":     "token",           // 登陆中心节点令牌
  "data_path": "./data"           // 数据存放目录
}
```

启动 `go run ampe.go` 或者 `go build ampe.go`+`./ampe`


## todo

进程监控报警、通知。 物理机监控报警

