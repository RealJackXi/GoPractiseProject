# idiomApp说明

#### 支持三中模式
1. amb + 模糊词：打印出包含模糊词的成语
2. acc + 词：    打印这个词的相关信息
3. poem + 模糊词：打印能查到模糊词的相关信息

### 逻辑
    - idiomsMap是全局变量
    - 初始化的时候，先加载本地的json文件数据到idiomsMap中
    
### 启动时
    - 启动gin的web服务来作为web服务
    
### 关闭时
    - 通过os.Singal来捕获终止信号，将gin服务关闭，将idiomMap写在本地
    