from typing import Annotated, TypedDict

Age = Annotated[int, "年龄，范围0-150"]


class Person(TypedDict):
    name: str
    age: int
    age2: Age


p = Person(name="z3", age=111, age2=188)
print(p)

# p = Person(name="z3",age="1111")
# print(p)

"""
一、核心原因 1：Annotated 本身不具备运行时校验能力
typing.Annotated的设计目的并不是在程序运行时对数据进行合法性校验（比如范围、格式检查），它的核心作用是：
为类型添加元数据（附加描述信息）：你这里的"年龄，范围0-150"就是元数据，仅用于说明、文档生成、静态分析工具识别等场景，不会被 Python 解释器在运行时解析和执行校验逻辑。
保留原始类型特性：Annotated[int, "年龄，范围0-150"]本质上还是int类型，Python 运行时只会校验它是否是int类型（这里 188 是合法int），不会关心附加的元数据内容。
简单说：Annotated只是给类型 “加注释”，不是给类型 “加校验规则”。

二、核心原因 2：Python 的类型提示（Type Hints）是静态的、仅供参考的（装饰性）
"""
