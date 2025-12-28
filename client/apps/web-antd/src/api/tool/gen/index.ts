import type { GenInfo } from './model';

import type { ID, IDS, PageQuery } from '#/api/common';

import { buildingQuery } from '#/api/helper';
import { pb, requestClient } from '#/api/request';
import { serializeImportTableData } from '#/views/tool/gen/gen-data-serializer';

const collection = pb.collection<GenInfo>('gen_table');

enum Api {
  batchGenCode = '/tool/gen/batchGenCode',
  columnList = '/tool/gen/column',
  dataSourceNames = '/tool/gen/getDataNames',
  download = '/tool/gen/download',
  genCode = '/tool/gen/genCode',
  generatedList = '/tool/gen/list',
  importTable = '/tool/gen/importTable',
  preview = '/tool/gen/preview',
  readyToGenList = '/tool/gen/db/list',
  root = '/tool/gen',
  syncDb = '/tool/gen/synchDb',
}

type TableItem = Record<any, any>;

// 查询代码生成列表
export function generatedList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return collection.getList(currentPage, pageSize, {
    filter,
    sort,
  });
}

// 修改代码生成业务
export function genInfo(tableId: ID) {
  return collection.getOne(`${tableId}`);
}

// 查询数据库列表
export function readyToGenList(): Promise<TableItem[]> {
  return requestClient.get('/system/collections');
}

// 查询数据表字段列表
export function columnList(tableId: ID) {
  return requestClient.get(`${Api.columnList}/${tableId}`);
}

async function checkTableExists(tableName: string) {
  try {
    const item = await collection.getFirstListItem(`name="${tableName}"`, {
      requestKey: tableName,
    });
    return !!item;
  } catch {
    return false;
  }
}

/**
 * 导入表结构（保存）
 * @param tables 表名称集合
 */
export function importTable(tables: TableItem | TableItem[]) {
  if (!Array.isArray(tables)) {
    tables = [tables];
  }
  return Promise.all(
    tables.map(async (table: TableItem) => {
      // 使用序列化函数处理数据
      const serializedData = serializeImportTableData(table);
      if (await checkTableExists(table.name)) return;
      return collection.create(serializedData, { requestKey: table.name });
    }),
  );
}

// 修改保存代码生成业务
export function editSave(data: GenInfo) {
  return collection.update(`${data.id}`, data);
}

// 删除代码生成
export function genRemove(tableIds: IDS) {
  return tableIds.map((tableId) => collection.delete(`${tableId}`));
}

// 预览代码
export function previewCode(tableId: ID) {
  return requestClient.get<{ [key: string]: string }>(
    `${Api.preview}/${tableId}`,
  );
}

// 生成代码（下载方式）
export function genDownload(tableId: ID) {
  return requestClient.get<Blob>(`${Api.download}/${tableId}`);
}

// 同步数据库
export async function syncDb(table: GenInfo) {
  const collection = await requestClient.get(
    `/system/collection/${table.name}`,
  );
  if (!collection) return;

  const serializedData = serializeImportTableData(collection);
  if (await checkTableExists(table.name)) {
    return collection.update(`${table.id}`, { fields: serializedData.fields });
  }
  return collection.create(serializedData);
}

// 批量生成代码
export function batchGenCode(tableIdStr: ID | IDS) {
  return requestClient.get<Blob>(Api.batchGenCode, {
    isTransformResponse: false,
    params: { tableIdStr },
    responseType: 'blob',
  });
}
