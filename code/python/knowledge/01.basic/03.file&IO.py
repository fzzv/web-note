with open("test.txt", "w", encoding="utf-8") as f:
    f.write("Hello")

with open("demo.txt", "w", encoding="utf-8") as f:
    f.write("第一行\n第二行\n第三行\n")
with open("demo.txt", "r", encoding="utf-8") as f:
    content = f.read()
    print(content)

with open("demo.txt", "r", encoding="utf-8") as f:
    line1 = f.readline()
    line2 = f.readline()
    line3 = f.readline()
    print(repr(line1))
    print(repr(line2))
    print(repr(line3))

with open("demo.txt", "r", encoding="utf-8") as f:
    line = f.readlines()
    print(line)
