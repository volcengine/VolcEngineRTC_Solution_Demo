# VRTC 产品Demo


### 获取 AppId 和临时 Token

可参考 RTC 接入指南获取 AppId 和临时 Token。临时 Token 的生成需要输入房间 ID 和用户 ID，这里输入的房间 ID 和用户 ID 需要在编译成功的 Demo 的登录页输入相同的，才可以正常进入房间。

### 配置 Demo 工程文件

1. 全局安装 node、yarn
2. 进入工程目录，修改 AppID。使用控制台获取的 AppID 覆盖 src 文件夹下 config.ts 里的 appId 值
3. 进入工程目录，修改 Token。临时 Token 覆盖 src 文件下 config.ts 里的 token 值

### 编译运行

1. 在命令行里输入 yarn
2. 运行项目 yarn start


## 项目结构

```shell
|-- Project
  |-- assets 图片、字体等资源
  |-- src 源代码
    |-- components Demo中复用的组件
    |-- layouts Demo页面结构
    |-- sdk 核心SDK与封装的公共库
    |-- lib 基于SDK封装的业务Util
    |-- locales 国际化文本库
    |-- models 业务数据模型层(dva)
    |-- pages 页面
    |-- app.less 全局样式
    |-- app.tsx 入口配置(dva)
    |-- config.ts 静态常量
    |-- app-interfaces.ts App全局类型定义
```