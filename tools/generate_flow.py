import sys
import pygraphviz as pgv
from graphviz2drawio import graphviz2drawio

def convert_dot_to_drawio(dot_path, output_drawio):
    try:
        # 加载 DOT 文件并用 Graphviz 自动布局
        G = pgv.AGraph(dot_path)
        G.layout(prog='dot')  # 可选 dot/neato/fdp/sfdp等
        dot_code = G.to_string()

        # 转换为 draw.io XML
        xml = graphviz2drawio.convert(dot_code)

        with open(output_drawio, "w", encoding="utf-8") as f:
            f.write(xml)

        print(f"✅ draw.io 文件生成成功: {output_drawio}")

    except Exception as e:
        print(f"❌ 转换失败: {e}")
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("❌ 用法: python generate_flow.py input.dot output.drawio")
        sys.exit(1)

    dot_path = sys.argv[1]
    output_drawio = sys.argv[2]
    convert_dot_to_drawio(dot_path, output_drawio)
