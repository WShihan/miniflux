中文

本仓库派生自 [miniflux/v2](https://github.com/miniflux/v2)，在官方版本基础上，添加一些额外功能，目前可用：

* 翻译标题：翻译订阅源标题到任意一门语言。

## 如何使用

### 翻译标题

你需要配置第三方翻译接口，目前支持[OpenAI](https://openai.com/)的ChatGPT和[百度智能云平台](https://cloud.baidu.com)的机器翻译（通用版），如果你没有权限访问上述服务，请先到对应官网申请。

在配置文件里添加`TRANSLATE_URL`，针对不同的翻译接口，该配置项的值有所不相同，

如果启用ChatGPT翻译，那么配置项的格式如下：

```plain
TRNASLATE_URL = chatgpt@url@key@model@language
```

说明：

> * Url: 接口地址
> * key：访问密钥
> * model：指定翻译模型
> * language：翻译到目标语言（需要用英文指定，比如中文：chinese）

e.g

```plain
TRANSLATE_URL=chatgpt@https://oa.api2d.net/v1/chat/completions@fk222771-TpBm4qmwaOyiyI6W3esffdfdDTuIq@gpt-3.5-turbo@chinese
```





百度智能云平台机器翻译通用版则为

```plain
TRANSLATE_URL = baidu_ml@clientid@secret@language
```

说明：

> * clientid：百度智能云平台应用id
> * secret：百度智能云平台应用密钥
> * language：翻译到目标语言（这一部分参考百度智能云平台翻译文档描述）

e.g

```
TRANSLATE_URL=baidu_ml@7mtS62xNRffffffVnP3VCaB@Yffffff8HptfCclcY5tGraONXpDEkn9@zh
```



### 效果

![image-20240528222541593](https://md-1301600412.cos.ap-nanjing.myqcloud.com/pic/typora/image-20240528222541593.png)

![image-20240528222448719](https://md-1301600412.cos.ap-nanjing.myqcloud.com/pic/typora/image-20240528222448719.png)

![image-20240528222635201](https://md-1301600412.cos.ap-nanjing.myqcloud.com/pic/typora/image-20240528222635201.png)

