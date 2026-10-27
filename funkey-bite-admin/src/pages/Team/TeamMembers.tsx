import React, { useState } from 'react';
import { Plus, Pencil, Trash2, X, UserSquare, Upload } from 'lucide-react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import {
  createTeamMember,
  deleteTeamMember,
  getTeamMembers,
  updateTeamMember,
  uploadTeamPhoto,
} from '../../api/adminApi';
import type { TeamMember, TeamMemberInput } from '../../types';

const MAX_PHOTO_BYTES = 2 * 1024 * 1024;
const inputClass = 'w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white';
const labelClass = 'block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1';

const initials = (name: string) =>
  name.split(/\s+/).filter(Boolean).slice(0, 2).map((part) => part[0].toUpperCase()).join('');

const Avatar: React.FC<{ name: string; photoUrl?: string; size?: string }> = ({ name, photoUrl, size = 'h-16 w-16' }) =>
  photoUrl ? (
    <img src={photoUrl} alt={name} className={`${size} rounded-full object-cover`} />
  ) : (
    <div className={`${size} rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center text-primary-700 dark:text-primary-300 font-semibold`}>
      {initials(name) || '?'}
    </div>
  );

interface MemberFormProps {
  member: TeamMember | null;
  nextOrder: number;
  onClose: () => void;
}

const MemberForm: React.FC<MemberFormProps> = ({ member, nextOrder, onClose }) => {
  const queryClient = useQueryClient();
  const [form, setForm] = useState<TeamMemberInput>({
    name: member?.name ?? '',
    role: member?.role ?? '',
    bio: member?.bio ?? '',
    photoUrl: member?.photoUrl,
    displayOrder: member?.displayOrder ?? nextOrder,
    isActive: member?.isActive ?? true,
  });
  const [isUploading, setIsUploading] = useState(false);

  const set = <K extends keyof TeamMemberInput>(key: K, value: TeamMemberInput[K]) =>
    setForm((prev) => ({ ...prev, [key]: value }));

  const error =
    form.name.trim().length < 2 ? 'Enter the person\'s name'
      : form.role.trim().length < 2 ? 'Enter their role'
        : '';

  const handlePhoto = async (file?: File) => {
    if (!file) return;
    if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
      toast.error('Use a JPEG, PNG or WebP photo');
      return;
    }
    if (file.size > MAX_PHOTO_BYTES) {
      toast.error('Photo must be 2 MB or smaller');
      return;
    }
    setIsUploading(true);
    try {
      const uploaded = await uploadTeamPhoto(file);
      set('photoUrl', uploaded.url || uploaded.key);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to upload photo');
    } finally {
      setIsUploading(false);
    }
  };

  const mutation = useMutation({
    mutationFn: () => (member ? updateTeamMember(member.id, form) : createTeamMember(form)),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['team-members'] });
      toast.success(member ? 'Team member updated' : 'Team member added');
      onClose();
    },
    onError: (err: Error) => toast.error(err.message || 'Failed to save team member'),
  });

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50" role="presentation">
      <form
        role="dialog"
        aria-modal="true"
        aria-labelledby="member-form-title"
        onSubmit={(e) => { e.preventDefault(); if (!error) mutation.mutate(); }}
        className="bg-white dark:bg-gray-800 rounded-xl w-full max-w-lg p-6 space-y-4"
      >
        <div className="flex items-center justify-between">
          <h2 id="member-form-title" className="text-lg font-semibold text-gray-900 dark:text-white">
            {member ? `Edit ${member.name}` : 'Add team member'}
          </h2>
          <button type="button" onClick={onClose} aria-label="Close" className="p-1 text-gray-400 hover:text-gray-600"><X className="h-5 w-5" /></button>
        </div>

        <div className="flex items-center gap-4">
          <Avatar name={form.name} photoUrl={form.photoUrl} />
          <div className="flex flex-col gap-2">
            <label className="inline-flex items-center gap-2 px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700">
              <Upload className="h-4 w-4" />
              {isUploading ? 'Uploading…' : form.photoUrl ? 'Change photo' : 'Upload photo'}
              <input type="file" accept="image/jpeg,image/png,image/webp" className="sr-only" disabled={isUploading} onChange={(e) => handlePhoto(e.target.files?.[0])} />
            </label>
            {form.photoUrl && (
              <button type="button" className="text-sm text-red-600 text-left" onClick={() => set('photoUrl', undefined)}>Remove photo</button>
            )}
          </div>
        </div>

        <div>
          <label htmlFor="member-name" className={labelClass}>Name *</label>
          <input id="member-name" className={inputClass} maxLength={120} value={form.name} onChange={(e) => set('name', e.target.value)} />
        </div>
        <div>
          <label htmlFor="member-role" className={labelClass}>Role *</label>
          <input id="member-role" className={inputClass} maxLength={120} value={form.role} onChange={(e) => set('role', e.target.value)} placeholder="e.g. Head Chef" />
        </div>
        <div>
          <label htmlFor="member-bio" className={labelClass}>Short bio</label>
          <textarea id="member-bio" className={inputClass} rows={3} maxLength={500} value={form.bio} onChange={(e) => set('bio', e.target.value)} />
          <p className="mt-1 text-xs text-gray-500">{form.bio.length}/500</p>
        </div>
        <div className="grid grid-cols-2 gap-4 items-end">
          <div>
            <label htmlFor="member-order" className={labelClass}>Display order</label>
            <input id="member-order" type="number" min={0} max={1000} className={inputClass} value={form.displayOrder} onChange={(e) => set('displayOrder', Number(e.target.value))} />
          </div>
          <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 pb-2">
            <input type="checkbox" checked={form.isActive} onChange={(e) => set('isActive', e.target.checked)} />
            Show on the website
          </label>
        </div>

        {error && <p className="text-sm text-red-600">{error}</p>}

        <div className="flex justify-end gap-3">
          <button type="button" onClick={onClose} className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700">Cancel</button>
          <button type="submit" disabled={!!error || isUploading || mutation.isPending} className="px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg disabled:opacity-50">
            {mutation.isPending ? 'Saving…' : 'Save'}
          </button>
        </div>
      </form>
    </div>
  );
};

