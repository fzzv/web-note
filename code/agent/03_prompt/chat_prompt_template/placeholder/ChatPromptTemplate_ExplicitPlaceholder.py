"""
如果我们不确定消息何时生成，也不确定要插入几条消息，比如在提示词中添加聊天历史记忆这种场景，
可以在ChatPromptTemplate添加MessagesPlaceholder占位符，在调用invoke时，在占位符处插入消息。

显式使用MessagesPlaceholder
"""

from langchain_core.messages import HumanMessage, AIMessage
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder

# 构建一个 ChatPromptTemplate，包含多种消息类型：
prompt = ChatPromptTemplate.from_messages(
    [
        # 添加一条系统消息，设定 AI 的角色或行为准则
        (
            "system",
            "你是一个资深的Python应用开发工程师，请认真回答我提出的Python相关的问题",
        ),
        # 插入 memory 占位符，用于填充历史对话记录（如多轮对话上下文）
        MessagesPlaceholder("memory"),
        # 添加一条用户问题消息，用变量 {question} 表示
        ("human", "{question}"),
    ]
)

# 调用 prompt.invoke 来格式化整个 Prompt 模板
# 传入的参数中：
# - memory：是一组历史消息，表示之前的对话内容（多轮上下文）
# - question：是当前用户的问题
prompt_value = prompt.invoke(
    {
        "memory": [
            # 用户第一轮说的话
            HumanMessage("我的名字叫Fan，是一名程序员"),
            # AI 第一轮的回应
            AIMessage("好的，Fan你好"),
        ],
        # 当前问题：结合上下文，测试模型是否记住了用户名字
        "question": "请问我的名字叫什么？",
    }
)

# 打印生成的完整 prompt 文本，格式化后的聊天记录
print(prompt_value.to_string())
# System: 你是一个资深的Python应用开发工程师，请认真回答我提出的Python相关的问题
# Human: 我的名字叫Fan，是一名程序员
# AI: 好的，Fan你好
# Human: 请问我的名字叫什么？
