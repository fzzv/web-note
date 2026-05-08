"""
from_messages
作用：将模板变量替换后，直接生成消息列表（List[BaseMessage]），
一般包含：SystemMessage``HumanMessage``AIMessage
常用场景：用于手动查看或调试 Prompt 的最终“消息结构”或者自己拼接进 Chain。

实例化时需要传入messages: Sequence[MessageLikeRepresentation]
messages 参数支持如下格式：
- tuple 构成的列表，格式为[(role, content)]
template = ChatPromptTemplate(
    [
        ("system", "你是一个AI开发工程师，你的名字是{name}。"),
        ("human", "你能帮我做什么?"),
        ("ai", "我能开发很多{thing}。"),
        ("human", "{user_input}"),
    ]
)
- dict 构成的列表，格式为[{“role”:... , “content”:...}]
chat_prompt = ChatPromptTemplate(
    [
        {"role": "system", "content": "你是AI助手，你的名字叫{name}。"},
        {"role": "user", "content": "请问：{question}"}
    ]
)
- Message 类构成的列表
"""

import os
from langchain.chat_models import init_chat_model
from langchain_core.prompts import ChatPromptTemplate
from dotenv import load_dotenv

load_dotenv()  # 加载环境变量

# 创建聊天提示模板，包含系统角色设定和用户问题格式
# 系统消息定义了AI助手的角色，人类消息定义了用户问题的格式
chat_prompt = ChatPromptTemplate.from_messages(
    [("system", "你是一个{role}，请回答我提出的问题"), ("human", "请回答:{question}")]
)

# 格式化聊天提示模板，填充角色和问题参数
# 参数role: 指定AI助手的角色身份
# 参数question: 用户提出的具体问题
# 返回值: 格式化后的消息列表
# prompt_value = chat_prompt.format_messages(role="python开发工程师", question="冒泡排序怎么写")
prompt_value = chat_prompt.format_messages(
    **{"role": "python开发工程师", "question": "堆排序怎么写"}
)
# 打印格式化后的提示消息
print(prompt_value)

print()
# 使用指定的角色和问题参数填充模板，生成具体的提示内容
# role: 指定AI扮演的角色
# question: 用户提出的具体问题
prompt_value2 = chat_prompt.invoke(
    {"role": "python开发工程师", "question": "堆排序怎么写"}
)
# 输出生成的提示内容
print(prompt_value2.to_string())

print()

prompt_value3 = chat_prompt.format(
    **{"role": "python开发工程师", "question": "快速排序怎么写"}
)
# 输出生成的提示内容
print(prompt_value3)


llm = init_chat_model(
    model="qwen-plus",
    model_provider="openai",
    api_key=os.getenv("QWEN_API_KEY"),
    base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",
)
print()
print("======================")

result = llm.invoke(prompt_value)
print(result)
print(result.content)
