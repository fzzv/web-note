import os
from langchain_deepseek import ChatDeepSeek
from dotenv import load_dotenv
from pydantic import SecretStr

load_dotenv()

api_key = os.getenv("deepseek-api")

# 初始化 deepseek
# 看看ChatDeepSeek类的源码，解释了为什么不写调用地址，chat_modesl.py源码第176行
model = ChatDeepSeek(
    model="deepseek-chat",
    temperature=0,
    max_tokens=None,
    timeout=None,
    max_retries=2,
    api_key=SecretStr(api_key) if api_key else None,
)

# 打印结果
print(model.invoke("什么是LangChain?20字以内回答，简洁"))
