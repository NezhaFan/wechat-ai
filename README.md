### 之前是 对接 OpenAI 的，太老了，重写了一下。

### 一、介绍
- 这是一个用于**公众号对接大模型聊天**的项目。
- 该项目最初版本是对接 OpenAI 开发的，彼时调用需要一些魔法且速度不快，现在已经没有此类问题
- 存在问题：微信限制，只能一问一答且15秒超时限制，如果15秒内不返回结果则无法主动推送。所以建议使用速度较快的模型。
- 体验。关注公众号`杠点杠`尝试提问，这仅是个人娱乐号，不推送。

### 二、特性
- [x] 优化微信被动回复超时问题。(微信是每次5秒，询问3次，即最大15秒)  
- [x] 支持参数调节，人物预设、滑动记录聊天次数、单次回复长度预估、温度等
- [x] 支持上下文。(在`chat`文件夹记录不同用户聊天内容，可以自己定期删除)
- [] 才知道公众号已经取消语音消息自动转文字功能，所以该功能不支持。

### 三、部署
- 拷贝配置文件 `config.yaml`
- 配置大模型 
  - 阿里百炼 
    - 申请Key https://bailian.console.aliyun.com/?apiKey=1#/api-key 
    - 模型列表  https://help.aliyun.com/zh/model-studio/getting-started/models 
  - 字节火山引擎
    - 申请Key https://console.volcengine.com/ark/region:ark+cn-beijing/apiKey
    - 模型列表 https://console.volcengine.com/ark/region:ark+cn-beijing/model
  - DeepSeek (不推荐，没有小模型，速度比较慢。非要使用可以用阿里或者字节的deepseek大模型)
    - 申请Key: https://platform.deepseek.com/api_keys
    - 模型： deepseek-reason (R1) 、deepseek-chat (V3)
- 配置微信公众号`令牌Token`：[微信公众平台](https://mp.weixin.qq.com/)->设置与开发->开发接口管理->基本配置->令牌(Token) 
  
- 部署服务。下载右侧 Releases 中的二进制文件与  `config.yaml` 同目录，直接执行即可。 (使用`nohup ./wechat-ai-amd64 >> ./data.log 2>&1 &` 后台运行)

- 配置公众号服务器地址(URL)。 填写 `http://服务器IP/wx`（该连接勿手动调用），设置明文方式传输，提交后，点击「启用」。 （初次启用可能要等一会生效）

