'use client';

import React from 'react';
import { ContentItemProps, ContentTypeEnum, StringContent, FormulaContent, ImageContent, HtmlContent, FigureContent, TableContent, BlankContent } from './types';
import { InlineMath, BlockMath } from 'react-katex';
import 'katex/dist/katex.min.css';

/**
 * 内容项组件 - 负责渲染不同类型的内容项
 */
export const ContentItem: React.FC<ContentItemProps> = ({
  item,
  index,
  isLast,
  totalItems,
  options
}) => {
  if (!item) return null;
  
  // 获取元数据和样式
  const metadata = item.metadata || {};
  const styles = metadata.style || [];
  
  // 检查是否为段落样式（但不改变元素本身的显示模式）
  const hasPStyle = styles.includes('p');
  
  // 渲染字符串类型
  const renderString = () => {
    if (item.type !== ContentTypeEnum.STRING) return null;
    const stringItem = item as StringContent;
    if (!stringItem.content) return null;
    
    // 应用样式
    let className = '';
    if (styles.includes('em') || styles.includes('i')) className += ' italic';
    if (styles.includes('strong') || styles.includes('b')) className += ' font-bold';
    if (styles.includes('sup')) className += ' align-super text-xs';
    
    return (
      <span className={className.trim()}>
        {stringItem.content}
      </span>
    );
  };
  
  // 渲染公式类型
  const renderFormula = () => {
    if (item.type !== ContentTypeEnum.FORMULA) return null;
    const formulaItem = item as FormulaContent;
    if (!formulaItem.content?.content) return null;
    
    const { content, display = 'inline' } = formulaItem.content;
    
    return display === 'block' ? (
      <BlockMath math={content} />
    ) : (
      <InlineMath math={content} />
    );
  };
  
  // 渲染图片类型
  const renderImage = () => {
    if (item.type !== ContentTypeEnum.IMAGE) return null;
    const imageItem = item as ImageContent;
    if (!imageItem.content?.content) return null;
    
    const { content, type, altText = '', dimensions } = imageItem.content;
    
    if (type === 'svg') {
      return <div dangerouslySetInnerHTML={{ __html: content }} />;
    }
    
    return (
      <img 
        src={content} 
        alt={altText} 
        className="max-w-full" 
        width={dimensions?.width} 
        height={dimensions?.height}
      />
    );
  };
  
  // 渲染HTML内容
  const renderHtml = () => {
    if (item.type !== ContentTypeEnum.HTML) return null;
    const htmlItem = item as HtmlContent;
    if (!htmlItem.content) return null;
    
    return <div dangerouslySetInnerHTML={{ __html: htmlItem.content }} />;
  };
  
  // 渲染图表内容
  const renderFigure = () => {
    if (item.type !== ContentTypeEnum.FIGURE) return null;
    const figureItem = item as FigureContent;
    
    // 添加调试日志
    console.log('Figure content:', figureItem);
    
    // 检查content的类型，处理可能的不同格式
    let figureContent = '';
    if (typeof figureItem.content === 'string') {
      figureContent = figureItem.content;
    } else if (figureItem.content?.content) {
      figureContent = figureItem.content.content;
    }
    
    if (!figureContent) {
      console.warn('Empty figure content for:', figureItem);
      return null;
    }
    
    // 添加额外的样式以确保SVG正确显示
    return (
      <div 
        className="figure-container w-full overflow-auto"
        style={{ maxWidth: '100%', margin: '10px 0' }}
        dangerouslySetInnerHTML={{ __html: figureContent }} 
      />
    );
  };
  
  // 渲染表格内容
  const renderTable = () => {
    if (item.type !== ContentTypeEnum.TABLE) return null;
    const tableItem = item as TableContent;
    if (!tableItem.content?.content) return null;
    
    return <div className="table-container" dangerouslySetInnerHTML={{ __html: tableItem.content.content }} />;
  };
  
  // 渲染空白填充
  const renderBlank = () => {
    // 检查是否有内容和自定义宽度
    let width = 'w-20'; // 默认宽度
    let customStyle = {};
    
    if (item.type === ContentTypeEnum.BLANK) {
      const blankItem = item as BlankContent;
      if (blankItem.content?.width) {
        // 如果提供了自定义宽度，删除默认宽度类
        width = '';
        // 添加自定义宽度到样式
        customStyle = { 
          ...customStyle, 
          width: typeof blankItem.content.width === 'number' 
            ? `${blankItem.content.width}px` 
            : blankItem.content.width 
        };
      }
      
      // 如果提供了自定义样式
      if (blankItem.content?.style) {
        customStyle = { ...customStyle, ...JSON.parse(blankItem.content.style) };
      }
    }
    
    return (
      <span 
        className={`inline-block ${width} border-b-2 border-black mx-1 align-baseline`} 
        style={{ 
          height: '0.15em', 
          marginBottom: '0.2em',
          ...customStyle
        }}
      ></span>
    );
  };
  
  // 渲染各种内容类型
  const renderContent = () => {
    switch (item.type) {
      case ContentTypeEnum.STRING:
        return renderString();
      case ContentTypeEnum.FORMULA:
        return renderFormula();
      case ContentTypeEnum.IMAGE:
        return renderImage();
      case ContentTypeEnum.HTML:
        return renderHtml();
      case ContentTypeEnum.FIGURE:
        return renderFigure();
      case ContentTypeEnum.TABLE:
        return renderTable();
      case ContentTypeEnum.BLANK:
        return renderBlank();
      case ContentTypeEnum.LINE_BREAK:
        return <br />;
      default:
        return <span className="text-gray-400">未知内容类型: {String(item.type)}</span>;
    }
  };
  
  // 决定是否在末尾添加空格
  const shouldAddSpace = !isLast && item.type !== ContentTypeEnum.LINE_BREAK;
  
  return (
    <>
      {/* 所有内容项都用inline-flex显示 */}
      <div className="content-item" style={{ display: 'inline-flex' }}>
        {renderContent()}
      </div>
      
      {/* 在元素之间添加空格，除非是最后一个元素或换行 */}
      {shouldAddSpace && <span> </span>}
      
      {/* 如果有p样式，在元素后添加换行 */}
      {hasPStyle && <br />}
    </>
  );
};

export default ContentItem; 