import pandas as pd
import json
import sys
 
df = pd.read_csv(sys.argv[1])

list = (
  df.loc[df["AssetBundleName"].str.contains("chara/still/", na=False), "AssetBundleName"]
    .str.extract(r"(\d+)$")[0]
    .to_list()
)

with open("list.json", "w", encoding="utf-8") as f:
    json.dump(list, f)