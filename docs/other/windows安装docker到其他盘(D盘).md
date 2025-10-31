# windows 安装docker 到其他盘

## 下载安装包

首先去官网下载安装包

[docker.com](https://www.docker.com)

## 进入命令行工具

比如将安装包放入D盘根目录，在根目录进入命令行

输入安装命令，安装到D盘下的Docker目录下（也可以自定义安装到其他目录）

```shell
"Docker Desktop Installer.exe"  install --installation-dir="D:\Docker"
```

## 设置文件存储路径

安装完成后，进入docker desktop的设置页面，修改 image 的路径到D盘

![image-20251031091325633](windows%E5%AE%89%E8%A3%85docker%E5%88%B0%E5%85%B6%E4%BB%96%E7%9B%98(D%E7%9B%98).assets/image-20251031091325633.png)