const TeamMembers: React.FC = () => {
  const queryClient = useQueryClient();
  const [editing, setEditing] = useState<TeamMember | null | 'new'>(null);
  const [confirmDelete, setConfirmDelete] = useState<TeamMember | null>(null);

  const { data: members = [], isLoading } = useQuery({
    queryKey: ['team-members'],
    queryFn: getTeamMembers,
  });

  const deletion = useMutation({
    mutationFn: (id: number) => deleteTeamMember(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['team-members'] });
      toast.success('Team member removed');
      setConfirmDelete(null);
    },
    onError: (err: Error) => toast.error(err.message || 'Failed to remove team member'),
  });

  const nextOrder = members.reduce((max, m) => Math.max(max, m.displayOrder + 1), 0);

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Team</h1>
          <p className="text-gray-600 dark:text-gray-400">
            People shown in the "Meet Our Team" section of the website's contact page. The section is hidden until you add someone.
          </p>
        </div>
        <button onClick={() => setEditing('new')} className="flex items-center gap-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg">
          <Plus className="h-4 w-4" /> Add Member
        </button>
      </div>

      {isLoading ? (
        <div className="p-8 text-center text-gray-500">Loading team…</div>
      ) : members.length === 0 ? (
        <div className="bg-white dark:bg-gray-800 rounded-xl p-8 border border-gray-200 dark:border-gray-700 text-center">
          <UserSquare className="h-12 w-12 text-gray-400 mx-auto mb-3" />
          <p className="text-gray-700 dark:text-gray-300">No team members yet.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {members.map((member) => (
            <article key={member.id} className={`bg-white dark:bg-gray-800 rounded-xl p-5 border ${member.isActive ? 'border-gray-200 dark:border-gray-700' : 'border-dashed border-gray-300 dark:border-gray-600 opacity-70'}`}>
              <div className="flex items-start gap-4">
                <Avatar name={member.name} photoUrl={member.photoUrl} />
                <div className="flex-1 min-w-0">
                  <h2 className="font-semibold text-gray-900 dark:text-white truncate">{member.name}</h2>
                  <p className="text-sm text-primary-600 dark:text-primary-400">{member.role}</p>
                  <span className={`inline-block mt-1 px-2 py-0.5 rounded-full text-xs ${member.isActive ? 'bg-green-100 text-green-800 dark:bg-green-900/40 dark:text-green-200' : 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300'}`}>
                    {member.isActive ? 'Shown' : 'Hidden'}
                  </span>
                </div>
              </div>
              {member.bio && <p className="mt-3 text-sm text-gray-600 dark:text-gray-300 line-clamp-3">{member.bio}</p>}
              <div className="mt-4 flex gap-2">
                <button onClick={() => setEditing(member)} className="flex items-center gap-1 px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700">
                  <Pencil className="h-4 w-4" /> Edit
                </button>
                <button onClick={() => setConfirmDelete(member)} className="flex items-center gap-1 px-3 py-1.5 text-sm border border-red-300 text-red-700 dark:text-red-300 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20">
                  <Trash2 className="h-4 w-4" /> Remove
                </button>
              </div>
            </article>
          ))}
        </div>
      )}

      {editing && (
        <MemberForm
          key={editing === 'new' ? 'new' : editing.id}
          member={editing === 'new' ? null : editing}
          nextOrder={nextOrder}
          onClose={() => setEditing(null)}
        />
      )}

      {confirmDelete && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50" role="presentation">
          <div role="alertdialog" aria-modal="true" aria-labelledby="remove-member-title" className="bg-white dark:bg-gray-800 rounded-xl w-full max-w-sm p-6 space-y-4">
            <h2 id="remove-member-title" className="text-lg font-semibold text-gray-900 dark:text-white">Remove {confirmDelete.name}?</h2>
            <p className="text-sm text-gray-600 dark:text-gray-400">They will be removed from the website. To keep them but stop showing them, edit and untick "Show on the website" instead.</p>
            <div className="flex justify-end gap-3">
              <button onClick={() => setConfirmDelete(null)} className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg">Cancel</button>
              <button onClick={() => deletion.mutate(confirmDelete.id)} disabled={deletion.isPending} className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg disabled:opacity-50">Remove</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default TeamMembers;
