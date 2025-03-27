'use client';

import React from 'react';
import BaseComponent from '../base/BaseComponent';
import { ContentMetadata, ContentRendererOptions, StringContent } from '../../types';

interface TextComponentProps {
  content: string;
  metadata?: ContentMetadata;
  options?: ContentRendererOptions;
}

/**
 * Component for rendering text content
 */
const TextComponent: React.FC<TextComponentProps> = ({
  content,
  metadata = {},
  options = {}
}) => {
  if (!content) return null;
  
  // Apply additional class names based on style metadata
  let additionalClasses = '';
  
  const styles = metadata.style || [];
  if (styles.includes('em') || styles.includes('i')) additionalClasses += ' italic';
  if (styles.includes('strong') || styles.includes('b')) additionalClasses += ' font-bold';
  if (styles.includes('sup')) additionalClasses += ' align-super text-xs';
  
  return (
    <BaseComponent metadata={metadata} className={additionalClasses}>
      {content}
    </BaseComponent>
  );
};

export default TextComponent; 