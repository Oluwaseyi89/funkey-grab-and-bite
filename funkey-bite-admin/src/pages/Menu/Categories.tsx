import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { ArrowDown, ArrowUp, Eye, EyeOff, FolderTree, Pencil, Plus } from 'lucide-react';
import toast from 'react-hot-toast';
import { createCategory, getAdminCategories, getMenuItems, updateCategory } from '../../api/adminApi';
import PageHeader from '../../components/ui/PageHeader';
import Sheet from '../../components/ui/Sheet';
import type { MenuCategory } from '../../types';

type Draft = { id?: number; name: string; description: string; isActive: boolean };

// Menu categories: add, rename, reorder, and show or hide them on the
// customer menu. There's no delete, so items are never left without one;
// hiding a category takes it off the menu.
const Categories: React.FC = () => {
  const queryClient = useQueryClient();
  const [draft, setDraft] = useState<Draft | null>(null);
  const [error, setError] = useState('');

  const { data: categories = [], isLoading, isError, refetch } = useQuery({
    queryKey: ['admin-categories'],
    queryFn: async () => (await getAdminCategories()) ?? [],
  });
  const { data: items = [] } = useQuery({
    queryKey: ['menu-items', 'all-for-counts'],
    queryFn: async () => (await getMenuItems({ limit: 1000 })) ?? [],
  });
  const countFor = (id: number) => items.filter((item) => item.categoryId === id).length;
  const sorted = [...categories].sort((a, b) => a.displayOrder - b.displayOrder || a.id - b.id);

  const refresh = () => {
    queryClient.invalidateQueries({ queryKey: ['admin-categories'] });
    queryClient.invalidateQueries({ queryKey: ['categories'] });
  };

  const save = useMutation({
    mutationFn: async (value: Draft) => {
      const body = { name: value.name.trim(), description: value.description.trim(), isActive: value.isActive };
      if (value.id) return updateCategory(value.id, body);
      const nextOrder = sorted.length ? Math.max(...sorted.map((c) => c.displayOrder)) + 1 : 1;
      return createCategory({ ...body, displayOrder: nextOrder });
    },
    onSuccess: (_d, value) => {
      toast.success(value.id ? 'Category saved' : 'Category added');
      setDraft(null);
      refresh();
    },
    onError: (e: Error) => setError(e.message || 'The category wasn’t saved.'),
  });

  const toggle = useMutation({
    mutationFn: (category: MenuCategory) => updateCategory(category.id, { isActive: !category.isActive }),
    onSuccess: (_d, category) => {
      toast.success(category.isActive ? `${category.name} is hidden from the menu` : `${category.name} is on the menu`);
      refresh();
    },
    onError: () => toast.error('That didn’t save. Please try again.'),
  });

  // Swaps places with the neighbour, then renumbers so the order is clean.
  const move = useMutation({
    mutationFn: async ({ index, direction }: { index: number; direction: -1 | 1 }) => {
      const order = [...sorted];
      const target = index + direction;
      [order[index], order[target]] = [order[target], order[index]];
      const changes = order
        .map((category, position) => ({ category, position: position + 1 }))
        .filter(({ category, position }) => category.displayOrder !== position);
      await Promise.all(changes.map(({ category, position }) => updateCategory(category.id, { displayOrder: position })));
    },
    onSuccess: refresh,
    onError: () => {
      toast.error('The new order wasn’t saved. Please try again.');
      refresh();
    },
  });

  const openEditor = (category?: MenuCategory) => {
    setError('');
    setDraft(category ? { id: category.id, name: category.name, description: category.description ?? '', isActive: category.isActive } : { name: '', description: '', isActive: true });
  };

  return (
    <div>
      <PageHeader
        title="Categories"
        description="The sections of your menu, in the order customers see them."
        actions={
          <button type="button" className="btn-primary" onClick={() => openEditor()}>
            <Plus className="h-4 w-4" aria-hidden="true" /> Add category
          </button>
        }
      />

      {isLoading ? (
        <div className="space-y-2">{[0, 1, 2].map((i) => <div key={i} className="skeleton h-16 rounded-2xl" />)}</div>
      ) : isError ? (
        <div className="card p-6 text-center" role="alert">
          <p>The categories didn’t load.</p>
          <button type="button" className="btn-secondary mt-3" onClick={() => refetch()}>Try again</button>
        </div>
      ) : sorted.length === 0 ? (
        <div className="card p-10 text-center">
          <FolderTree className="mx-auto h-10 w-10 text-gray-300" aria-hidden="true" />
          <p className="mt-3 font-medium">No categories yet</p>
          <p className="mt-1 text-sm text-gray-500">Add one, e.g. “Rice dishes” or “Drinks”, then put menu items in it.</p>
          <button type="button" className="btn-primary mt-4" onClick={() => openEditor()}>Add your first category</button>
        </div>
      ) : (
        <ol className="card divide-y divide-gray-100 dark:divide-gray-800">
          {sorted.map((category, index) => (
            <li key={category.id} className={`flex items-center gap-2 p-3 ${category.isActive ? '' : 'bg-gray-50/70 dark:bg-gray-800/30'}`}>
              <div className="flex flex-col">
                <button type="button" className="icon-btn h-9 w-9" disabled={index === 0 || move.isPending} onClick={() => move.mutate({ index, direction: -1 })} aria-label={`Move ${category.name} up`}>
                  <ArrowUp className="h-4 w-4" aria-hidden="true" />
                </button>
                <button type="button" className="icon-btn h-9 w-9" disabled={index === sorted.length - 1 || move.isPending} onClick={() => move.mutate({ index, direction: 1 })} aria-label={`Move ${category.name} down`}>
                  <ArrowDown className="h-4 w-4" aria-hidden="true" />
                </button>
              </div>
              <div className="min-w-0 flex-1">
                <p className={`truncate font-semibold ${category.isActive ? 'text-gray-900 dark:text-white' : 'text-gray-500'}`}>
                  {category.name}
                  {!category.isActive && <span className="badge ml-2 bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-300">Hidden</span>}
                </p>
                <p className="truncate text-sm text-gray-500">
                  <Link to={`/menu?category=${category.id}`} className="hover:underline">{countFor(category.id)} items</Link>
                  {category.description ? ` · ${category.description}` : ''}
                </p>
              </div>
              <button type="button" className="icon-btn" onClick={() => toggle.mutate(category)} disabled={toggle.isPending} aria-label={category.isActive ? `Hide ${category.name} from the menu` : `Show ${category.name} on the menu`}>
                {category.isActive ? <Eye className="h-5 w-5" aria-hidden="true" /> : <EyeOff className="h-5 w-5" aria-hidden="true" />}
              </button>
              <button type="button" className="icon-btn" onClick={() => openEditor(category)} aria-label={`Edit ${category.name}`}>
                <Pencil className="h-5 w-5" aria-hidden="true" />
              </button>
            </li>
          ))}
        </ol>
      )}

      <Sheet
        open={!!draft}
        title={draft?.id ? 'Edit category' : 'Add category'}
        onClose={() => setDraft(null)}
        footer={
          <div className="grid grid-cols-2 gap-3">
            <button type="button" className="btn-secondary" onClick={() => setDraft(null)}>Cancel</button>
            <button type="submit" form="category-form" className="btn-primary" disabled={save.isPending}>
              {save.isPending ? 'Saving…' : 'Save'}
            </button>
          </div>
        }
      >
        {draft && (
          <form
            id="category-form"
            className="space-y-4"
            onSubmit={(e) => {
              e.preventDefault();
              if (!draft.name.trim()) {
                setError('Give the category a name.');
                return;
              }
              save.mutate(draft);
            }}
          >
            <div>
              <label htmlFor="category-name" className="form-label">Name</label>
              <input id="category-name" className="input-primary" value={draft.name} maxLength={100} onChange={(e) => setDraft({ ...draft, name: e.target.value })} aria-invalid={error && !draft.name.trim() ? 'true' : undefined} aria-describedby={error ? 'category-error' : undefined} autoFocus />
            </div>
            <div>
              <label htmlFor="category-description" className="form-label">Description <span className="font-normal text-gray-400">(optional)</span></label>
              <textarea id="category-description" className="input-primary resize-none" rows={2} maxLength={500} value={draft.description} onChange={(e) => setDraft({ ...draft, description: e.target.value })} placeholder="Shown under the category on the menu" />
            </div>
            <label className="flex min-h-[44px] items-center gap-3 text-sm">
              <input type="checkbox" className="h-5 w-5 accent-primary-600" checked={draft.isActive} onChange={(e) => setDraft({ ...draft, isActive: e.target.checked })} />
              Show on the customer menu
            </label>
            {error && <p id="category-error" className="text-sm text-red-600" role="alert">{error}</p>}
          </form>
        )}
      </Sheet>
    </div>
  );
};

export default Categories;
