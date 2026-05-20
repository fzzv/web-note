# LangChain1.0+版本使用方式,目前主流

# 1.导入依赖
import os
from dotenv import load_dotenv
from langchain.chat_models import init_chat_model

#通过 python-dotenv 库读取 env 文件中的环境变量，并加载到当前运行的环境中
load_dotenv(encoding='utf-8')

# 2.实例化模型
# 什么是关键字参数 k1=v1 ,k2 = v2
model = init_chat_model(
    model="qwen3.5-35b-a3b",
    model_provider="openai",
    api_key=os.getenv("QWEN_API_KEY"),
    base_url="https://dashscope.aliyuncs.com/compatible-mode/v1"
)

# 3.调用模型
print(model.invoke("你是谁").content)

