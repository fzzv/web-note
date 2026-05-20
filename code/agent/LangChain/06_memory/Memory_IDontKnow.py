from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import PromptTemplate
from langchain.chat_models import init_chat_model
import os
from dotenv import load_dotenv

load_dotenv()

# 设置本地模型
llm = init_chat_model(
    model="mimo-v2.5-pro",
    model_provider="openai",
    api_key=os.getenv("XIAOMI_API_KEY"),
    temperature=0.0,
    base_url="https://token-plan-cn.xiaomimimo.com/v1",
)

prompt = PromptTemplate.from_template("请回答我的问题：{question}")
# 创建字符串输出解析器
parser = StrOutputParser()

# 构建链式调用
chain = prompt | llm | parser

# 执行链式调用
print(chain.invoke({"question": "我叫张三，你叫什么?"}))

print(chain.invoke({"question": "你知道我是谁吗?"}))

"""
你好，张三！我叫通义千问（Qwen），是阿里云研发的超大规模语言模型。很高兴认识你！😊
我不知道你是谁。我是一个AI助手，没有能力识别或获取用户的身份信息。如果你有任何问题需要帮助，我会很乐意为你提供支持！

解释：
我们刚刚在本地程序，前一轮对话告诉大语言模型的信息，下一轮就被“遗忘了”。
但如果我们使用使用 Qwen 聊天时，它能记住多轮对话中的内容，这Qwen网页版实现了历史记忆功能
"""
