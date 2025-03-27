'use client';

import React from 'react';
import BaseComponent from '../base/BaseComponent';
import { ContentMetadata, ContentRendererOptions } from '../../types';

interface BlankContent {
  width?: string | number;
  style?: string;
}

interface BlankComponentProps {
  content?: BlankContent;
  metadata?: ContentMetadata;
  options?: ContentRendererOptions;
}

/**
 * Component for rendering blank spaces for fill-in questions
 */
const BlankComponent: React.FC<BlankComponentProps> = ({
  content,
  metadata = {},
  options = {}
}) => {
  // Define default width and custom styles
  let width = 'w-20'; // Default width class
  let customStyle: React.CSSProperties = {};
  
  // Apply custom width if provided
  if (content?.width) {
    width = ''; // Remove default width class when custom width is used
    customStyle.width = typeof content.width === 'number' 
      ? `${content.width}px` 
      : content.width;
  }
  
  // Apply custom style if provided
  if (content?.style) {
    try {
      const parsedStyle = JSON.parse(content.style);
      customStyle = { ...customStyle, ...parsedStyle };
    } catch (e) {
      console.warn('Invalid custom style JSON for blank component:', content.style);
    }
  }
  
  return (
    <BaseComponent metadata={metadata} className="inline-flex align-baseline">
      <span 
        className={`${width} border-b-2 border-black mx-1`} 
        style={{ 
          height: '0.15em', 
          marginBottom: '0.2em',
          ...customStyle
        }}
      />
    </BaseComponent>
  );
};

export default BlankComponent; 