migrations/init.sql
-- Users table (extends Keycloak user info)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    keycloak_user_id VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(100),
    bio TEXT,
    avatar_url VARCHAR(500),
    location VARCHAR(100),
    website VARCHAR(500),
    musical_genres TEXT[], -- Array of genres
    instruments TEXT[], -- Array of instruments
    experience_level VARCHAR(20), -- beginner, intermediate, advanced, professional
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User social profiles
CREATE TABLE IF NOT EXISTS user_social_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL, -- spotify, soundcloud, youtube, etc.
    profile_url VARCHAR(500) NOT NULL,
    username VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Organizations (like GitLab groups/organizations)
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    avatar_url VARCHAR(500),
    website VARCHAR(500),
    visibility VARCHAR(20) DEFAULT 'public', -- public, internal, private
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Organization members
CREATE TABLE IF NOT EXISTS organization_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL, -- owner, maintainer, developer, guest
    access_level INTEGER NOT NULL, -- 50=owner, 40=maintainer, 30=developer, 10=guest
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, user_id)
);

-- Music projects
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    organization_id UUID REFERENCES organizations(id),
    created_by UUID REFERENCES users(id) NOT NULL,
    visibility VARCHAR(20) DEFAULT 'public',
    
    -- Music specific fields
    genre VARCHAR(50),
    tempo INTEGER, -- BPM
    key_signature VARCHAR(10), -- C major, A minor, etc.
    time_signature VARCHAR(10), -- 4/4, 3/4, etc.
    
    -- Project status
    status VARCHAR(20) DEFAULT 'active', -- active, completed, archived
    license VARCHAR(50), -- Creative Commons, etc.
    
    -- Statistics
    stars_count INTEGER DEFAULT 0,
    forks_count INTEGER DEFAULT 0,
    collaborators_count INTEGER DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, slug)
);

-- Project members/collaborators
CREATE TABLE IF NOT EXISTS project_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL, -- producer, musician, mixer, vocalist, etc.
    access_level INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, user_id)
);

-- Project stars (like GitLab stars)
CREATE TABLE IF NOT EXISTS project_stars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, user_id)
);

-- Audio tracks/stems
CREATE TABLE IF NOT EXISTS audio_tracks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    file_path VARCHAR(1000) NOT NULL,
    file_size BIGINT,
    duration_seconds FLOAT,
    format VARCHAR(10), -- wav, mp3, flac, etc.
    sample_rate INTEGER,
    bit_depth INTEGER,
    channels INTEGER, -- mono=1, stereo=2, etc.
    
    -- Track metadata
    track_type VARCHAR(20), -- vocals, drums, bass, guitar, synth, mix, master
    instrument VARCHAR(50),
    
    -- Version control
    version_number INTEGER DEFAULT 1,
    parent_track_id UUID REFERENCES audio_tracks(id),
    
    created_by UUID REFERENCES users(id) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Track versions (like commits)
CREATE TABLE IF NOT EXISTS track_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    track_id UUID REFERENCES audio_tracks(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    commit_message TEXT,
    changes_description TEXT,
    file_path VARCHAR(1000) NOT NULL,
    created_by UUID REFERENCES users(id) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Collaboration requests (like GitLab merge requests)
CREATE TABLE IF NOT EXISTS collaboration_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Source and target
    source_branch VARCHAR(100), -- feature branch
    target_branch VARCHAR(100) DEFAULT 'main',
    
    -- Status
    status VARCHAR(20) DEFAULT 'open', -- open, merged, closed, draft
    
    created_by UUID REFERENCES users(id) NOT NULL,
    assigned_to UUID REFERENCES users(id),
    merged_by UUID REFERENCES users(id),
    merged_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Reviews and feedback
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    collaboration_request_id UUID REFERENCES collaboration_requests(id) ON DELETE CASCADE,
    reviewer_id UUID REFERENCES users(id) NOT NULL,
    status VARCHAR(20) NOT NULL, -- approved, changes_requested, commented
    feedback TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Review comments (on specific tracks/timestamps)
CREATE TABLE IF NOT EXISTS review_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id UUID REFERENCES reviews(id) ON DELETE CASCADE,
    track_id UUID REFERENCES audio_tracks(id),
    timestamp_seconds FLOAT, -- Comment at specific time in track
    comment TEXT NOT NULL,
    created_by UUID REFERENCES users(id) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Issues (bugs, feature requests, feedback)
CREATE TABLE IF NOT EXISTS issues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    issue_type VARCHAR(20) NOT NULL, -- bug, feedback, feature_request, question
    status VARCHAR(20) DEFAULT 'open', -- open, closed, in_progress
    priority VARCHAR(10) DEFAULT 'medium', -- low, medium, high, critical
    
    created_by UUID REFERENCES users(id) NOT NULL,
    assigned_to UUID REFERENCES users(id),
    closed_by UUID REFERENCES users(id),
    closed_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Issue comments
CREATE TABLE IF NOT EXISTS issue_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issue_id UUID REFERENCES issues(id) ON DELETE CASCADE,
    comment TEXT NOT NULL,
    created_by UUID REFERENCES users(id) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Activity feed (like GitLab activity)
CREATE TABLE IF NOT EXISTS activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) NOT NULL,
    project_id UUID REFERENCES projects(id),
    organization_id UUID REFERENCES organizations(id),
    
    activity_type VARCHAR(50) NOT NULL, -- track_uploaded, project_created, collaboration_request, etc.
    title VARCHAR(255) NOT NULL,
    description TEXT,
    metadata JSONB, -- Flexible data for different activity types
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Notifications
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_id UUID REFERENCES users(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id),
    
    notification_type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT,
    
    is_read BOOLEAN DEFAULT false,
    is_emailed BOOLEAN DEFAULT false,
    
    related_entity_type VARCHAR(50), -- project, collaboration_request, issue
    related_entity_id UUID,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- File storage tracking
CREATE TABLE IF NOT EXISTS file_storage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_path VARCHAR(1000) UNIQUE NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    checksum VARCHAR(255), -- For integrity verification
    
    -- Storage metadata
    storage_provider VARCHAR(50) DEFAULT 'local', -- local, s3, gcs, etc.
    storage_bucket VARCHAR(100),
    
    uploaded_by UUID REFERENCES users(id) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Performance indexes
CREATE INDEX IF NOT EXISTS idx_users_keycloak_id ON users(keycloak_user_id);
CREATE INDEX IF NOT EXISTS idx_projects_organization ON projects(organization_id);
CREATE INDEX IF NOT EXISTS idx_projects_created_by ON projects(created_by);
CREATE INDEX IF NOT EXISTS idx_audio_tracks_project ON audio_tracks(project_id);
CREATE INDEX IF NOT EXISTS idx_activities_user_created ON activities(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_recipient_read ON notifications(recipient_id, is_read);