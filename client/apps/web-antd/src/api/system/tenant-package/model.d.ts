import type { BaseCollectionModel } from 'pocketbase';

/**
 * @description 租户套餐
 * @param id id
 * @param package_name 名称
 * @param menu_ids 菜单id  格式为[1,2,3] 返回为string 提交为数组
 * @param remark 备注
 * @param menu_check_strictly 是否关联父节点
 * @param status 状态
 */
export interface TenantPackage extends BaseCollectionModel {
  id: string;
  package_name: string;
  menu_ids: string[];
  remark?: string;
  menu_check_strictly?: boolean;
  status?: string;
}
