from langchain_openai import ChatOpenAI
from langchain.messages import AIMessage, HumanMessage, SystemMessage
import os
from dotenv import load_dotenv

load_dotenv()
api_key = os.getenv("ZHIPU_API_KEY")
llm = ChatOpenAI(
    temperature=0.6,
    model="glm-5.1",
    api_key=api_key,
    base_url="https://open.bigmodel.cn/api/paas/v4/",
)


# 创建消息
messages = [
    AIMessage(content="你好！"),
    SystemMessage(content="你是一个诗人"),
    HumanMessage(content="写一首四行的关于人工智能的短诗。"),
]

# 调用模型
response = llm.invoke(messages)
print(response.content)
