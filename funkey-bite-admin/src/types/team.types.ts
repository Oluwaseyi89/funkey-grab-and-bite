// A person shown in the team section of the website's contact page.
export interface TeamMember {
  id: number;
  name: string;
  role: string;
  bio: string;
  photoUrl?: string;
  displayOrder: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface TeamMemberInput {
  name: string;
  role: string;
  bio: string;
  photoUrl?: string;
  displayOrder: number;
  isActive: boolean;
}
