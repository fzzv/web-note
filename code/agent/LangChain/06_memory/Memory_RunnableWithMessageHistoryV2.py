"""
可持续记忆（RunnableWithMessageHistory）
"""

from langchain.chat_models import init_chat_model
from langchain_core.chat_history import InMemoryChatMessageHistory  # 内存型消息记录
from langchain_core.runnables.history import RunnableWithMessageHistory
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain_core.output_parsers import StrOutputParser
import os
from dotenv import load_dotenv

load_dotenv()

# 设置本地模型
llm = init_chat_model(
    model="mimo-v2.5-pro",
    model_provider="openai",
    api_key=os.getenv("XIAOMI_API_KEY"),
    base_url="https://token-plan-cn.xiaomimimo.com/v1",
)


# 定义全局的“会话存储”，用来保存每个 session 的聊天历史
#    （真实项目中可改为 Redis、SQLite 等）
store = {}


def get_session_history(session_id: str):
    """
    根据 session_id 获取对应的历史消息对象。
    如果不存在则创建一个新的 InMemoryChatMessageHistory。
    """
    if session_id not in store:
        store[session_id] = InMemoryChatMessageHistory()
    return store[session_id]


# 定义 Prompt 模板
#     - system: 给模型设定角色
#     - MessagesPlaceholder: 历史消息将注入这里
#     - human: 当前用户输入
prompt = ChatPromptTemplate.from_messages(
    [
        ("system", "你是一个友好的中文助理，会根据上下文回答问题。"),
        MessagesPlaceholder("history"),
        ("human", "{question}"),
    ]
)


# 构建基本链：Prompt → LLM → 输出解析
memory_chain = prompt | llm | StrOutputParser()

# -----------------------------------------------------
# 将链包装为支持记忆的版本
with_history = RunnableWithMessageHistory(
    memory_chain,  # 原始链
    get_session_history,  # 获取历史函数
    input_messages_key="question",  # 对应 prompt 输入的 key
    history_messages_key="history",  # 对应 MessagesPlaceholder 的变量名
)

# -----------------------------------------------------
# 模拟一个会话，用 session_id 区分不同用户
config = {"configurable": {"session_id": "user-001"}}

# 第一次提问：告诉模型“我叫张三”
print("用户：我叫张三。")
print("AI：", with_history.invoke({"question": "我叫张三。"}, config))

# 第二次提问：让模型回忆前面的对话
print("\n 用户：我叫什么？")
print("AI：", with_history.invoke({"question": "我叫什么？"}, config))
