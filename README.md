### YQBlog 是一个基于语雀文档开放接口开发的极简博客系统

### 技术原理
+ 通过语雀开放接口获取文章信息并渲染。
+ 通过将语雀文章数据缓存在内存中并持久化实现加速
+ 通过实现倒排索引实现文章全文搜索功能。
+ 通过反向代理CDN实现语雀静态资源的展示。
+ 通过acme实现自动https证书


### 设计初衷
+ 折腾各种博客系统，各种模板眼花缭乱，很多时候在选择博客模板上花费了巨量的时间。希望回归博客的本质：专注于写作而非专注于博客本身。
+ 语雀的编辑功能是目前写作体验最好的。但是不支持域名绑定， 且官方的博客知识库只能单知识库发布，想要分类只能靠目录。而实际场景需要多知识库来对内容进行分类。
+ 希望文章完成写作即发布，无需构建deploy过程。无需繁重的各种配置、样式适配。
+ 不用关心闹心的图床问题，直接粘贴图片发布到语雀，会自动反向代理解决跨域图片加载问题。
+ YQBlog一次部署好之后，可以遗忘， 专注的创作即可。博客内容写作即所得。

### 功能介绍
+ [Demo](https://wangxun.tech/)
+ 基于语雀文档开放API实现， 在语雀中编辑文档，博客即所得。
+ 语雀的文章搜索很难用，实现了一个倒排索引提供精确的全文检索功能。
+ 自动维护https证书，通过autoSSL参数控制https开关
+ 支持主题， 虽然目前仅有唯一的一个极简(懒)主题，其实YQBlog的主题很容易做， 参照themes中的四个html文件的实现即可。欢迎贡献。
+ 微信公众号排班：复制粘贴文章内容至公众号编辑框即可。
+ 对接vssue留言板插件


### 为什么不用 xxx ？？？
+ wordpress
    + 过于臃肿
    + markdown编辑体验差是硬伤，编辑器支持的markdown样式发布后会变成另外一个样子。
    + 速度慢。
    
+ hexo等静态博客系统
    + 折腾成本较高
    + 编辑完文章需要构建发布
    + 让人头痛的图床问题
    + 管理文章心累
    + 编辑体验不好
  
+ github pages
  + 慢
  + 需要hexo等静态博客方案，所以具有其所有缺点。
  
### 开始使用

+ clone 代码
```bigquery
git pull https://github.com/lllwan/YQBlog.git
```
+ 使用docker-compose

  1. 准备配置文件，参照：config.yaml.example
     ```bigquery
     cp config.yaml.example config.yaml
     vim config.yaml```
    
  2. 编辑docker-compose.yml
  3. 运行
    ```bigquery
    docker-compose up -d --build
    ```
  
+ 手动编译
  ```bigquery
  go build *.go -o YQBlog
  ./YQBlog
  ```
  
### 配置文件
+ token可以在语雀--> 账户设置--> Token中创建， 链接：https://www.yuque.com/settings/tokens
+ token需要读取知识库和文档的权限

```bigquery
---
yuque:
  api: "https://www.yuque.com/api/v2"
  token: "you yuque token"
  user: "yuque id"
  # 语雀仓库ID， 参阅：https://www.yuque.com/yuque/developer/repo
  repos:
    - name: "运维笔记"
      repo: "kkgfxm"   # 语雀仓库Id
    - name: "云原生"
      repo: "ooa19f"
    - name: "DIY搞事情"
      repo: "bua6cb"
    - name: "开开脑洞"
      repo: "ussmi8"
blog:
  title: "WangXun`s Blog"
  subtitle: "大道至简"
  keywords: "页面keywords"
  avatar: "https://blog-download-1251192068.cos.ap-shanghai.myqcloud.com/background.jpg"
  description: "页面description"
  author: "WangXun"
  # 友情链接
  link:
      - name: "google"
      - link: "https://www.google.com"
  # 留言板插件vssue的配置， 获取步骤参考：https://vssue.js.org/zh/guide/github.html#%E5%88%9B%E5%BB%BA%E4%B8%80%E4%B8%AA%E6%96%B0%E7%9A%84-oauth-app
  vssue:
    owner: "lllwan"
    repo: "blog"
    clientId: "xxxxx"
    clientSecret: "xxxxx"

manage:
  # 是否开启自动https证书维护， 此功能需要博客可以直接被外网访问。
  autoSSL: false
  # http端口
  httpPort: 80
  # https监听端口
  httpsPort: 443
  # 博客主题，目前只有默认主题，无需修改。
  theme: default
```

### TODO
+ 文章搜索
+ RSS订阅