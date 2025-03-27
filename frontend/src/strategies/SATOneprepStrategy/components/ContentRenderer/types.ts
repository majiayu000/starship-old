/**
 * 内容渲染器类型定义
 */

// 内容类型枚举
export enum ContentTypeEnum {
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

// 内容元数据
export interface ContentMetadata {
  style?: string[];
  content_type?: string;
  [key: string]: any;
}

// 基础内容项接口
export interface ContentItemBase {
  type: ContentTypeEnum;
  metadata?: ContentMetadata;
}

// 字符串内容
export interface StringContent extends ContentItemBase {
  type: ContentTypeEnum.STRING;
  content: string;
}

// 公式内容
export interface FormulaContent extends ContentItemBase {
  type: ContentTypeEnum.FORMULA;
  content: {
    content: string;
    display?: 'inline' | 'block';
  };
}

// 图片内容
export interface ImageContent extends ContentItemBase {
  type: ContentTypeEnum.IMAGE;
  content: {
    content: string;
    type: 'base64' | 'url' | 'svg';
    altText?: string;
    dimensions?: {
      width?: number | string;
      height?: number | string;
    };
  };
}

// HTML内容
export interface HtmlContent extends ContentItemBase {
  type: ContentTypeEnum.HTML;
  content: string;
}

// 空白填充内容
export interface BlankContent extends ContentItemBase {
  type: ContentTypeEnum.BLANK;
  content?: {
    width?: string | number;
    style?: string;
  };
}

// 图表内容
export interface FigureContent extends ContentItemBase {
  type: ContentTypeEnum.FIGURE;
  content: {
    content: string;
    type?: string;
    caption?: string;
  } | string;
}

// 表格内容
export interface TableContent extends ContentItemBase {
  type: ContentTypeEnum.TABLE;
  content: {
    content: string;
    caption?: string;
  };
}

// 换行内容
export interface LineBreakContent extends ContentItemBase {
  type: ContentTypeEnum.LINE_BREAK;
}

// 段落容器
export interface ParagraphContainer extends ContentItemBase {
  type: ContentTypeEnum.PARAGRAPH;
  content: ContentItemType[];
}

// 内容项类型联合
export type ContentItemType = 
  | StringContent
  | FormulaContent
  | ImageContent
  | HtmlContent
  | BlankContent
  | FigureContent
  | TableContent
  | LineBreakContent
  | ParagraphContainer;

// 内容渲染器属性
export interface ContentRendererProps {
  content: ContentItemType[] | ContentItemType | string | null;
  options?: ContentRendererOptions;
}

// 内容渲染器选项
export interface ContentRendererOptions {
  showSpecialSymbols?: boolean;
  mergeConsecutiveStrings?: boolean;
  convertOldRootRadicand?: boolean;
  renderMode?: 'default' | 'compact' | 'expanded';
  enableParagraphGrouping?: boolean;
  [key: string]: any;
}

// 内容项组件属性
export interface ContentItemProps {
  item: ContentItemType;
  index: number;
  isLast: boolean;
  totalItems: number;
  options: ContentRendererOptions;
}

// 特定内容组件属性
export interface TypedContentProps<T extends ContentItemType> {
  content: T extends { content: any } ? T['content'] : undefined;
  metadata?: ContentMetadata;
  options: ContentRendererOptions;
}

// 内容组件属性
export interface ContentComponentProps {
  content: any;
  metadata?: ContentMetadata;
  options?: ContentRendererOptions;
} 