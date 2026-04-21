# Windows Ollama 安装

[ollama官网](https://ollama.com)

[下载地址](https://ollama.com/download)

产品定位

![image-20260421215034895](ollama%E5%AE%89%E8%A3%85.assets/image-20260421215034895.png)

## 自定义安装路径

在有安装包目录下的命令行中，执行安装到指定目录

```bash
OllamaSetup.exe /DIR=D:\application\Ollama
```

## 手动创建大模型仓库目录

新建环境变量

![image-20260421221059822](ollama%E5%AE%89%E8%A3%85.assets/image-20260421221059822.png)

大模型存放的对应目录

## 测试是否安装成功

```bash
# 第一种
netstat -ano | findstr 11434
# 第二种
ollama --version
```

## Ollama 常用命令速查表  

| 命令                                  | 一句话说明                         |
| ------------------------------------- | ---------------------------------- |
| `ollama pull llama3`                  | 下载指定模型（例：llama3）。       |
| `ollama run llama3`                   | 启动并进入该模型交互对话。         |
| `ollama list`                         | 列出本机已下载的所有模型。         |
| `ollama rm llama3`                    | 删除不再需要的模型以节省磁盘。     |
| `ollama cp llama3 my-llama3`          | 本地复制/重命名模型。              |
| `ollama show llama3`                  | 查看模型详细信息（参数、大小等）。 |
| `ollama create my-model -f Modelfile` | 用自定义 Modelfile 构建新模型。    |
| `ollama serve`                        | 启动后台服务，供 API 调用。        |
| `ollama ps`                           | 查看当前正在运行的模型进程。       |
| `ollama stop llama3`                  | 停止正在运行的模型。               |

比如下载运行 `qwen3.5:9b` 模型

```bash
ollama run qwen3.5:9b
```

![image-20260421225016652](ollama%E5%AE%89%E8%A3%85.assets/image-20260421225016652.png)