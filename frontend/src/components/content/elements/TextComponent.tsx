'use client';

import React from 'react';
import BaseComponent from '../base/BaseComponent';
import { TextComponentProps } from '../types';

/**
 * 文本内容组件 - 用于渲染字符串内容
 */
const TextComponent: React.FC<TextComponentProps> = ({
  content,
  metadata = {},
  options = {}
}) => {
  if (!content) return null;
  
  // 应用基于元数据的额外样式类
  let additionalClasses = '';
  
  const styles = metadata.style || [];
  if (styles.includes('em') || styles.includes('i')) additionalClasses += ' italic';
  if (styles.includes('strong') || styles.includes('b')) additionalClasses += ' font-bold';
  if (styles.includes('sup')) additionalClasses += ' align-super text-xs';
  if (styles.includes('sub')) additionalClasses += ' align-sub text-xs';
  
  return (
    <BaseComponent 
      metadata={metadata} 
      className={`text-component ${additionalClasses}`}
      testId="text-component"
    >
      {content}
    </BaseComponent>
  );
};

// 使用 React.memo 优化渲染性能
export default React.memo(TextComponent); 