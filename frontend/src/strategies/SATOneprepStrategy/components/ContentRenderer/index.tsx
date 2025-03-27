'use client';

import React, { useEffect } from 'react';
import { useContentAdapter } from './hooks/useContentAdapter';
import { ContentRendererProps, ContentTypeEnum } from './types';
import { ContentItemFactory, initializeContentComponents } from '../../../../factories/content/registerComponents';
import ParagraphContainer from '../../../../components/content/containers/ParagraphContainer';
import './styles.css';

/**
 * SAT Oneprep 内容渲染器组件
 * 负责根据特定策略渲染各种类型的内容
 */
export const ContentRenderer: React.FC<ContentRendererProps> = ({ 
  content,
  options = {}
}) => {
  // 初始化内容组件
  useEffect(() => {
    initializeContentComponents();
  }, []);
  
  // 如果内容为空，返回空内容提示
  if (!content) {
    return <div className="text-center text-gray-500">无内容</div>;
  }
  
  // 默认启用段落分组，但禁用字符串合并
  const finalOptions = {
    mergeConsecutiveStrings: false,
    enableParagraphGrouping: true,
    ...options
  };
  
  // 使用适配器处理内容，使其符合标准格式
  const adaptedContent = useContentAdapter(content, finalOptions);
  
  // 渲染处理后的内容
  return (
    <div className="sat-content">
      {adaptedContent.map((item, index) => {
        // 处理段落容器
        if (item.type === ContentTypeEnum.PARAGRAPH) {
          const paragraphItem = item as any; // 类型断言处理段落项目
          return (
            <ParagraphContainer
              key={`paragraph-${index}`}
              content={paragraphItem.content || []}
              metadata={paragraphItem.metadata}
              options={finalOptions}
            />
          );
        }
        
        // 处理常规内容项
        return (
          <ContentItemFactory 
            key={`item-${index}`}
            item={item as any}
            options={finalOptions}
          />
        );
      })}
    </div>
  );
};

export default ContentRenderer; 