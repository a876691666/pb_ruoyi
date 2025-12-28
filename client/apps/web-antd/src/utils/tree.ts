/*
 * 通用树构建工具：将平铺数据通过 parentId 关系构造成树，并按 weight 排序。
 */
export interface TreeBuildOptions<T, ID = number | string> {
  getId: (item: T) => ID;
  getParentId?: (item: T) => ID | null | undefined;
  /** 根节点的父ID标识，默认 '0' */
  rootPid?: ID | null | undefined;
  /** 可选：从原始项生成用于展示的 label */
  getLabel?: (item: T) => string;
  /** 可选：用于排序的权重，越小越靠前 */
  getWeight?: (item: T) => number | undefined;
  /** 可选：为节点追加额外字段（如 icon、menuType、key 等） */
  assign?: (item: T, node: any) => void;
}

export type TreeNodeBase<N = any> = {
  children?: N[];
  id: any;
  key?: any;
  label?: string;
  parentId?: any;
  weight?: number;
};

export function buildTree<T, N extends TreeNodeBase<N> = any>(
  items: T[],
  options: TreeBuildOptions<T, any>,
): N[] {
  const {
    getId,
    getParentId,
    rootPid = '0',
    getLabel,
    getWeight,
    assign,
  } = options;

  const flat: N[] = items.map((item) => {
    const node: any = {
      id: getId(item),
      parentId: getParentId ? (getParentId(item) ?? rootPid) : rootPid,
      children: [],
    };
    if (getLabel) node.label = getLabel(item);
    if (getWeight) node.weight = getWeight(item);
    if (assign) assign(item, node);
    return node as N;
  });

  const map = new Map<string, N>(flat.map((n: any) => [String(n.id), n]));
  const roots: N[] = [];

  for (const node of flat as any as N[]) {
    const pid = String(((node as any).parentId ?? rootPid) as any);
    const parent = map.get(pid);
    if (!pid || pid === String(rootPid) || !parent) {
      roots.push(node);
    } else {
      (parent as any).children = (parent as any).children || [];
      (parent as any).children.push(node);
    }
  }

  const sortTree = (nodes: N[]) => {
    nodes.sort(
      (a: any, b: any) =>
        ((a?.weight ?? 0) as number) - ((b?.weight ?? 0) as number),
    );
    nodes.forEach((n: any) => {
      if (n.children && n.children.length > 0) sortTree(n.children);
    });
  };
  sortTree(roots);

  return roots;
}
