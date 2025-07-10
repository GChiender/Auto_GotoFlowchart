import drawpyo
from drawpyo.diagram_types import TreeDiagram, NodeObject

def create_simple_tree(output_path="layout.drawio"):
    tree = TreeDiagram(
        file_name=output_path,
        direction="down",
        level_spacing=80,
        item_spacing=20
    )

    root = NodeObject(tree=tree, value="开始")
    step1 = NodeObject(tree=tree, value="赋值", tree_parent=root)
    cond = NodeObject(tree=tree, value="判断", tree_parent=step1)
    step2 = NodeObject(tree=tree, value="循环体", tree_parent=cond)
    end = NodeObject(tree=tree, value="结束", tree_parent=step2)

    tree.write()
    print(f"✅ 自动布局生成完成: {output_path}")

if __name__ == "__main__":
    create_simple_tree()
