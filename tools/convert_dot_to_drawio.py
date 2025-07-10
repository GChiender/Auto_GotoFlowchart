from graphviz2drawio import graphviz2drawio
import argparse

def convert_dot_to_drawio(dot_path, output_drawio):
    with open(dot_path, "r", encoding="utf-8") as f:
        dot_content = f.read()

    xml_content = graphviz2drawio.convert(dot_content)

    with open(output_drawio, "w", encoding="utf-8") as f:
        f.write(xml_content)

    print(f"✅ draw.io 文件生成成功：{output_drawio}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Convert .dot to .drawio")
    parser.add_argument("--input", "-i", required=True, help="Path to input .dot file")
    parser.add_argument("--output", "-o", required=True, help="Path to output .drawio file")
    args = parser.parse_args()

    convert_dot_to_drawio(args.input, args.output)
