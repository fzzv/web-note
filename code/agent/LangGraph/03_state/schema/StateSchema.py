"""
LangGraph 图输入输出模式和私有状态传递演示

该演示展示了：
1. 如何定义图的输入和输出模式
"""

from langgraph.graph import StateGraph, START, END
from typing_extensions import TypedDict


# 定义输入状态模式
class InputState(TypedDict):
    question: str


# 定义输出状态模式
class OutputState(TypedDict):
    answer: str


# 定义整体状态模式，结合输入和输出
class OverallState(InputState, OutputState):
    pass


# 定义处理节点
def answer_node(state: InputState):
    """
    处理输入并生成答案的节点
    Args:
        state: 输入状态
    Returns:
        dict: 包含答案的字典
    """
    print(f"执行 answer_node 节点:")
    print(f"  输入: {state}")

    # 示例答案
    answer = "再见" if "bye" in state["question"].lower() else "你好"
    result = {"answer": answer, "question": state["question"]}

    print(f"  输出: {result}")
    return result


def demo_input_output_schema():
    """演示输入输出模式"""
    print("=== 演示输入输出模式 ===")

    # 使用指定的输入和输出模式构建图
    builder = StateGraph(
        OverallState, input_schema=InputState, output_schema=OutputState
    )
    builder.add_edge(START, "answer_node")  # 定义起始边
    builder.add_node("answer_node", answer_node)  # 添加答案节点
    builder.add_edge("answer_node", END)  # 定义结束边
    graph = builder.compile()  # 编译图

    # 使用输入调用图并打印结果
    result = graph.invoke({"question": "你好"})
    print(f"图调用结果: {result}")
    # 打印图的ascii可视化结构
    print(graph.get_graph().print_ascii())
    print()


def main():
    """主函数"""
    print("=== LangGraph 图输入输出模式===\n")

    # 演示输入输出模式
    demo_input_output_schema()

    print("=== 演示完成 ===")


if __name__ == "__main__":
    main()
