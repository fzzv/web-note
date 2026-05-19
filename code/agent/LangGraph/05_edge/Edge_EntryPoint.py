"""
LangGraph入口点演示

入口点定义了图开始执行的第一个节点。
"""

from typing_extensions import TypedDict
from langgraph.graph import StateGraph, START, END


# 定义状态
class AtguiguState(TypedDict):
    value: int
    step: str


# 定义节点函数
def node_a(state: AtguiguState) -> dict:
    """节点A"""
    print("执行节点A")
    print("state[value]:" + str(state["value"]))
    print("state[step]:" + str(state["step"]))
    return {"value": state["value"] + 1, "step": "A执行完毕"}


def node_b(state: AtguiguState) -> dict:
    """节点B"""
    print("执行节点B")
    return {"value": state["value"] * 2, "step": "B执行完毕"}


def main():
    """演示入口点"""
    print("=== 入口点演示 ===")

    # 创建图
    builder = StateGraph(AtguiguState)

    # 添加节点
    builder.add_node("node_a", node_a)
    builder.add_node("node_b", node_b)

    """
    set_entry_point(node_id) 和 set_finish_point(node_id) 是 LangGraph 为「图对象」提供的配置方法，
    核心作用是将你自定义的业务节点，和内置的 START/END 特殊节点做 “自动绑定”，简化图的入口 / 出口边的定义，
    本质是语法糖（底层还是帮你执行了 add_edge(START, 入口节点) / add_edge(出口节点, END)）

    set_entry_point(node_id)
        图的实际执行入口是 node_id 这个自定义节点，底层会自动创建一条边 add_edge(START, node_id)，
        无需你手动写这条边
    set_finish_point(node_id)
        当流程走到 node_id 这个自定义节点时视为流程结束，底层会自动创建一条边 add_edge(node_id, END)，
        无需你手动写这条边
    """
    builder.set_entry_point("node_a")
    builder.add_edge("node_a", "node_b")
    builder.set_finish_point("node_b")

    # 编译图
    graph = builder.compile()
    # 执行图
    result = graph.invoke({"value": 0, "step": "hello"})
    print(f"执行结果: {result}\n")

    print()
    # 打印图的ascii可视化结构
    print(graph.get_graph().print_ascii())
    print("=================================")
    print()
    # 打印图的可视化结构，生成更加美观的Mermaid 代码，通过processon 编辑器查看
    print(graph.get_graph().draw_mermaid())


if __name__ == "__main__":
    main()
