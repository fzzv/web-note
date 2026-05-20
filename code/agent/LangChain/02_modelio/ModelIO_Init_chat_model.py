# 1.导入依赖
import os
from langchain.chat_models import init_chat_model
from dotenv import load_dotenv

load_dotenv()

# 2.实例化模型
model = init_chat_model(
    model="mimo-v2.5",
    model_provider="openai",  # 1.0版本需要指定模型提供商，默认支持的提供商可以省略 model_provider
    api_key=os.getenv("XIAOMI_API_KEY"),
    base_url="https://token-plan-cn.xiaomimimo.com/v1",
)

# 3.调用模型
print(model.invoke("你是谁").content)
