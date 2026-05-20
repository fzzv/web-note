"""
invoke() 是 LangChain Expression Language（LCEL 的统一执行入口，用于执行任意可运行对象（Runnable ）。返回的是一个 PromptValue 对象，
可以用 .to_string() 或 .to_messages() 查看内容
"""

from langchain_core.prompts import PromptTemplate

# 创建一个PromptTemplate对象，用于生成格式化的提示词模板
# 模板中包含两个占位符：{role}表示角色，{question}表示问题
template = PromptTemplate.from_template(
    "你是一个专业的{role}工程师，请回答我的问题给出回答，我的问题是：{question}"
)

# 使用invoke方法填充模板中的占位符，生成具体的提示词
# 参数：字典类型，包含role和question两个键值对
# 返回值：PromptValue对象，包含了格式化后的提示词
prompt = template.invoke({"role": "python开发", "question": "冒泡排序怎么写？"})

# 打印PromptValue对象及其类型
print(
    prompt
)  # text='你是一个专业的python开发工程师，请回答我的问题给出回答，我的问题是：冒泡排序怎么写？'
print(type(prompt))  # <class 'langchain_core.prompt_values.StringPromptValue'>
print()

# 将PromptValue对象转换为字符串并打印
# to_string()方法将PromptValue转换为可读的字符串格式
print(
    prompt.to_string()
)  # 你是一个专业的python开发工程师，请回答我的问题给出回答，我的问题是：冒泡排序怎么写？
print(type(prompt.to_string()))  # <class 'str'>
print()

print(
    prompt.to_messages()
)  # [HumanMessage(content='你是一个专业的python开发工程师，请回答我的问题给出回答，我的问题是：冒泡排序怎么写？', additional_kwargs={}, response_metadata={})]
print(type(prompt.to_messages()))  # <class 'list'>
