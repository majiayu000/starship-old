'use client';

import React from 'react';
import { ContentMetadata } from '../../types';

interface BaseComponentProps {
  metadata?: ContentMetadata;
  children: React.ReactNode;
  className?: string;
}

/**
 * Base component that wraps all content elements with consistent styling
 */
const BaseComponent: React.FC<BaseComponentProps> = ({ 
  metadata = {}, 
  children,
  className = ''
}) => {
  const { style = [], content_type } = metadata;
  
  // Generate class names from style array
  const styleClasses = style.map(s => `style-${s}`).join(' ');
  const contentTypeClass = content_type ? `content-type-${content_type}` : '';
  
  // Check if this component should be rendered as a block
  const isBlock = style.includes('p');
  
  return (
    <div 
      className={`sat-component ${styleClasses} ${contentTypeClass} ${isBlock ? 'block-component' : 'inline-component'} ${className}`}
    >
      {children}
    </div>
  );
};

export default BaseComponent; 