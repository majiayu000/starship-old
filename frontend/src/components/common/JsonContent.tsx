'use client';

import React from 'react';

interface JsonContentProps {
  content: any;
}

const JsonContent: React.FC<JsonContentProps> = ({ content }) => {
  if (!content) return <div>无内容</div>;
  
  // 简单实现，实际可能需要更复杂的处理逻辑
  if (typeof content === 'string') {
    try {
      // 尝试解析JSON字符串
      const parsed = JSON.parse(content);
      return renderJson(parsed);
    } catch {
      // 如果不是JSON，直接显示
      return <div className="text-content">{content}</div>;
    }
  }
  
  return renderJson(content);
};

// 根据JSON结构和数据类型渲染不同的UI
const renderJson = (data: any): React.ReactNode => {
  if (Array.isArray(data)) {
    // 检查是否所有元素都是图形类型
    const allGraphics = data.every(
      item => item && typeof item === 'object' && 
      (item.type === 'figure' || item.type === 'image')
    );
    
    if (allGraphics) {
      return (
        <div className="flex flex-wrap justify-center items-start">
          {data.map((item, index) => (
            <div key={index} className="mx-2 my-2 flex-shrink-0">
              {renderJson(item)}
            </div>
          ))}
        </div>
      );
    }
    
    return (
      <div className="json-array">
        {data.map((item, index) => (
          <div key={index} className="json-array-item">
            {renderJson(item)}
          </div>
        ))}
      </div>
    );
  }
  
  if (data && typeof data === 'object') {
    // 特殊处理某些已知结构
    if (data.type === 'string' && data.content) {
      return <div className="text-content">{data.content}</div>;
    }
    
    if (data.type === 'figure' && data.content) {
      // 处理SVG内容
      if (data.content.type === 'svg' || data.content.type === 'html') {
        const content = data.content.content;
        // 安全处理内容，去除可能的脚本
        const safeContent = typeof content === 'string' 
          ? content.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
                   .replace(/on\w+="[^"]*"/gi, '')
          : '';
        
        return (
          <figure className="image inline-block max-w-md mx-2 my-2 align-top">
            <div dangerouslySetInnerHTML={{ __html: safeContent }} />
            {data.content.caption && (
              <figcaption className="text-sm text-gray-500 mt-1 text-center">
                {data.content.caption}
              </figcaption>
            )}
          </figure>
        );
      }
      return <div className="text-error">不支持的图形类型</div>;
    }
    
    if (data.type === 'image' && data.content && data.content.content) {
      // 检查图片类型
      if (data.content.type === 'base64') {
        // 验证base64数据格式
        const isValidBase64 = typeof data.content.content === 'string' && 
                             data.content.content.startsWith('data:image');
        
        if (isValidBase64) {
          return (
            <img 
              src={data.content.content} 
              alt={data.content.altText || "题目图片"} 
              style={{ 
                maxWidth: '100%', 
                display: 'block',
                margin: data.content.style?.includes('center') ? '0 auto' : 'initial' 
              }} 
            />
          );
        } else {
          return <div className="text-error">无效的base64图片格式</div>;
        }
      }
      
      // 处理URL类型图片或其他类型
      return <img src={data.content.content} alt="题目图片" style={{ maxWidth: '100%' }} />;
    }
    
    // 默认显示为JSON
    return (
      <div className="json-object">
        {Object.entries(data).map(([key, value]) => (
          <div key={key} className="json-property">
            <strong>{key}: </strong>
            {renderJson(value)}
          </div>
        ))}
      </div>
    );
  }
  
  // 基本类型
  return <div className="json-value">{String(data)}</div>;
};

export default JsonContent; 