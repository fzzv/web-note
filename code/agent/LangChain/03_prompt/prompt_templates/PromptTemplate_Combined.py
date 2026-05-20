"""
组合提示词模板

通过将多个子提示（Prompt）按一定逻辑顺序或层级组合起来，形成一个复杂任务的整体 Prompt。
例如实现多消息对话、多阶段任务、多输入源组合等场景。
尤其在，AI产品，你一言我一语，构建最后提示词，有用
"""

from langchain_core.prompts import PromptTemplate

# 创建一个PromptTemplate模板，用于生成介绍某个主题的提示词
# 该模板包含两个占位符：topic（主题）和length（字数限制）
# template1 = PromptTemplate.from_template("请用一句话介绍{topic}，要求通俗易懂,内容不超过{length}个字")
template1 = (
    PromptTemplate.from_template("请用一句话介绍{topic}，要求通俗易懂\n")
    + "内容不超过{length}个字"
)
# 使用format方法填充模板中的占位符，生成具体的提示词
prompt1 = template1.format(topic="LangChain", length=100)
print(prompt1)  # 请用一句话介绍LangChain，要求通俗易懂 内容不超过100个字

# 分别创建两个独立的PromptTemplate模板
prompt_a = PromptTemplate.from_template("请用一句话介绍{topic}，要求通俗易懂\n")
prompt_b = PromptTemplate.from_template("内容不超过{length}个字")
# 将两个模板进行拼接组合
prompt_all = prompt_a + prompt_b
# 填充组合后模板的占位符，生成最终的提示词
prompt2 = prompt_all.format(topic="LangChain", length=200)
print(prompt2)  # 请用一句话介绍LangChain，要求通俗易懂 内容不超过200个字
