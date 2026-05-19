"""
LangGraph Reducer函数演示 - add_messages Reducer（消息列表专用）
"""

from typing import Annotated, List

from langchain_core.messages import HumanMessage, AIMessage
from typing_extensions import TypedDict
from langgraph.graph import StateGraph, START, END
from langgraph.graph.message import add_messages


# 2. add_messages Reducer（消息列表专用）
class AddMessagesState(TypedDict):
    """
    引入的 Annotated 类型，它允许给类型添加额外的元数据。
    messages: Annotated[List, add_messages]
    表示:
    - messages 我的状态里只有一个字段叫 messages，类型是是 List列表类型,
    - add_messages  这里的 add_messages 是一个函数，用于修改 messages 列表
                    每当节点返回对 messages 的“局部更新”时，
                    请用 add_messages 规约器把它合并到旧列表上（追加，而不是覆盖）
    总结：
    节点永远只 return 增量字典，不用手动把旧列表读出来再拼接。
    add_messages 在后台帮你完成“追加”动作；如果换成默认 reducer，旧消息会被整份替换掉
    """

    messages: Annotated[List, add_messages]


def chat_node_1(state: AddMessagesState) -> dict:
    return {"messages": [("assistant", "Hello from node 1")]}


def chat_node_2(state: AddMessagesState) -> dict:
    return {"messages": [("assistant", "Hello from node 2")]}


def run_demo():
    print("2. add_messages Reducer（消息列表专用）演示:")
    builder = StateGraph(AddMessagesState)
    builder.add_node("chat1", chat_node_1)
    builder.add_node("chat2", chat_node_2)

    builder.add_edge(START, "chat1")
    builder.add_edge(START, "chat2")  # 并行执行
    builder.add_edge("chat1", END)
    builder.add_edge("chat2", END)
    graph = builder.compile()

    result = graph.invoke({"messages": [("user", "Hi there!")]})
    print("初始状态: {{'messages': [('user', 'Hi there!')]}}")
    print(f"执行结果: {result}\n")

    print("*" * 60)

    # 打印图的ascii可视化结构
    print(graph.get_graph().print_ascii())


if __name__ == "__main__":
    run_demo()

# 2. add_messages Reducer（消息列表专用）演示:
# 初始状态: {{'messages': [('user', 'Hi there!')]}}
# 执行结果: {
#     'messages': [
#         HumanMessage(content='Hi there!', additional_kwargs={}, response_metadata={}, id='f029a089-9898-4b99-8882-1a01aac04cfc'),
#         AIMessage(content='Hello from node 1', additional_kwargs={}, response_metadata={}, id='abc8bef1-ca86-4b5b-912b-e00cc4b60cff', tool_calls=[], invalid_tool_calls=[]),
#         AIMessage(content='Hello from node 2', additional_kwargs={}, response_metadata={}, id='c9c12df7-bc15-4332-b5b6-15a248a9b524', tool_calls=[], invalid_tool_calls=[])
#     ]
# }
