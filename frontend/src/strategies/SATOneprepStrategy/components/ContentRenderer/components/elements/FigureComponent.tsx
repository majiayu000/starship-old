'use client';

import React from 'react';
import BaseComponent from '../base/BaseComponent';
import { ContentMetadata, ContentRendererOptions } from '../../types';
import ImageComponent from './ImageComponent';

interface FigureContent {
  content: string | { type: string; content: string };
  caption?: string;
}

interface FigureComponentProps {
  content: FigureContent;
  metadata?: ContentMetadata;
  options?: ContentRendererOptions;
}

/**
 * Component for rendering figures with captions
 */
const FigureComponent: React.FC<FigureComponentProps> = ({
  content,
  metadata = {},
  options = {}
}) => {
  if (!content) return null;
  
  // Extract the figure content - handle both string and object formats
  let figureContent = '';
  let contentType = 'svg';
  
  if (typeof content.content === 'string') {
    figureContent = content.content;
  } else if (content.content && typeof content.content === 'object') {
    figureContent = content.content.content;
    contentType = content.content.type || 'svg';
  }
  
  if (!figureContent) {
    console.warn('Empty figure content:', content);
    return null;
  }
  
  // Use ImageComponent for rendering the actual figure
  const imageContent = {
    content: figureContent,
    type: contentType as 'svg' | 'base64' | 'url',
    altText: content.caption || 'Figure'
  };
  
  return (
    <BaseComponent metadata={metadata} className="figure-component my-4">
      <figure className="w-full">
        <div className="figure-content">
          <ImageComponent 
            content={imageContent}
            metadata={{ ...metadata, style: [...(metadata.style || []), 'center'] }}
          />
        </div>
        
        {content.caption && (
          <figcaption className="text-sm text-gray-500 mt-2 text-center">
            {content.caption}
          </figcaption>
        )}
      </figure>
    </BaseComponent>
  );
};

export default FigureComponent; 