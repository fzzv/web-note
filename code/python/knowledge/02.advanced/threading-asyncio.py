from concurrent.futures import ThreadPoolExecutor
from pathlib import Path
import hashlib


def compute_hash(filepath):
    """计算文件的 MD5 哈希"""
    h = hashlib.md5()
    with open(filepath, "rb") as f:
        for chunk in iter(lambda: f.read(8192), b""):
            h.update(chunk)
    return filepath, h.hexdigest()


def batch_hash(directory, pattern="*"):
    """批量计算目录中文件的哈希"""
    files = [f for f in Path(directory).rglob(pattern) if f.is_file()]
    if not files:
        print("没有找到文件")
        return {}

    results = {}
    with ThreadPoolExecutor(max_workers=8) as executor:
        futures = {executor.submit(compute_hash, f): f for f in files}
        for future in futures:
            filepath, hash_value = future.result()
            results[str(filepath)] = hash_value
            print(f"  {filepath.name}: {hash_value[:16]}...")

    return results


hashes = batch_hash(".", "*.py")
print(f"\n共处理 {len(hashes)} 个文件")
