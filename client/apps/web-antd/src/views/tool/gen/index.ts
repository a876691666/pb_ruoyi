import type { GenInfo } from '#/api/tool/gen/model';

import { BlobWriter, TextReader, ZipWriter } from '@zip.js/zip.js';
import velocityjs from 'velocityjs';

import * as GenConstants from './gen.constant';

// 使用 Vite 的 import.meta.glob 导入模板文件
const templateModules = import.meta.glob('/template/**/*.vm', {
  query: '?raw',
  eager: true,
  import: 'default',
}) as Record<string, string>;

// 匹配所有换行符到#之间的内容，并且替换为空，兼容crlf和lf
function replaceSpace(content: string) {
  return content.replaceAll(/\r?\n\s*#/g, '\n#');
}

// 传入：xx{aa}x{c} 和 { aa: 'bb', c: 1 } 返回 xxbbx1
function replaceStr(content: string, options: { [key: string]: string }) {
  return content.replaceAll(/\{([^}]+)\}/g, (match, key) => {
    return options[key] || match;
  });
}

function ignoreField(field: string) {
  return [
    'createBy',
    'createTime',
    'delFlag',
    'updateBy',
    'updateTime',
  ].includes(field);
}

function getValidatorDecorator(javaType: string) {
  if (javaType === 'Boolean') return `@IsBoolean()`;
  if (javaType === 'String') return `@IsString()`;
  if (javaType === 'Date') return `@IsDate()`;
  if (javaType === 'Number') return `@IsNumber()`;
  if ([`BigDecimal`, `Double`, `Float`, `Integer`, `Long`].includes(javaType))
    return `@IsNumber()`;
  return ``;
}

function getTsType(javaType: string) {
  if (javaType === 'Boolean') return `boolean`;
  if (javaType === 'String') return `string`;
  if (javaType === 'Date') return `Date`;
  if (javaType === 'Number') return `number`;
  if ([`BigDecimal`, `Double`, `Float`, `Integer`, `Long`].includes(javaType))
    return `number`;
  return `any`;
}

function getUiTsType(column: any) {
  if (column.type === 'bool') return `boolean`;
  if (column.type === 'file') return `blob`;
  if (column.type === 'number') return `number`;
  if (column.type === 'json') return `any`;
  return `string`;
}

function getFieldComment(column: any) {
  const comment = column.comment || column.name;
  const index = comment.indexOf('（');
  if (index !== -1) {
    return comment.slice(0, Math.max(0, index));
  }
  return comment;
}

function isDate(column: any) {
  return (
    ['autodate', 'datetime'].includes(column.type.toLowerCase()) ||
    ['datetime'].includes(column.htmlType.toLowerCase())
  );
}

function getBigintType(ColumnType: string) {
  if (ColumnType === 'bigint') return 'string';
  return 'number';
}

function getBool(str: string) {
  console.log(str);
  return str === '1' ? 'true' : 'false';
}

function convertToSnakeCase(str: string) {
  return str
    .replaceAll(/[A-Z]/g, (letter) => `_${letter.toLowerCase()}`)
    .replace(/^_/, '');
}

/**
 * 将小驼峰命名（camelCase）转换为下划线分割（snake_case）
 * @param {string} str - 输入的字符串
 * @returns {string} - 下划线分割的字符串
 */
export function camelToSnake(str?: string): string {
  return (str || '').replaceAll(/([a-z0-9])([A-Z])/g, '$1_$2').toLowerCase();
}

