from typing import Annotated, List
from typing_extensions import TypedDict
from langgraph.graph import StateGraph, START, END
from langgraph.graph.message import add_messages
import operator


class ChatState(TypedDict):
    messages: Annotated[list, add_messages]  # 消息历史
    tags: Annotated[List[str], operator.add]  # 标签列表
    score: Annotated[float, operator.add]  # 累计分数


def process_user_message(state: ChatState) -> dict:
    user_message = state["messages"][-1]  # 获取最新消息
    return {
        "messages": [("assistant", f"Echo: {user_message.content}")],
        "tags": ["processed"],
        "score": 1.0,
    }


def add_sentiment_tag(state: ChatState) -> dict:
    return {"tags": ["positive"], "score": 0.5}


# 构建图
builder = StateGraph(ChatState)
builder.add_node("process", process_user_message)
builder.add_node("sentiment", add_sentiment_tag)

builder.add_edge(START, "process")
builder.add_edge(START, "sentiment")
builder.add_edge("process", END)
builder.add_edge("sentiment", END)

graph = builder.compile()

# 使用示例 -使用正确的消息格式
result = graph.invoke(
    {
        "messages": [{"role": "user", "content": "Hello, how are you?"}],
        "tags": ["greeting"],
        "score": 0.0,
    }
)

print(result)

# {
#     "messages": [
#         HumanMessage(
#             content="Hello, how are you?",
#             additional_kwargs={},
#             response_metadata={},
#             id="259efc74-5e90-4c01-a703-700e118b922e",
#         ),
#         AIMessage(
#             content="Echo: Hello, how are you?",
#             additional_kwargs={},
#             response_metadata={},
#             id="662e4e7f-d8a2-42c0-879e-c88d62cfc47d",
#             tool_calls=[],
#             invalid_tool_calls=[],
#         ),
#     ],
#     "tags": ["greeting", "processed", "positive"],
#     "score": 1.5,
# }
