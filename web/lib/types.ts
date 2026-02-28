// Auth
export interface RegisterRequest {
  email: string;
  password: string;
  fullName: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface TokenResponse {
  accessToken: string;
  refreshToken: string;
}

export interface UserResponse {
  id: string;
  email: string;
  fullName: string;
  avatar: string;
  createdAt: string;
}

// Workspace
export interface CreateWorkspaceRequest {
  name: string;
  slug: string;
}

export interface WorkspaceResponse {
  id: string;
  name: string;
  slug: string;
  ownerId: string;
  createdAt: string;
  updatedAt?: string;
}

export interface AddMembersRequest {
  userIds: string[];
  role: "Admin" | "Editor" | "Viewer";
}

export interface WorkspaceMemberResponse {
  id: string;
  userId: string;
  role: string;
  joinedAt: string;
}

// Document
export interface CreateDocRequest {
  workspaceId: string;
  parentId?: string;
  title: string;
  emoji?: string;
}

export interface DocumentResponse {
  id: string;
  title: string;
  emoji: string;
  parentId: string;
  authorId: string;
  createdAt: string;
}

export interface DocTreeItem {
  id: string;
  title: string;
  parentId: string;
  children: DocTreeItem[];
}

// API
export interface ApiResponse<T> {
  data: T;
  message?: string;
}

export interface ApiError {
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
}
