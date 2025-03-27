'use client';

import React from 'react';
import { contentRegistry } from './ContentComponentFactory';
import { ContentType } from '../../components/content/types';

// 导入标准组件
import TextComponent from '../../components/content/elements/TextComponent';
import FormulaComponent from '../../components/content/elements/FormulaComponent';
import ImageComponent from '../../components/content/elements/ImageComponent';
import BlankComponent from '../../components/content/elements/BlankComponent';
import LineBreakComponent from '../../components/content/elements/LineBreakComponent';
import HtmlComponent from '../../components/content/elements/HtmlComponent';

// 默认组件注册函数
export function registerDefaultComponents(): void {
  // 注册基础内容类型
  contentRegistry.register(ContentType.STRING, TextComponent, 10);
  contentRegistry.register(ContentType.FORMULA, FormulaComponent, 10);
  contentRegistry.register(ContentType.IMAGE, ImageComponent, 10);
  contentRegistry.register(ContentType.BLANK, BlankComponent, 10);
  contentRegistry.register(ContentType.LINE_BREAK, LineBreakComponent, 10);
  contentRegistry.register(ContentType.HTML, HtmlComponent, 10);
  
  // 注册默认内容处理组件
  const DefaultComponent: React.FC<any> = ({ content, metadata }) => (
    <div className="unknown-content p-2 border border-gray-300 bg-gray-50 rounded">
      <div className="text-xs text-gray-500 mb-1">未识别的内容类型</div>
      <pre className="text-xs overflow-auto">{JSON.stringify({ content, metadata }, null, 2)}</pre>
    </div>
  );
  
  contentRegistry.registerDefaultComponent(DefaultComponent);
  
  console.log('已注册默认内容组件');
}

// 自动注册组件 (在应用初始化时调用)
export function initializeContentComponents(): void {
  // 检查是否已注册组件
  if (contentRegistry.hasType(ContentType.STRING)) {
    console.log('内容组件已注册，跳过初始化');
    return;
  }
  
  // 注册默认组件
  registerDefaultComponents();
}

// 导出单例实例
export { contentRegistry }; 