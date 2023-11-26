# Github Proxy 
一个简单的Github代理,支持页面代理和git命令代理 暂时只支持```git clone```不支持```git push```
# 编译方法
```
go build -o ghproxy main.go
```
# 使用方法
```
ghproxy -port=[运行的端口] -proxy.http=[HTTP代理地址 如果没有可以不指定] -github.index=[为true时会展示github主页]
```
