# pip install langchain-community
# pip install dashscope
# 如果报错，强制重装 cffi（解决模块损坏、缓存残留问题，优先推荐）
# pip install --upgrade --force-reinstall cffi

import os
from langchain_community.chat_models.tongyi import ChatTongyi
from langchain_core.messages import HumanMessage
from dotenv import load_dotenv

load_dotenv()

chatLLM = ChatTongyi(
    model="qwen-plus",
    api_key=os.getenv("QWEN_API_KEY"),
    streaming=True,
    # other params...
)
# 打印结果
print(chatLLM.invoke("你是谁"))

print("*" * 60)

res = chatLLM.stream([HumanMessage(content="你好，你是谁")], streaming=True)
for r in res:
    print("chat resp:", r.content)
