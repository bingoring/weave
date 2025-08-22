export interface WeaveContent {
  type: string;
  data: Record<string, any>;
}

export interface Weave {
  id: string;
  user_id: string;
  channel_id: string;
  title: string;
  cover_image?: string;
  content: WeaveContent;
  version: number;
  parent_weave_id?: string;
  is_published: boolean;
  is_featured: boolean;
  view_count: number;
  like_count: number;
  fork_count: number;
  comment_count: number;
  created_at: string;
  updated_at: string;
  user?: User;
  channel?: Channel;
}

export interface Channel {
  id: string;
  name: string;
  slug: string;
  description?: string;
  cover_image?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface WeaveVersion {
  id: string;
  weave_id: string;
  version: number;
  title: string;
  content: WeaveContent;
  change_log?: string;
  created_at: string;
}

export interface CreateWeaveRequest {
  title: string;
  channel_id: string;
  content: WeaveContent;
  cover_image?: string;
}

export interface UpdateWeaveRequest {
  title?: string;
  content?: WeaveContent;
  cover_image?: string;
  change_log?: string;
}

export interface WeaveFilters {
  channel_id?: string;
  user_id?: string;
  is_published?: boolean;
  is_featured?: boolean;
  search?: string;
  tags?: string[];
}

export interface PaginatedWeavesResponse {
  weaves: Weave[];
  page: number;
  limit: number;
  total: number;
  has_next: boolean;
  has_prev: boolean;
}

// Recipe-specific content structure
export interface RecipeContent {
  ingredients: Array<{
    name: string;
    amount: string;
    unit: string;
  }>;
  instructions: Array<{
    step: number;
    description: string;
    duration?: number;
  }>;
  prep_time?: number;
  cook_time?: number;
  servings?: number;
  difficulty?: 'easy' | 'medium' | 'hard';
  cuisine?: string;
  dietary_restrictions?: string[];
}

// Travel plan content structure
export interface TravelPlanContent {
  destination: string;
  duration: number;
  itinerary: Array<{
    day: number;
    activities: Array<{
      time: string;
      activity: string;
      location: string;
      duration?: number;
      cost?: number;
      notes?: string;
    }>;
  }>;
  budget?: {
    accommodation: number;
    transport: number;
    food: number;
    activities: number;
    misc: number;
  };
  tips?: string[];
}

// Workout routine content structure
export interface WorkoutContent {
  type: 'strength' | 'cardio' | 'yoga' | 'mixed';
  duration: number;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  equipment: string[];
  exercises: Array<{
    name: string;
    sets?: number;
    reps?: number;
    duration?: number;
    weight?: number;
    rest?: number;
    instructions: string;
    target_muscles: string[];
  }>;
  warmup?: string[];
  cooldown?: string[];
}