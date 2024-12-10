import re

xxlist = [[
    "",
    "广州卢腾贸易有限公司\n地址： 广州市白云区同和街同沙路283号之二2009铺TJ2415",
    "在业",
    "张旭",
    "2500万元人民币",
    "2022-11-06",
        "91440111MAC1KPXFX9",
        "广州市白云区同和街同沙路283号之二2009铺TJ2415",
        "-",
        "-",
        "笔记关注"
    ]]


for data in xxlist:
    text=data[1]
    
    
    # 提取公司名称和地址的正则表达式
    pattern = r"(.*?公司)\n地址：\s*(.*)"
    
    # 使用正则表达式匹配
    match = re.search(pattern, text)
    
    if match:
        company_name = match.group(1)  # 提取公司名称
        address = match.group(2)       # 提取地址
        print(f"公司名称: {company_name}")
        print(f"地址: {address}")
    else:
        print("未能匹配到公司名称和地址")

