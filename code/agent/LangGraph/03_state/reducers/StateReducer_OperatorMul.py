"""
LangGraph Reducer函数演示 - operator.mul Reducer（数值相乘）
"""

import operator
from typing import Annotated
from typing_extensions import TypedDict
from langgraph.graph import StateGraph, START, END


# 4. operator.mul Reducer（数值相乘）
class MultiplyState(TypedDict):
    factor: Annotated[float, operator.mul]


def multiplier(state: MultiplyState) -> dict:
    return {"factor": 2.0}


def run_demo():
    """
    这不是bug，是设计决策：LangGraph选择用类型默认值初始化状态字段
    对于不同操作，需要不同处理：
    加法：恒等元是0.0，所以operator.add可以直接用
    乘法：恒等元是1.0，需要特殊处理初始的0.0
    自定义reducer是标准做法：复杂的业务逻辑都应该使用自定义reducer

    这是LangGraph使用中的一个常见陷阱！建议在使用乘法、除法等非加法操作时，总是使用自定义reducer来处理初始值问题

        在执行初始阶段（我们定义的第一个node前），会默认调用一次reducer（后面自定义reducer案例中进行了打印验证），
        用默认值与invoke传递的值进行计算：
        此案例中，invoke中传递了一个默认值5.0，由于会默认调用一次reducer，
        执行的计算是： 0.0（float默认值） * 5.0(invoke传递的初始值) = 0.0
        导致后续乘法结果一直都是0

        初始默认值: factor = 0.0
        invoke传入: factor = 5.0
        reducer计算: 0.0 * 5.0 = 0.0
        然后才执行你的multiplier节点...
    operator.mul作为 LangGraph 归约器的执行逻辑是：
    最终值 = 初始值 * 增量值1 * 增量值2 * ...
    归约器会迭代节点返回的增量值并依次相乘，若直接返回单个数值（非可迭代），会被判定为「无增量」，最终按初始值 * 空处理，而乘法中「空累积」的默认结果是乘法单位元 0.0。

        解决方案： 使用自定义reducer
    """
    print("4. operator.mul Reducer（数值相乘）演示:")
    builder = StateGraph(MultiplyState)
    builder.add_node("multiplier", multiplier)
    builder.add_edge(START, "multiplier")
    builder.add_edge("multiplier", END)
    graph = builder.compile()

    result = graph.invoke({"factor": 5.0})
    print(f"初始状态: {{'factor': 5.0}}")
    print(f"执行结果: {result}\n")


if __name__ == "__main__":
    run_demo()

# 4. operator.mul Reducer（数值相乘）演示:
# 初始状态: {'factor': 5.0}
# 执行结果: {'factor': 0.0}
