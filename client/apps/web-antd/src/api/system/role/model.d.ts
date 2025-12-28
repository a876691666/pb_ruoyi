import type { BaseCollectionModel } from 'pocketbase';

export interface Role extends BaseCollectionModel {
  id: string;
  role_name: string;
  role_key: string;
  role_sort: number;
  data_scope: string;
  menu_check_strictly: boolean;
  dept_check_strictly: boolean;
  status: string;
  remark: string;
  create_time: string;
  // 用户是否存在此角色标识 默认不存在
  flag: boolean;
}

export interface DeptOption {
  id: number;
  parentId: number;
  label: string;
  weight: number;
  children: DeptOption[];
  key: string; // 实际上不存在 ide报错
}

export interface DeptResp {
  checkedKeys: number[];
  depts: DeptOption[];
}
