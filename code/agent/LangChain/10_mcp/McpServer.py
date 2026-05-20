import json
import os
import httpx
from mcp.server.fastmcp import FastMCP
from loguru import logger
from dotenv import load_dotenv

load_dotenv()


#  该案例跑在python3.12.7及以下
# 一个MCPServer对外暴露多个，一系列TOOLCalling工具类的集合


# 创建FastMCP实例，用于启动天气服务器SSE服务
mcp = FastMCP("WeatherServerSSE", host="127.0.0.1", port=8000)


@mcp.tool()
def get_weather(city: str) -> str:
    """
    查询指定城市的即时天气信息。
    参数 city: 城市英文名，如 Beijing
    返回: OpenWeather API 的 JSON 字符串
    """
    url = "https://api.openweathermap.org/data/2.5/weather"
    params = {
        "q": city,
        "appid": os.getenv("OPENWEATHER_API_KEY"),
        "units": "metric",
        "lang": "zh_cn",
    }
    resp = httpx.get(url, params=params, timeout=10)
    data = resp.json()
    logger.info(f"查询 {city} 天气结果：{data}")
    return json.dumps(data, ensure_ascii=False)


if __name__ == "__main__":
    logger.info("启动 MCP SSE 天气服务器，监听 http://127.0.0.1:8000/sse")
    # 运行MCP客户端，使用Server-Sent Events(SSE)作为传输协议
    mcp.run(transport="sse")
    # mcp.run(transport="stdio")


"""
核心重点：202 Accepted 状态码的意义（结合 MCP SSE 场景）
HTTP 202 状态码和常见的200 OK有本质区别，适配 MCP SSE 的流式处理特性：
200 OK：请求已处理完成，服务端立即返回最终结果（适合一次性请求 - 响应的场景，比如普通接口查询）；
202 Accepted：请求已接收并受理，服务端会在后台处理（比如调用天气工具、执行 MCP 指令），
处理完成后通过 SSE 流将结果推送给客户端（适合耗时 / 流式处理的场景，这正是 MCP SSE 服务的设计逻辑）。
"""
