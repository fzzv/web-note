import os
from langchain.chat_models import init_chat_model
from langgraph.checkpoint.memory import InMemorySaver
from langchain.agents import create_agent
from dotenv import load_dotenv

load_dotenv()

# ==========定义大模型 ==========
llm = init_chat_model(
    model="qwen-plus",
    model_provider="openai",
    api_key=os.getenv("QWEN_API_KEY"),
    temperature=0.0,
    base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",
)

# 定义短期记忆使用内存（生产可以换 RedisSaver/PostgresSaver）
checkpointer = InMemorySaver()
agent = create_agent(model=llm, checkpointer=checkpointer)
# 多轮对话配置，同一 thread_id 即同一会话
config = {"configurable": {"thread_id": "user-001"}}

msg1 = agent.invoke(
    {"messages": [("user", "你好，我叫张三，喜欢足球，60字内简洁回复")]}, config
)
msg1["messages"][-1].pretty_print()

# 6. 第二轮（继续同一 thread）
msg2 = agent.invoke({"messages": [("user", "我叫什么？我喜欢做什么？")]}, config)
msg2["messages"][-1].pretty_print()
