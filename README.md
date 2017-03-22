### todo
* 暂时用redis缓存，实际项目中会使用mysql保存冷数据，游戏角色超过一个星期没有上线，从缓存保存到mysql中并删除；缓存中的数据每天凌晨定时保存到mysql

### 通讯架构
* 架构分为三部分，客户端，网关，游戏区服务器，以及其他逻辑服务器（例如pvp服务器）；
* 消息头为16个字节；id，消息源或者目标id，4byte；seq，消息序列号，2byte；ret，消息返回结果，2byte；cmd，消息命令字，2byte；unlen，未压缩前长度；len，消息体长度，2byte，理论上不要超过65535字节.


### 项目说明
* Gate - 网关服务器，用于转发客户端与游戏服务器之间的消息，验证玩家数据，负载均衡，广播，数据统计
* Login - 登陆服务器，玩家通过登陆服务器获取token，再使用token与Gate登陆，并提供第三方登陆和充值功能
* Center - 游戏区服务器，提供每个区服的逻辑
* Client - 客户端测试程序
* Tool - 工具集合，

### Log
* 后期会完善，目前该框架已承受30万人同时在线，并且不断完善中，本框架适合于一般的游戏server，后期会开源一个高可用游戏框架。不过需要比较长的时间去完善。
* 后期框架介绍：

   分为两个版本：
   1. 简单版本：用于小型游戏，50万人以下的游戏，逻辑简单，交叉性不强。该框架主要是用于游戏初期使用，分为Gate、Login、Center以及其他逻辑服务器，并且带有后台管理以及监控系统。
   2. 大型版本：达到单实例3000K长连接，每秒消息下发量20k-50kQPS
   
      dispatcher service
      room service
      register service
      coordinator service
      saver service
      center service
      agent service
