import os
import subprocess
import sys

# 项目根目录
PROJECT_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))

# 输入输出目录
TARGET_DIR = os.path.join(PROJECT_ROOT, "target")
OUTPUT_DIR = os.path.join(PROJECT_ROOT, "output")
CONVERTER = os.path.join(PROJECT_ROOT, "tools", "convert_dot_to_drawio.py")

def run_go_parser(input_file, dot_file):
    result = subprocess.run(["go", "run", "main.go"], cwd=PROJECT_ROOT)
    return result.returncode == 0

def convert_to_drawio(dot_file, drawio_file):
    result = subprocess.run(["python3", CONVERTER, "--input", dot_file, "--output", drawio_file])
    return result.returncode == 0

def open_drawio(drawio_file):
    try:
        subprocess.Popen(["open", drawio_file])  # macOS
    except Exception as e:
        print("❌ 无法打开文件:", e)

def main():
    os.makedirs(OUTPUT_DIR, exist_ok=True)

    for filename in os.listdir(TARGET_DIR):
        if filename.endswith(".go"):
            base = os.path.splitext(filename)[0]
            input_go = os.path.join(TARGET_DIR, filename)
            dot_dir = os.path.join(OUTPUT_DIR, "dot")
            drawio_dir = os.path.join(OUTPUT_DIR, "drawio")
            os.makedirs(dot_dir, exist_ok=True)
            os.makedirs(drawio_dir, exist_ok=True)
            output_dot = os.path.join(dot_dir, f"{base}.dot")
            output_drawio = os.path.join(drawio_dir, f"{base}.drawio")

            print(f"🔧 处理 {filename}...")

            # 设置环境变量：让 main.go 知道当前目标文件
            os.environ["GO_INPUT_FILE"] = input_go
            if not run_go_parser(input_go, output_dot):
                print(f"❌ 生成 DOT 失败: {filename}")
                continue

            if not convert_to_drawio(output_dot, output_drawio):
                print(f"❌ 转换 drawio 失败: {output_dot}")
                continue

            open_drawio(output_drawio)

    print("✅ 批量处理完成！")

if __name__ == "__main__":
    main()
