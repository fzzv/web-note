"""
partial()方法可以格式化部分变量，并且继续返回一个模板，通常在部分提示词模板场景下使用
"""

from langchain_core.prompts import PromptTemplate

# 创建模板对象，定义提示词模板格式
# 模板包含两个占位符：role（角色）和 question（问题）
template = PromptTemplate.from_template(
    "你是一个专业的{role}工程师，请回答我的问题给出回答，我的问题是：{question}"
)

# 使用partial方法固定role参数为"python开发"
# 返回一个新的模板对象，其中role参数已被绑定
partial = template.partial(role="python开发")

# 打印partial对象及其类型信息
print(
    partial
)  # input_variables=['question'] input_types={} partial_variables={'role': 'python开发'} template='你是一个专业的{role}工程师，请回答我的问题给出回答，我的问题是：{question}'
print(type(partial))  # <class 'langchain_core.prompts.prompt.PromptTemplate'>
print()

# 使用format方法填充question参数，生成最终的提示词字符串
# 此时所有占位符都已填充完毕，返回完整的提示词文本
prompt = partial.format(question="冒泡排序怎么写？")

# 输出生成的提示词
print(
    prompt
)  # 你是一个专业的python开发工程师，请回答我的问题给出回答，我的问题是：冒泡排序怎么写
print(type(prompt))  # <class 'str'>
