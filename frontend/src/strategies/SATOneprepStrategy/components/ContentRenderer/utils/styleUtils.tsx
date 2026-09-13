import React, { ReactNode } from 'react';
import { sanitizeHtml } from '@/lib/sanitizeHtml';

/**
 * 应用样式到内容
 * @param content 要应用样式的React节点
 * @param styles 样式数组
 * @returns 应用样式后的React节点
 */
export const applyStyles = (content: ReactNode, styles?: string[]): ReactNode => {
  if (!styles || !styles.length) return content;
  
  return styles.reduce((styledContent, style) => {
    switch(style) {
      case 'p': 
        // p样式由父组件处理，这里不处理
        return styledContent;
      case 'em': 
      case 'i': 
        return <em>{styledContent}</em>;
      case 'strong': 
      case 'b': 
        return <strong>{styledContent}</strong>;
      case 'sup': 
        return <sup>{styledContent}</sup>;
      default:
        // 处理html-*前缀
        if (style.startsWith('html-')) {
          const actualStyle = style.substring(5);
          const styleObj: Record<string, string> = {};
          
          // 简单解析HTML样式
          const styleProps = actualStyle.split(';');
          styleProps.forEach(prop => {
            const [key, value] = prop.split(':').map(s => s.trim());
            if (key && value) {
              // 转换为驼峰式
              const camelKey = key.replace(/-([a-z])/g, (_, letter) => letter.toUpperCase());
              styleObj[camelKey] = value;
            }
          });
          
          return <span style={styleObj}>{styledContent}</span>;
        }
        return styledContent;
    }
  }, content);
};

/**
 * 处理HTML内容
 * @param html HTML字符串
 * @returns 处理后的React节点
 */
export const processHTML = (html: string): React.ReactNode => {
  try {
    // 安全检查：如果HTML内容为空或不是字符串，返回空内容
    if (!html || typeof html !== 'string') {
      return <span>无内容</span>;
    }
    
    // 处理特殊的HTML格式
    if (html.includes('old-root-radicand')) {
      // 提取根号内部的内容 - 简单的正则表达式实现
      const rootContent = html.match(/old-root-radicand[^>]*>([^<]+)<\/span>/);
      const innerContent = rootContent ? rootContent[1].trim() : '';
      
      // 这里为了简单，我们直接返回渲染好的内容
      // 实际应用中，你可能想要返回一个特殊的对象供后续处理
      return (
        <span className="question-formula root-formula">
          √({innerContent})
        </span>
      );
    }
    
    const safeHtml = sanitizeHtml(html);
    
    // 默认返回HTML内容
    return <div dangerouslySetInnerHTML={{ __html: safeHtml }} />;
  } catch (error) {
    console.error('HTML处理错误:', error);
    return <span className="error">HTML处理错误</span>;
  }
};
