from typing import Annotated
from typing_extensions import TypedDict
from langgraph.graph import StateGraph, START, END


def MyOperatorMul(current: float, update: float) -> float:
    """自定义乘法reducer，处理初始值为1.0"""
    # 如果是第一次调用，current会是默认值0.0
    if current == 0.0:
        # 对于乘法，恒等元应该是1.0或者 return 1.0 * update
        print(f"current:{current}")
        print(f"update:{update}")
        return 1.0 * update
    return current * update


class MultiplyState(TypedDict):
    factor: Annotated[float, MyOperatorMul]


def multiplier(state: MultiplyState) -> dict:
    return {"factor": 2.0}


def run_demo():
    print("使用自定义reducer解决乘法问题:")
    builder = StateGraph(MultiplyState)
    builder.add_node("multiplier", multiplier)
    builder.add_edge(START, "multiplier")
    builder.add_edge("multiplier", END)
    graph = builder.compile()

    result = graph.invoke({"factor": 5.0})
    print("初始状态: {{'factor': 5.0}}")
    print(f"执行结果: {result}")  # 应该是 {'factor': 10.0}
    print("解释: 5.0 * 2.0 = 10.0\n")


if __name__ == "__main__":
    run_demo()

# 使用自定义reducer解决乘法问题:
# current:0.0
# update:5.0
# 初始状态: {{'factor': 5.0}}
# 执行结果: {'factor': 10.0}
# 解释: 5.0 * 2.0 = 10.0
