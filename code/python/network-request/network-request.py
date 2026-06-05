import requests

# response = requests.get("https://httpbin.org/get")
# print(response.status_code)

# response = requests.get("https://httpbin.org/get", params={"name": "Alice", "age": 25})
# print(response.status_code)
# print(response.json()["args"])

# response = requests.post(
#     "https://httpbin.org/post", data={"username": "alice", "password": "123"}
# )
# print(response.json()["form"])

# response = requests.post(
#     "https://httpbin.org/post", json={"name": "Alice", "items": [1, 2, 3]}
# )
# print(response.json()["json"])

# headers = {
#     "User-Agent": "MyApp/1.0",
#     "Accept": "application/json",
#     "Authorization": "Bearer your_token_here",
# }
# response = requests.get("https://httpbin.org/headers", headers=headers)
# print(response.json()["headers"]["User-Agent"])


# ===== 文件下载 =====
def download_file(url, filepath, chunk_size=512):
    """流式下载大文件"""
    with requests.get(url, stream=True) as r:
        r.raise_for_status()
        total = int(r.headers.get("Content-Length", 0))
        downloaded = 0
        with open(filepath, "wb") as f:
            for chunk in r.iter_content(chunk_size=chunk_size):
                f.write(chunk)
                downloaded += len(chunk)
                if total:
                    percent = downloaded / total * 100
                    print(f"\r下载进度：{percent:.1f}%", end="", flush=True)
        print("\n下载完成")


# download_file("https://httpbin.org/stream-bytes/10485760?chunk_size=1024", "file.zip")
