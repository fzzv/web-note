import os
import json
import httpx
from typing import TypedDict

from langchain.agents import create_agent
from langchain_core.tools import tool
from langchain_openai import ChatOpenAI
from pydantic import SecretStr
from dotenv import load_dotenv

load_dotenv()


# 1.Tool 定义
@tool
def get_weather(loc: str):
    """
    查询即时天气函数

    :param loc: 必要参数，字符串类型，用于表示查询天气的具体城市名称。
                注意，中国的城市需要用对应城市的英文名称代替，例如如果需要查询北京市天气，
                则 loc 参数需要输入 'Beijing'/'shanghai'。
    :return: OpenWeather API 查询即时天气的结果。
    """
    url = "https://api.openweathermap.org/data/2.5/weather"
    params = {
        "q": loc,
        "appid": os.getenv("OPENWEATHER_API_KEY"),
        "units": "metric",
        "lang": "zh_cn",
    }
    response = httpx.get(url, params=params, timeout=30)
    data = response.json()
    # print(json.dumps(data, ensure_ascii=False, indent=2))
    return json.dumps(data, ensure_ascii=False)


# 2 结构化输出（推荐）
class WeatherCompareOutput(TypedDict):
    beijing_temp: float
    shanghai_temp: float
    hotter_city: str
    summary: str


# 3 模型（OpenAI Compatible）
api_key = os.getenv("QWEN_API_KEY")
model = ChatOpenAI(
    model="qwen-plus",
    api_key=SecretStr(api_key) if api_key else None,
    base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",
)


# 4 创建Agent
agent = create_agent(
    model=model,
    tools=[get_weather],
    system_prompt=(
        "你是天气助手。"
        "当用户询问多个城市天气时，"
        "你需要分别调用工具获取数据，并进行比较分析。"
    ),
    response_format=WeatherCompareOutput,
)

# 5 调用Agent
result = agent.invoke({"input": "请问今天北京和上海的天气怎么样，哪个城市更热？"})
print(result)

print()

print(json.dumps(result["structured_response"], ensure_ascii=False, indent=2))
