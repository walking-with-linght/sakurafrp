

基于[樱花内网穿透](https://github.com/ZeroDream-CN/SakuraPanel)改版而来

服务器使用[frpv0.60.0](https://github.com/fatedier/frp/releases/tag/v0.60.0)稳定版本

开源万岁。

基于**go version go1.22.0 linux/amd64**编译

需要改动如下：

```
sakura\modules\download.php
```

download.demo.com改成自己的下载链接



```
frp\pkg\config\flags.go 
frp\pkg\config\v1\server.go
```

frpapi.demo.com改成面板所在服务器域名



```
sakura\frps-config\frps_for_sakura.ini
sakura\api\index.php
```

http://frp.demo.com/api/ 改成自己的面板api地址
api_token 改成自己的api_token

要改的基本和原版一样，只是把api地址改成自己的