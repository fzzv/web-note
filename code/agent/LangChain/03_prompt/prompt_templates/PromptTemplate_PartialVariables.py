# 方式3：部分提示词模板(partial_variables),实例化过程中指定 partial_variables 参数
from langchain_core.prompts import PromptTemplate
from datetime import datetime
import time

# 1 实例化过程中指定 partial_variables 参数
# 创建一个包含时间变量的模板，时间变量使用partial_variables预设为当前时间,然后格式化问题生成最终提示词
template1 = PromptTemplate.from_template(
    "现在时间是：{time},请对我的问题给出答案，我的问题是：{question}",
    partial_variables={"time": datetime.now().strftime("%Y-%m-%d %H:%M:%S")}
)

prompt1 = template1.format(question="今天是几号？")
print(prompt1)

time.sleep(2)  # 程序暂停 2 秒，期间不执行任何代码

# 2 使用 partial 方法指定默认值
template2 = PromptTemplate.from_template("现在时间是：{time},请对我的问题给出答案，我的问题是：{question}")
# 使用 partial 方法指定默认值
partial = template2.partial(time=datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
prompt2 = partial.format(question="今天是几号？")
print(prompt2)


template3 = PromptTemplate(
    template="{foo} {bar}",
    input_variables=["foo", "bar"],
    partial_variables={"foo": "hello"},  # 预先定义部分变量foo值为hello
)

prompt = template3.format(foo="li4",bar="world")
print(prompt)  # li4 world

prompt = template3.format(bar="world")
print(prompt)  # hello world