import type { BaseCollectionModel } from 'pocketbase';

/**
 * @description: Post interface
 */
export interface Post extends BaseCollectionModel {
  id: string;
  post_code: string;
  post_name: string;
  post_sort: number;
  status: string;
  remark: string;
  create_time: string;
}
