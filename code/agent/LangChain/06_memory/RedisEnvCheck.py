# pip install redis==5.3.1

# 尝试导入 redis 包
import redis

# 验证包版本（无报错即为导入成功）
print(redis.__version__)

# 极简 redis 导入测试脚本
try:
    # 导入 redis 包
    import redis
    print("✅ redis 包导入成功！")
    print(f"✅ redis 包版本：{redis.__version__}")
except ModuleNotFoundError:
    print("❌ 未找到 redis 包，请先安装！")
except Exception as e:
    print(f"❌ redis 包导入异常：{e}")

"""
5.3.1
✅ redis 包导入成功！
✅ redis 包版本：5.3.1
"""