"""
LangGraph 简单案例HelloWorld：
构建一个最小的有向图，流程是：START → 模型节点 → END

LangGraph的灵魂：State(状态) + Nodes(节点) + Edges(边) + Graph(图)
"""

import uuid
from typing import TypedDict, Annotated, List
from langgraph.graph import StateGraph, START, END
from langgraph.graph.message import add_messages
import os
from langchain.chat_models import init_chat_model
from langchain_core.messages import HumanMessage
from dotenv import load_dotenv

load_dotenv()


# ========== 1. 定义状态（State） ==========
# 存储对话消息
class AtguiguState(TypedDict):
    # messages 是一个消息列表，Annotated + add_messages 表示支持自动追加消息
    messages: Annotated[List, add_messages]


# ========== 2. 定义大模型 ==========
llm = init_chat_model(
    model="qwen-plus",
    model_provider="openai",
    api_key=os.getenv("QWEN_API_KEY"),
    base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",
)


# ========== 3. 定义节点函数 ==========
# 节点：调用大模型，并把回复加入到 state["messages"] 里
def model_node(state: AtguiguState):
    reply = llm.invoke(state["messages"])  # 输入历史消息，调用模型
    return {"messages": [reply]}  # 返回新消息，自动加到 state


# ========== 4. 构建图结构 ==========
graph = StateGraph(AtguiguState)  # 初始化图，指定 State 类型

graph.add_node("model", model_node)  # 添加一个节点，名字叫 "model"

graph.add_edge(START, "model")  # 从 START 到 "model"
graph.add_edge("model", END)  # 从 "model" 到 END
# 打印图的边和节点信息
# print(graph.edges)
print()
# print(graph.nodes)

# ========== 5. 编译==========
app = graph.compile()

# ========== 6. 运行 ==========
# result = app.invoke({"messages": [HumanMessage(content="请用一句话解释什么是 LangGraph。")]})
result = app.invoke({"messages": "请用一句话解释什么是 LangGraph。"})

# 打印模型的最后一条回复
print("模型回答：", result["messages"][-1].content)

print()
# =========================
# 1. 打印图的ascii可视化结构
print(app.get_graph().print_ascii())
print("=" * 50)

# 2. 打印图的Mermaid代码可视化结构并通过https://www.processon.com/mermaid编辑器查看
print(app.get_graph().draw_mermaid())
print("=" * 50)

# 3. 生成 PNG并写入文件
png_bytes = app.get_graph().draw_mermaid_png()
output_path = "langgraph" + str(uuid.uuid4())[:8] + ".png"
with open(output_path, "wb") as f:
    f.write(png_bytes)
print(f"图片已生成：{output_path}")
