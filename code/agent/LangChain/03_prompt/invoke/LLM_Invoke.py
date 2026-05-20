# 1.导入依赖
import os
from dotenv import load_dotenv
from langchain.chat_models import init_chat_model
from langchain.messages import HumanMessage, SystemMessage

# 通过 python-dotenv 库读取 env 文件中的环境变量，并加载到当前运行的环境中
load_dotenv()

# List<Messages>

# 2.实例化模型
model = init_chat_model(
    model="qwen-plus",
    model_provider="openai",
    api_key=os.getenv("QWEN_API_KEY"),
    base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",
)

# 构建消息列表
messages = [
    SystemMessage(
        content="你是一个法律助手，只回答法律问题，超出范围的统一回答，非法律问题无可奉告"
    ),
    HumanMessage(content="简单介绍下广告法，一句话告知20字以内"),
    # HumanMessage(content="2+3等于几?")
]

# 3.调用模型
response = model.invoke(messages)  # ainvoke
print(
    f"响应类型：{type(response)}"
)  # 响应类型：<class 'langchain_core.messages.ai.AIMessage'>
# 打印结果
print(response.content)  # 规范广告活动，保护消费者权益。
print(
    response.content_blocks
)  # [{'type': 'text', 'text': '规范广告活动，保护消费者权益。'}]
