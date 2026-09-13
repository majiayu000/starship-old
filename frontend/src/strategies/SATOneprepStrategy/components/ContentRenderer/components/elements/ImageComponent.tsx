'use client';

import React from 'react';
import BaseComponent from '../base/BaseComponent';
import { ContentMetadata, ContentRendererOptions } from '../../types';
import { sanitizeSvg } from '@/lib/sanitizeHtml';

interface ImageContent {
  content: string;
  type: 'base64' | 'url' | 'svg';
  altText?: string;
  dimensions?: {
    width?: number | string;
    height?: number | string;
  };
}

interface ImageComponentProps {
  content: ImageContent;
  metadata?: ContentMetadata;
  options?: ContentRendererOptions;
}

/**
 * Component for rendering different types of images (base64, URL, SVG)
 */
const ImageComponent: React.FC<ImageComponentProps> = ({
  content,
  metadata = {},
  options = {}
}) => {
  if (!content || !content.content) return null;
  
  const { content: imageContent, type, altText = '', dimensions } = content;
  
  // Handle SVG content
  if (type === 'svg') {
    return (
      <BaseComponent metadata={metadata}>
        <div 
          className="svg-container" 
          dangerouslySetInnerHTML={{ __html: sanitizeSvg(imageContent) }}
        />
      </BaseComponent>
    );
  }
  
  // Verify base64 format if needed
  if (type === 'base64' && !imageContent.startsWith('data:')) {
    console.warn('Invalid base64 image format:', imageContent.substring(0, 20) + '...');
    return null;
  }
  
  // Determine if image should be centered based on metadata
  const isCentered = (metadata.style || []).includes('center');
  const containerStyle = isCentered ? { textAlign: 'center' as const } : {};
  
  return (
    <BaseComponent metadata={metadata}>
      <div className="image-container" style={containerStyle}>
        <img 
          src={imageContent} 
          alt={altText} 
          className="max-w-full" 
          width={dimensions?.width} 
          height={dimensions?.height}
        />
      </div>
    </BaseComponent>
  );
};

export default ImageComponent; 