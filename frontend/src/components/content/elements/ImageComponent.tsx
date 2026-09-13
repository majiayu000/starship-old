'use client';

import React, { useState } from 'react';
import BaseComponent from '../base/BaseComponent';
import { ImageComponentProps } from '../types';
import { sanitizeSvg } from '@/lib/sanitizeHtml';

/**
 * 图片内容组件 - 用于渲染不同类型的图片
 */
const ImageComponent: React.FC<ImageComponentProps> = ({
  content,
  metadata = {},
  options = {}
}) => {
  const [isLoaded, setIsLoaded] = useState(false);
  const [hasError, setHasError] = useState(false);
  
  if (!content || !content.content) return null;
  
  const { content: imageContent, type, altText = '', dimensions } = content;
  
  // 处理SVG内容
  if (type === 'svg') {
    const sanitizedSvg = sanitizeSvg(imageContent);
    
    return (
      <BaseComponent 
        metadata={metadata} 
        className="image-component svg-image"
        testId="svg-component"
      >
        <div 
          className="svg-container" 
          dangerouslySetInnerHTML={{ __html: sanitizedSvg }}
          aria-label={altText || '图表内容'}
        />
      </BaseComponent>
    );
  }
  
  // 验证base64格式
  if (type === 'base64' && !isValidBase64(imageContent)) {
    console.warn('无效的base64图片格式:', imageContent.substring(0, 20) + '...');
    
    return (
      <BaseComponent 
        metadata={metadata} 
        className="image-component image-error"
        testId="image-error"
      >
        <div className="image-error-container text-red-500 text-sm">
          图片格式无效
        </div>
      </BaseComponent>
    );
  }
  
  // 确定图片是否应居中
  const isCentered = (metadata.style || []).includes('center');
  const containerStyle: React.CSSProperties = isCentered 
    ? { textAlign: 'center' } 
    : {};
  
  // 构建图片尺寸属性
  const imgProps: React.ImgHTMLAttributes<HTMLImageElement> = {
    src: imageContent,
    alt: altText || '题目图片',
    className: `max-w-full ${hasError ? 'hidden' : ''} ${isLoaded ? 'opacity-100' : 'opacity-0'}`,
    onLoad: () => setIsLoaded(true),
    onError: () => setHasError(true),
    style: { 
      transition: 'opacity 0.3s', 
      maxWidth: '100%',
      ...dimensions
    }
  };
  
  return (
    <BaseComponent 
      metadata={metadata} 
      className="image-component"
      testId="image-component"
    >
      <div className="image-container" style={containerStyle}>
        {!isLoaded && !hasError && (
          <div className="image-placeholder bg-gray-100 animate-pulse" 
            style={{ 
              width: dimensions?.width || '100%', 
              height: dimensions?.height || '200px',
              maxWidth: '100%'
            }}
          />
        )}
        
        {hasError ? (
          <div className="image-error-container p-4 border border-red-300 rounded text-center">
            <p className="text-red-500">图片加载失败</p>
            {options.showAltTextOnError && altText && (
              <p className="text-gray-600 mt-2">{altText}</p>
            )}
          </div>
        ) : (
          <img {...imgProps} />
        )}
      </div>
    </BaseComponent>
  );
};

/**
 * 验证Base64图片格式
 */
function isValidBase64(str: string): boolean {
  return typeof str === 'string' && (
    str.startsWith('data:image/') || 
    str.startsWith('data:application/octet-stream;base64,')
  );
}

// 使用 React.memo 优化渲染性能
export default React.memo(ImageComponent); 