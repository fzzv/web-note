"""
StreamLLMTokens.py

使用messages流模式，从图中的任何部分（包括节点、工具、子图或任务）逐token流式传输大型语言模型（LLM）的输出。
messages模式的流式输出是一个元组(message_chunk, metadata)，其中：
- message_chunk：来自大语言模型（LLM）的令牌或消息片段。
- metadata：一个包含图节点和大语言模型调用详情的字典元数据。

"""

from typing import TypedDict
from langgraph.graph import StateGraph, START
from langchain.chat_models import init_chat_model
import os
from dotenv import load_dotenv

load_dotenv()


class State(TypedDict):
    query: str
    answer: str


def node(state: State):
    print("开始调用node节点")

    model = init_chat_model(
        model="qwen-plus",
        model_provider="openai",
        api_key=os.getenv("QWEN_API_KEY"),
        base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",
    )

    llm_result = model.invoke([("user", state["query"])])
    print("llm invoke结束", end="\n\n")

    return {"answer": llm_result}


def main():
    graph = (
        StateGraph(state_schema=State).add_node(node).add_edge(START, "node").compile()
    )

    inputs = {"query": "帮我生成一个100字的小学生作文，主题为我的一天"}

    # stream_mode="messages"从任何调用了大语言模型的图节点流式传输二元组（大语言模型token，元数据）。
    """messages模式的流式输出是一个元组(message_chunk, metadata)，其中：
    - message_chunk：来自大语言模型（LLM）的令牌或消息片段。
    - metadata：一个包含图节点和大语言模型调用详情的字典元数据。"""
    for chunk, meta_data in graph.stream(inputs, stream_mode="messages"):
        # print(f"type of chunk:{type(chunk)}")
        print(chunk.content, end="")
        # print(meta_data, end="")
        # print(chunk, end="")


if __name__ == "__main__":
    main()