// 处理模板文件列表
const templateList = Object.entries(templateModules).map(
  ([filePath, content]) => {
    const fileTplCategory =
      filePath.match(/\/(\w+)_index\.vue/)?.[1] ||
      filePath.match(/\/(\w+)_popup\.vue/)?.[1] ||
      '';
    // filePath 格式: /template/xxx/yyy.zzz.vm
    // 移除 /template/ 前缀和 .vm 后缀
    const relativePath = filePath
      .replace(/^\/template\//, '')
      .replace(/\.vm$/, '')
      .replace(/\/\w+_index\.vue/, '/index.vue')
      .replace(/\/\w+_popup\.vue/, '/popup.vue');
    const name = relativePath;

    const previewName = filePath
      .replace(/^\/template\//, '')
      .replace(/\/\w+_index\.vue/, '/index.vue')
      .replace(/\/\w+_popup\.vue/, '/popup.vue');

    return [
      name,
      previewName,
      replaceSpace(content as string),
      fileTplCategory,
    ] as [string, string, string, string];
  },
);

/**
 * 校验数组是否包含指定值
 *
 * @param arr 数组
 * @param targetValue 值
 * @return 是否包含
 */

export function arraysContains(array: string[], value: string): boolean {
  return array.includes(value);
}

/**
 * 将字符串转换为小驼峰命名法
 * @param {string} str - 输入的字符串，使用下划线分隔
 * @returns {string} - 转换后的小驼峰命名法字符串
 */
export function toCamelCase(str: string): string {
  return str.replaceAll(/_([a-z])/g, (_match, letter) => letter.toUpperCase());
}

/**
 * 将字符串转换为大驼峰命名法
 * @param {string} str - 输入的字符串，使用下划线分隔
 * @returns {string} - 转换后的大驼峰命名法字符串
 */
export function toPascalCase(str: string) {
  return str?.[0]?.toUpperCase() + toCamelCase(str).slice(1);
}
/**
 * 将字符串转换为驼峰命名法
 * @param {string} str - 输入的字符串，使用下划线分隔
 * @returns {string} - 转换后的驼峰命名法字符串
 */
export function convertToCamelCase(str?: string): string {
  return (
    str?.replaceAll(/_([a-z])/g, (_, letter) => letter.toUpperCase()) || ''
  );
}

/**
 * 将字符串的首字母大写
 * @param {string} str - 输入的字符串
 * @returns {string} - 首字母大写后的字符串
 */
export function capitalize(str: string): string {
  return str.charAt(0).toUpperCase() + str.slice(1);
}
export function gen(options: any, isVM = false) {
  const result: {
    [name: string]: string;
  } = {};
  for (const [name, previewName, content, fileTplCategory] of templateList) {
    if (fileTplCategory && fileTplCategory !== options.tpl_category) {
      continue;
    }
    result[replaceStr(isVM ? previewName : name, options)] = velocityjs.render(
      content,
      {
        GenConstants,
        getValidatorDecorator,
        getUiTsType,
        getFieldComment,
        isDate,
        getBigintType,
        ignoreField,
        getTsType,
        getBool,
        convertToSnakeCase,
        ...options,
      },
    );
  }

  return result;
}

/**
 * 将文件打包成 zip 并返回 Blob
 * @param files - 文件对象，键为文件路径,值为文件内容（字符串）
 * @returns Promise<Blob> - zip 文件的 Blob 对象
 */
export async function createZipBlob(files: {
  [filePath: string]: string;
}): Promise<Blob> {
  const blobWriter = new BlobWriter('application/zip');
  const zipWriter = new ZipWriter(blobWriter);

  // 将所有文件添加到 zip
  for (const [filePath, content] of Object.entries(files)) {
    await zipWriter.add(filePath, new TextReader(content));
  }

  // 关闭 zip 写入器并获取 blob
  await zipWriter.close();
  return blobWriter.getData();
}

/**
 * 根据 GenInfo 数据生成完整代码（用于下载）
 * @param data - GenInfo 类型的数据
 * @returns 生成的代码文件对象
 */
export function generateCodeFromGenInfo(data: GenInfo, isVM = false) {
  // 从 fields 中找到主键列
  const primaryColumn = data.fields?.find((col: any) => col.pk) || null;
  const primaryKey = primaryColumn?.name || 'id';

  // 处理 columns 数据
  const columns = data.fields || [];
  const info = {
    primaryColumn,
    primaryKey,
    permissionPrefix: `${data.module_name}:${data.business_name}`,
    ...data,
    columns,

    functionName: convertToCamelCase(data.function_name),
    FunctionName: capitalize(convertToCamelCase(data.function_name)),
    _functionName: data.function_name,

    moduleName: convertToCamelCase(data.module_name),
    ModuleName: capitalize(convertToCamelCase(data.module_name)),
    _moduleName: data.module_name,
    module_name: convertToSnakeCase(data.module_name),

    businessName: convertToCamelCase(data.business_name),
    BusinessName: capitalize(convertToCamelCase(data.business_name)),
    _businessName: data.business_name,
    business_name: convertToSnakeCase(data.business_name),

    popupComponent: convertToCamelCase(data.options?.popup_type),
    PopupComponent: capitalize(convertToCamelCase(data.options?.popup_type)),
    _popupComponent: data.options?.popup_type,
    popup_component: convertToSnakeCase(data.options?.popup_type),

    formComponent: convertToCamelCase(data.options?.form_type),
    FormComponent: capitalize(convertToCamelCase(data.options?.form_type)),
    _formComponent: data.options?.form_type,
    form_component: convertToSnakeCase(data.options?.form_type),
  };

  return gen(info, isVM);
}
