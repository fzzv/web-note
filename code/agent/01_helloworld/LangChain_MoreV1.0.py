# LangChain1.0+版本使用方式,目前主流,多模型共存


# 1.导入依赖
from dotenv import load_dotenv
from langchain.chat_models import init_chat_model
import os
# .env文件读取
load_dotenv()

# 2.实例化模型
# model = init_chat_model(
#     model="qwen-plus",
#     model_provider="openai",
#     api_key=os.getenv("QWEN_API_KEY"),
#     base_url="https://dashscope.aliyuncs.com/compatible-mode/v1"
# )
#
# # 3.调用模型
# print(model.invoke("你是谁").content)
#
# print("*" * 70)

# 4.实例化模型v2
"""
说明：
model="deepseek-chat" 和 base_url="https://api.deepseek.com" 
刚好匹配默认的 model_provider（如 deepseek），因此无需显式传入，函数内部做了智能推导
如果切换成其他模型（如 OpenAI），若默认值不匹配，就需要显式指定 model_provider="openai"。
"""
model = init_chat_model(
    model="deepseek-chat", # deepseek-chat 对应 DeepSeek-V3.2 的非思考模式
    model_provider="deepseek",
    api_key=os.getenv("deepseek-api"),
    base_url="https://api.deepseek.com"
)

# 5.调用模型v2
print(model.__dict__)
print(model.invoke("你是谁").content)