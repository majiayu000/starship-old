'use client';

import React from 'react';
import { ContentMetadata } from '../base/BaseComponent';
import { ContentItem, ContentRenderOptions, ParagraphContainerProps } from '../types';
import { ContentErrorBoundary } from '../ErrorBoundary';

/**
 * 段落容器组件 - 用于将内容项分组为段落
 */
const ParagraphContainer: React.FC<ParagraphContainerProps> = ({
  content = [],
  metadata = {},
  options = {},
  children
}) => {
  const styles = metadata?.style || [];
  let containerClasses = 'paragraph-container mb-4';
  
  // 应用段落特定样式
  if (styles.includes('indent')) containerClasses += ' pl-8';
  if (styles.includes('center')) containerClasses += ' text-center';
  if (styles.includes('right')) containerClasses += ' text-right';
  if (styles.includes('justify')) containerClasses += ' text-justify';
  
  // 应用其他样式类
  if (styles.includes('bordered')) containerClasses += ' border border-gray-200 p-4 rounded';
  if (styles.includes('shaded')) containerClasses += ' bg-gray-50 p-4 rounded';
  
  // 使用子元素或渲染内容项
  return (
    <ContentErrorBoundary>
      <div 
        className={containerClasses}
        data-testid="paragraph-container"
        data-content-type={metadata.content_type}
      >
        {children || (
          <ParagraphContentRenderer 
            items={content} 
            options={options} 
          />
        )}
      </div>
    </ContentErrorBoundary>
  );
};

/**
 * 段落内容渲染器组件 - 负责渲染段落内的内容项
 */
interface ParagraphContentRendererProps {
  items: ContentItem[];
  options: ContentRenderOptions;
}

const ParagraphContentRenderer: React.FC<ParagraphContentRendererProps> = ({
  items,
  options
}) => {
  // 如果没有内容项，返回空
  if (!items || items.length === 0) {
    return null;
  }
  
  // 这里需要实现内容项的渲染
  // 我们会依赖工厂组件 ContentComponentFactory 来实现
  // 但在架构设计阶段，我们先返回一个占位符
  return (
    <div className="paragraph-content">
      {items.map((item, index) => (
        <div key={index} className="paragraph-item">
          {/* 这里将由工厂组件实现 */}
          <pre className="text-xs text-gray-500">
            {JSON.stringify(item, null, 2)}
          </pre>
        </div>
      ))}
    </div>
  );
};

// 使用 React.memo 优化渲染性能
export default React.memo(ParagraphContainer); 