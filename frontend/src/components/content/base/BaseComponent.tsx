'use client';

import React from 'react';
import { ContentErrorBoundary } from '../ErrorBoundary';

// 内容元数据类型
export interface ContentMetadata {
  style?: string[];
  content_type?: string;
  [key: string]: any;
}

// 基础内容组件属性
export interface BaseComponentProps {
  metadata?: ContentMetadata;
  children: React.ReactNode;
  className?: string;
  testId?: string;
}

/**
 * 基础内容组件 - 所有内容元素的基础组件
 * 提供通用样式处理和错误边界
 */
const BaseComponent: React.FC<BaseComponentProps> = ({ 
  metadata = {}, 
  children,
  className = '',
  testId
}) => {
  const { style = [], content_type } = metadata;
  
  // 生成样式类名
  const styleClasses = style.map(s => `style-${s}`).join(' ');
  const contentTypeClass = content_type ? `content-type-${content_type}` : '';
  
  // 检查是否应渲染为块级元素
  const isBlock = style.includes('p');
  
  return (
    <ContentErrorBoundary>
      <div 
        className={`sat-component ${styleClasses} ${contentTypeClass} ${isBlock ? 'block-component' : 'inline-component'} ${className}`}
        data-testid={testId}
        data-content-type={content_type}
      >
        {children}
      </div>
    </ContentErrorBoundary>
  );
};

// 使用 React.memo 优化渲染性能
export default React.memo(BaseComponent); 