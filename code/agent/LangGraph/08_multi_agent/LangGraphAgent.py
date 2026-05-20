import os
from langchain.chat_models import init_chat_model
from langchain_core.messages import HumanMessage
from langchain.agents import create_agent
from dotenv import load_dotenv

load_dotenv()


# ========== 1. 获取指定城市的天气信息 ==========
def get_weather(city: str) -> str:
    """获取指定城市的天气信息。
    Args:
        city: 城市名称
    Returns:
        返回该城市的天气描述
    """
    return f"今天{city}是晴天，仅做测试，固定写死"


# ========== 2. 定义大模型 ==========
llm = init_chat_model(
    model="qwen-plus",
    model_provider="openai",
    api_key=os.getenv("QWEN_API_KEY"),
    base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",
)

# 使用LangGraph提供的API创建Agent
agent = create_agent(
    model=llm,  # 添加模型
    tools=[get_weather],  # 添加工具
)
print("agent底层本质是个什么对象: " + str(type(agent)))
human_message = HumanMessage(content="今天深圳天气怎么样？")
response = agent.invoke({"messages": [human_message]})
# print(response)
print()
print("模型回答：", response["messages"][-1].content)
print()
response["messages"][-1].pretty_print()


# 使用stream方法进行流式调用
# 这里stream_mode参数有四种选项：
# - messages：流式输出大语言模型回复的token
# - updates : 流式输出每个工具调用的每个步骤。
# - values : 一次输出到所有的chunk。默认值。
# - custom : 自定义输出。主要是可以在工具内部使用get_stream_writer获取输入流，添加自定义的内容。
# for chunk in agent.stream(
#         {"messages": [{"role": "user", "content": "请问北京今天天气如何？"}]},
#         stream_mode="values",
# ):
#     chunk["messages"][-1].pretty_print()
