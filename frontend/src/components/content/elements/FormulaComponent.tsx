'use client';

import React, { useEffect, useState } from 'react';
import BaseComponent from '../base/BaseComponent';
import { FormulaComponentProps } from '../types';
import { InlineMath, BlockMath } from 'react-katex';
import 'katex/dist/katex.min.css';

/**
 * 公式内容组件 - 用于渲染数学公式
 */
const FormulaComponent: React.FC<FormulaComponentProps> = ({
  content,
  metadata = {},
  options = {}
}) => {
  const [error, setError] = useState<string | null>(null);
  
  // 如果没有内容，返回空
  if (!content || !content.content) return null;
  
  const { content: formulaContent, display = 'inline' } = content;
  
  // 判断是否为块级公式
  const isBlock = display === 'block' || (metadata.style || []).includes('p');
  
  // 处理渲染错误
  useEffect(() => {
    // 重置错误状态
    setError(null);
  }, [formulaContent]);
  
  // 渲染公式
  const renderFormula = () => {
    try {
      return isBlock ? (
        <BlockMath math={formulaContent} />
      ) : (
        <InlineMath math={formulaContent} />
      );
    } catch (e) {
      const errorMessage = e instanceof Error ? e.message : '公式渲染错误';
      setError(errorMessage);
      console.error('公式渲染失败:', errorMessage);
      
      // 返回原始公式文本作为回退
      return (
        <code className="formula-error text-red-500">
          {formulaContent}
        </code>
      );
    }
  };
  
  return (
    <BaseComponent 
      metadata={metadata} 
      className={`formula-component ${isBlock ? 'formula-block' : 'formula-inline'}`}
      testId="formula-component"
    >
      <div className="formula-container">
        {error ? (
          <div className="formula-error text-red-500">
            <code>{formulaContent}</code>
            {options.showErrors && <p className="text-xs mt-1">{error}</p>}
          </div>
        ) : renderFormula()}
      </div>
    </BaseComponent>
  );
};

// 使用 React.memo 优化渲染性能
export default React.memo(FormulaComponent); 