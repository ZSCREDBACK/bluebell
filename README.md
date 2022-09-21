# Bluebell

> 基于 goScaffold 项目，类似于论坛。

---

> goland 默认推送到gitlab
>
> 不需要额外设置

---

> 本地使用git bash推送至github

```shell
# git add .
# git commit -m "XXX"
# 如果之前使用gitlab推送过了,这里只需要push即可
git push https://github.com/ZSCREDBACK/bluebell.git
```

## 在线接口文档

> 在线调试。

```bash
http://localhost:12345/swagger/index.html
```

## 离线接口文档

> 导入 Postman 使用。

- `bluebell\src\docs\Bluebell.postman_collection.json`

## 压力测试

```ps1
go-wrk.exe -t8 -c100 -d30s --latency http://127.0.0.1:12345/api/v1/posts?size=10
```

## 打包镜像

```shell
docker build -t bluebell_0921 .
docker run -d --name bluebell -p 12345:12345 bluebell_0921
```

## 前端配置

```bash
# 前端下载
https://github.com/Q1mi/bluebell/tree/master/bluebell_frontend
```

```nginx
# ...
    # bluebell项目动静分离
	server {
        listen       12345;
        server_name  _;

        access_log   D:/Linux/nginx-1.22.0/logs/bluebell-access.log;
        error_log    D:/Linux/nginx-1.22.0/logs/bluebell-error.log;

		# 静态文件请求
        location ~ .*\.(gif|jpg|jpeg|png|js|css|eot|ttf|woff|svg|otf)$ {
            access_log off;
            expires    1d;
            root       D:/Linux/nginx-1.22.0/html/bluebell;
        }

        # index.html 页面请求
        # 因为是单页面应用这里使用 try_files 处理一下，避免刷新页面时出现404的问题
        location / {
            root D:/Linux/nginx-1.22.0/html/bluebell;
            index index.html;
            try_files $uri $uri/ /index.html;
        }

		# 凡是以/api开头的请求,都转发给bluebell后端进行处理
        location /api {
            proxy_pass                 http://192.168.118.116:12345;
            proxy_redirect             off;
            proxy_set_header           Host             $host;
            proxy_set_header           X-Real-IP        $remote_addr;
            proxy_set_header           X-Forwarded-For  $proxy_add_x_forwarded_for;
        }
    }
# ...
```