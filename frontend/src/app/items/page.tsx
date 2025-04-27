'use client';

import React, { useState } from 'react';
import { GenericItem } from '@/models/generic-item';
import ItemList from '@/components/items/ItemList';
import ItemDetail from '@/components/items/ItemDetail';
import CreateItemForm from '@/components/items/CreateItemForm';
import { Button } from '@/components/ui/button';
import { PlusIcon } from 'lucide-react';

export default function ItemsPage() {
  // State
  const [selectedItem, setSelectedItem] = useState<GenericItem | null>(null);
  const [showCreateForm, setShowCreateForm] = useState<boolean>(false);

  // Handle item selection
  const handleItemSelect = (item: GenericItem) => {
    setSelectedItem(item);
    setShowCreateForm(false);
  };

  // Handle item update
  const handleItemUpdated = (updatedItem: GenericItem) => {
    setSelectedItem(updatedItem);
  };

  // Handle item creation
  const handleItemCreated = () => {
    setShowCreateForm(false);
    setSelectedItem(null);
  };

  // Handle create button click
  const handleCreateClick = () => {
    setSelectedItem(null);
    setShowCreateForm(true);
  };

  // Handle close
  const handleClose = () => {
    setSelectedItem(null);
    setShowCreateForm(false);
  };

  return (
    <div className="container mx-auto py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Items</h1>
        <Button onClick={handleCreateClick}>
          <PlusIcon className="h-4 w-4 mr-2" />
          Create Item
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className={`${selectedItem || showCreateForm ? 'hidden md:block' : ''} md:col-span-2`}>
          <ItemList onItemSelect={handleItemSelect} />
        </div>

        {selectedItem && (
          <div className="md:col-span-1">
            <ItemDetail 
              item={selectedItem} 
              onItemUpdated={handleItemUpdated}
              onClose={handleClose}
            />
          </div>
        )}

        {showCreateForm && (
          <div className="md:col-span-1">
            <CreateItemForm 
              onItemCreated={handleItemCreated}
              onCancel={handleClose}
            />
          </div>
        )}
      </div>
    </div>
  );
}
