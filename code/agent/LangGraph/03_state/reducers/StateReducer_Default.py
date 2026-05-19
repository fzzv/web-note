"""
如果未明确指定reducer函数，则默认对该键的更新是覆盖行为。
LangGraph Reducer函数演示 - 默认Reducer（覆盖更新）

直接覆盖：
如果没有为状态字段指定 Reducer，默认会覆盖更新。
也就是说，后执行的节点返回的值会直接覆盖先执行节点的值，
即下一个节点的State数据是上一个节点的返回。
"""

from typing import List
from typing_extensions import TypedDict
from langgraph.graph import StateGraph, START, END


# 1. 默认Reducer（覆盖更新）
# 未指定合并策略，默认覆盖，上一个节点的返回是下一个节点的值
class DefaultReducerState(TypedDict):
    foo: int
    bar: List[str]


def node_default_1(state: DefaultReducerState) -> dict:
    print(state["foo"])
    print(state["bar"])
    return {"foo": 22}


def node_default_2(state: DefaultReducerState) -> dict:
    print()
    print(state["foo"])
    print(state["bar"])
    return {"bar": ["bye1", "bye2", "bye3"]}


def main():
    print("1. 默认Reducer（覆盖更新）演示:\n")
    builder = StateGraph(DefaultReducerState)

    builder.add_node("node1", node_default_1)
    builder.add_node("node2", node_default_2)

    builder.add_edge(START, "node1")
    builder.add_edge("node1", "node2")
    builder.add_edge("node2", END)

    graph = builder.compile()

    result = graph.invoke(input={"foo": 1, "bar": ["hi"]})
    # print(f"初始状态: {{'foo': 1, 'bar': ['hi']}}")
    print(f"执行结果: {result}\n")


if __name__ == "__main__":
    main()

# 1. 默认Reducer（覆盖更新）演示:

# 1
# ['hi']

# 22
# ['hi']
# 执行结果: {'foo': 22, 'bar': ['bye1', 'bye2', 'bye3']}
