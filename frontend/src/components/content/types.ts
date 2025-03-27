import { ContentMetadata } from './base/BaseComponent';

/**
 * 内容类型枚举
 */
export enum ContentType {
  STRING = 'string',
  FORMULA = 'formula',
  IMAGE = 'image',
  HTML = 'html',
  BLANK = 'blank',
  FIGURE = 'figure',
  TABLE = 'table',
  LINE_BREAK = 'line_break',
  PARAGRAPH = 'paragraph'
}

/**
 * 渲染模式枚举
 */
export enum RenderMode {
  DEFAULT = 'default',
  COMPACT = 'compact',
  EXPANDED = 'expanded'
}

/**
 * 内容渲染选项
 */
export interface ContentRenderOptions {
  showSpecialSymbols?: boolean;
  mergeConsecutiveStrings?: boolean;
  convertOldRootRadicand?: boolean;
  renderMode?: RenderMode;
  enableParagraphGrouping?: boolean;
  [key: string]: any;
}

/**
 * 基础内容项接口
 */
export interface ContentItem {
  type: ContentType | string;
  metadata?: ContentMetadata;
  content?: any;
}

/**
 * 内容组件基础属性
 */
export interface ContentComponentProps<T = any> {
  content: T;
  metadata?: ContentMetadata;
  options?: ContentRenderOptions;
}

/**
 * 字符串内容组件属性
 */
export interface TextComponentProps extends ContentComponentProps<string> {}

/**
 * 公式内容组件属性
 */
export interface FormulaComponentProps extends ContentComponentProps<{
  content: string;
  display?: 'inline' | 'block';
}> {}

/**
 * 图片内容组件属性
 */
export interface ImageComponentProps extends ContentComponentProps<{
  content: string;
  type: 'base64' | 'url' | 'svg';
  altText?: string;
  dimensions?: {
    width?: number | string;
    height?: number | string;
  };
}> {}

/**
 * HTML内容组件属性
 */
export interface HtmlComponentProps extends ContentComponentProps<string> {}

/**
 * 空白填充组件属性
 */
export interface BlankComponentProps extends ContentComponentProps<{
  width?: string | number;
  style?: string;
}> {}

/**
 * 图表组件属性
 */
export interface FigureComponentProps extends ContentComponentProps<{
  content: string;
  type?: string;
  caption?: string;
}> {}

/**
 * 表格组件属性
 */
export interface TableComponentProps extends ContentComponentProps<{
  content: string;
  caption?: string;
}> {}

/**
 * 换行组件属性
 */
export interface LineBreakComponentProps extends ContentComponentProps {}

/**
 * 段落容器组件属性
 */
export interface ParagraphContainerProps extends ContentComponentProps<ContentItem[]> {
  children?: React.ReactNode;
} 