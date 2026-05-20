"""
"placeholder" 是 ("placeholder", "{memory}") 的简写语法，
等价于 MessagesPlaceholder("memory")。

隐式使用MessagesPlaceholder
"""

from langchain_core.messages import HumanMessage, AIMessage
from langchain_core.prompts import ChatPromptTemplate

# 使用 ChatPromptTemplate 构建一个多角色对话提示模板
prompt = ChatPromptTemplate.from_messages(
    [
        # 占位符，用于插入对话“记忆”内容，即之前的聊天记录（历史上下文）
        ("placeholder", "{memory}"),
        # 系统消息，用于设定 AI 的角色 —— 是一个资深的 Python 应用开发工程师
        (
            "system",
            "你是一个资深的Python应用开发工程师，请认真回答我提出的Python相关的问题",
        ),
        # 用户当前提问，使用变量 {question} 进行动态填充
        ("human", "{question}"),
    ]
)

# 使用 invoke 方法传入上下文变量，生成格式化后的对话 prompt 内容
prompt_value = prompt.invoke(
    {
        # memory：是之前的对话上下文，会被插入到 {memory} 的位置
        "memory": [
            # 用户第一轮对话
            HumanMessage("我的名字叫亮仔，是一名程序员"),
            # AI 第一轮回答
            AIMessage("好的，亮仔你好"),
        ],
        # 当前的问题，将替换模板中的 {question}
        "question": "请问我的名字叫什么？",
    }
)
# 使用 .to_string() 将格式化后的对话链转换成纯文本字符串，方便查看输出
print(prompt_value.to_string())
# Human: 我的名字叫亮仔，是一名程序员
# AI: 好的，亮仔你好
# System: 你是一个资深的Python应用开发工程师，请认真回答我提出的Python相关的问题
# Human: 请问我的名字叫什么？
