'use client';

import React, { useState } from 'react';
import { GenericItem, StatusUpdate } from '@/models/generic-item';
import { genericApi } from '@/api/generic-api';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { toast } from 'sonner';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

interface ItemDetailProps {
  item: GenericItem;
  onItemUpdated?: (updatedItem: GenericItem) => void;
  onClose?: () => void;
}

export default function ItemDetail({ item, onItemUpdated, onClose }: ItemDetailProps) {
  // State
  const [loading, setLoading] = useState<boolean>(false);
  const [showEditDialog, setShowEditDialog] = useState<boolean>(false);
  const [showStatusDialog, setShowStatusDialog] = useState<boolean>(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState<boolean>(false);

  // Form state
  const [formData, setFormData] = useState<{
    name: string;
    description: string;
  }>({
    name: item.name,
    description: item.description,
  });

  const [statusUpdate, setStatusUpdate] = useState<StatusUpdate>({
    status: item.status,
    comment: '',
  });

  // Handle form input change
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  // Handle status change
  const handleStatusChange = (status: string) => {
    setStatusUpdate({
      ...statusUpdate,
      status,
    });
  };

  // Handle comment change
  const handleCommentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setStatusUpdate({
      ...statusUpdate,
      comment: e.target.value,
    });
  };

  // Update item
  const handleUpdate = async () => {
    setLoading(true);
    try {
      const updatedItem = await genericApi.update(item.id, {
        name: formData.name,
        description: formData.description,
        status: item.status,
      });

      toast.success('Item updated successfully');
      setShowEditDialog(false);

      if (onItemUpdated) {
        onItemUpdated(updatedItem);
      }
    } catch (error) {
      console.error('Error updating item:', error);
      toast.error('Failed to update item');
    } finally {
      setLoading(false);
    }
  };

  // Update item status
  const handleUpdateStatus = async () => {
    setLoading(true);
    try {
      await genericApi.updateStatus(item.id, statusUpdate);

      toast.success('Status updated successfully');
      setShowStatusDialog(false);

      if (onItemUpdated) {
        // Fetch the updated item
        const updatedItem = await genericApi.getById(item.id);
        onItemUpdated(updatedItem);
      }
    } catch (error) {
      console.error('Error updating status:', error);
      toast.error('Failed to update status');
    } finally {
      setLoading(false);
    }
  };

  // Delete item
  const handleDelete = async () => {
    setLoading(true);
    try {
      await genericApi.delete(item.id);

      toast.success('Item deleted successfully');
      setShowDeleteDialog(false);

      if (onClose) {
        onClose();
      }
    } catch (error) {
      console.error('Error deleting item:', error);
      toast.error('Failed to delete item');
    } finally {
      setLoading(false);
    }
  };

  // Render status badge
  const renderStatusBadge = (status: string) => {
    let variant: 'default' | 'secondary' | 'destructive' | 'outline' = 'default';

    switch (status) {
      case 'active':
        variant = 'default';
        break;
      case 'draft':
        variant = 'secondary';
        break;
      case 'inactive':
        variant = 'outline';
        break;
      case 'archived':
        variant = 'destructive';
        break;
    }

    return <Badge variant={variant}>{status}</Badge>;
  };

  // Render the edit form
  if (showEditDialog) {
    return (
      <Card className="w-full">
        <CardHeader>
          <CardTitle>Edit Item</CardTitle>
          <CardDescription>Make changes to the item details.</CardDescription>
        </CardHeader>

        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              name="name"
              value={formData.name}
              onChange={handleInputChange}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              name="description"
              value={formData.description}
              onChange={handleInputChange}
              rows={4}
            />
          </div>
        </CardContent>

        <CardFooter className="flex justify-between">
          <Button variant="outline" onClick={() => setShowEditDialog(false)}>
            Cancel
          </Button>
          <Button onClick={handleUpdate} disabled={loading}>
            {loading ? 'Saving...' : 'Save Changes'}
          </Button>
        </CardFooter>
      </Card>
    );
  }

  // Render the status update form
  if (showStatusDialog) {
    return (
      <Card className="w-full">
        <CardHeader>
          <CardTitle>Update Status</CardTitle>
          <CardDescription>Change the status of this item.</CardDescription>
        </CardHeader>

        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="status">Status</Label>
            <Select
              value={statusUpdate.status}
              onValueChange={handleStatusChange}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="draft">Draft</SelectItem>
                <SelectItem value="active">Active</SelectItem>
                <SelectItem value="inactive">Inactive</SelectItem>
                <SelectItem value="archived">Archived</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="comment">Comment (Optional)</Label>
            <Textarea
              id="comment"
              value={statusUpdate.comment}
              onChange={handleCommentChange}
              placeholder="Add a comment about this status change"
              rows={3}
            />
          </div>
        </CardContent>

        <CardFooter className="flex justify-between">
          <Button variant="outline" onClick={() => setShowStatusDialog(false)}>
            Cancel
          </Button>
          <Button onClick={handleUpdateStatus} disabled={loading}>
            {loading ? 'Updating...' : 'Update Status'}
          </Button>
        </CardFooter>
      </Card>
    );
  }

  // Render the delete confirmation
  if (showDeleteDialog) {
    return (
      <Card className="w-full">
        <CardHeader>
          <CardTitle>Delete Item</CardTitle>
          <CardDescription>
            Are you sure you want to delete this item? This action cannot be undone.
          </CardDescription>
        </CardHeader>

        <CardFooter className="flex justify-between">
          <Button variant="outline" onClick={() => setShowDeleteDialog(false)}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={handleDelete} disabled={loading}>
            {loading ? 'Deleting...' : 'Delete'}
          </Button>
        </CardFooter>
      </Card>
    );
  }

  // Render the item details
  return (
    <Card className="w-full">
      <CardHeader>
        <div className="flex justify-between items-start">
          <div>
            <CardTitle>{item.name}</CardTitle>
            <CardDescription>ID: {item.id}</CardDescription>
          </div>
          {renderStatusBadge(item.status)}
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        <div>
          <h3 className="text-sm font-medium">Description</h3>
          <p className="mt-1">{item.description || 'No description provided'}</p>
        </div>

        <Separator />

        <div className="grid grid-cols-2 gap-4">
          <div>
            <h3 className="text-sm font-medium">Created</h3>
            <p className="mt-1">{new Date(item.createdAt).toLocaleString()}</p>
          </div>
          <div>
            <h3 className="text-sm font-medium">Last Updated</h3>
            <p className="mt-1">{new Date(item.updatedAt).toLocaleString()}</p>
          </div>
        </div>
      </CardContent>

      <CardFooter className="flex justify-between">
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>

        <div className="flex gap-2">
          <Button variant="outline" onClick={() => setShowEditDialog(true)}>
            Edit
          </Button>
          <Button onClick={() => setShowStatusDialog(true)}>
            Change Status
          </Button>
          <Button variant="destructive" onClick={() => setShowDeleteDialog(true)}>
            Delete
          </Button>
        </div>
      </CardFooter>
    </Card>
  );
}
