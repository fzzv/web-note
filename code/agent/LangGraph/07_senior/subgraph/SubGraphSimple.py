"""将子图作为节点添加到父图"""

from langgraph.graph import StateGraph, START, END
from typing import TypedDict


# 1. 状态定义（统一字段名，避免执行时KeyError）
class ParentState(TypedDict):
    parent_messages: list  # 与子图共享数据


class SubgraphState(TypedDict):
    parent_messages: list  # 与父图共享的数据
    sub_message: str  # 子图私有数据


# 2. 定义子图节点函数
def subgraph_node(state: SubgraphState) -> SubgraphState:
    """子图节点处理逻辑：修改共享数据+设置私有数据"""
    # 向共享的parent_messages中添加内容
    state["parent_messages"].append("message from subgraph updateO(∩_∩)O")
    # 设置子图私有数据
    state["sub_message"] = "subgraph private message"
    # print(state["sub_message"])
    return state


# 3. 定义父图节点函数
def parent_node(state: ParentState) -> ParentState:
    """父图初始节点：初始化共享数据"""
    if not state.get("parent_messages"):
        state["parent_messages"] = []
    state["parent_messages"].append("message from 父亲 node")
    return state


# 4. 构建子图
def build_subgraph() -> StateGraph:
    """构建并返回编译后的子图"""
    sub_builder = StateGraph(SubgraphState)
    sub_builder.add_node("sub_node", subgraph_node)
    sub_builder.add_edge(START, "sub_node")
    sub_builder.add_edge("sub_node", END)  # 子图执行完指向结束
    return sub_builder.compile()


# 5. 构建父图
def build_parent_graph(compiled_subgraph) -> StateGraph:
    """构建并返回编译后的父图"""
    builder = StateGraph(ParentState)
    # 添加父图初始节点
    builder.add_node("parent_node", parent_node)
    # 将子图作为节点添加到父图，添加子图添加为父图的节点
    builder.add_node("subgraph_node", compiled_subgraph)
    # 父图执行流程：START -> parent_node -> subgraph_node -> END
    builder.add_edge(START, "parent_node")
    builder.add_edge("parent_node", "subgraph_node")  # 将子图作为节点添加到父图
    builder.add_edge("subgraph_node", END)
    return builder.compile()


# 6. 主方法（程序入口）
def main():
    """主函数：执行父图并输出结果"""
    # 构建子图
    compiled_subgraph = build_subgraph()
    # 构建父图
    parent_graph = build_parent_graph(compiled_subgraph)

    # 执行父图，先初始
    initial_state = {"parent_messages": ["我是父消息"]}
    print("初始状态：", initial_state)

    # 执行父图并获取最终状态
    # 父图执行时会自动调用子图，子图可修改共享的parent_messages，
    # 私有sub_message仅在子图内有效（父图最终状态不会显示，因为父图状态定义中无该字段）
    final_state = parent_graph.invoke(initial_state)
    print("\n执行后最终状态：", final_state)


# 程序入口
if __name__ == "__main__":
    main()
