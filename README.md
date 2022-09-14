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