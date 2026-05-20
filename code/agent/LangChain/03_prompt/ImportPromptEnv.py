from langchain_core.prompts import PromptTemplate
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.prompts import FewShotPromptTemplate
from langchain_core.prompts import (
    ChatMessagePromptTemplate,
    SystemMessagePromptTemplate,
    AIMessagePromptTemplate,
    HumanMessagePromptTemplate,
)

from langchain.messages import SystemMessage, HumanMessage, AIMessage, ToolMessage


messages = [
    SystemMessage(content="你是一位乐于助人的智能小助手"),
    HumanMessage(content="你好，请你介绍一下你自己"),
    AIMessage(content="我是一名人工智能助手，请问您有什么想问的嘛?"),
    # ToolMessage - 用于工具调用场景
    ToolMessage(
        tool_call_id="call_abc123",  # 关联的工具调用ID
        content='{"population": 21540000, "area": "16410平方公里"}',  # 工具执行结果
    )
]

print(messages)