import UnityPy
import os
import sys
import json

def extract(input_path, output_path):
  try:
    env = UnityPy.load(input_path)

    for obj in env.objects:
      if obj.type.name in ["Texture2D"]:
        container_path = obj.container if hasattr(obj, "container") and obj.container else ""
        parts = container_path.replace('\\', '/').split("/")
        if len(parts) >= 2:
          del parts[-2]
        fix_container_path = os.path.join(output_path, *parts[:-1])
        os.makedirs(fix_container_path, exist_ok=True)
        
        data = obj.read()
        img = data.image
        des = os.path.join(fix_container_path, parts[-1].capitalize())

        img.save(des)
        print(f"Extracted: {str(parts[-1])}")
      
      if obj.type.name == "TextAsset":
        container_path = obj.container if hasattr(obj, "container") and obj.container else ""
        parts = container_path.replace('\\', '/').split("/")
        if len(parts) >= 2:
          del parts[-2]
        fix_container_path = os.path.join(output_path, *parts[:-1])
        os.makedirs(fix_container_path, exist_ok=True)

        data = obj.parse_as_object()
        text = data.m_Script
        des = os.path.join(fix_container_path, parts[-1].capitalize())
        
        if des.lower().endswith(".json"):
          try:
            parsed = json.loads(text)
            with open(des, "w", encoding="utf-8") as f:
              json.dump(parsed, f, ensure_ascii=False, indent=2)
          except json.JSONDecodeError:
            with open(des, "w", encoding="utf-8") as f:
              f.write(text)
        else:
          cleaned_text = "\n".join(
            line for line in text.splitlines()
            if line.strip()
          )   
          with open(des, "w", encoding="utf-8") as f:
            f.write(cleaned_text)
  except Exception as e:
    print(f"(Error) Failed to handle Unity file: {str(e)}")

inputs = sys.argv[1:]
output_folder = "./extracted"

for file in inputs:
  extract(file, output_folder)