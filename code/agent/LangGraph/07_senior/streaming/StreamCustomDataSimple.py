"""
StreamCustomDataSimple.py

要从LangGraph节点或工具内部发送自定义用户定义数据，请遵循以下步骤：
- 使用get_stream_writer访问流写入器并发送自定义数据。
- 调用.stream()或.astream()时，设置stream_mode="custom"以在流中获取自定义数据。
你可以组合多种模式（例如["updates", "custom"]），但至少有一种模式必须是"custom"。

LangGraph 自定义数据流式传输演示
展示如何从节点内部发送自定义用户定义数据
"""

from typing import TypedDict
from langgraph.config import get_stream_writer
from langgraph.graph import StateGraph, START, END


class State(TypedDict):
    query: str
    answer: str


def node(state: State):
    # Get the stream writer to send custom data
    writer = get_stream_writer()
    # Emit a custom key-value pair (e.g., progress update)
    writer({"custom_key": "学习Agent，O(∩_∩)O"})
    return {"answer": "some data"}


graph = (
    StateGraph(State)
    .add_node(node)
    .add_edge(START, "node")
    .add_edge("node", END)
    .compile()
)

# Set stream_mode="custom" to receive the custom data in the stream
# for chunk in graph.stream({"query": "example"}, stream_mode=["custom"]):
#     print(chunk)
#
# for chunk in graph.stream({"query": "example"}, stream_mode=["updates", "custom"]):
#     print(chunk)
#
for chunk in graph.stream({"query": "example"}, stream_mode=["values", "custom"]):
    print(chunk)
