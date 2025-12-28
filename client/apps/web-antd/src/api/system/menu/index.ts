import type { TenantPackage } from '../tenant-package/model';
import type { Menu, MenuOption } from './model';

import type { ID, IDS, PageQuery } from '#/api/common';

import { buildingQuery } from '#/api/helper';
import { Ors, pb } from '#/api/request';
import { buildTree } from '#/utils/tree';

const menuCollection = pb.collection<Menu>('menu');
const roleMenuCollection = pb.collection('role_menu');
const tenantPackageCollection = pb.collection('tenant_package');

/**
 * 菜单列表（树形结构，获取全部数据）
 * @param params 参数
 * @returns 列表
 */
export function menuList(params?: PageQuery) {
  const { filter, sort } = buildingQuery(params);
  return menuCollection.getFullList({
    filter,
    sort,
  });
}

export function menuPageList(params?: PageQuery) {
  const { sort } = buildingQuery(params);
  return menuCollection.getFullList({
    filter: Ors(["menu_type = 'M'", "menu_type = 'C'"]),
    sort,
  });
}

/**
 * 菜单详情
 * @param id 菜单id
 * @returns 菜单详情
 */
export function menuInfo(id: ID) {
  return menuCollection.getOne(`${id}`);
}

/**
 * 菜单新增
 * @param data 参数
 */
export function menuAdd(data: Partial<Menu>) {
  return menuCollection.create(data);
}

/**
 * 菜单更新
 * @param data 参数
 */
export function menuUpdate(data: Partial<Menu>) {
  return menuCollection.update(`${data.id}`, data);
}

/**
 * 菜单删除
 * @param menuIds ids
 */
export function menuRemove(menuIds: IDS) {
  return Promise.all(menuIds.map((id) => menuCollection.delete(`${id}`)));
}

/**
 * 返回对应角色的菜单
 * @param roleId id
 * @returns resp
 */
export async function roleMenuTreeSelect(roleId: ID) {
  // 1) 查询当前角色已勾选的菜单ID
  const roleMenus = await roleMenuCollection.getFullList<{
    menu: string;
  }>({
    filter: `role="${roleId}"`,
    fields: 'menu',
  });
  const checkedKeys = roleMenus.map((m) => m.menu);

  const menus = await menuTreeSelect();

  return { checkedKeys, menus };
}

/**
 * 下拉框使用  返回所有的菜单
 * @returns []
 */
export async function menuTreeSelect(filter: string = `status="0"`) {
  const menus = await menuCollection.getFullList<{
    icon?: string;
    id: string;
    menu_name: string;
    menu_type: string;
    order_num?: number;
    parent_id?: null | string;
    status: string;
  }>({
    filter,
    fields: 'id,parent_id,menu_name,menu_type,icon,order_num,status',
  });

  const roots = buildTree<any, MenuOption>(menus, {
    getId: (m) => m.id,
    getParentId: (m) => m.parent_id ?? '0',
    rootPid: '0',
    getLabel: (m) => m.menu_name,
    getWeight: (m) => m.order_num ?? 0,
    assign: (m, node: any) => {
      node.menuType = m.menu_type;
      node.icon = m.icon ?? '#';
      node.key = m.id;
    },
  });

  return roots as unknown as MenuOption[];
}

/**
 * 租户套餐使用
 * @param packageId packageId
 * @returns resp
 */
export async function tenantPackageMenuTreeSelect(packageId: ID) {
  const checkedKeys =
    `${packageId}` === '0'
      ? []
      : await tenantPackageCollection
          .getOne<TenantPackage>(`${packageId}`)
          .then((tp) => {
            return tp.menu_ids;
          });

  const menus = await menuTreeSelect(`status="0"&&menu_name!~"租户"`);

  return { checkedKeys, menus };
}
