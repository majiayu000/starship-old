'use client';

import React from 'react';
import { ContentItem } from './ContentItem';
import { ParagraphContainer as ParagraphContainerType, ContentRendererOptions, ContentTypeEnum } from './types';

/**
 * 段落容器组件 - 处理段落级别的布局和渲染
 */
interface ParagraphContainerProps {
  paragraph: ParagraphContainerType;
  options: ContentRendererOptions;
}

const ParagraphContainer: React.FC<ParagraphContainerProps> = ({ 
  paragraph, 
  options 
}) => {
  // 检查段落中是否包含图表等特殊内容
  const hasSpecialContent = paragraph.content.some(item => 
    item.type === ContentTypeEnum.FIGURE || 
    item.type === ContentTypeEnum.TABLE || 
    item.type === ContentTypeEnum.IMAGE
  );

  // 对于包含特殊内容的段落，使用不同的布局
  if (hasSpecialContent) {
    return (
      <div 
        className="paragraph-container-special"
        style={{ 
          display: 'block',
          margin: '10px 0',
          width: '100%',
          wordWrap: 'break-word',
          overflowWrap: 'break-word'
        }}
      >
        {paragraph.content.map((item, index) => (
          <ContentItem
            key={`${paragraph.id}-item-${index}`}
            item={item}
            index={index}
            isLast={index === paragraph.content.length - 1}
            totalItems={paragraph.content.length}
            options={options}
          />
        ))}
      </div>
    );
  }

  // 常规段落使用现有的flex布局
  return (
    <div 
      className="paragraph-container"
      style={{ 
        display: 'flex',
        flexWrap: 'wrap',
        alignItems: 'center',
        margin: 0,
        width: '100%',
        wordWrap: 'break-word',
        overflowWrap: 'break-word'
      }}
    >
      {paragraph.content.map((item, index) => {
        // 检查前一个元素是否有换行
        const prevItem = index > 0 ? paragraph.content[index - 1] : null;
        const prevHasLineBreak = prevItem && (
          prevItem.type === ContentTypeEnum.LINE_BREAK || 
          prevItem.metadata?.style?.includes('p')
        );
        
        // 计算是否需要添加间隔
        const needSpaceBefore = index > 0 && !prevHasLineBreak;
        
        return (
          <React.Fragment key={`${paragraph.id}-item-${index}`}>
            {/* 在元素之前添加间隔，但在段落开始或换行后不添加 */}
            {needSpaceBefore && (
              <div style={{ display: 'inline-block', width: '0.3rem' }}></div>
            )}
            
            <ContentItem
              item={item}
              index={index}
              isLast={index === paragraph.content.length - 1}
              totalItems={paragraph.content.length}
              options={options}
            />
          </React.Fragment>
        );
      })}
    </div>
  );
};

export default ParagraphContainer; 