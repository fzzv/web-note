"""
模型参数演示
"""

import os
from langchain.chat_models import init_chat_model
from dotenv import load_dotenv

load_dotenv()

model = init_chat_model(
    model="mimo-v2.5",
    model_provider="openai",
    api_key=os.getenv("XIAOMI_API_KEY"),
    base_url="https://token-plan-cn.xiaomimimo.com/v1",
    temperature=1.0,
)

# 3.调用模型
for x in range(3):
    print(model.invoke("写一句关于春天的词,14字以内").content)
