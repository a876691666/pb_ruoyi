import type { BaseCollectionModel } from 'pocketbase';

export interface Menu extends BaseCollectionModel {
  create_by?: any;
  create_time: string;
  update_by?: any;
  update_time?: any;
  remark?: any;
  id: string;
  menu_name: string;
  parent_name?: string;
  parent_id: string;
  order_num: number;
  path: string;
  component?: string;
  query: string;
  is_frame: string;
  is_cache: string;
  menu_type: string;
  visible: string;
  status: string;
  perms: string;
  icon: string;
  children: Menu[];
}

/**
 * @description 菜单信息
 * @param label 菜单名称
 */
export interface MenuOption {
  id: number;
  parent_id: string;
  label: string;
  weight: number;
  children: MenuOption[];
  key: string; // 实际上不存在 ide报错
  menu_type: string;
  icon: string;
}

/**
 * @description 菜单返回
 * @param checkedKeys 选中的菜单id
 * @param menus 菜单信息
 */
export interface MenuResp {
  checked_keys: number[];
  menus: MenuOption[];
}

/**
 * 菜单表单查询
 */
export interface MenuQuery {
  menu_name?: string;
  visible?: string;
  status?: string;
}
